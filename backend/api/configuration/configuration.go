package configuration

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const ImportDir = "./uploads/import/"
const ExportDir = "./uploads/export/"

var Version = "0.0.1"

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AllowEmptyEnv(true)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".\\backend\\api")
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
