package service

import (
	"chat/internal/repository"
	"context"

	pb "github.com/hardsmile98/messager/sdk/chat/v1"
	commonpb "github.com/hardsmile98/messager/sdk/common/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	return nil, status.Errorf(codes.Unimplemented, "method CreatePrivateChat not implemented")
}

func (s *ChatService) GetUserChats(ctx context.Context, req *pb.GetUserChatsRequest) (*pb.GetUserChatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserChats not implemented")
}

func (s *ChatService) GetChatInfo(ctx context.Context, req *pb.GetChatInfoRequest) (*commonpb.ChatInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChatInfo not implemented")
}
