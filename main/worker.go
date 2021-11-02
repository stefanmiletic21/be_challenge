// +build worker

package main

import (
	"fmt"
	"github.com/GlassNode/be_challenge/connectors"
	"github.com/GlassNode/be_challenge/data_providers"
	"github.com/GlassNode/be_challenge/services"
	_ "github.com/lib/pq"
)

func main() {
	dbConnection, err := connectors.NewDBConnector()
	if err != nil {
		fmt.Println(err)
		return
	}
	dataProvider := data_providers.NewDataProvider(dbConnection)
	err = dataProvider.MakeAggregatedTable()
	if err != nil {
		fmt.Println(err)
	}

	migrationService := services.NewMigrationService(dataProvider)
	migrationService.Start()

}
