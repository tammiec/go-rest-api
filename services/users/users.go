package users

import (
	"log"

	users "github.com/tammiec/go-rest-api/dals/users"
	model "github.com/tammiec/go-rest-api/models/user"
)

type Deps struct {
	Users users.Users
}

type Config struct{}

type UsersService interface {
	List() ([]*model.UserResponse, error)
	Get(request *model.UserRequest) (*model.UserResponse, error)
	Create(request *model.UserRequest) (*model.UserResponse, error)
	Update(request *model.UserRequest) (*model.UserResponse, error)
	Delete(request *model.UserRequest) (*model.UserResponse, error)
}

type UsersServiceImpl struct {
	deps *Deps
}

func New(deps *Deps, config *Config) UsersService {
	return &UsersServiceImpl{deps: deps}
}

func (impl *UsersServiceImpl) List() ([]*model.UserResponse, error) {
	users, err := impl.deps.Users.GetUsers()
	if err != nil {
		log.Println("Could not get users. ", err)
		return []*model.UserResponse{}, err
	}
	return users, nil
}

func (impl *UsersServiceImpl) Get(request *model.UserRequest) (*model.UserResponse, error) {
	user, err := impl.deps.Users.GetUser(*request.Id)
	if err != nil {
		log.Println("Could not get user. ", err)
		return &model.UserResponse{}, err
	}
	return user, nil
}

func (impl *UsersServiceImpl) Create(request *model.UserRequest) (*model.UserResponse, error) {
	user, err := impl.deps.Users.CreateUser(*request.Name, *request.Email, *request.Password)
	if err != nil {
		log.Println("Could not create user. ", err)
		return &model.UserResponse{}, err
	}
	return user, nil
}

func (impl *UsersServiceImpl) Update(request *model.UserRequest) (*model.UserResponse, error) {
	user, err := impl.deps.Users.UpdateUser(*request.Id, *request.Name, *request.Email, *request.Password)
	if err != nil {
		log.Println("Could not create user. ", err)
		return &model.UserResponse{}, err
	}
	return user, nil
}

func (impl *UsersServiceImpl) Delete(request *model.UserRequest) (*model.UserResponse, error) {
	user, err := impl.deps.Users.DeleteUser(*request.Id)
	if err != nil {
		log.Println("Could not delete user. ", err)
		return &model.UserResponse{}, err
	}
	return user, nil
}
