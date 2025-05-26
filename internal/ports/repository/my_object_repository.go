package repository

import "github.com/cesarfreire/go-boilerplate/internal/domain/entity"

type MyObjectRepository interface {
	GetAllObjects() ([]entity.MyObject, error)
	GetObjectByID(id int64) (entity.MyObject, error)
}
