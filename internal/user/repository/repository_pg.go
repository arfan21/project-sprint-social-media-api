package userrepo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/arfan21/project-sprint-social-media-api/internal/entity"
	"github.com/arfan21/project-sprint-social-media-api/internal/model"
	"github.com/arfan21/project-sprint-social-media-api/pkg/constant"
	dbpostgres "github.com/arfan21/project-sprint-social-media-api/pkg/db/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	db dbpostgres.Queryer
}

func New(db dbpostgres.Queryer) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) Begin(ctx context.Context) (tx pgx.Tx, err error) {
	return r.db.Begin(ctx)
}

func (r Repository) WithTx(tx pgx.Tx) *Repository {
	r.db = tx
	return &r
}

func (r Repository) Create(ctx context.Context, data entity.User) (err error) {
	query := `
		INSERT INTO users (id, email, phone, name, password)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = r.db.Exec(ctx, query, data.ID, data.Email, data.Phone, data.Name, data.Password)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == constant.ErrSQLUniqueViolation {
				err = constant.ErrEmailOrPhoneAlreadyRegistered
			}
		}

		err = fmt.Errorf("user.repository.Create: failed to create user: %w", err)
		return
	}

	return
}

func (r Repository) GetByCredential(ctx context.Context, credentialType, credentialValue string) (data entity.User, err error) {
	credType := "email"
	if credentialType == "phone" {
		credType = "phone"
	}
	query := `
		SELECT id, name, password, email, phone
		FROM users
		WHERE ` + credType + ` = $1
	`

	err = r.db.QueryRow(ctx, query, credentialValue).Scan(
		&data.ID,
		&data.Name,
		&data.Password,
		&data.Email,
		&data.Phone,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrUserNotFound
		}

		err = fmt.Errorf("user.repository.GetByCredential: failed to get user by email: %w", err)

		return
	}

	return
}

func (r Repository) GetByID(ctx context.Context, id string) (data entity.User, err error) {
	query := `
		SELECT id, name, email, phone
		FROM users
		WHERE id = $1
	`

	err = r.db.QueryRow(ctx, query, id).Scan(
		&data.ID,
		&data.Name,
		&data.Email,
		&data.Phone,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrUserNotFound
		}

		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == constant.ErrSQLInvalidUUID {
				err = constant.ErrUserNotFound
			}
		}

		err = fmt.Errorf("user.repository.GetByID: failed to get user by id: %w", err)

		return
	}

	return
}

func (r Repository) AddFriend(ctx context.Context, userIdAdder, userIdAdded string) (err error) {
	query := `
		INSERT INTO friends (userIdAdder, userIdAdded)
		VALUES ($1, $2)
	`

	_, err = r.db.Exec(ctx, query, userIdAdder, userIdAdded)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == constant.ErrSQLUniqueViolation {
				err = constant.ErrFriendUseralreadyAdded
			}
		}

		err = fmt.Errorf("user.repository.AddFriend: failed to add friend: %w", err)
		return
	}

	return
}

func (r Repository) DeleteFriend(ctx context.Context, userIdAdder, userIdAdded string) (err error) {
	query := `
		DELETE FROM friends
		WHERE( userIdAdder = $1 AND userIdAdded = $2 ) OR ( userIdAdder = $2 AND userIdAdded = $1)
	`

	cmd, err := r.db.Exec(ctx, query, userIdAdder, userIdAdded)
	if err != nil {
		err = fmt.Errorf("user.repository.DeleteFriend: failed to delete friend: %w", err)
		return
	}

	if cmd.RowsAffected() == 0 {
		err = fmt.Errorf("user.repository.DeleteFriend: failed to delete friend: %w", constant.ErrFriendUserNotAdded)
		return
	}

	return
}

// queryGetListWithFilter is a helper function to get list of user with filter
// the query is expected to be joined with friends table
// where table users is alias as u and friends is alias as fr
func (r Repository) queryGetListWithFilter(ctx context.Context, query string, groupByCols []string, filter model.UserGetListRequest) (rows pgx.Rows, err error) {
	arrArgs := []interface{}{}
	whereQuery := ""
	andStatement := " AND "

	if filter.OnlyFriend && filter.UserID != "" {
		arrArgs = append(arrArgs, filter.UserID)
		whereQuery += fmt.Sprintf("(u.id != $%d) %s", len(arrArgs), andStatement)

	}

	if len(filter.UserIDs) > 0 {
		arrArgs = append(arrArgs, filter.UserIDs)
		whereQuery += fmt.Sprintf("u.id = ANY ($%d) %s", len(arrArgs), andStatement)
	}

	if filter.OnlyFriend {
		arrArgs = append(arrArgs, filter.UserID)
		whereQuery += fmt.Sprintf("(fr.userIdAdder = $%d OR fr.userIdAdded = $%d) %s", len(arrArgs), len(arrArgs), andStatement)
	}

	if filter.Search != "" {
		arrArgs = append(arrArgs, "%"+strings.ToLower(filter.Search)+"%")
		whereQuery += fmt.Sprintf("(LOWER(u.name) LIKE $%d) %s", len(arrArgs), andStatement)
	}

	if lenArgs := len(arrArgs); lenArgs > 0 {
		whereQuery = "WHERE " + whereQuery[:len(whereQuery)-len(andStatement)] + " "
	}

	query += whereQuery

	if len(groupByCols) > 0 {
		colsStr := strings.Join(groupByCols, ", ")
		query += fmt.Sprintf("GROUP BY %s ", colsStr)
	}

	if !filter.DisableOrder {
		sortBy := "id"
		if filter.SortBy != "" && filter.SortBy != "createdAt" {
			sortBy = "friendCount"

		}

		query += fmt.Sprintf("ORDER BY %s ", sortBy)

		orderBy := "DESC"
		if filter.OrderBy != "" && filter.OrderBy != "desc" {
			orderBy = "ASC"
		}
		query += fmt.Sprintf("%s ", orderBy)
	}

	if !filter.DisableOffset {
		arrArgs = append(arrArgs, filter.Limit)
		query += fmt.Sprintf("LIMIT $%d ", len(arrArgs))

		arrArgs = append(arrArgs, filter.Offset)
		query += fmt.Sprintf("OFFSET $%d ", len(arrArgs))
	}
	return r.db.Query(ctx, query, arrArgs...)
}

func (r Repository) GetList(ctx context.Context, filter model.UserGetListRequest) (data []entity.User, err error) {
	query := `
		SELECT u.id, u.name, u.imageurl, u.createdat, COUNT(fr.useridadder) AS friendCount
		FROM users u
		LEFT JOIN friends fr ON (fr.useridadder = u.id OR fr.useridadded = u.id)
	`

	rows, err := r.queryGetListWithFilter(ctx, query, []string{"u.id", "u.name", "u.imageurl", "u.createdat"}, filter)
	if err != nil {
		err = fmt.Errorf("user.repository.GetList: failed to get list of user: %w", err)
		return
	}

	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.ID, &user.Name, &user.ImageUrl, &user.CreatedAt, &user.FriendCount)
		if err != nil {
			err = fmt.Errorf("user.repository.GetList: failed to scan user: %w", err)
			return
		}

		data = append(data, user)
	}

	return
}

func (r Repository) GetCountList(ctx context.Context, filter model.UserGetListRequest) (count int, err error) {
	query := `
		SELECT COUNT(DISTINCT u.id)
		FROM users u
		LEFT JOIN friends fr ON (fr.useridadder = u.id OR fr.useridadded = u.id)
	`
	filter.DisableOffset = true
	filter.DisableOrder = true
	rows, err := r.queryGetListWithFilter(ctx, query, []string{}, filter)
	if err != nil {
		err = fmt.Errorf("user.repository.GetCountList: failed to get count list of user: %w", err)
		return
	}

	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			err = fmt.Errorf("user.repository.GetCountList: failed to scan count: %w", err)
			return
		}
	}

	return
}

func (r Repository) IsFriend(ctx context.Context, userIdAdder, userIdAdded string) (isFriend bool, err error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM friends
			WHERE (userIdAdder = $1 AND userIdAdded = $2) OR (userIdAdder = $2 AND userIdAdded = $1)
		)
	`

	err = r.db.QueryRow(ctx, query, userIdAdder, userIdAdded).Scan(&isFriend)
	if err != nil {
		err = fmt.Errorf("user.repository.IsFriend: failed to check is friend: %w", err)
		return
	}

	return
}

