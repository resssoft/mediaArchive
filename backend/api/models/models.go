package models

type Flag struct {
	Name   string `bson:"name"`
	Status bool   `bson:"status"`
}
