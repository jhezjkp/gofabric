package main

import (
	"github.com/xssdoctor/gofabric/cli"
	"github.com/xssdoctor/gofabric/db"
	"github.com/xssdoctor/gofabric/utils"
)

func main() {
	err := db.InitDB() // initialize the database, including creating tables and populating the database. If the database already exists, it will not be overwritten, but the tables will be created if they do not exist.
	if err != nil {
		utils.LogError(err)
	}
	cli.Cli()

}