package handler

import (
	"errors"
	"gateway/internal/transport/http/dto"
	"gateway/internal/transport/http/middleware"
	"gateway/internal/transport/http/response"
	"gateway/internal/validation"
	"net/http"

	"github.com/go-chi/chi/v5"
	pb "github.com/hardsmile98/messager/sdk/chat/v1"
)

type ChatHandler struct {
	client pb.ChatServiceClient
}

func NewChatHandler(client pb.ChatServiceClient) *ChatHandler {
	return &ChatHandler{client}
}

func (h *ChatHandler) CreatePrivateChat(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r)

	var req dto.CreatePrivateChatRequest

	if err := validation.DecodeAndValidate(r, &req); err != nil {
		response.RequestError(w, err)
		return
	}

	resp, err := h.client.CreatePrivateChat(r.Context(), &pb.CreatePrivateChatRequest{
		InitiatorId:  userID,
		TargetUserId: req.TargetUserID,
	})
	if err != nil {
		response.GRPCError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, resp)
}

func (h *ChatHandler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r)

	resp, err := h.client.GetUserChats(r.Context(), &pb.GetUserChatsRequest{
		UserId: userID,
	})

	if err != nil {
		response.GRPCError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, resp)
}

func (h *ChatHandler) GetChatInfo(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r)
	chatID := chi.URLParam(r, "chat_id")

	if chatID == "" {
		response.RequestError(w, errors.New("chat_id is required"))
		return
	}

	resp, err := h.client.GetChatInfo(r.Context(), &pb.GetChatInfoRequest{
		ChatId: chatID,
		UserId: userID,
	})

	if err != nil {
		response.GRPCError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, resp)
}
