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

const configCollectionName = "config"

type ConfigRepository interface {
	Add(models.Config) error
	GetItemByID(string) (models.Config, error)
	List(models.DataFilter) ([]models.Config, error)
	Update(string, models.Config) error
}

type configRepo struct {
	dbApp      database.MongoClientApplication
	collection *mongo.Collection
}

func NewConfigRepo(db database.MongoClientApplication) ConfigRepository {
	collection := db.GetCollection(configCollectionName)
	db.CreateUniqueIndex(collection, "name", "group", "perm")
	return &configRepo{
		dbApp:      db,
		collection: collection,
	}
}

func (r *configRepo) Add(item models.Config) error {
	item.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(r.dbApp.GetContext(), item)
	if err != nil {
		log.Error().AnErr("Insert item error", err).Send()
		return err
	}
	return nil
}

func (r *configRepo) getByField(name string, value interface{}) (models.Config, error) {

	item := models.Config{}
	filter := bson.M{name: value}
	err := r.collection.FindOne(r.dbApp.GetContext(), filter).Decode(&item)
	if err != nil {
		log.Error().AnErr("item read error", err).Send()
		return item, err
	}
	return item, nil
}

func (r *configRepo) getByFields(values map[string]interface{}) (models.Config, error) {
	var findQuery []bson.M
	item := models.Config{}
	//filter := bson.M{values...}
	err := r.collection.FindOne(r.dbApp.GetContext(), findQuery).Decode(&item)
	if err != nil {
		log.Error().AnErr("item read error", err).Send()
		return item, err
	}
	return item, nil
}

func (r *configRepo) List(filter models.DataFilter) ([]models.Config, error) {
	options := options.Find()
	items := make([]models.Config, 0)
	mongoFilter := bson.D(filter.Data)
	cur, err := r.collection.Find(r.dbApp.GetContext(), mongoFilter, options)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	for cur.Next(r.dbApp.GetContext()) {
		var item models.Config
		err := cur.Decode(&item)
		if err != nil {
			log.Error().Err(err).Send()
			continue
		}
		items = append(items, item)
	}

	log.Info().Interface("configs", items).Send()
	if err := cur.Err(); err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}
	cur.Close(r.dbApp.GetContext())
	return items, nil
}

func (r *configRepo) GetItemByID(id string) (models.Config, error) {
	item, err := r.getByField("_id", id)
	return item, err
}

func (r *configRepo) Update(id string, item models.Config) error {
	log.Info().Msg(id)
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.UpdateOne(
		r.dbApp.GetContext(),
		bson.M{"_id": idObj},
		bson.D{
			{"$set", item},
		})
	if err != nil {
		log.Error().AnErr("update error", err).Send()
		return err
	}
	return nil
}