func (r Repository) GetListMap(ctx context.Context, filter model.UserGetListRequest) (data map[string]entity.User, err error) {
	query := `
		SELECT u.id, u.name, u.imageurl, u.createdat, COUNT(fr.useridadder) AS friendCount
		FROM users u
		LEFT JOIN friends fr ON (fr.useridadder = u.id OR fr.useridadded = u.id)
	`

	rows, err := r.queryGetListWithFilter(ctx, query, []string{"u.id", "u.name", "u.imageurl", "u.createdat"}, filter)
	if err != nil {
		err = fmt.Errorf("user.repository.GetList: failed to get list of user: %w", err)
		return
	}

	data = make(map[string]entity.User)

	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.ID, &user.Name, &user.ImageUrl, &user.CreatedAt, &user.FriendCount)
		if err != nil {
			err = fmt.Errorf("user.repository.GetList: failed to scan user: %w", err)
			return
		}

		data[user.ID.String()] = user
	}

	return
}

func (r Repository) UpdatePhone(ctx context.Context, userId, phone string) (err error) {
	query := `
		UPDATE users
		SET phone = $1
		WHERE id = $2
	`

	_, err = r.db.Exec(ctx, query, phone, userId)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == constant.ErrSQLUniqueViolation {
				err = constant.ErrPhoneAlreadyRegistered
			}
		}

		err = fmt.Errorf("user.repository.UpdatePhone: failed to update phone: %w", err)
		return
	}

	return
}

func (r Repository) UpdateEmail(ctx context.Context, userId, email string) (err error) {
	query := `
		UPDATE users
		SET email = $1
		WHERE id = $2
	`

	_, err = r.db.Exec(ctx, query, email, userId)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == constant.ErrSQLUniqueViolation {
				err = constant.ErrEmailAlreadyRegistered
			}
		}

		err = fmt.Errorf("user.repository.UpdateEmail: failed to update email: %w", err)
		return
	}

	return
}

func (r Repository) UpdateProfile(ctx context.Context, data entity.User) (err error) {
	query := `
		UPDATE users
	`

	arrArgs := []interface{}{}
	comma := ", "
	setQuery := ""
	arrArgs = append(arrArgs, data.Name)
	setQuery += fmt.Sprintf("name = $%d%s", len(arrArgs), comma)

	if data.ImageUrl.Valid {
		arrArgs = append(arrArgs, data.ImageUrl)
		setQuery += fmt.Sprintf("imageurl = $%d%s", len(arrArgs), comma)
	}

	if len(arrArgs) > 0 {
		setQuery = "SET " + setQuery[:len(setQuery)-len(comma)] + " "
	}

	query += setQuery
	arrArgs = append(arrArgs, data.ID)
	query += fmt.Sprintf("WHERE id = $%d", len(arrArgs))
	fmt.Println(query)
	_, err = r.db.Exec(ctx, query, arrArgs...)
	if err != nil {
		err = fmt.Errorf("user.repository.UpdateProfile: failed to update profile: %w", err)
		return
	}

	return
}
