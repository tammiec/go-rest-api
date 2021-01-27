package users

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mock_users_dal "github.com/tammiec/go-rest-api/generatedmocks/dals/users"
	model "github.com/tammiec/go-rest-api/models/user"
)

func TestList_Success(t *testing.T) {
	usersMock := mock_users_dal.NewMockUsers(gomock.NewController(t))

	expectedResponse := []*model.UserResponse{
		&model.UserResponse{
			Id:    1,
			Name:  "test",
			Email: "test@test.com",
		},
	}

	usersMock.EXPECT().GetUsers().Return(expectedResponse, nil).Times(1)

	users := impl{&Deps{Users: usersMock}}

	result, err := users.List()

	require.NoError(t, err)
	require.Equal(t, expectedResponse, result)
}

func TestList_Error(t *testing.T) {
	usersMock := mock_users_dal.NewMockUsers(gomock.NewController(t))

	usersMock.EXPECT().GetUsers().Return([]*model.UserResponse{}, errors.New("test")).Times(1)

	users := impl{&Deps{Users: usersMock}}

	_, err := users.List()

	require.Error(t, err)
}

func TestGet_Success(t *testing.T) {
	usersMock := mock_users_dal.NewMockUsers(gomock.NewController(t))

	expectedResponse := &model.UserResponse{
		Id:    1,
		Name:  "test",
		Email: "test@test.com",
	}

	usersMock.EXPECT().GetUser(gomock.Any()).Return(expectedResponse, nil).Times(1)

	users := impl{&Deps{Users: usersMock}}

	id := 1
	result, err := users.Get(&model.UserRequest{Id: &id})

	require.NoError(t, err)
	require.Equal(t, expectedResponse, result)
}

func TestGet_Error(t *testing.T) {
	usersMock := mock_users_dal.NewMockUsers(gomock.NewController(t))

	usersMock.EXPECT().GetUser(gomock.Any()).Return(&model.UserResponse{}, errors.New("test")).Times(1)

	users := impl{&Deps{Users: usersMock}}

	id := 1
	_, err := users.Get(&model.UserRequest{Id: &id})

	require.Error(t, err)
}
