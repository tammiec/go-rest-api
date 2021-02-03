package users

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mock_render "github.com/tammiec/go-rest-api/generatedmocks/clients/render"
	mock_users_service "github.com/tammiec/go-rest-api/generatedmocks/services/users"
	model "github.com/tammiec/go-rest-api/models/user"
)

func TestUsers_ListSuccess(t *testing.T) {
	usersServiceMock := mock_users_service.NewMockUsersService(gomock.NewController(t))
	renderMock := mock_render.NewMockRender(gomock.NewController(t))

	req, err := http.NewRequest("GET", "/users", nil)
	require.NoError(t, err)

	gomock.InOrder(
		usersServiceMock.EXPECT().List().Return([]*model.UserResponse{}, nil),
		renderMock.EXPECT().JSON(gomock.Any(), gomock.Any(), gomock.Any()),
	)

	deps := Deps{
		UsersService: usersServiceMock,
		Render:       renderMock,
	}

	usersHandler := HandlerImpl{Deps: &deps}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersHandler.List)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUsers_GetSuccess(t *testing.T) {
	usersServiceMock := mock_users_service.NewMockUsersService(gomock.NewController(t))
	renderMock := mock_render.NewMockRender(gomock.NewController(t))

	req, err := http.NewRequest("GET", "/users/1", nil)
	require.NoError(t, err)

	id := 1

	expectedResponse := &model.UserResponse{
		Id:    1,
		Name:  "test",
		Email: "test",
	}
	gomock.InOrder(
		renderMock.EXPECT().JSON(gomock.Any(), gomock.Any(), gomock.Any()),
		usersServiceMock.EXPECT().Get(&model.UserRequest{Id: &id}).Return(expectedResponse, nil),
	)

	deps := Deps{
		UsersService: usersServiceMock,
		Render:       renderMock,
	}

	usersHandler := HandlerImpl{Deps: &deps}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersHandler.Get)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUsers_CreateSuccess(t *testing.T) {
	usersServiceMock := mock_users_service.NewMockUsersService(gomock.NewController(t))
	renderMock := mock_render.NewMockRender(gomock.NewController(t))

	req, err := http.NewRequest("POST", "/users?name=test&email=test&password=test", nil)
	require.NoError(t, err)

	testString := "test"

	request := &model.UserRequest{
		Name:     &testString,
		Email:    &testString,
		Password: &testString,
	}

	expectedResponse := &model.UserResponse{
		Id:    1,
		Name:  "test",
		Email: "test",
	}
	gomock.InOrder(
		usersServiceMock.EXPECT().Create(request).Return(expectedResponse, nil),
		renderMock.EXPECT().JSON(gomock.Any(), gomock.Any(), gomock.Any()),
	)

	deps := Deps{
		UsersService: usersServiceMock,
		Render:       renderMock,
	}

	usersHandler := HandlerImpl{Deps: &deps}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersHandler.Create)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUsers_UpdateSuccess(t *testing.T) {
	usersServiceMock := mock_users_service.NewMockUsersService(gomock.NewController(t))
	renderMock := mock_render.NewMockRender(gomock.NewController(t))

	req, err := http.NewRequest("PUT", "/users/1?name=test&email=test&password=test", nil)
	require.NoError(t, err)

	id := 1
	testString := "test"

	request := &model.UserRequest{
		Id:       &id,
		Name:     &testString,
		Email:    &testString,
		Password: &testString,
	}

	expectedResponse := &model.UserResponse{
		Id:    1,
		Name:  "test",
		Email: "test",
	}
	gomock.InOrder(
		renderMock.EXPECT().JSON(gomock.Any(), gomock.Any(), gomock.Any()),
		usersServiceMock.EXPECT().Update(request).Return(expectedResponse, nil),
	)

	deps := Deps{
		UsersService: usersServiceMock,
		Render:       renderMock,
	}

	usersHandler := HandlerImpl{Deps: &deps}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersHandler.Update)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUsers_DeleteSuccess(t *testing.T) {
	usersServiceMock := mock_users_service.NewMockUsersService(gomock.NewController(t))
	renderMock := mock_render.NewMockRender(gomock.NewController(t))

	req, err := http.NewRequest("DELETE", "/users/1", nil)
	require.NoError(t, err)

	id := 1

	expectedResponse := &model.UserResponse{
		Id:    1,
		Name:  "test",
		Email: "test",
	}
	gomock.InOrder(
		renderMock.EXPECT().JSON(gomock.Any(), gomock.Any(), gomock.Any()),
		usersServiceMock.EXPECT().Delete(&model.UserRequest{Id: &id}).Return(expectedResponse, nil),
	)

	deps := Deps{
		UsersService: usersServiceMock,
		Render:       renderMock,
	}

	usersHandler := HandlerImpl{Deps: &deps}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersHandler.Delete)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}
