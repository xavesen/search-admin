package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/xavesen/search-admin/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	client 			*mongo.Client
	database 		*mongo.Database
	usersCollecton	*mongo.Collection
}

func NewMongoStorage(ctx context.Context, addr string, db string, user string, password string) (*MongoStorage, error) {
	clientCreds := options.Credential{
		Username: user,
		Password: password,
		AuthSource: db,
	}
	clientOpts := options.Client()
	clientOpts.SetAuth(clientCreds)

	newClient, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		 return nil, err
	}

	if err = newClient.Ping(ctx, nil); err != nil {
		return nil, err
	}

	appDb := newClient.Database("search_app")
	usersCol := appDb.Collection("users")

	newStorage := &MongoStorage{
		client: newClient,
		database: appDb,
		usersCollecton: usersCol,
	}

	return newStorage, nil
}

func getOid(supposedOid interface{}) (string, bool) {
	if oid, ok := supposedOid.(primitive.ObjectID); ok {
		return oid.Hex(), true
	}

	return "", false
}

func (s *MongoStorage) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	result, err := s.usersCollecton.InsertOne(ctx, user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	id, ok := getOid(result.InsertedID)

	if !ok {
		return nil, errors.New("db did not return object id")
	}

	user.Id = id

	return user, nil
}