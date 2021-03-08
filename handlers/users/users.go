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
	request, err := parseRequest(w, r, true, []string{})
	if err != nil {
		switch err {
		case errs.ErrBadRequestParams:
			respondWithError(impl.Deps.Render, w, http.StatusBadRequest, "Bad Request - Missing Params")
		case errs.ErrInvalidId:
			respondWithError(impl.Deps.Render, w, http.StatusBadRequest, "Bad Request - Invalid ID")
		default:
			respondWithError(impl.Deps.Render, w, http.StatusInternalServerError, "Something went wrong. Please try again.")
		}
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
	request, err := parseRequest(w, r, false, []string{"name", "email", "password"})
	if err != nil {
		switch err {
		case errs.ErrBadRequestParams:
			respondWithError(impl.Deps.Render, w, http.StatusBadRequest, "Bad Request - Missing Params")
		default:
			respondWithError(impl.Deps.Render, w, http.StatusInternalServerError, "Something went wrong. Please try again.")
		}
		return
	}
	result, err := impl.Deps.UsersService.Create(request)
	if err != nil {
		respondWithError(impl.Deps.Render, w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	impl.Deps.Render.JSON(w, http.StatusCreated, result)
}

func (impl *HandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	request, err := parseRequest(w, r, true, []string{"name", "email", "password"})
	if err != nil {
		switch err {
		case errs.ErrBadRequestParams:
			respondWithError(impl.Deps.Render, w, http.StatusBadRequest, "Bad Request - Missing Params")
		case errs.ErrInvalidId:
			respondWithError(impl.Deps.Render, w, http.StatusBadRequest, "Bad Request - Invalid ID")
		default:
			respondWithError(impl.Deps.Render, w, http.StatusInternalServerError, "Something went wrong. Please try again.")
		}
		return
	}
	result, err := impl.Deps.UsersService.Update(request)
	if err != nil {
		respondWithError(impl.Deps.Render, w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	impl.Deps.Render.JSON(w, http.StatusCreated, result)
}

func (impl *HandlerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	request, err := parseRequest(w, r, true, []string{})
	if err != nil {
		switch err {
		case errs.ErrInvalidId:
			respondWithError(impl.Deps.Render, w, http.StatusBadRequest, "Bad Request - Invalid ID")
		default:
			respondWithError(impl.Deps.Render, w, http.StatusInternalServerError, "Something went wrong. Please try again.")
		}
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

func parseRequest(w http.ResponseWriter, r *http.Request, getId bool, paramKeys []string) (*model.UserRequest, error) {
	urlParams := r.URL.Query()
	parsedParams, err := parseUrlParams(urlParams, paramKeys)
	if err != nil {
		return nil, err
	}
	userRequest := model.UserRequest{
		Name:     parsedParams["name"],
		Email:    parsedParams["email"],
		Password: parsedParams["password"],
	}
	if getId {
		userRequest.Id, err = validateId(mux.Vars(r)["id"], w)
		if err != nil {
			return nil, err
		}
	}
	return &userRequest, nil
}

func validateId(idString string, w http.ResponseWriter) (*int, error) {
	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, errs.ErrInvalidId
	}
	return &id, nil
}

func parseUrlParams(urlParams map[string][]string, paramKeys []string) (map[string]*string, error) {
	parsedParams := make(map[string]*string)
	for _, key := range paramKeys {
		val, ok := urlParams[key]
		if !ok {
			return nil, errs.ErrBadRequestParams
		}
		parsedParams[key] = &val[0]
	}
	return parsedParams, nil
}
