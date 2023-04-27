package utils

import (
	"errors"
	"net/http"
	"github.com/go-chi/render"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	ErrMsg string `json:"errorMsg,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(errMsg string) render.Renderer {
	err := errors.New(errMsg)
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		ErrMsg:         err.Error(),
	}
}

func ErrServerError(errMsg string) render.Renderer {
	err := errors.New(errMsg)
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		ErrMsg:         err.Error(),
	}
}

func ErrNotFoundError(errMsg string) render.Renderer {
	err := errors.New(errMsg)
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 404,
		ErrMsg:         err.Error(),
	}
}

func ErrCustomError(errCode int, errMsg string) render.Renderer {
	err := errors.New(errMsg)
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: errCode,
		ErrMsg:         err.Error(),
	}
}
