package main

import (
	"github.com/keiwi/storage"
	"github.com/keiwi/utils/log"
	"github.com/keiwi/utils/log/handlers/cli"
	"github.com/keiwi/utils/log/handlers/file"
)

func main() {
	log.Log = log.NewLogger(log.DEBUG, []log.Reporter{
		cli.NewCli(),
		file.NewFile("./logs", "%date%_storage.log"),
	})

	storage := storage.Storage{}

	storage.StartStorage()

}
