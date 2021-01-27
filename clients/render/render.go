package render

import (
	"io"
	"net/http"

	"github.com/unrolled/render"
)

type Deps struct {
	Render *render.Render
}

type Config struct{}

type Render interface {
	JSON(w io.Writer, status int, v interface{})
}

type impl struct {
	deps *Deps
}

func New(deps *Deps, config *Config) Render {
	return &impl{deps: deps}
}

func (impl *impl) JSON(w io.Writer, status int, v interface{}) {
	err := impl.deps.Render.JSON(w, status, v)

	if err != nil {
		// nolint
		impl.deps.Render.Text(w, http.StatusInternalServerError, "Failed to write response")
	}
}
