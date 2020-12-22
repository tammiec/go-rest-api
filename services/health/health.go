package health

import (
	"context"
	"log"
)

type Pingable interface {
	Ping(ctx context.Context) error
	GetName() string
}

type Deps struct {
	Users Pingable
}

type Config struct{}

type Health interface {
	Ping(ctx context.Context) *HealthResponse
}

type HealthImpl struct {
	deps *Deps
}

func New(deps *Deps, config *Config) Health {
	return &HealthImpl{deps: deps}
}

type HealthResponse struct {
	UsersAvailable bool `json:"users"`
}

func RunPing(ctx context.Context, pingable Pingable) bool {
	err := pingable.Ping(ctx)

	if err != nil {
		log.Printf("Error while pinging: %v", pingable.GetName())
		log.Printf("Error: %v", err.Error())
		return false
	}

	return true
}

func (impl *HealthImpl) Ping(ctx context.Context) *HealthResponse {
	return &HealthResponse{
		UsersAvailable: RunPing(ctx, impl.deps.Users),
	}
}
