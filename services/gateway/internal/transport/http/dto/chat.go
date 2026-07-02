package dto

type CreatePrivateChatRequest struct {
	TargetUserID string `json:"target_user_id" validate:"required,uuid"`
}

type GetUserChatsRequest struct {
	PageSize  int32  `json:"page_size" validate:"required,min=1,max=100"`
	PageToken string `json:"page_token" validate:"required"`
}
