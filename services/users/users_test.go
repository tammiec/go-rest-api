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
			Email: "test",
		},
	}

	usersMock.EXPECT().List().Return(expectedResponse, nil).Times(1)

	users := impl{&Deps{Users: usersMock}}

	result, err := users.List()

	require.NoError(t, err)
	require.Equal(t, expectedResponse, result)
}

func TestList_Error(t *testing.T) {
	usersMock := mock_users_dal.NewMockUsers(gomock.NewController(t))

	usersMock.EXPECT().List().Return([]*model.UserResponse{}, errors.New("test")).Times(1)

	users := impl{&Deps{Users: usersMock}}

	_, err := users.List()

	require.Error(t, err)
}

func TestGet_Success(t *testing.T) {
	usersMock := mock_users_dal.NewMockUsers(gomock.NewController(t))

	expectedResponse := &model.UserResponse{
		Id:    1,
		Name:  "test",
		Email: "test",
	}

	usersMock.EXPECT().Get(gomock.Any()).Return(expectedResponse, nil).Times(1)

	users := impl{&Deps{Users: usersMock}}

	id := 1
	result, err := users.Get(&model.UserRequest{Id: &id})

	require.NoError(t, err)
	require.Equal(t, expectedResponse, result)
}

func TestGet_Error(t *testing.T) {
	usersMock := mock_users_dal.NewMockUsers(gomock.NewController(t))

	usersMock.EXPECT().Get(gomock.Any()).Return(&model.UserResponse{}, errors.New("test")).Times(1)

	users := impl{&Deps{Users: usersMock}}

	id := 1
	_, err := users.Get(&model.UserRequest{Id: &id})

	require.Error(t, err)
}

func TestCreate_Success(t *testing.T) {
	usersMock := mock_users_dal.NewMockUsers(gomock.NewController(t))

	expectedResponse := &model.UserResponse{
		Id:    1,
		Name:  "test",
		Email: "test",
	}

	usersMock.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedResponse, nil).Times(1)

	users := impl{&Deps{Users: usersMock}}

	testString := "test"
	result, err := users.Create(&model.UserRequest{
		Name:     &testString,
		Email:    &testString,
		Password: &testString,
	})

	require.NoError(t, err)
	require.Equal(t, expectedResponse, result)
}

func TestCreate_Error(t *testing.T) {
	usersMock := mock_users_dal.NewMockUsers(gomock.NewController(t))

	usersMock.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.UserResponse{}, errors.New("test")).Times(1)

	users := impl{&Deps{Users: usersMock}}

	testString := "test"
	_, err := users.Create(&model.UserRequest{
		Name:     &testString,
		Email:    &testString,
		Password: &testString,
	})

	require.Error(t, err)
}

func TestUpdate_Success(t *testing.T) {
	usersMock := mock_users_dal.NewMockUsers(gomock.NewController(t))

	expectedResponse := &model.UserResponse{
		Id:    1,
		Name:  "test",
		Email: "test",
	}

	usersMock.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedResponse, nil).Times(1)

	users := impl{&Deps{Users: usersMock}}

	id := 1
	testString := "test"
	result, err := users.Update(&model.UserRequest{
		Id:       &id,
		Name:     &testString,
		Email:    &testString,
		Password: &testString,
	})

	require.NoError(t, err)
	require.Equal(t, expectedResponse, result)
}

func TestUpdate_Error(t *testing.T) {
	usersMock := mock_users_dal.NewMockUsers(gomock.NewController(t))

	usersMock.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.UserResponse{}, errors.New("test")).Times(1)

	users := impl{&Deps{Users: usersMock}}

	id := 1
	testString := "test"
	_, err := users.Update(&model.UserRequest{
		Id:       &id,
		Name:     &testString,
		Email:    &testString,
		Password: &testString,
	})

	require.Error(t, err)
}

func TestDelete_Success(t *testing.T) {
	usersMock := mock_users_dal.NewMockUsers(gomock.NewController(t))

	expectedResponse := &model.UserResponse{
		Id:    1,
		Name:  "test",
		Email: "test",
	}

	usersMock.EXPECT().Delete(gomock.Any()).Return(expectedResponse, nil).Times(1)

	users := impl{&Deps{Users: usersMock}}

	id := 1
	testString := "test"
	result, err := users.Delete(&model.UserRequest{
		Id:       &id,
		Name:     &testString,
		Email:    &testString,
		Password: &testString,
	})

	require.NoError(t, err)
	require.Equal(t, expectedResponse, result)
}

func TestDelete_Error(t *testing.T) {
	usersMock := mock_users_dal.NewMockUsers(gomock.NewController(t))

	usersMock.EXPECT().Delete(gomock.Any()).Return(&model.UserResponse{}, errors.New("test")).Times(1)

	users := impl{&Deps{Users: usersMock}}

	id := 1
	_, err := users.Delete(&model.UserRequest{Id: &id})

	require.Error(t, err)
}
