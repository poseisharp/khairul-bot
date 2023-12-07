package interface_repositories

import "github.com/poseisharp/khairul-bot/internal/domain/entities"

type ServerRepository interface {
	FindOne(id string) (entities.Server, error)
	FindAll() ([]entities.Server, error)
	Store(server entities.Server) error
	Update(server entities.Server) error
	Delete(id string) error
}
