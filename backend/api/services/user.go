package services

import (
	"errors"
	config "github.com/resssoft/mediaArchive/configuration"
	"github.com/resssoft/mediaArchive/models"
	"github.com/resssoft/mediaArchive/repositories"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Add(models.User, models.User) error
	Update(models.User, models.User) error
	List(string, interface{}) ([]*models.User, error)
	hashPassword(string) (string, error)
	CheckUserPassword(models.User, string) bool
}

type userApp struct {
	repo repositories.UserRepository
}

func NewUserApp(repo repositories.UserRepository) IUserService {
	app := &userApp{
		repo: repo,
	}
	return app
}

func (r *userApp) Add(user, author models.User) error {
	if !author.Role.CheckPerm(models.OwnerPerm) {
		return errors.New("you can`t create users")
	}
	var err error
	user.Password, err = r.hashPassword(user.Password)
	if err != nil {
		return err
	}
	err = r.repo.Add(user)
	if err != nil {
		return err
	}
	return err
}

func (r *userApp) Update(user, author models.User) error {
	if user.Id != author.Id || !author.Role.CheckPerm(models.OwnerPerm) {
		return errors.New("you can`t update users")
	}
	var err error
	err = r.repo.Add(user)
	if err != nil {
		return err
	}
	return err
}

func (r *userApp) List(name string, value interface{}) ([]*models.User, error) {
	//TODO: add limit
	users, err := r.repo.List(name, value)
	if err != nil {
		return users, err
	}
	return users, err
}

func (r *userApp) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password+config.PasswordSalt()), 14)
	return string(bytes), err
}

func (r *userApp) CheckUserPassword(user models.User, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password+config.PasswordSalt()),
		[]byte(password+config.PasswordSalt()),
	)
	return err == nil
}
