package configuration

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"strings"
)

var Version = "0.0.1"

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AllowEmptyEnv(true)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Info().Msg("Unable to read config file")
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
