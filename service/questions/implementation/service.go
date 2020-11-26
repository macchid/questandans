package implementation

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gofrs/uuid"
	qsrv "github.com/macchid/questandans/questions"
	"github.com/macchid/questandans/utils"
)

type service struct {
	repository qsrv.Repository
	logger     log.Logger
}

/*
NewService instantiates a new questions.Service instance to be used by the application.
The newly created questions.Service stores a Repository and a Logger instance.
*/
func NewService(rep qsrv.Repository, logger log.Logger) qsrv.Service {
	return &service{
		repository: rep,
		logger:     logger,
	}
}

func (s *service) Add(ctx context.Context, q qsrv.Question) (string, error) {
	logger, logEnd := utils.LogStart(s.logger, "Service::Add")
	defer logEnd(time.Now())

	uuid4, _ := uuid.NewV4()

	q.ID = uuid4.String()

	err := s.repository.CreateQuestion(ctx, q)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", qsrv.ErrCmdRepository
	}

	return q.ID, nil
}

func (s *service) Read(ctx context.Context) ([]*qsrv.Question, error) {
	logger, logEnd := utils.LogStart(s.logger, "Service::Read")
	defer logEnd(time.Now())

	questions, err := s.repository.GetQuestions(ctx)
	if err != nil {
		level.Error(logger).Log("err", err)
		if err == qsrv.ErrQuestionNotFound {
			return questions, err
		}

		return nil, err
	}

	return questions, nil
}

func (s *service) ReadOwned(ctx context.Context, user string) ([]*qsrv.Question, error) {
	logger, logEnd := utils.LogStart(s.logger, "Service::ReadOwned")
	defer logEnd(time.Now())

	questions, err := s.repository.GetQuestionsByUser(ctx, user)
	if err != nil {
		level.Error(logger).Log("err", err)
		if err == qsrv.ErrQuestionNotFound {
			return questions, err
		}

		return nil, err
	}

	return questions, nil
}

func (s *service) Peek(ctx context.Context, id string) (*qsrv.Question, error) {
	logger, logEnd := utils.LogStart(s.logger, "Service::Peek")
	defer logEnd(time.Now())

	question, err := s.repository.GetQuestionByID(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	return question, nil
}

func (s *service) Change(ctx context.Context, id, statement string) (*qsrv.Question, error) {
	logger, logEnd := utils.LogStart(s.logger, "Service::Change")
	defer logEnd(time.Now())

	question, err := s.repository.GetQuestionByID(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	question.Statement = statement
	if err := s.repository.ChangeQuestion(ctx, *question); err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	return question, nil
}

func (s *service) Answer(ctx context.Context, id, answer string) (*qsrv.Question, error) {
	logger, logEnd := utils.LogStart(s.logger, "Service::Answer")
	defer logEnd(time.Now())

	question, err := s.repository.GetQuestionByID(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	question.Answer = answer
	if err := s.repository.ChangeQuestion(ctx, *question); err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	return question, nil
}

func (s *service) Delete(ctx context.Context, id string) (*qsrv.Question, error) {
	logger, logEnd := utils.LogStart(s.logger, "Service::Delete")
	defer logEnd(time.Now())

	question, err := s.repository.GetQuestionByID(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	if err := s.repository.DeleteQuestion(ctx, id); err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	return question, nil
}

func (s *service) Logger() log.Logger {
	return s.logger
}
