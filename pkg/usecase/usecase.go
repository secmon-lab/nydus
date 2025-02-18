package usecase

import (
	"github.com/secmon-lab/nydus/pkg/adapter"
)

type UseCase struct {
	clients *adapter.Clients
}

func New(clients *adapter.Clients) *UseCase {
	return &UseCase{
		clients: clients,
	}
}
