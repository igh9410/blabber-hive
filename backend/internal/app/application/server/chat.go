package server

import (
	"context"

	"github.com/igh9410/blabber-hive/backend/internal/api"
	"github.com/igh9410/blabber-hive/backend/pkg/pointerutil"
)

// ChatServiceCreateChatRoom implements api.StrictServerInterface.
func (a *API) ChatServiceCreateChatRoom(ctx context.Context, request api.ChatServiceCreateChatRoomRequestObject) (api.ChatServiceCreateChatRoomResponseObject, error) {
	chatRoomName := request.Body.Name

	chatRoom, err := a.chatSerivce.CreateChatRoom(ctx, chatRoomName)
	if err != nil {
		return nil, err
	}

	return &api.ChatServiceCreateChatRoom200JSONResponse{
		ChatRoom: &api.ChatRoom{
			Id:        pointerutil.StringToPtr(chatRoom.ID.String()),
			Name:      chatRoom.Name,
			CreatedAt: pointerutil.TimeToPtr(chatRoom.CreatedAt),
		},
	}, nil

}
