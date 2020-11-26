package questions

import (
	"context"
	"errors"
)

type Question struct {
	ID        string `json:"id,omitempty" bson:"id"`
	Who       string `json:"who" bson:"who"`
	Statement string `json:"statement" bson:"statement"`
	Answer    string `json:"answer,omitempty" bson:"answer"`
}

type Repository interface {
	CreateQuestion(ctx context.Context, q Question) error
	GetQuestions(ctx context.Context) ([]*Question, error)
	GetQuestionsByUser(ctx context.Context, username string) ([]*Question, error)
	GetQuestionByID(ctx context.Context, id string) (*Question, error)
	ChangeQuestion(ctx context.Context, q Question) error
	DeleteQuestion(ctx context.Context, id string) error
}

var (
	ErrCmdRepository    = errors.New("Unable to command repository")
	ErrQueryRepository  = errors.New("Unable to query repository")
	ErrQuestionNotFound = errors.New("Question not found")
)
