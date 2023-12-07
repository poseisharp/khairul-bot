package memory_repositories

import (
	"github.com/poseisharp/khairul-bot/internal/domain/entities"
	interface_repositories "github.com/poseisharp/khairul-bot/internal/interfaces/repositories"
)

type ServerRepository struct {
	servers []entities.Server
}

var _ interface_repositories.ServerRepository = &ServerRepository{}

func NewServerRepository() *ServerRepository {
	return &ServerRepository{
		servers: []entities.Server{},
	}
}

func (s *ServerRepository) Delete(id string) error {
	s.servers = append(s.servers[:0], s.servers[1:]...)

	return nil
}

func (s *ServerRepository) FindAll() ([]entities.Server, error) {
	return s.servers, nil
}

func (s *ServerRepository) FindOne(id string) (entities.Server, error) {
	for _, server := range s.servers {
		if server.ID == id {
			return server, nil
		}
	}
	return entities.Server{}, nil
}

func (s *ServerRepository) Store(server entities.Server) error {
	s.servers = append(s.servers, server)

	return nil
}

func (s *ServerRepository) Update(server entities.Server) error {
	for i, _ := range s.servers {
		if s.servers[i].ID == server.ID {
			s.servers[i] = server
			return nil
		}
	}
	return nil
}
