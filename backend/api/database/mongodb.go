package database

import (
	"context"
	"errors"
	config "github.com/resssoft/mediaArchive/configuration"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

var (
	onceMongoAction sync.Once
	mongoClient     *mongo.Client
	mongoContext    context.Context
)

type MongoClientApplication interface {
	GetCollection(string) *mongo.Collection
	GetContext() context.Context
	CreateUniqueIndex(*mongo.Collection, ...string)
	CreateIndexWithTimeout(*mongo.Collection, string, int32)
}

type mongoClientOriginal struct {
	client  *mongo.Client
	context context.Context
	dbName  string
}

func ProvideMongo() (MongoClientApplication, error) {
	onceMongoAction.Do(func() {
		configureMongo(config.MongoUrl())
	})
	if mongoClient == nil || mongoContext == nil {
		return &mongoClientOriginal{}, errors.New("mongo client or context is empty")
	}
	return &mongoClientOriginal{
		client:  mongoClient,
		context: mongoContext,
		dbName:  config.MongoUrl(),
	}, nil
}

func (r *mongoClientOriginal) GetCollection(collection string) *mongo.Collection {
	return mongoClient.Database(config.MongoDbName()).Collection(collection)
}

func (r *mongoClientOriginal) GetContext() context.Context {
	return r.context
}

func (r *mongoClientOriginal) CreateUniqueIndex(collection *mongo.Collection, keyList ...string) {
	keys := bson.M{}
	for _, key := range keyList {
		keys[key] = 1
	}
	indexName, err := collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    keys,
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Error().AnErr("create index", err).Send()
	}
	log.Info().Interface("indexName", indexName).Send()
}

func (r *mongoClientOriginal) CreateIndexWithTimeout(collection *mongo.Collection, key string, expireTimeSeconds int32) {
	indexName, err := collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				key: 1,
			},
			Options: options.Index().SetExpireAfterSeconds(expireTimeSeconds),
		},
	)
	if err != nil {
		log.Error().AnErr("create index", err).Send()
	}
	log.Info().Interface("indexName", indexName).Send()
}

func configureMongo(address string) {
	var err error
	log.Info().Msg(config.MongoUrl())
	mongoContext = context.Background()
	clientOptions := options.Client().ApplyURI(address)
	mongoClient, err = mongo.Connect(mongoContext, clientOptions)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot connect to mongo")
	}
	err = mongoClient.Ping(mongoContext, nil)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot connect to mongo")
	}
}
