package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/log/level"
	qsrv "github.com/macchid/questandans/questions"
	"github.com/macchid/questandans/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *MongoDBRepo) CreateQuestion(ctx context.Context, q qsrv.Question) error {
	logger, logEnd := utils.LogStart(r.logger, "MongoDBRepo::CreateQuestion")
	defer logEnd(time.Now())
	_, err := r.questions.InsertOne(ctx, q)
	if err != nil {
		level.Error(logger).Log("err", err)
		return qsrv.ErrCmdRepository
	}

	return nil
}

func (r *MongoDBRepo) GetQuestions(ctx context.Context) ([]*qsrv.Question, error) {
	_, logEnd := utils.LogStart(r.logger, "MongoDBRepo::GetQuestions")
	defer logEnd(time.Now())
	filter := bson.D{{}}
	return r.search(ctx, filter)
}

func (r *MongoDBRepo) GetQuestionsByUser(ctx context.Context, username string) ([]*qsrv.Question, error) {
	_, logEnd := utils.LogStart(r.logger, "MongoDBRepo::GetQuestionsByUser")
	defer logEnd(time.Now())

	filter := bson.D{bson.E{Key: "who", Value: username}}
	return r.search(ctx, filter)
}

func (r *MongoDBRepo) GetQuestionByID(ctx context.Context, id string) (*qsrv.Question, error) {
	logger, logEnd := utils.LogStart(r.logger, "MongoDBRepo::GetQuestionByID")
	defer logEnd(time.Now())
	filter := bson.D{bson.E{Key: "id", Value: id}}
	q := &qsrv.Question{}

	cur := r.questions.FindOne(ctx, filter)
	if err := cur.Err(); err != nil {
		level.Error(logger).Log("err", err)
		if err == mongo.ErrNoDocuments {
			return nil, qsrv.ErrQuestionNotFound
		}
	}

	err := cur.Decode(q)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, qsrv.ErrQueryRepository
	}

	return q, nil

}

func (r *MongoDBRepo) ChangeQuestion(ctx context.Context, q qsrv.Question) error {
	logger, logEnd := utils.LogStart(r.logger, "MongoDBRepo::ChangeQuestion")
	defer logEnd(time.Now())

	d, err := toDoc(q)
	if err != nil {
		level.Error(logger).Log("err", err)
		return qsrv.ErrCmdRepository
	}

	filter := bson.D{bson.E{Key: "id", Value: q.ID}}

	_, err = r.questions.ReplaceOne(ctx, filter, d)
	if err != nil {
		level.Error(logger).Log("err", err)
		return qsrv.ErrCmdRepository
	}

	return nil
}

func (r *MongoDBRepo) DeleteQuestion(ctx context.Context, id string) error {
	logger, logEnd := utils.LogStart(r.logger, "MongoDBRepo::DeleteQuestion")
	defer logEnd(time.Now())
	filter := bson.D{bson.E{Key: "id", Value: id}}

	res, err := r.questions.DeleteOne(ctx, filter)
	if err != nil {
		level.Error(logger).Log("err", err)
		return qsrv.ErrCmdRepository
	}

	if res.DeletedCount == 0 {
		level.Error(logger).Log("err", fmt.Errorf("No question was found with filter %v", filter))
		return qsrv.ErrQuestionNotFound
	}

	return nil
}

func (r *MongoDBRepo) search(ctx context.Context, filter interface{}) ([]*qsrv.Question, error) {
	logger, logEnd := utils.LogStart(r.logger, "MongoDBRepo::search")
	defer logEnd(time.Now())

	qs := make([]*qsrv.Question, 1)

	cur, err := r.questions.Find(ctx, filter)
	if err != nil {
		level.Error(logger).Log("err", err)
		return qs, qsrv.ErrQueryRepository
	}

	for cur.Next(ctx) {
		var q qsrv.Question
		err := cur.Decode(&q)
		if err != nil {
			level.Error(logger).Log("err", err)
			return qs, qsrv.ErrQueryRepository
		}

		qs = append(qs, &q)
	}

	if err := cur.Err(); err != nil {
		level.Error(logger).Log("err", err)
		return qs, qsrv.ErrQueryRepository
	}

	cur.Close(ctx)
	if len(qs) == 0 {
		level.Error(logger).Log("err", fmt.Errorf("No question was found with filter %v", filter))
		return qs, qsrv.ErrQuestionNotFound
	}

	return qs, nil
}
