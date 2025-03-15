package response

import (
	"github.com/go-chi/render"
	"net/http"
	"strings"
)

type Object struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type HealthCheck struct {
	Commit   string            `json:"commit"`
	Database map[string]string `json:"database"`
	Version  string            `json:"version"`
}

func OK(w http.ResponseWriter, r *http.Request, data any) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, data)
}

func NoContent(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)

	v := Object{
		Success: true,
	}

	render.JSON(w, r, v)
}

func BadRequest(w http.ResponseWriter, r *http.Request, err error, data any) {
	render.Status(r, http.StatusBadRequest)

	v := Object{
		Success: false,
		Data:    data,
		Message: err.Error(),
	}
	render.JSON(w, r, v)
}

func NotFound(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusNotFound)

	v := Object{
		Success: false,
		Message: err.Error(),
	}
	render.JSON(w, r, v)
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusInternalServerError)

	v := Object{
		Success: false,
		Message: err.Error(),
	}

	if err != nil && strings.Contains(err.Error(), "context deadline exceeded") {
		switch r.Header.Get("Language") {
		case "RUS":
			v.Message = "Превышено время ожидания запроса"
		case "KAZ":
			v.Message = "Сұраудың күту уақыты асып кетті"
		default:
			v.Message = "Request timeout exceeded"
		}
	}

	render.JSON(w, r, v)
}
