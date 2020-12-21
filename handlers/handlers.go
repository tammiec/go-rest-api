package handlers

import (
	"github.com/tammiec/go-rest-api/clients/render"
	"github.com/tammiec/go-rest-api/handlers/health"
	"github.com/tammiec/go-rest-api/handlers/users"
	healthService "github.com/tammiec/go-rest-api/services/health"
	usersService "github.com/tammiec/go-rest-api/services/users"
)

type Deps struct {
	Health healthService.Health
	Users  usersService.UsersService
	Render render.Render
}

type Handlers struct {
	Health health.Handler
	Users  users.Handler
}

func New(deps *Deps) *Handlers {
	healthDeps := health.Deps{
		Health: deps.Health,
		Render: deps.Render,
	}
	healthMiddle := health.HandlerImpl{Deps: &healthDeps}

	usersDeps := users.Deps{
		UsersService: deps.Users,
		Render:       deps.Render,
	}
	usersMiddle := users.HandlerImpl{Deps: &usersDeps}

	return &Handlers{
		Health: &healthMiddle,
		Users:  &usersMiddle,
	}
}
