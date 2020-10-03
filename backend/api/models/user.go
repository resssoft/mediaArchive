package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserPermission string

const (
	OwnerPerm          UserPermission = "owner"
	FileUploadPerm     UserPermission = "file-upload"
	CoubPreloadPerm    UserPermission = "coub-preload"
	WebPagePreviewPerm UserPermission = "web-page-preview"
	FileSize100mbPerm  UserPermission = "files-100mb"
	FileSize500mbPerm  UserPermission = "files-500mb"
	VoiceAnalysePerm   UserPermission = "voice-analyse"
	SmsPerm            UserPermission = "sms"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id"`
	Email    string             `bson:"email"`
	Phone    string             `bson:"phone"`
	Password string             `bson:"password"`
	Settings UserSettings       `bson:"settings"`
	Role     UserRole           `bson:"role"`
}

type UserData struct {
	Id       string       `bson:"id"`
	Email    string       `bson:"email"`
	Settings UserSettings `bson:"settings"`
	Role     UserRole     `bson:"role"`
}

type UserFriends struct {
	UserId string `bson:"user_id"`
}

type UserSettings struct {
	Id                string `bson:"id"`
	CoubProcessing    bool   `bson:"coub_processing"`
	VimeoProcessing   bool   `bson:"vimeo_processing"`
	YoutubeProcessing bool   `bson:"youtube_processing"`
	Language          string `bson:"language"`
}

type UserRole struct {
	Id          int              `bson:"id"`
	Name        string           `bson:"name"`
	Permissions []UserPermission `bson:"permissions"`
}

type UserNotice struct {
	UserId  string    `bson:"user_id"`
	Created time.Time `bson:"created"`
	Message string    `bson:"message"`
	Type    string    `bson:"type"`
	Source  string    `bson:"source"`
}

func (r *UserRole) CheckPerm(perm UserPermission) bool {
	for _, item := range r.Permissions {
		if item == perm {
			return true
		}
	}
	return false
}
