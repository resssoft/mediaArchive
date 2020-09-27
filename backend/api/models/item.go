package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	ItemForMusic         ItemAssignment = "music"
	ItemForMusicAndVideo ItemAssignment = "musicAndVideo"
	ItemForVideo         ItemAssignment = "video"
	ItemForNothing       ItemAssignment = ""
)

type ItemAssignment string

type Item struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             `bson:"name"`
	Assignment  ItemAssignment     `bson:"assignment"`
	Image       string             `bson:"image"`
	URL         string             `bson:"url"`
	Tags        []string           `bson:"tags"`
	Categories  []string           `bson:"group"`
	Cache       string             `bson:"cache"`
	Media       []ItemMedia        `bson:"media"`
	Error       string             `bson:"error"`
	Service     string             `bson:"Service"`
	ServiceData interface{}        `bson:"service_data"`
}

type ItemMedia struct {
	Type      string `bson:"type"`
	RemoteUrl string `bson:"remoteUrl"`
	LocalUrl  string `bson:"localUrl"`
}
