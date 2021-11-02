// +build server-api

package main

import (
	"github.com/GlassNode/be_challenge/connectors"
	"github.com/GlassNode/be_challenge/handlers"
	"github.com/GlassNode/be_challenge/services"
	_ "github.com/lib/pq"
	"net/http"
	"fmt"
	"github.com/GlassNode/be_challenge/data_providers"

)

func main() {
	dbConnection, err := connectors.NewDBConnector()
	if err != nil {
		fmt.Println(err)
		return
	}
	dataProvider := data_providers.NewDataProvider(dbConnection)
	apiService := services.NewAPIService(dataProvider)
	handler := handlers.NewHandler(apiService)

	http.HandleFunc("/lastHour", handler.HandleLastHour)
	http.HandleFunc("/range", handler.HandleRange)
	http.ListenAndServe(":8080", nil)
}
