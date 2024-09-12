package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	models "github.com/miselaytes-anton/airy/internal/models"
)

// ServerEnv represents the environment containing server dependencies.
type Server struct {
	Router interface {
		HandlerFunc(string, string, http.HandlerFunc)
	}
	Measurements models.MeasurementModelInterface
	Events       models.EventModelInterface
	LogError     *log.Logger
	LogInfo      *log.Logger
}

type ResponseError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

// StartServer starts the http server.
func (s Server) routes() {
	s.Router.HandlerFunc(http.MethodGet, "/api/graphs", s.handleGraphs())
	s.Router.HandlerFunc(http.MethodGet, "/api/events", s.handleEventsList())
	s.Router.HandlerFunc(http.MethodPost, "/api/events", s.handleEventsCreate())
	s.Router.HandlerFunc(http.MethodPatch, "/api/events/:id", s.handleEventsUpdate())
	s.Router.HandlerFunc(http.MethodGet, "/api/measurements", s.handleMeasurements())
}

func (s Server) jsonError(w http.ResponseWriter, err error, code int) {
	var responseError ResponseError

	if code == http.StatusInternalServerError {
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
		s.LogError.Output(2, trace)

		responseError = ResponseError{
			Status: http.StatusText(code),
			Error:  "internal server error occured",
		}
	} else {
		responseError = ResponseError{
			Status: http.StatusText(code),
			Error:  err.Error(),
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(responseError)
}

func (s Server) readJson(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func lowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

func (s Server) jsonValidationError(w http.ResponseWriter, err error) {
	var formattedErrors []string

	for _, err := range err.(validator.ValidationErrors) {
		formattedError := fmt.Sprintf("%s did not pass validation rules: %s %s", lowerFirst(err.Field()), err.Tag(), lowerFirst(err.Param()))
		formattedErrors = append(formattedErrors, strings.TrimSpace(formattedError))
	}

	err = errors.New(strings.Join(formattedErrors, ", "))

	s.jsonError(w, err, http.StatusBadRequest)
}
