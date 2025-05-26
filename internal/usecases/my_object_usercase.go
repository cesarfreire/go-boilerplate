package usecases

import (
	"github.com/cesarfreire/go-boilerplate/internal/domain/entity"
	"github.com/cesarfreire/go-boilerplate/internal/ports/repository"
)

type MyObjectUseCase struct {
	myObjectRepo repository.MyObjectRepository
}

func NewMyObjectUseCase(myObjectRepo repository.MyObjectRepository) *MyObjectUseCase {
	return &MyObjectUseCase{
		myObjectRepo: myObjectRepo,
	}
}

// GetAllObjects retrieves all objects using the repository.
func (uc *MyObjectUseCase) GetAllObjects() ([]entity.MyObject, error) {
	objects, err := uc.myObjectRepo.GetAllObjects()
	if err != nil {
		return nil, err
	}
	return objects, nil
}

// GetObjectByID retrieves an object by its ID using the repository.
func (uc *MyObjectUseCase) GetObjectByID(id int64) (entity.MyObject, error) {
	object, err := uc.myObjectRepo.GetObjectByID(id)
	if err != nil {
		return entity.MyObject{}, err
	}
	return object, nil
}
