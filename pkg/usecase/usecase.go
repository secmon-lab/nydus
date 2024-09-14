package usecase

import (
	"github.com/secmon-as-code/nydus/pkg/adapter"
)

type UseCase struct {
	clients *adapter.Clients
}

func New(clients *adapter.Clients) *UseCase {
	return &UseCase{
		clients: clients,
	}
}
