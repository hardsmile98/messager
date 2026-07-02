package dto

type CreatePrivateChatRequest struct {
	TargetUserID string `json:"target_user_id" validate:"required,uuid"`
}
