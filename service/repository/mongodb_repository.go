package repository

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//NeQuestionsRepository establishes a connection to MongoDB and creates a new Repository
func NewRepository(ctx context.Context, parent log.Logger, uri string) (*MongoDBRepo, error) {
	logger := log.With(parent, "method", "NewQuestionsRepository")

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &MongoDBRepo{
		questions: client.Database("qanda").Collection("questions"),
		logger:    parent,
	}, nil
}

type MongoDBRepo struct {
	questions *mongo.Collection
	logger    log.Logger
}

func toDoc(v interface{}) (*bson.D, error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}

	doc := bson.D{}
	err = bson.Unmarshal(data, &doc)
	return &doc, nil
}
