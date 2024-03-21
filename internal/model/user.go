package model

// {
// 	"credentialType": "phone | email", // not null, should in enum
// 	"credentialType": "+621.... | n@n.co", // not null
// 	// 👆 if credentialType == email, value should be in email format,
// 	// otherwise, phone number should start with "international calling code" (including the "+" prefix)
// 	// with minLength=7 and maxLength=13 (including the "international calling code" with the "+" and only
// 	// applicable with credentialType == phone)
// 	"name": "namadepan namabelakang", // not null, minLength 5, maxLength 50, name can be duplicate with others
// 	"password": "" // not null, minLength 5, maxLength 15
// }

type UserRegisterRequest struct {
	CredentialType  string `json:"credentialType" validate:"required,oneof=phone email"`
	CredentialValue string `json:"credentialValue" validate:"required,emailorphone=CredentialType"`
	Name            string `json:"name" validate:"required,min=5,max=50"`
	Password        string `json:"password" validate:"required,min=5,max=15"`
}

type UserLoginRequest struct {
	CredentialType  string `json:"credentialType" validate:"required,oneof=phone email"`
	CredentialValue string `json:"credentialValue" validate:"required,emailorphone=CredentialType"`
	Password        string `json:"password" validate:"required,min=5,max=15"`
}

type UserLoginResponse struct {
	Phone       *string `json:"phone,omitempty"`
	Email       *string `json:"email,omitempty"`
	Name        string  `json:"name"`
	AccessToken string  `json:"accessToken"`
}

type FriendRequest struct {
	UserIDAdder string `json:"-" validate:"required"`
	UserID      string `json:"userId" validate:"required"`
}

type UserGetListRequest struct {
	Limit         int      `query:"limit" validate:"omitempty,gte=0"`
	Offset        int      `query:"offset" validate:"omitempty,gte=0"`
	SortBy        string   `query:"sortBy" validate:"omitempty,oneof=friendCount createdAt"`
	OrderBy       string   `query:"orderBy" validate:"omitempty,oneof=asc desc"`
	OnlyFriend    bool     `query:"onlyFriend"`
	Search        string   `query:"search"`
	UserID        string   `query:"-"`
	UserIDs       []string `query:"-"`
	DisableOffset bool     `query:"-"`
	DisableOrder  bool     `query:"-"`
}

type UserResponse struct {
	UserID      string `json:"userId"`
	Name        string `json:"name"`
	ImageUrl    string `json:"imageUrl"`
	FriendCount int    `json:"friendCount"`
	CreatedAt   string `json:"createdAt,omitempty"`
}
