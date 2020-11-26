package transport

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	qsrv "github.com/macchid/questandans/questions"
	"github.com/macchid/questandans/utils"
)

type Endpoints struct {
	Create    endpoint.Endpoint
	GetAll    endpoint.Endpoint
	GetByID   endpoint.Endpoint
	GetByUser endpoint.Endpoint
	Change    endpoint.Endpoint
	Answer    endpoint.Endpoint
	Delete    endpoint.Endpoint
}

func MakeEndpoints(s qsrv.Service) Endpoints {
	return Endpoints{
		Create:    logEndpoint(s.Logger(), "Endpoints::Create")(makeCreateEndpoint(s)),
		GetAll:    logEndpoint(s.Logger(), "Endpoints::GetAll")(makeGetAllEndpoint(s)),
		GetByID:   logEndpoint(s.Logger(), "Endpoints::GetByID")(makeGetByIDEndpoint(s)),
		GetByUser: logEndpoint(s.Logger(), "Endpoints::GetByUser")(makeGetByUserEndpoint(s)),
		Change:    logEndpoint(s.Logger(), "Endpoints::Change")(makeChangeEndpoint(s)),
		Answer:    logEndpoint(s.Logger(), "Endpoints::Answer")(makeAnswerEndpoint(s)),
		Delete:    logEndpoint(s.Logger(), "Endpoints::Delete")(makeDeleteEndpoint(s)),
	}
}

func logEndpoint(logger log.Logger, method string) func(endpoint.Endpoint) endpoint.Endpoint {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			_, logEnd := utils.LogStart(logger, method)
			defer logEnd(time.Now())

			return next(ctx, request)
		}
	}
}

func makeCreateEndpoint(s qsrv.Service) endpoint.Endpoint {
	_, logEnd := utils.LogStart(s.Logger(), "MakeEndpoints::makeCreateEndpoint")
	defer logEnd(time.Now())

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		id, err := s.Add(ctx, req.Question)
		return CreateResponse{ID: id, Err: err}, nil
	}
}

func makeGetAllEndpoint(s qsrv.Service) endpoint.Endpoint {
	_, logEnd := utils.LogStart(s.Logger(), "MakeEndpoints::makeGetAllEndpoint")
	defer logEnd(time.Now())

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		qs, err := s.Read(ctx)
		return GetAllResponse{Count: len(qs), Questions: qs, Err: err}, nil
	}
}

func makeGetByIDEndpoint(s qsrv.Service) endpoint.Endpoint {
	_, logEnd := utils.LogStart(s.Logger(), "MakeEndpoints::makeGetByIDEndpoint")
	defer logEnd(time.Now())

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		q, err := s.Peek(ctx, req.ID)
		return GetByIDResponse{Question: q, Err: err}, nil
	}
}

func makeGetByUserEndpoint(s qsrv.Service) endpoint.Endpoint {
	_, logEnd := utils.LogStart(s.Logger(), "MakeEndpoints::makeGetByUserEndpoint")
	defer logEnd(time.Now())

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByUserRequest)
		qs, err := s.ReadOwned(ctx, req.User)
		return GetByUserResponse{Count: len(qs), Questions: qs, Err: err}, nil
	}
}

func makeChangeEndpoint(s qsrv.Service) endpoint.Endpoint {
	_, logEnd := utils.LogStart(s.Logger(), "MakeEndpoints::makeChangeEndpoint")
	defer logEnd(time.Now())

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ChangeQuestionRequest)
		q, err := s.Change(ctx, req.ID, req.Statement)
		return ChangeQuestionResponse{Question: q, Err: err}, nil
	}
}

func makeAnswerEndpoint(s qsrv.Service) endpoint.Endpoint {
	_, logEnd := utils.LogStart(s.Logger(), "MakeEndpoints::makeAnswerEndpoint")
	defer logEnd(time.Now())

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AnswerQuestionRequest)
		q, err := s.Answer(ctx, req.ID, req.Answer)
		return ChangeQuestionResponse{Question: q, Err: err}, nil
	}
}

func makeDeleteEndpoint(s qsrv.Service) endpoint.Endpoint {
	_, logEnd := utils.LogStart(s.Logger(), "MakeEndpoints::makeDeleteEndpoint")
	defer logEnd(time.Now())

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteQuestionRequest)
		q, err := s.Delete(ctx, req.ID)
		return DeleteQuestionResponse{Question: q, Err: err}, nil
	}
}
