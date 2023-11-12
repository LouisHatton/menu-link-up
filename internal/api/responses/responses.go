package responses

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

type HttpResponse struct {
	Err        error `json:"-"` // low-level runtime error
	StatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *HttpResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

func NotFoundResponse(name string) render.Renderer {
	return &HttpResponse{
		StatusCode: http.StatusNotFound,
		StatusText: fmt.Sprintf("%s not found", name),
	}
}

func ErrInvalidRequest(err error) render.Renderer {
	return &HttpResponse{
		Err:        err,
		StatusCode: http.StatusBadRequest,
		StatusText: "Invalid request",
		ErrorText:  err.Error(),
	}
}

func ErrUnauthorised() render.Renderer {
	return &HttpResponse{
		StatusCode: http.StatusUnauthorized,
		StatusText: "Unauthroised",
	}
}

func ErrForbidden() render.Renderer {
	return &HttpResponse{
		StatusCode: http.StatusForbidden,
		StatusText: "Forbidden",
	}
}

func ErrInternalServerError() render.Renderer {
	return &HttpResponse{
		StatusCode: http.StatusInternalServerError,
		StatusText: "Internal Server Error",
	}
}

func Accepted() render.Renderer {
	return &HttpResponse{
		StatusCode: http.StatusAccepted,
		StatusText: "accepted",
	}
}
