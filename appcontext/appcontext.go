package appcontext

import (
	renderClient "github.com/tammiec/go-rest-api/clients/render"
	config "github.com/tammiec/go-rest-api/config"
	users "github.com/tammiec/go-rest-api/dals/users"
	"github.com/tammiec/go-rest-api/services/health"
	usersService "github.com/tammiec/go-rest-api/services/users"
	"github.com/unrolled/render"
)

type Clients struct {
	Render renderClient.Render
}

type DALs struct {
	Users users.Users
}

type Services struct {
	Health health.Health
	Users  usersService.UsersService
}

func BuildClients(config *config.Config) *Clients {
	_render := renderClient.New(&renderClient.Deps{Render: render.New(render.Options{})}, &renderClient.Config{})
	return &Clients{Render: _render}
}

func BuildDALs(clients *Clients, config *config.Config) *DALs {
	usersConfig := &users.Config{Url: config.DB.Url}
	usersDeps := &users.Deps{}
	usersDB := users.New(usersDeps, usersConfig)

	return &DALs{Users: usersDB}
}

func BuildServices(dals *DALs) *Services {
	usersDeps := &usersService.Deps{Users: dals.Users}
	usersService := usersService.New(usersDeps, &usersService.Config{})

	healthDeps := &health.Deps{Users: dals.Users}
	healthService := health.New(healthDeps)

	return &Services{Health: healthService, Users: usersService}
}

type AppContext struct {
	Clients  *Clients
	DALs     *DALs
	Services *Services
	Config   *config.Config
}

func NewAppContext(config *config.Config) *AppContext {
	builtClients := BuildClients(config)
	builtDALs := BuildDALs(builtClients, config)
	builtServices := BuildServices(builtDALs)

	return &AppContext{DALs: builtDALs, Services: builtServices, Config: config, Clients: builtClients}
}
