package auth

import (
	"encoding/json"
	"github.com/mileusna/useragent"
	config "github.com/resssoft/mediaArchive/configuration"
	"github.com/resssoft/mediaArchive/database"
	"github.com/resssoft/mediaArchive/models"
	"github.com/twinj/uuid"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net"
	"time"
)

const sessionCollectionName = "session"

var mongoDbApp database.MongoClientApplication
var collection *mongo.Collection

type Session struct {
	ID        string    `bson:"token"`
	UserId    string    `bson:"user_id"`
	UserLang  string    `bson:"user_lang"`
	Perms     string    `bson:"perms"`
	IP        net.IP    `bson:"ip"`
	OS        string    `bson:"os"`
	Device    string    `bson:"device"`
	Browser   string    `bson:"browser"`
	UserAgent string    `bson:"user_agent"`
	ExpiredAt time.Time `bson:"expired_at"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func init() {
	var err error
	mongoDbApp, err = database.ProvideMongo()
	if err != nil {
		return
	}
	collection = mongoDbApp.GetCollection(sessionCollectionName)
	mongoDbApp.CreateUniqueIndex(collection, "token")
	mongoDbApp.CreateIndexWithTimeout(collection, "expired_at", 0)
}

func (s *Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Session) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *Session) IsValid(ctx *fasthttp.RequestCtx) bool {
	if time.Now().After(s.ExpiredAt) {
		return false
	}
	userAgent := ua.Parse(string(ctx.UserAgent()))
	return s.OS == userAgent.OS &&
		s.Browser == userAgent.Name &&
		s.Device == userAgent.Device
}

func NewSessionByRequest(ctx *fasthttp.RequestCtx, user models.User) *Session {
	now := time.Now()
	userAgent := ua.Parse(string(ctx.UserAgent()))
	return &Session{
		ID:        uuid.NewV4().String(),
		UserId:    user.Id.Hex(),
		Perms:     user.Role.PermsToString(),
		UserLang:  user.Lang,
		IP:        ctx.RemoteIP(),
		OS:        userAgent.OS,
		Device:    userAgent.Device,
		Browser:   userAgent.Name,
		UserAgent: userAgent.String,
		ExpiredAt: now.Add(config.JwtRtExpires()),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func SaveSession(session *Session) error {
	_, err := collection.InsertOne(mongoDbApp.GetContext(), session)
	return err
}

func GetSession(sessionId string) (*Session, error) {
	session := new(Session)
	filter := bson.M{"token": sessionId}
	err := collection.FindOne(mongoDbApp.GetContext(), filter).Decode(&session)
	return session, err
}

func UpdateSessionExpiration(session *Session, expiration time.Duration) (*Session, error) {
	now := time.Now()
	session.UpdatedAt = now
	session.ExpiredAt = now.Add(expiration)
	filter := bson.M{"token": session.ID}
	update := bson.M{"$set": bson.M{"updated_at": session.UpdatedAt, "expired_at": session.ExpiredAt}}
	_, err := collection.UpdateOne(mongoDbApp.GetContext(), filter, update)
	return session, err
}

func DeleteSessionById(sessionId string) error {
	filter := bson.M{"token": sessionId}
	_, err := collection.DeleteOne(mongoDbApp.GetContext(), filter)
	return err
}
