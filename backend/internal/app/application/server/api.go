package server

import (
	"github.com/igh9410/blabber-hive/backend/internal/api"
	"github.com/igh9410/blabber-hive/backend/internal/app/application/service"
)

// API implements the api.StrictServerInterface
var _ api.StrictServerInterface = (*API)(nil)

type API struct {
	chatSerivce service.ChatService
}

func NewAPI(chatService service.ChatService) api.StrictServerInterface {

	return &API{
		chatSerivce: chatService,
	}
}
