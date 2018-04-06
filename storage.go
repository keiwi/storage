package storage

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/keiwi/storage/database"
	"github.com/keiwi/utils/log"
	"github.com/keiwi/utils/log/handlers/cli"
	"github.com/keiwi/utils/log/handlers/file"
	"github.com/nats-io/go-nats"
	"github.com/spf13/viper"
)

var (
	configType string
)

type Storage struct{}

func (s Storage) StartStorage() {
	ReadConfig()

	log.Info("Trying to connect to database")
	db, err := database.NewDatabase(
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.ip"),
		viper.GetString("database.port"),
		viper.GetString("database.database"),
	)
	if err != nil {
		log.WithError(err).Error("error when creating database")
		return
	}
	log.Info("Successfully connected to the database")

	log.Info("Trying to listen to events")
	db.Listen()
	log.Info("Storage is now fully initialized")

	// Wait here until CTRL-C or other term signal is received.
	log.Info("Storage is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	db.Close()
}

// ReadConfig will try to find the config and read, if config file
// does not exists it will create one with default options
func ReadConfig() {
	configType = os.Getenv("KeiwiConfigType")
	if configType == "" {
		configType = "json"
	}
	viper.SetConfigType(configType)

	viper.SetConfigFile("config." + configType)
	viper.AddConfigPath(".")

	viper.SetDefault("log.dir", "./logs")
	viper.SetDefault("log.syntax", "%date%_storage.log")
	viper.SetDefault("log.level", "info")

	viper.SetDefault("database.username", "admin")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.ip", "127.0.0.1")
	viper.SetDefault("database.port", "27017")
	viper.SetDefault("database.database", "keiwi")

	viper.SetDefault("nats.url", nats.DefaultURL)

	if err := viper.ReadInConfig(); err != nil {
		log.Debug("Config file not found, saving default")
		if err = viper.WriteConfigAs("config." + configType); err != nil {
			log.WithField("error", err.Error()).Fatal("Can't save default config")
		}
	}

	level := strings.ToLower(viper.GetString("log.level"))
	log.Log = log.NewLogger(log.GetLevelFromString(level), []log.Reporter{
		cli.NewCli(),
		file.NewFile(viper.GetString("log.dir"), viper.GetString("log.syntax")),
	})
}
