package questions

import (
	"context"

	"github.com/go-kit/kit/log"
)

type Service interface {
	Add(ctx context.Context, q Question) (string, error)                 // Create a new question.
	Read(ctx context.Context) ([]*Question, error)                       // Read all the questions
	ReadOwned(ctx context.Context, user string) ([]*Question, error)     // Read all the questions for a given owner
	Peek(ctx context.Context, id string) (*Question, error)              // Read the question with ID
	Change(ctx context.Context, id, statement string) (*Question, error) // Change the question statement
	Answer(ctx context.Context, id, answer string) (*Question, error)    // Answer the question.
	Delete(ctx context.Context, id string) (*Question, error)            // Delete the question
	Logger() log.Logger
}
