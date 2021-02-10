package users

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tammiec/go-rest-api/clients/render"
	"github.com/tammiec/go-rest-api/errs"
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
	request := parseRequest(w, r, impl.Deps.Render, true, []string{})
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
	request := parseRequest(w, r, impl.Deps.Render, false, []string{"name", "email", "password"})
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
	request := parseRequest(w, r, impl.Deps.Render, true, []string{"name", "email", "password"})
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
	request := parseRequest(w, r, impl.Deps.Render, true, []string{})
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

func parseRequest(w http.ResponseWriter, r *http.Request, render render.Render, getId bool, paramKeys []string) *model.UserRequest {
	urlParams := r.URL.Query()
	parsedParams, err := parseUrlParams(urlParams, paramKeys, render)
	if err != nil {
		respondWithError(render, w, http.StatusBadRequest, "Bad Request - Missing Params")
		return nil
	}
	userRequest := model.UserRequest{
		Name:     parsedParams["name"],
		Email:    parsedParams["email"],
		Password: parsedParams["password"],
	}
	if getId {
		userRequest.Id, err = validateId(mux.Vars(r)["id"], w)
	}
	if err != nil {
		respondWithError(render, w, http.StatusBadRequest, "Bad Request - Invalid User ID")
		return nil
	}
	return &userRequest
}

func validateId(idString string, w http.ResponseWriter) (*int, error) {
	id, err := strconv.Atoi(idString)
	return &id, err
}

func parseUrlParams(urlParams map[string][]string, paramKeys []string, render render.Render) (map[string]*string, error) {
	parsedParams := make(map[string]*string)
	for _, key := range paramKeys {
		val, ok := urlParams[key]
		if ok == false {
			return nil, errs.ErrBadRequest
		}
		parsedParams[key] = &val[0]
	}
	return parsedParams, nil
}
