package health

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mock_users "github.com/tammiec/go-rest-api/generatedmocks/dals/users"
)

func TestPingSuccess(t *testing.T) {
	usersMock := mock_users.NewMockUsers(gomock.NewController(t))

	deps := &Deps{Users: usersMock}
	healthService := New(deps)

	ctx := context.TODO()

	usersMock.EXPECT().Ping(ctx).Return(nil)

	response := healthService.Ping(ctx)

	expected := &HealthResponse{
		UsersAvailable: true,
	}

	require.Equal(t, expected, response)
}

func TestPingFail(t *testing.T) {
	usersMock := mock_users.NewMockUsers(gomock.NewController(t))

	deps := &Deps{Users: usersMock}
	healthService := New(deps)

	ctx := context.TODO()

	usersMock.EXPECT().Ping(ctx).Return(errors.New("test"))
	usersMock.EXPECT().GetName().Return("Users")

	response := healthService.Ping(ctx)

	expected := &HealthResponse{
		UsersAvailable: false,
	}

	require.Equal(t, expected, response)
}
