package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	ID    primitive.ObjectID `bson:"_id"`
	Name  string             `bson:"name"`
	Type  string             `bson:"type"`
	Image string             `bson:"image"`
	URL   string             `bson:"url"`
	Tags  []string           `bson:"tags"`
	Group string             `bson:"group"`
	File  string             `bson:"file"`
}
