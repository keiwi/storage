package main

import (
	"flag"
	"os"
	"text/template"
)

var dbTmpl *template.Template
var modelTmpl *template.Template

var databaseTemplate = flag.String("database-template", "gen/db.tpl", "path for db file")
var modelTemplate = flag.String("model-template", "gen/model.tpl", "path for model file")
var templateOutput = flag.String("template-out", "database/", "path for output of generated code")

type Generate struct {
	ID   string
	Name string
}

var generate = []Generate{
	{ID: "alerts", Name: "Alert"},
	{ID: "alert_options", Name: "AlertOption"},
	{ID: "checks", Name: "Check"},
	{ID: "clients", Name: "Client"},
	{ID: "commands", Name: "Command"},
	{ID: "groups", Name: "Group"},
	{ID: "servers", Name: "Server"},
	{ID: "uploads", Name: "Upload"},
	{ID: "users", Name: "User"},
}

func init() {
	flag.Parse()
	dbTmpl = template.Must(template.ParseFiles(*databaseTemplate))
	modelTmpl = template.Must(template.ParseFiles(*modelTemplate))
}

func main() {
	f, err := os.OpenFile(*templateOutput+"db.go", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}

	err = dbTmpl.Execute(f, generate)
	if err != nil {
		f.Close()
		panic(err)
	}
	f.Close()

	for _, gen := range generate {
		genFile, err := os.OpenFile(*templateOutput+gen.ID+".go", os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}

		err = modelTmpl.Execute(genFile, gen)
		if err != nil {
			genFile.Close()
			panic(err)
		}
		genFile.Close()
	}
}
