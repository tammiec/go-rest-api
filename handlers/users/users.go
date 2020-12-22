package users

import (
	"net/http"
	"strconv"

	"github.com/dyninc/qstring"
	"github.com/gorilla/mux"
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
	Deps *Deps
}

func New(deps *Deps, config *Config) Handler {
	return &HandlerImpl{Deps: deps}
}

func (impl *HandlerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	result, err := impl.Deps.UsersService.GetAll()
	if err != nil {
		respondWithError(impl.Deps.Render, w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	impl.Deps.Render.JSON(w, http.StatusOK, result)
}

func (impl *HandlerImpl) Get(w http.ResponseWriter, r *http.Request) {
	request := parseRequest(w, r, impl.Deps.Render, true)
	if request == nil {
		return
	}
	result, err := impl.Deps.UsersService.Get(request)
	if err != nil {
		respondWithError(impl.Deps.Render, w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	impl.Deps.Render.JSON(w, http.StatusOK, result)
}

func (impl *HandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	request := parseRequest(w, r, impl.Deps.Render, false)
	if request == nil {
		return
	}
	result, err := impl.Deps.UsersService.Create(request)
	if err != nil {
		respondWithError(impl.Deps.Render, w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	impl.Deps.Render.JSON(w, http.StatusOK, result)
}

func (impl *HandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	request := parseRequest(w, r, impl.Deps.Render, true)
	if request == nil {
		return
	}
	result, err := impl.Deps.UsersService.Update(request)
	if err != nil {
		respondWithError(impl.Deps.Render, w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	impl.Deps.Render.JSON(w, http.StatusOK, result)
}

func (impl *HandlerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	request := parseRequest(w, r, impl.Deps.Render, true)
	if request == nil {
		return
	}
	result, err := impl.Deps.UsersService.Delete(request)
	if err != nil {
		respondWithError(impl.Deps.Render, w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	impl.Deps.Render.JSON(w, http.StatusOK, result)
}

func respondWithError(render render.Render, w http.ResponseWriter, status int, message string) {
	response := utils.Status{
		Status:  status,
		Message: message,
	}
	render.JSON(w, response.Status, response)
}

func parseRequest(w http.ResponseWriter, r *http.Request, render render.Render, getId bool) *model.UserRequest {
	userRequest := &model.UserRequest{}
	var err error
	if getId {
		userRequest.Id, err = validateId(mux.Vars(r)["id"], w)
	}
	err = qstring.Unmarshal(r.URL.Query(), userRequest)
	if err != nil {
		respondWithError(render, w, http.StatusInternalServerError, "Internal Server Error")
		return nil
	}
	return userRequest
}

func validateId(idString string, w http.ResponseWriter) (*int, error) {
	id, err := strconv.Atoi(idString)
	return &id, err
}
