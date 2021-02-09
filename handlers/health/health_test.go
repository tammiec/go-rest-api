package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mock_render "github.com/tammiec/go-rest-api/generatedmocks/clients/render"
	mock_health "github.com/tammiec/go-rest-api/generatedmocks/services/health"
	"github.com/tammiec/go-rest-api/handlers/utils"
	"github.com/tammiec/go-rest-api/services/health"
)

func TestHealthPing_OK(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/health", nil)
	require.NoError(t, err)

	healthMock := mock_health.NewMockHealth(gomock.NewController(t))

	healthResponse := &health.HealthResponse{
		UsersAvailable: true,
	}

	handlerResponse := utils.Status{
		Status:  http.StatusOK,
		Message: "Healthy",
		Failed:  make([]string, 0),
	}

	healthMock.EXPECT().Ping(gomock.Any()).Return(healthResponse).Times(1)

	renderCtrl := gomock.NewController(t)
	renderMock := mock_render.NewMockRender(renderCtrl)

	renderMock.EXPECT().JSON(gomock.Any(), gomock.Eq(http.StatusOK), gomock.Eq(handlerResponse)).Times(1)

	deps := Deps{
		Health: healthMock,
		Render: renderMock,
	}

	healthHandler := HandlerImpl{Deps: &deps}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler.Ping)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	// Response is asserted by the renderMock
}

func TestHealthPing_Fail(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/health", nil)
	require.NoError(t, err)

	healthMock := mock_health.NewMockHealth(gomock.NewController(t))

	healthResponse := &health.HealthResponse{
		UsersAvailable: false,
	}

	handlerResponse := utils.Status{
		Status:  http.StatusServiceUnavailable,
		Message: "Service Unavailable",
		Failed:  []string{"Users"},
	}

	healthMock.EXPECT().Ping(gomock.Any()).Return(healthResponse).Times(1)

	renderCtrl := gomock.NewController(t)
	renderMock := mock_render.NewMockRender(renderCtrl)

	renderMock.EXPECT().JSON(gomock.Any(), gomock.Eq(http.StatusServiceUnavailable), gomock.Eq(handlerResponse)).Times(1)

	deps := Deps{
		Health: healthMock,
		Render: renderMock,
	}

	healthHandler := HandlerImpl{Deps: &deps}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler.Ping)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	// Response is asserted by the renderMock
}
