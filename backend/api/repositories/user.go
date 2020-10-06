package repositories

import (
	"github.com/resssoft/mediaArchive/database"
	"github.com/resssoft/mediaArchive/models"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userCollectionName = "user"

type UserRepository interface {
	Add(models.User) error
	Update(models.User) error
	GetByID(string) (models.User, error)
	List(string, interface{}) ([]*models.User, error)
	GetByEmail(string) (models.User, error)
}

type userRepo struct {
	dbApp      database.MongoClientApplication
	collection *mongo.Collection
}

func NewUserRepo(db database.MongoClientApplication) UserRepository {
	collection := db.GetCollection(userCollectionName)
	db.CreateUniqueIndex(collection, "email")
	return &userRepo{
		dbApp:      db,
		collection: collection,
	}
}

func (r *userRepo) Add(user models.User) error {
	user.Id = primitive.NewObjectID()
	_, err := r.collection.InsertOne(r.dbApp.GetContext(), user)
	if err != nil {
		log.Error().AnErr("Insert user error", err).Send()
		return err
	}
	return nil
}

func (r *userRepo) Update(user models.User) error {
	filter := bson.M{"_id": user.Id}
	_, err := r.collection.UpdateOne(r.dbApp.GetContext(), filter, user)
	if err != nil {
		log.Error().AnErr("Update user error", err).Send()
		return err
	}
	return nil
}

func (r *userRepo) getByField(name string, value interface{}) (models.User, error) {
	user := models.User{}
	filter := bson.M{name: value}
	err := r.collection.FindOne(r.dbApp.GetContext(), filter).Decode(&user)
	if err != nil {
		log.Error().AnErr("user read error", err).Send()
		return user, err
	}
	return user, nil
}

func (r *userRepo) List(name string, value interface{}) ([]*models.User, error) {
	options := options.Find()
	users := make([]*models.User, 0)
	//filter := bson.M{name: value}
	filter := bson.M{}
	cur, err := r.collection.Find(r.dbApp.GetContext(), filter, options)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	for cur.Next(r.dbApp.GetContext()) {
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			log.Error().Err(err).Send()
			continue
		}
		users = append(users, &user)
	}

	log.Info().Interface("users", users).Send()
	if err := cur.Err(); err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}
	cur.Close(r.dbApp.GetContext())
	return users, nil
}

func (r *userRepo) GetByID(id string) (models.User, error) {
	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, err
	}
	user, err := r.getByField("_id", ObjId)
	return user, err
}

func (r *userRepo) GetByEmail(email string) (models.User, error) {
	user, err := r.getByField("email", email)
	return user, err
}
