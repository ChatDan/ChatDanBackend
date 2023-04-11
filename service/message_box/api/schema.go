package api

import "ChatDanBackend/common"

// Box

type BoxCommonResponse struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
	Title   string `json:"title"`
}

type BoxCreateRequest struct {
	Title string `json:"title" query:"title" validate:"required"`
}

type BoxCreateResponse struct {
	BoxCommonResponse
}

type BoxListRequest struct {
	common.PageRequest
	Title string `json:"title"`
	Owner int    `json:"owner"`
}

type BoxListResponse struct {
	BoxCommonResponse
}

type BoxGetResponse struct {
	BoxCommonResponse
	Posts []string `json:"posts"`
}

type BoxModifyRequest struct {
	Title string `json:"title" query:"title" validate:"required"`
}

type BoxModifyResponse struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
	Title   string `json:"title"`
}

type BoxDeleteResponse struct {
	Message string `json:"message"`
}

// Post

type PostCommonResponse struct {
	ID         string `json:"id"`
	PosterID   string `json:"poster_id"`
	Content    string `json:"content"`
	Visibility string `json:"visibility"`
}

type PostCreateRequest struct {
	MessageBoxID string `json:"message_box_id" query:"message_box_id" validate:"required"`
	Content      string `json:"content" query:"content" validate:"required"`
	Visibility   string `json:"visibility" query:"visibility" validate:"required,oneof=public private"`
}

type PostCreateResponse struct {
	PostCommonResponse
}

type PostListRequest struct {
	common.PageRequest
	MessageBoxID string `json:"message_box_id" query:"message_box_id" validate:"required"`
}

type PostListResponse struct {
	PostCommonResponse
}

type PostGetResponse struct {
	PostCommonResponse
	Channels []string `json:"channels"`
}

type PostDeleteResponse struct {
	Message string `json:"message"`
}
