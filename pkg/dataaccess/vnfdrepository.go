package dataaccess

import (
	"github.com/vishwanathj/protovnfdparser/pkg/models"
)

// VnfdRepository defines what methods we need from the database
type VnfdRepository interface {
	InsertVnfd(vnfd *models.Vnfd) error
	FindVnfdByID(id string) (*models.Vnfd, error)
	FindVnfdByName(name string) (*models.Vnfd, error)
	GetVnfds(start string, limit int) ([]models.Vnfd, int, error)
}
