package raydium

import (
	"context"
	"github.com/google/uuid"
	"github.com/mr-tron/base58"
	"io"
	"solanaindexer/internal/constants"
	"solanaindexer/internal/db/models"
	"solanaindexer/internal/db/repository"
	"solanaindexer/internal/geyser"
	"solanaindexer/internal/geyser/proto"
	"solanaindexer/internal/logger"
	"solanaindexer/internal/utils"
	"time"
)

type RaydiumIndexer struct {
	envVariables *utils.Config
	repository.CoinRepository
}

func (r *RaydiumIndexer) StartRaydiumIndexer() error {
	conn, err := geyser.StartGeyser(r.envVariables.Grpc)
	if err != nil {
		logger.Errorf("Error %v", err)
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
			AccountInclude:  []string{constants.PumpFunMigration.String()},
			AccountRequired: []string{constants.RaydiumV4Contract.String()},
		}}, Commitment: &commitmentLevel}

	ctx := context.Background()

	stream, err := client.Subscribe(ctx)
	if err != nil {
		logger.Errorf("Error on subscribing to grpc %v", err)
		return err
	}
	err = stream.Send(&subscription)
	if err != nil {
		logger.Errorf("Error on grpc streaming %v", err)
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			logger.Errorf("Error on grpc listening %v", err)
			return err
		}

		if resp != nil && resp.GetTransaction() != nil {
			accountKeys := resp.GetTransaction().Transaction.Transaction.Message.AccountKeys
			logMessages := resp.GetTransaction().GetTransaction().GetMeta().GetLogMessages()
			poolMigrated := utils.ContainsSubstring(logMessages, "initialize2: InitializeInstruction2")

			if poolMigrated {
				coinAddress := base58.Encode(accountKeys[18])
				poolId := base58.Encode(accountKeys[2])
				pool1 := base58.Encode(accountKeys[5])
				pool2 := base58.Encode(accountKeys[6])
				block := resp.GetTransaction().Slot
				signature := base58.Encode(resp.GetTransaction().GetTransaction().Signature)

				coinData := models.RaydiumCoin{
					MigratedAt:  time.Now().Format("Jan 2 15:04:05.00"),
					CoinAddress: coinAddress,
					PoolId:      poolId,
					Pool1:       pool1,
					Pool2:       pool2,
					Block:       block,
					Signature:   signature,
				}

				err := r.InsertCoin(ctx, coinData)
				if err != nil {
					logger.Errorf("Error on inserting to database %v", err)
					return err
				}
				logger.Infof("Added %v from raydiumMigration to db", coinAddress)
			}
		}
	}
}

func NewRaydiumIndexer() *RaydiumIndexer {
	envVariables := utils.LoadEnvVariables()
	coinRepository := repository.NewIndexerRepository("raydiumIndexer")
	return &RaydiumIndexer{envVariables: envVariables, CoinRepository: coinRepository}
}
