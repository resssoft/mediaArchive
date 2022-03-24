package configuration

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

const (
	ImportDir          = "./uploads/import/"
	ExportDir          = "./uploads/export/"
	TranslationDir     = "./translations/"
	DateTimeFormat     = "2006-01-02T15:04:05Z07:00"
	DateTimeFlatFormat = "20060102150405"
	Version            = "0.0.1001"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AllowEmptyEnv(true)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	//TODO:change this to .
	viper.AddConfigPath(".")
	viper.AddConfigPath("./backend/api")
	err := viper.ReadInConfig()
	if err != nil {
		log.Info().Msg("Unable to read config file")
	}
	err = os.MkdirAll(ImportDir, 0666)
	if err != nil {
		log.Error().AnErr("Cant create import dir", err).Msg(ImportDir)
	}
	err = os.MkdirAll(ExportDir, 0666)
	if err != nil {
		log.Error().AnErr("Cant create export dir", err).Msg(ExportDir)
	}
}

func GetString(key string) string {
	return viper.GetString(key)
}

func MongoUrl() string {
	return fmt.Sprintf("mongodb://%s:%s@%s",
		viper.GetString("db.mongoDb.user"),
		viper.GetString("db.mongoDb.password"),
		viper.GetString("db.mongoDb.url"),
	)
}

func MongoDbName() string {
	return viper.GetString("db.mongoDb.dbname")
}

func ApiUrl() string {
	return viper.GetString("address.api")
}

func PasswordSalt() string {
	return viper.GetString("security.salt")
}

func JwtAtExpires() time.Duration {
	return viper.GetDuration("Jwt.AtExpires")
}

func JwtRtExpires() time.Duration {
	return viper.GetDuration("Jwt.RtExpires")
}

func JwtSecretAccess() []byte {
	return []byte(viper.GetString("Jwt.SecretAccess"))
}

func AdminLogin() string {
	return viper.GetString("admin.login")
}

func AdminPassword() string {
	return viper.GetString("admin.password")
}
