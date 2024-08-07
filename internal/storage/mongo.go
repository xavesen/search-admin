package storage

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/xavesen/search-admin/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	client 			*mongo.Client
	database 		*mongo.Database
	usersCollection	*mongo.Collection
}

func NewMongoStorage(ctx context.Context, addr string, db string, user string, password string) (*MongoStorage, error) {
	// FIXME: actually use addr!
	log.Infof("Initializing client and connecting mongo db %s on %s with user %s", db, addr, user)

	clientCreds := options.Credential{
		Username: user,
		Password: password,
		AuthSource: db,
	}
	clientOpts := options.Client()
	clientOpts.SetAuth(clientCreds)

	log.Debug("Initializing mongo client")
	newClient, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Errorf("Error while initializing mongo client for db %s on %s with user %s: %s", db, addr, user, err.Error())
		return nil, err
	}

	log.Debug("Connecting mongo db")
	if err = newClient.Ping(ctx, nil); err != nil {
		log.Errorf("Error while connecting mongo db %s on %s with user %s: %s", db, addr, user, err.Error())
		return nil, err
	}

	log.Debug("Initializing db and collections")
	appDb := newClient.Database(db)
	usersCol := appDb.Collection("users")

	newStorage := &MongoStorage{
		client: newClient,
		database: appDb,
		usersCollection: usersCol,
	}

	log.Info("Successfully initialized and connected mongo db")
	return newStorage, nil
}

func getOid(supposedOid interface{}) (string, bool) {
	log.Debug("Getting object id")
	if oid, ok := supposedOid.(primitive.ObjectID); ok {
		hexOid := oid.Hex()
		log.Debugf("Got object id %s", hexOid)
		return hexOid, true
	}

	return "", false
}

func (s *MongoStorage) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	log.Debugf("Inserting user %s to db", user)

	result, err := s.usersCollection.InsertOne(ctx, user)
	if err != nil {
		log.Errorf("Error inserting user %s to db: %s", user, err.Error())
		return nil, err
	}

	id, ok := getOid(result.InsertedID)
	if !ok {
		log.Errorf("Unable to get oid from interface returned by db after trying to insert user %s", user)
		return nil, errors.New("db did not return object id")
	}

	user.Id = id

	log.Debugf("Successfully inserted user %s to db", user)
	return user, nil
}

func (s *MongoStorage) GetUser(ctx context.Context, id string) (*models.User, error) {
	log.Debugf("Searching for user with id %s in db", id)
	var user *models.User

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warningf("Error converting id string %s to object id while searching for user in db: %s", id, err.Error())
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: oid}}

	if err := s.usersCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Warningf("Tried to find in db non-existent user with id %s ", id)
		} else {
			log.Errorf("Error searching for user with id %s in db: %s", id, err.Error())
		}
		return nil, err
	}

	log.Debugf("Successfully found user with id %s in db: %s", id, user)
	return user, nil
}

func (s *MongoStorage) GetAllUsers(ctx context.Context) ([]models.User, error) {
	log.Debug("Getting all users from db")
	users := []models.User{}

	filter := bson.D{{}}
	cur, err := s.usersCollection.Find(ctx, filter)
	if err != nil {
		log.Errorf("Error finding all users in db: %s", err.Error())
		return users, err
	}

	if err = cur.All(ctx, &users); err != nil {
		log.Errorf("Error iterating and decoding all users from db: %s", err.Error())
		return users, err
	}

	log.Debug("Successfully got all users from db")
	if log.IsLevelEnabled(log.TraceLevel) {
		usersString := ""
		for _, user := range users {
			usersString = usersString + user.String() + ", "
		}
		log.Tracef("Users from db: [%s]", usersString)
	}

	return users, nil
}

func (s *MongoStorage) DeleteUser(ctx context.Context, id string) error {
	log.Debugf("Deleting user with id %s", id)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warningf("Error converting id string %s to object id while deleting user from db: %s", id, err.Error())
		return err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	
	result, err := s.usersCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Errorf("Error deleting user with id %s from db: %s", id, err.Error())
		return err
	} else if result.DeletedCount < 1 {
		log.Warningf("Tried to delete from db non-existent user with id %s ", id)
		return mongo.ErrNoDocuments
	}

	log.Debugf("Successfully deleted user with id %s from db", id)
	return nil
}