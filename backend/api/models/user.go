package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
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
	Lang     string             `bson:"lang"`
	Phone    string             `bson:"phone"`
	Password string             `bson:"password"`
	Settings UserSettings       `bson:"settings"`
	Role     UserRole           `bson:"role"`
}

type UserData struct {
	Id       string       `json:"id"`
	Email    string       `json:"email"`
	Lang     string       `json:"lang"`
	Settings UserSettings `json:"settings"`
	Role     UserRole     `json:"role"`
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

func (r *UserRole) PermsToString() string {
	permissions := make([]string, 0)
	for _, item := range r.Permissions {
		permissions = append(permissions, string(item))
	}
	return strings.Join(permissions, ",")
}

func (r *UserRole) PermsFromString(perms string) {
	permissions := make([]UserPermission, 0)
	for _, item := range perms {
		permissions = append(permissions, UserPermission(item))
	}
	r.Permissions = permissions
}

func (u *User) Data() UserData {
	return UserData{
		Id:       u.Id.Hex(),
		Lang:     u.Lang,
		Email:    u.Email,
		Settings: u.Settings,
		Role:     u.Role,
	}
}
