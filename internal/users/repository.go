package users

import "github.com/LouisHatton/menu-link-up/internal/common/repository"

type Repository interface {
	repository.CrudRepository[User]
}
