package services

import (
	"errors"
	config "github.com/resssoft/mediaArchive/configuration"
	"github.com/resssoft/mediaArchive/models"
	"github.com/resssoft/mediaArchive/repositories"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Add(models.User, string) error
	Update(models.User, string) error
	List(string, interface{}, string) ([]*models.User, error)
	hashPassword(string) (string, error)
	checkUserPassword(models.User, string) bool
	CheckCredentials(string, string) (models.User, error)
	Get(string, string) (models.UserData, error)
	AddOwner()
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

func (r *userApp) Add(user models.User, authorId string) error {
	author, err := r.repo.GetByID(authorId)
	if err != nil {
		return err
	}
	if !author.Role.CheckPerm(models.OwnerPerm) {
		return errors.New("you can`t create users")
	}
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

func (r *userApp) Update(user models.User, authorId string) error {
	author, err := r.repo.GetByID(authorId)
	if err != nil {
		return err
	}
	if user.Id != author.Id || !author.Role.CheckPerm(models.OwnerPerm) {
		return errors.New("you can`t update users")
	}
	err = r.repo.Add(user)
	if err != nil {
		return err
	}
	return err
}

func (r *userApp) List(name string, value interface{}, authorId string) ([]*models.User, error) {
	author, err := r.repo.GetByID(authorId)
	if err != nil {
		return nil, err
	}
	if !author.Role.CheckPerm(models.OwnerPerm) {
		return nil, errors.New("permission denied")
	}
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

func (r *userApp) checkUserPassword(user models.User, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password+config.PasswordSalt()),
		[]byte(password+config.PasswordSalt()),
	)
	return err == nil
}

func (r *userApp) CheckCredentials(email, password string) (models.User, error) {
	user, err := r.repo.GetByEmail(email)
	if err != nil {
		return models.User{}, err
	}
	if !r.checkUserPassword(user, password) {
		return models.User{}, errors.New("user or password is incorrect")
	}
	return user, nil
}

func (r *userApp) Get(userId, authorId string) (models.UserData, error) {
	author, err := r.repo.GetByID(authorId)
	if err != nil {
		return models.UserData{}, err
	}
	if !author.Role.CheckPerm(models.OwnerPerm) {
		return models.UserData{}, errors.New("you can`t create users")
	}
	user, err := r.repo.GetByID(userId)
	if err != nil {
		return models.UserData{}, err
	}
	return user.Data(), err
}

func (r *userApp) AddOwner() {
	password, err := r.hashPassword(config.AdminPassword())
	if err != nil {
		log.Err(err).Send()
	}
	err = r.repo.Add(models.User{
		Id:       primitive.NewObjectID(),
		Email:    config.AdminLogin(),
		Password: password,
		Role:     models.UserRole{Permissions: []models.UserPermission{models.OwnerPerm}},
	})
	if err != nil {
		log.Err(err).Send()
	}
}
