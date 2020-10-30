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

const itemGroupCollectionName = "item_group"

type ItemGroupRepository interface {
	Add(models.ItemGroup) error
	GetItemByID(string) (models.ItemGroup, error)
	List(...models.DataFilter) ([]*models.ItemGroup, error)
}

type itemGroupRepo struct {
	dbApp      database.MongoClientApplication
	collection *mongo.Collection
}

func NewItemGroupRepo(db database.MongoClientApplication) ItemGroupRepository {
	collection := db.GetCollection(itemGroupCollectionName)
	return &itemGroupRepo{
		dbApp:      db,
		collection: collection,
	}
}

func (r *itemGroupRepo) Add(itemGroup models.ItemGroup) error {
	itemGroup.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(r.dbApp.GetContext(), itemGroup)
	if err != nil {
		log.Error().AnErr("Insert itemGroup error", err).Send()
		return err
	}
	return nil
}

func (r *itemGroupRepo) getByField(name string, value interface{}) (models.ItemGroup, error) {

	itemGroup := models.ItemGroup{}
	filter := bson.M{name: value}
	err := r.collection.FindOne(r.dbApp.GetContext(), filter).Decode(&itemGroup)
	if err != nil {
		log.Error().AnErr("itemGroup read error", err).Send()
		return itemGroup, err
	}
	return itemGroup, nil
}

func (r *itemGroupRepo) List(filter ...models.DataFilter) ([]*models.ItemGroup, error) {
	options := options.Find()
	itemGroups := make([]*models.ItemGroup, 0)
	//filter := bson.M{name: value}
	mongoFilter := bson.M{}
	cur, err := r.collection.Find(r.dbApp.GetContext(), mongoFilter, options)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	for cur.Next(r.dbApp.GetContext()) {
		var item models.ItemGroup
		err := cur.Decode(&item)
		if err != nil {
			log.Error().Err(err).Send()
			continue
		}
		itemGroups = append(itemGroups, &item)
	}

	log.Info().Interface("itemGroups", itemGroups).Send()
	if err := cur.Err(); err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}
	cur.Close(r.dbApp.GetContext())
	return itemGroups, nil
}

func (r *itemGroupRepo) GetItemByID(id string) (models.ItemGroup, error) {
	itemGroup, err := r.getByField("_id", id)
	return itemGroup, err
}
