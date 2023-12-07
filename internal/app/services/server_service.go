package services

import (
	"github.com/poseisharp/khairul-bot/internal/domain/entities"
	interface_repositories "github.com/poseisharp/khairul-bot/internal/interfaces/repositories"
)

type ServerService struct {
	serverRepository interface_repositories.ServerRepository
}

func NewServerService(serverRepository interface_repositories.ServerRepository) *ServerService {
	return &ServerService{
		serverRepository: serverRepository,
	}
}

func (s *ServerService) GetServer(id string) (*entities.Server, error) {
	server, err := s.serverRepository.FindOne(id)
	if err != nil {
		return nil, err
	}

	return &server, nil
}

func (s *ServerService) GetServers() ([]entities.Server, error) {
	servers, err := s.serverRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (s *ServerService) CreateServer(server entities.Server) error {
	return s.serverRepository.Store(server)
}

func (s *ServerService) UpdateServer(server entities.Server) error {
	return s.serverRepository.Update(server)
}

func (s *ServerService) DeleteServer(id string) error {
	return s.serverRepository.Delete(id)
}
