package main

import (
	"solanaindexer/internal/config"
	"solanaindexer/internal/indexer/global"
	"solanaindexer/internal/logger"
)

func main() {
	logger.Init()
	err := config.ConnectDB()
	if err != nil {
		logger.Errorf("Error while connecting to database: %v", err)
	}
	ch := make(chan int)
	global.StartGlobalIndexer()
	<-ch
}
