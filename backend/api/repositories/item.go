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

const itemCollectionName = "item"

type ItemRepository interface {
	Add(models.Item) error
	GetItemByID(string) (models.Item, error)
	List(models.DataFilter) ([]*models.Item, error)
}

type itemRepo struct {
	dbApp      database.MongoClientApplication
	collection *mongo.Collection
}

func NewItemRepo(db database.MongoClientApplication) ItemRepository {
	collection := db.GetCollection(itemCollectionName)
	return &itemRepo{
		dbApp:      db,
		collection: collection,
	}
}

func (r *itemRepo) Add(item models.Item) error {
	item.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(r.dbApp.GetContext(), item)
	if err != nil {
		log.Error().AnErr("Insert item error", err).Send()
		return err
	}
	return nil
}

func (r *itemRepo) getByField(name string, value interface{}) (models.Item, error) {

	item := models.Item{}
	filter := bson.M{name: value}
	err := r.collection.FindOne(r.dbApp.GetContext(), filter).Decode(&item)
	if err != nil {
		log.Error().AnErr("item read error", err).Send()
		return item, err
	}
	return item, nil
}

func (r *itemRepo) List(filter models.DataFilter) ([]*models.Item, error) {
	options := options.Find()
	items := make([]*models.Item, 0)
	//filter := bson.M{name: value}
	mongoFilter := bson.D(filter.Data)
	cur, err := r.collection.Find(r.dbApp.GetContext(), mongoFilter, options)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	for cur.Next(r.dbApp.GetContext()) {
		var item models.Item
		err := cur.Decode(&item)
		if err != nil {
			log.Error().Err(err).Send()
			continue
		}
		items = append(items, &item)
	}

	log.Info().Interface("items", items).Send()
	if err := cur.Err(); err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}
	cur.Close(r.dbApp.GetContext())
	return items, nil
}

func (r *itemRepo) GetItemByID(id string) (models.Item, error) {
	item, err := r.getByField("_id", id)
	return item, err
}
