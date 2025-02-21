package global

import (
	"solanaindexer/internal/indexer/pumpfun"
	"solanaindexer/internal/indexer/raydium"
	"solanaindexer/internal/logger"
	"solanaindexer/internal/utils"
)

type GlobalIndexer struct {
	pumpfunIndexer *pumpfun.PumpfunIndexer
	raydiumIndexer *raydium.RaydiumIndexer
}

func (g *GlobalIndexer) StartIndexer() {

	go func() {

		if err := g.pumpfunIndexer.StartPumpfunIndexer(); err != nil {
			logger.Errorf("Error on PumpfunIndexer: %v", err)
		}
	}()

	go func() {

		if err := g.raydiumIndexer.StartRaydiumIndexer(); err != nil {
			logger.Errorf("Error on RaydiumIndexer: %v", err)
		}
	}()
}

func StartGlobalIndexer(config *utils.Config) {
	pfIndexer := pumpfun.NewPumpfunIndexer(config)
	rdIndexer := raydium.NewRaydiumIndexer(config)
	globalIndexer := GlobalIndexer{pumpfunIndexer: pfIndexer, raydiumIndexer: rdIndexer}
	globalIndexer.StartIndexer()
}
