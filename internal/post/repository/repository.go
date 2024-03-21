package postrepo

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

func (r Repository) Create(ctx context.Context, data entity.Post) (err error) {
	query := `
		INSERT INTO posts (id, userId, body, tags)
		VALUES ($1, $2, $3, $4)
	`

	_, err = r.db.Exec(ctx, query, data.ID, data.UserID, data.Body, data.Tags)
	if err != nil {
		err = fmt.Errorf("post.repository.Create: failed to create post: %w", err)
		return
	}

	return
}

func (r Repository) GetByID(ctx context.Context, id string) (data entity.Post, err error) {
	query := `
		SELECT id, userId, body, tags, createdAt, updatedAt
		FROM posts
		WHERE id = $1
	`

	err = r.db.QueryRow(ctx, query, id).Scan(&data.ID, &data.UserID, &data.Body, &data.Tags, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrPostNotFound
		}

		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == constant.ErrSQLInvalidUUID {
				err = constant.ErrPostNotFound
			}
		}

		err = fmt.Errorf("post.repository.GetByID: failed to get post by id: %w", err)
		return
	}

	return
}

func (r Repository) CreateComment(ctx context.Context, data entity.PostComment) (err error) {
	query := `
		INSERT INTO post_comments (id, postId, userId, comment)
		VALUES ($1, $2, $3, $4)
	`

	_, err = r.db.Exec(ctx, query, data.ID, data.PostID, data.UserID, data.Comment)
	if err != nil {
		err = fmt.Errorf("post.repository.CreateComment: failed to create comment: %w", err)
		return
	}

	return
}

// queryGetListWithFilter is a helper function to get post of user with filter
// the query is expected to be joined with post_comments and friends table
// where table posts as p, post_comments as pc, and friends as f
func (r Repository) queryGetListWithFilter(ctx context.Context, query string, filter model.PostGetListRequest) (rows pgx.Rows, err error) {
	arrArgs := []interface{}{}
	whereQuery := ""
	andStatement := " AND "

	if filter.Search != "" {
		arrArgs = append(arrArgs, "%"+strings.ToLower(filter.Search)+"%")
		whereQuery += fmt.Sprintf("(LOWER(p.body) LIKE $%d) %s", len(arrArgs), andStatement)
	}

	if len(filter.SearchTags) > 0 {
		arrArgs = append(arrArgs, filter.SearchTags)
		whereQuery += fmt.Sprintf("tags && $%d %s", len(arrArgs), andStatement)
	}

	// only friend post or self post
	if filter.UserID != "" {
		arrArgs = append(arrArgs, filter.UserID)
		whereQuery += fmt.Sprintf("(($%d IN (f.useridadder, f.useridadded) AND p.userId != $%d)  OR p.userId = $%d) %s", len(arrArgs), len(arrArgs), len(arrArgs), andStatement)
	}

	if lenArgs := len(arrArgs); lenArgs > 0 {
		whereQuery = "WHERE " + whereQuery[:len(whereQuery)-len(andStatement)] + " "
	}

	query += whereQuery

	if !filter.DisableOrder {
		query += "ORDER BY p.id DESC "
	}

	if !filter.DisableOffset {
		arrArgs = append(arrArgs, filter.Limit)
		query += fmt.Sprintf("LIMIT $%d ", len(arrArgs))

		arrArgs = append(arrArgs, filter.Offset)
		query += fmt.Sprintf("OFFSET $%d ", len(arrArgs))
	}

	return r.db.Query(ctx, query, arrArgs...)
}

func (r Repository) GetList(ctx context.Context, filter model.PostGetListRequest) (res []entity.Post, userIDs []string, err error) {
	query := `
		SELECT
			p.id, p.userId, p.body, p.tags, p.createdAt, pc.userId as commentUserId, pc.comment, pc.createdAt as commentCreatedAt
		FROM posts p
		LEFT JOIN post_comments pc ON p.id = pc.postId
		INNER JOIN friends f ON (p.userId = f.useridadder OR p.userId = f.useridadded)
	`

	rows, err := r.queryGetListWithFilter(ctx, query, filter)
	if err != nil {
		err = fmt.Errorf("post.repository.GetList: failed to get list post: %w", err)
		return
	}

	postMap := make(map[string]struct{})
	userIdUnique := make(map[string]struct{})
	i := 0
	for rows.Next() {
		var post entity.Post
		var comment entity.PostCommentNullable

		err = rows.Scan(&post.ID, &post.UserID, &post.Body, &post.Tags, &post.CreatedAt, &comment.UserID, &comment.Comment, &comment.CreatedAt)
		if err != nil {
			err = fmt.Errorf("post.repository.GetList: failed to scan rows: %w", err)
			return
		}

		if _, ok := postMap[post.ID.String()]; !ok {
			userIdUnique[post.UserID.String()] = struct{}{}
			postMap[post.ID.String()] = struct{}{}
			res = append(res, post)
			i = len(res) - 1
		}

		currentPost := res[i]
		if comment.UserID.Valid {
			currentPost.Comments = append(currentPost.Comments, comment)
			res[i] = currentPost

			userIdUnique[comment.UserID.UUID.String()] = struct{}{}
		}

	}

	userIDs = make([]string, len(userIdUnique))
	i = 0
	for k := range userIdUnique {
		userIDs[i] = k
		i++
	}

	return
}

func (r Repository) GetCountList(ctx context.Context, filter model.PostGetListRequest) (count int, err error) {
	query := `
		SELECT COUNT(DISTINCT p.id)
		FROM posts p
		LEFT JOIN post_comments pc ON p.id = pc.postId
		INNER JOIN friends f ON (p.userId = f.useridadder OR p.userId = f.useridadded)
	`
	filter.DisableOffset = true
	filter.DisableOrder = true
	rows, err := r.queryGetListWithFilter(ctx, query, filter)
	if err != nil {
		err = fmt.Errorf("post.repository.GetCountList: failed to get count list post: %w", err)
		return
	}

	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			err = fmt.Errorf("post.repository.GetCountList: failed to scan count: %w", err)
			return
		}
	}

	return
}
