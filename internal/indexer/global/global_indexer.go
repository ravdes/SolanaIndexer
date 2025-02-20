package global

import (
	"solanaindexer/internal/indexer/pumpfun"
	"solanaindexer/internal/indexer/raydium"
	"solanaindexer/internal/logger"
	"sync"
)

type GlobalIndexer struct {
	pumpfunIndexer *pumpfun.PumpfunIndexer
	raydiumIndexer *raydium.RaydiumIndexer
}

func (g *GlobalIndexer) StartIndexer() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := g.pumpfunIndexer.StartPumpfunIndexer(); err != nil {
			logger.Errorf("Error on PumpfunIndexer: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := g.raydiumIndexer.StartRaydiumIndexer(); err != nil {
			logger.Errorf("Error on RaydiumIndexer: %v", err)
		}
	}()

	wg.Wait()
}

func StartGlobalIndexer() {
	pfIndexer := pumpfun.NewPumpfunIndexer()
	rdIndexer := raydium.NewRaydiumIndexer()
	globalIndexer := GlobalIndexer{pumpfunIndexer: pfIndexer, raydiumIndexer: rdIndexer}
	globalIndexer.StartIndexer()
}
