package transport

import (
	qsrv "github.com/macchid/questandans/questions"
)

type CreateRequest struct {
	Question qsrv.Question
}

type CreateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

type GetAllResponse struct {
	Count     int `json:"count"`
	Questions []*qsrv.Question
	Err       error `json:"error,omitempty"`
}

type GetByIDRequest struct {
	ID string
}

type GetByIDResponse struct {
	Question *qsrv.Question
	Err      error `json:"error,omitempty"`
}

type GetByUserRequest struct {
	User string `json:"string"`
}

type GetByUserResponse struct {
	Count     int `json:"count"`
	Questions []*qsrv.Question
	Err       error `json:"error,omitempty"`
}

type AnswerQuestionRequest struct {
	ID     string
	Answer string `json:"answer"`
}

type AnswerQuestionResponse struct {
	Question *qsrv.Question
	Err      error `json:"error,omitempty"`
}

type ChangeQuestionRequest struct {
	ID        string
	Statement string `json:"statement"`
}

type ChangeQuestionResponse struct {
	Question *qsrv.Question
	Err      error `json:"error,omitempty"`
}

type DeleteQuestionRequest struct {
	ID string
}

type DeleteQuestionResponse struct {
	Question *qsrv.Question
	Err      error `json:"error,omitempty"`
}

func (r CreateResponse) Error() error {
	return r.Err
}

func (r GetAllResponse) Error() error {
	return r.Err
}

func (r GetByIDResponse) Error() error {
	return r.Err
}

func (r GetByUserResponse) Error() error {
	return r.Err
}

func (r AnswerQuestionResponse) Error() error {
	return r.Err
}

func (r ChangeQuestionResponse) Error() error {
	return r.Err
}

func (r DeleteQuestionResponse) Error() error {
	return r.Err
}
