package service

import (
	apperrors "chat/internal/errors"
	"chat/internal/repository"
	"context"

	pb "github.com/hardsmile98/messager/sdk/chat/v1"
	commonpb "github.com/hardsmile98/messager/sdk/common/v1"
)

type ChatService struct {
	pb.UnimplementedChatServiceServer
	chatRepo        repository.ChatRepository
	participantRepo repository.ParticipantRepository
}

func NewChatService(
	chatRepo repository.ChatRepository,
	participantRepo repository.ParticipantRepository,
) *ChatService {
	return &ChatService{
		chatRepo:        chatRepo,
		participantRepo: participantRepo,
	}
}

func (s *ChatService) CreatePrivateChat(ctx context.Context, req *pb.CreatePrivateChatRequest) (*commonpb.Chat, error) {
	if req.InitiatorId == "" || req.TargetUserId == "" {
		return nil, apperrors.InvalidArgument("initiator_id and target_user_id are required")
	}

	if req.InitiatorId == req.TargetUserId {
		return nil, apperrors.InvalidArgument("cannot chat with yourself")
	}

	exists, err := s.chatRepo.ExistsPrivateChat(ctx, req.InitiatorId, req.TargetUserId)

	if err != nil {
		return nil, apperrors.InternalError(ctx, "failed to check if chat exists", err)
	}

	if exists {
		return nil, apperrors.AlreadyExists("chat already exists")
	}

	chatID, err := s.chatRepo.CreatePrivateChat(ctx, req.InitiatorId, req.TargetUserId)

	if err != nil {
		return nil, apperrors.InternalError(ctx, "failed to create chat", err)
	}

	return &commonpb.Chat{
		Id:   chatID,
		Type: commonpb.ChatType_CHAT_TYPE_PRIVATE,
	}, nil
}

func (s *ChatService) GetUserChats(ctx context.Context, req *pb.GetUserChatsRequest) (*pb.GetUserChatsResponse, error) {
	if req.UserId == "" {
		return nil, apperrors.InvalidArgument("user_id is required")
	}

	rows, err := s.chatRepo.GetChatsByUserID(ctx, req.UserId)

	if err != nil {
		return nil, apperrors.InternalError(ctx, "failed to get chats", err)
	}

	chats := make([]*commonpb.ChatInfo, 0, len(rows))

	for _, chat := range rows {
		chats = append(chats, &commonpb.ChatInfo{
			Chat: &commonpb.Chat{
				Id:   chat.ID,
				Type: commonpb.ChatType_CHAT_TYPE_PRIVATE,
				Members: []*commonpb.User{
					{Id: chat.UserID, Username: chat.Username},
				},
			},
		})
	}

	return &pb.GetUserChatsResponse{
		Chats: chats,
	}, nil
}

func (s *ChatService) GetChatInfo(ctx context.Context, req *pb.GetChatInfoRequest) (*commonpb.ChatInfo, error) {
	if req.ChatId == "" || req.UserId == "" {
		return nil, apperrors.InvalidArgument("chat_id and user_id are required")
	}

	ok, err := s.participantRepo.IsParticipant(ctx, req.ChatId, req.UserId)

	if err != nil {
		return nil, apperrors.InternalError(ctx, "failed to check if user is participant", err)
	}

	if !ok {
		return nil, apperrors.NotFound("user is not a participant of the chat")
	}

	companionID, err := s.chatRepo.GetCompanionID(ctx, req.ChatId, req.UserId)

	if err != nil {
		return nil, apperrors.InternalError(ctx, "failed to get companion id", err)
	}

	return &commonpb.ChatInfo{
		Chat: &commonpb.Chat{
			Id:   req.ChatId,
			Type: commonpb.ChatType_CHAT_TYPE_PRIVATE,
			Members: []*commonpb.User{
				{Id: companionID},
			},
		},
	}, nil
}
