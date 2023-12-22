package services

import (
	"github.com/poseisharp/khairul-bot/internal/domain/aggregates"
	"github.com/poseisharp/khairul-bot/internal/persistent/repositories"
)

type ServerService struct {
	serverRepository *repositories.ServerRepository
}

func NewServerService(serverRepository *repositories.ServerRepository) *ServerService {
	return &ServerService{
		serverRepository: serverRepository,
	}
}

func (s *ServerService) GetServer(id string) (*aggregates.Server, error) {
	server, err := s.serverRepository.FindOne(id)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (s *ServerService) GetServers() ([]aggregates.Server, error) {
	servers, err := s.serverRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (s *ServerService) CreateServerIfNotExists(server aggregates.Server) error {
	if _, err := s.serverRepository.FindOne(server.ID); err != nil {
		return s.serverRepository.Store(server)
	}

	return nil
}

func (s *ServerService) UpdateServer(server aggregates.Server) error {
	return s.serverRepository.Update(server)
}

func (s *ServerService) DeleteServer(id string) error {
	return s.serverRepository.Delete(id)
}
