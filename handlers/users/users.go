package users

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	List(w http.ResponseWriter, r *http.Request)
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

func (impl *HandlerImpl) List(w http.ResponseWriter, r *http.Request) {
	result, err := impl.Deps.UsersService.List()
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
	urlParams := r.URL.Query()
	userRequest := model.UserRequest{
		Name:     parseUrlParam(urlParams, "name", render, w),
		Email:    parseUrlParam(urlParams, "email", render, w),
		Password: parseUrlParam(urlParams, "password", render, w),
	}
	log.Println(r.URL.Query())
	var err error
	if getId {
		userRequest.Id, err = validateId(mux.Vars(r)["id"], w)
	}
	if err != nil {
		log.Println(err)
		respondWithError(render, w, http.StatusInternalServerError, "Internal Server Error")
		return nil
	}
	return &userRequest
}

func validateId(idString string, w http.ResponseWriter) (*int, error) {
	id, err := strconv.Atoi(idString)
	return &id, err
}

func parseUrlParam(urlParams map[string][]string, key string, render render.Render, w http.ResponseWriter) *string {
	var value *string
	val, ok := urlParams[key]
	if ok == false {
		respondWithError(render, w, http.StatusBadRequest, fmt.Sprintf("Bad Request - Missing Parameter %s", key))
		return value
	} else {
		value = &val[0]
	}
	log.Println(key, value, ok)
	return value
}
