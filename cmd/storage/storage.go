package main

import (
	"github.com/keiwi/storage"
	"github.com/keiwi/utils/log"
	"github.com/keiwi/utils/log/handlers/cli"
	"github.com/keiwi/utils/log/handlers/file"
)

func main() {
	fileConfig := file.Config{Folder: "./logs", Filename: "%date%_storage.log"}
	log.Log = log.NewLogger(log.DEBUG, []log.Reporter{
		cli.NewCli(),
		file.NewFile(&fileConfig),
	})

	storage := storage.Storage{}
	storage.StartStorage()

}
