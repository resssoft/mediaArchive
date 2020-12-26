package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ItemGroup struct {
	ID         primitive.ObjectID `bson:"_id"`
	Code       string             `bson:"code"`
	ParentCode string             `bson:"parent_code"`
	Name       string             `bson:"name"`
	UserID     string             `bson:"user_id"`
	Sort       int                `bson:"sort"`
	Flags      []Flag             `bson:"flags"`
	Depth      int                `bson:"depth"`
}

type ItemGroupFlat struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Depth int    `json:"depth"`
}

func (ig *ItemGroup) ToFlat() ItemGroupFlat {
	return ItemGroupFlat{
		Code:  ig.Code,
		Name:  ig.Name,
		Depth: ig.Depth,
	}
}
