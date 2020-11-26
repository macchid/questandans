package implementation

import (
	"context"
	"fmt"

	qsrv "github.com/macchid/questandans/questions"
)

func newFakeRepository() (qsrv.Repository, error) {
	fr := make(fakeRepository, 1)
	return &fr, nil
}

type fakeRepository []*qsrv.Question

func (rep *fakeRepository) CreateQuestion(ctx context.Context, q qsrv.Question) error {
	newrep := append(*rep, &q)
	rep = &newrep
	return nil
}

func (rep *fakeRepository) GetQuestions(ctx context.Context) ([]*qsrv.Question, error) {
	return *rep, nil
}

func (rep *fakeRepository) GetQuestionsByUser(ctx context.Context, username string) ([]*qsrv.Question, error) {
	byUser := make([]*qsrv.Question, 1)

	for _, q := range *rep {
		if q.Who == username {
			byUser = append(byUser, q)
		}
	}

	return byUser, nil

}

func (rep *fakeRepository) GetQuestionByID(ctx context.Context, id string) (*qsrv.Question, error) {
	for _, q := range *rep {
		if q.ID == id {
			return q, nil
		}
	}

	return nil, fmt.Errorf("Question with ID %v not found", id)
}

func (rep *fakeRepository) ChangeQuestion(ctx context.Context, newq qsrv.Question) error {
	for i, q := range *rep {
		if newq.ID == q.ID {
			(*rep)[i] = &newq
			return nil
		}
	}

	return fmt.Errorf("Question with ID %v not found", newq.ID)
}

func (rep *fakeRepository) DeleteQuestion(ctx context.Context, id string) error {
	for i, q := range *rep {
		if q.ID == id {
			newrep := append((*rep)[:i], (*rep)[i+1:]...)
			rep = &newrep
			return nil
		}
	}

	return fmt.Errorf("Question with ID %v not found", id)
}
