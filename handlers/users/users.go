package users

import (
	"net/http"

	"github.com/dyninc/qstring"
	"github.com/tammiec/go-rest-api/clients/render"
	"github.com/tammiec/go-rest-api/handlers/utils"
	model "github.com/tammiec/go-rest-api/models/user"
	usersService "github.com/tammiec/go-rest-api/services/users"
)

type Deps struct {
	UsersService usersService.UsersService
	Render       render.Render
}

type Config struct{}

type Handler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type HandlerImpl struct {
	deps *Deps
}

func New(deps *Deps, config *Config) Handler {
	return &HandlerImpl{deps: deps}
}

func (impl *HandlerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	result, err := impl.deps.UsersService.GetAll()
	if err != nil {
		respondWithError(impl.deps.Render, w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	impl.deps.Render.JSON(w, http.StatusOK, result)
}

func (impl *HandlerImpl) Get(w http.ResponseWriter, r *http.Request) {
	request := parseRequest(w, r, impl.deps.Render)
	if request == nil {
		return
	}
	result, err := impl.deps.UsersService.Get(request)
	if err != nil {
		respondWithError(impl.deps.Render, w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	impl.deps.Render.JSON(w, http.StatusOK, result)
}

func (impl *HandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	request := parseRequest(w, r, impl.deps.Render)
	if request == nil {
		return
	}
	result, err := impl.deps.UsersService.Create(request)
	if err != nil {
		respondWithError(impl.deps.Render, w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	impl.deps.Render.JSON(w, http.StatusOK, result)
}

func (impl *HandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	request := parseRequest(w, r, impl.deps.Render)
	if request == nil {
		return
	}
	result, err := impl.deps.UsersService.Update(request)
	if err != nil {
		respondWithError(impl.deps.Render, w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	impl.deps.Render.JSON(w, http.StatusOK, result)
}

func (impl *HandlerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	request := parseRequest(w, r, impl.deps.Render)
	if request == nil {
		return
	}
	result, err := impl.deps.UsersService.Delete(request)
	if err != nil {
		respondWithError(impl.deps.Render, w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	impl.deps.Render.JSON(w, http.StatusOK, result)
}

func respondWithError(render render.Render, w http.ResponseWriter, status int, message string) {
	response := utils.Status{
		Status:  status,
		Message: message,
	}
	render.JSON(w, response.Status, response)
}

func parseRequest(w http.ResponseWriter, r *http.Request, render render.Render) *model.UserRequest {
	userRequest := &model.UserRequest{}
	err := qstring.Unmarshal(r.URL.Query(), userRequest)
	if err != nil {
		respondWithError(render, w, http.StatusInternalServerError, "Internal Server Error")
		return nil
	}
	return userRequest
}