package model

type PostRequest struct {
	PostInHtml string   `json:"postInHtml" validate:"required,min=3,max=500"`
	Tags       []string `json:"tags" validate:"required,dive,required"`
	UserID     string   `json:"-" validate:"required"`
}

type PostCommentRequest struct {
	PostID  string `json:"postId" validate:"required"`
	Comment string `json:"comment" validate:"required,min=2,max=500"`
	UserID  string `json:"-" validate:"required"`
}

type PostGetListRequest struct {
	UserID        string   `query:"-" validate:"required"`
	Limit         int      `query:"limit" validate:"omitempty,gte=0"`
	Offset        int      `query:"offset" validate:"omitempty,gte=0"`
	Search        string   `query:"search"`
	SearchTags    []string `query:"searchTag"`
	DisableOffset bool     `query:"-"`
	DisableOrder  bool     `query:"-"`
}

type PostListResponse struct {
	PostID   string                `json:"postId"`
	Post     PostResponse          `json:"post"`
	Creator  UserResponse          `json:"creator"`
	Comments []PostCommentResponse `json:"comments"`
}

type PostResponse struct {
	PostInHtml string   `json:"postInHtml"`
	Tags       []string `json:"tags"`
	CreatedAt  string   `json:"createdAt"`
}

type PostCommentResponse struct {
	Comment   string       `json:"comment"`
	CreatedAt string       `json:"createdAt"`
	Creator   UserResponse `json:"creator"`
}
