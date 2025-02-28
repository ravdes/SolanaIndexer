package pumpfun

import (
	"context"
	"github.com/google/uuid"
	"github.com/mr-tron/base58"
	"io"
	"slices"
	"solanaindexer/internal/constants"
	"solanaindexer/internal/db/models"
	"solanaindexer/internal/db/repository"
	"solanaindexer/internal/geyser"
	"solanaindexer/internal/geyser/proto"
	"solanaindexer/internal/logger"
	"solanaindexer/internal/utils"
	"time"
)

type PumpfunIndexer struct {
	*utils.Config
	repository.CoinRepository
}

func (p *PumpfunIndexer) StartPumpfunIndexer() error {
	conn, err := geyser.StartGeyser(p.Grpc)
	if err != nil {
		logger.Errorf("Error starting Geyser: %v", err)
		return err
	}
	client := proto.NewGeyserClient(conn)

	var subscription proto.SubscribeRequest
	f := false
	commitmentLevel := proto.CommitmentLevel_PROCESSED
	name := uuid.NewString()
	name = name[:len(name)-4]
	subscription = proto.SubscribeRequest{Transactions: map[string]*proto.SubscribeRequestFilterTransactions{
		name: {
			Vote:            &f,
			Failed:          &f,
			Signature:       nil,
			AccountInclude:  []string{constants.PumpFunContract},
			AccountExclude:  nil,
			AccountRequired: nil,
		}}, Commitment: &commitmentLevel}

	ctx := context.Background()

	stream, err := client.Subscribe(ctx)
	if err != nil {
		logger.Errorf("Error subscribing to gRPC: %v", err)
		return err
	}
	err = stream.Send(&subscription)
	if err != nil {
		logger.Errorf("Error sending subscription: %v", err)
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			logger.Errorf("Error receiving stream: %v", err)
			return err
		}

		if resp != nil && resp.GetTransaction() != nil {
			parsedTx := resp.GetTransaction()
			accountKeys := parsedTx.Transaction.Transaction.Message.AccountKeys
			logMessages := parsedTx.GetTransaction().GetMeta().GetLogMessages()
			newCoinCreated := slices.Contains(logMessages, "Program log: Instruction: InitializeMint2")

			if newCoinCreated {
				creator := base58.Encode(accountKeys[0])
				coinAddress := base58.Encode(accountKeys[1])
				bondingCurve := base58.Encode(accountKeys[2])
				associatedBondingCurve := base58.Encode(accountKeys[3])
				block := parsedTx.Slot
				signature := base58.Encode(parsedTx.GetTransaction().Signature)

				coinData := models.PumpfunCoin{
					CreatedAt:              time.Now().Format("Jan 2 15:04:05.00"),
					CoinAddress:            coinAddress,
					Creator:                creator,
					BondingCurve:           bondingCurve,
					AssociatedBondingCurve: associatedBondingCurve,
					Block:                  block,
					Signature:              signature,
				}

				err = p.InsertCoin(ctx, coinData)
				if err != nil {
					logger.Errorf("Error inserting coin into database: %v", err)
					return err
				}
				logger.Infof("Added %v from pumpfun creation to db", coinAddress)
			}
		}
	}
}

func NewPumpfunIndexer(config *utils.Config) *PumpfunIndexer {
	coinRepository := repository.NewIndexerRepository("pumpfunIndexer")
	return &PumpfunIndexer{Config: config, CoinRepository: coinRepository}
}
