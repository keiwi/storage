package storage

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/keiwi/storage/database"
	"github.com/keiwi/utils"
)

type Storage struct{}

func (s Storage) StartStorage() {
	utils.Log.Info("Trying to connect to database")
	db, err := database.NewDatabase("admin", "", "127.0.0.1", "27017", "keiwi")
	if err != nil {
		utils.Log.WithError(err).Error("error when creating database")
		return
	}
	utils.Log.Info("Successfully connected to the database")

	utils.Log.Info("Trying to listen to events")
	db.Listen()
	utils.Log.Info("Storage is now fully initialized")

	// Wait here until CTRL-C or other term signal is received.
	utils.Log.Info("Storage is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	db.Close()
}
