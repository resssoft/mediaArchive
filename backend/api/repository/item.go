package repository

import (
	"fmt"
	"github.com/resssoft/mediaArchive/database"
	"github.com/resssoft/mediaArchive/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const itemCollectionName = "item"

type ItemRepository interface {
	Add(model.Item) error
	GetItemByID(string) (model.Item, error)
	List(string, interface{}) ([]*model.Item, error)
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

func (r *itemRepo) Add(item model.Item) error {
	log.Info().Interface("new item", item).Send()
	item.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(r.dbApp.GetContext(), item)
	if err != nil {
		log.Error().AnErr("Insert item error", err).Send()
		return err
	}
	return nil
}

func (r *itemRepo) getByField(name string, value interface{}) (model.Item, error) {

	item := model.Item{}
	filter := bson.M{name: value}
	err := r.collection.FindOne(r.dbApp.GetContext(), filter).Decode(&item)
	if err != nil {
		log.Error().AnErr("item read error", err).Send()
		return item, err
	}
	return item, nil
}

func (r *itemRepo) List(name string, value interface{}) ([]*model.Item, error) {
	options := options.Find()
	items := make([]*model.Item, 0)
	//filter := bson.M{name: value}
	filter := bson.M{}
	cur, err := r.collection.Find(r.dbApp.GetContext(), filter, options)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	for cur.Next(r.dbApp.GetContext()) {
		var item model.Item
		err := cur.Decode(&item)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		log.Info().Interface("item", item).Send()
		items = append(items, &item)
	}

	log.Info().Interface("items", items).Send()
	if err := cur.Err(); err != nil {
		log.Fatal().Err(err).Send()
	}
	cur.Close(r.dbApp.GetContext())

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", items)
	return items, nil
}

func (r *itemRepo) GetItemByID(id string) (model.Item, error) {
	item, err := r.getByField("_id", id)
	return item, err
}
