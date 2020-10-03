package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	ItemForMusic         ItemAssignment = "music"
	ItemForMusicAndVideo ItemAssignment = "musicAndVideo"
	ItemForVideo         ItemAssignment = "video"
	ItemForPicture       ItemAssignment = "picture"
	ItemForWebPage       ItemAssignment = "webPage"
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
	Categories  []string           `bson:"categories"`
	Groups      []string           `bson:"groups"`
	Cache       string             `bson:"cache"`
	Media       []ItemMedia        `bson:"media"`
	Error       string             `bson:"error"`
	Favorite    bool               `bson:"favorite"`
	Viewed      bool               `bson:"viewed"`
	Service     string             `bson:"service"`
	ServiceData interface{}        `bson:"service_data"`
}

type ItemMedia struct {
	Type      string `bson:"type"`
	RemoteUrl string `bson:"remoteUrl"`
	LocalUrl  string `bson:"localUrl"`
	HashSum   string `bson:"hashSum"`
	HashAlg   string `bson:"hashAlg"`
	Size      int64  `bson:"size"`
}
