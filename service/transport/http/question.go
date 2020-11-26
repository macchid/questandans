package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	qsrv "github.com/macchid/questandans/questions"
	"github.com/macchid/questandans/transport"
)

func NewService(root *mux.Router, svcEndpoints transport.Endpoints, logger log.Logger) *mux.Router {

	options := []kithttp.ServerOption{

		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeErrorResponse),
	}

	qroutes := root.PathPrefix("/api/v1/questions").Subrouter()
	qroutes.Path("/").Methods(http.MethodPost).Handler(
		kithttp.NewServer(
			svcEndpoints.Create,
			decodeCreateRequest,
			encodeResponse,
			options...,
		),
	)

	qroutes.Path("/").Methods(http.MethodGet).Handler(
		kithttp.NewServer(
			svcEndpoints.GetAll,
			decodeGetAllRequest,
			encodeResponse,
			options...,
		),
	)

	qroutes.Path("/{id}").Methods(http.MethodGet).Handler(
		kithttp.NewServer(
			svcEndpoints.GetByID,
			decodeGetByIDRequest,
			encodeResponse,
			options...,
		),
	)

	qroutes.Path("/users/{username}").Methods(http.MethodGet).Handler(
		kithttp.NewServer(
			svcEndpoints.GetByUser,
			decodeGetByUserRequest,
			encodeResponse,
			options...,
		),
	)

	qroutes.Path("/{id}/answer").Methods(http.MethodPost).Handler(
		kithttp.NewServer(
			svcEndpoints.Answer,
			decodeAnswerRequest,
			encodeResponse,
			options...,
		),
	)

	qroutes.Path("/{id}").Methods(http.MethodPatch).Handler(
		kithttp.NewServer(
			svcEndpoints.Change,
			decodeChangeRequest,
			encodeResponse,
			options...,
		),
	)

	qroutes.Path("/{id}").Methods(http.MethodDelete).Handler(
		kithttp.NewServer(
			svcEndpoints.Delete,
			decodeDeleteRequest,
			encodeResponse,
			options...,
		),
	)

	return root
}

func extractURI(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, "http-url", r.URL.Path)
}

func decodeCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req transport.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req.Question); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetAllRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetByIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	return transport.GetByIDRequest{ID: id}, nil
}

func decodeGetByUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	username := vars["username"]
	return transport.GetByUserRequest{User: username}, nil
}

func decodeAnswerRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req transport.AnswerQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	id := vars["id"]
	req.ID = id

	return req, nil
}

func decodeChangeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req transport.ChangeQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	id := vars["id"]
	req.ID = id

	return req, nil
}

func decodeDeleteRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	return transport.DeleteQuestionRequest{ID: id}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if e, ok := resp.(erroer); ok && e.Error() != nil {
		encodeErrorResponse(ctx, e.Error(), w)
		return nil
	}

	return json.NewEncoder(w).Encode(resp)
}

type erroer interface {
	Error() error
}

func encodeErrorResponse(ctx context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case qsrv.ErrQuestionNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
