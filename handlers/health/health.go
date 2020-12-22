package health

import (
	"context"
	"net/http"
	"time"

	"github.com/tammiec/go-rest-api/clients/render"
	"github.com/tammiec/go-rest-api/handlers/utils"
	"github.com/tammiec/go-rest-api/services/health"
)

type Deps struct {
	Health health.Health
	Render render.Render
}

type HandlerImpl struct {
	Deps *Deps
}

type Handler interface {
	Ping(w http.ResponseWriter, r *http.Request)
}

func (impl *HandlerImpl) Ping(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	result := impl.Deps.Health.Ping(ctx)

	unavailableSystems := make([]string, 0)

	if !result.UsersAvailable {
		unavailableSystems = append(unavailableSystems, "Users")
	}

	if len(unavailableSystems) > 0 {
		response := utils.Status{
			Status:  http.StatusServiceUnavailable,
			Message: "Service Unavailable",
			Failed:  unavailableSystems,
		}

		//nolint
		impl.Deps.Render.JSON(w, http.StatusServiceUnavailable, response)
		return
	}

	//nolint
	impl.Deps.Render.JSON(w, http.StatusOK, utils.Status{
		Status:  http.StatusOK,
		Message: "Healthy",
		Failed:  unavailableSystems,
	})
}
