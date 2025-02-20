package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"solanaindexer/internal/config"
	"solanaindexer/internal/db/models"
	"solanaindexer/internal/logger"
)

type CoinRepository interface {
	InsertCoin(ctx context.Context, coin interface{}) error
	GetCoin(ctx context.Context, coinAddress string) (interface{}, error)
	UpdateCoin(ctx context.Context, coinAddress string, coin interface{}) error
	DeleteCoin(ctx context.Context, coinAddress string) error
}
type IndexerRepository struct {
	collection *mongo.Collection
}

func (r *IndexerRepository) InsertCoin(ctx context.Context, model interface{}) error {
	_, err := r.collection.InsertOne(ctx, model)
	return err
}

func (r *IndexerRepository) GetCoin(ctx context.Context, coinAddress string) (interface{}, error) {
	var coin models.PumpfunCoin
	err := r.collection.FindOne(ctx, bson.M{"coinAddress": coinAddress}).Decode(&coin)
	if err != nil {
		return nil, err
	}
	return &coin, nil
}

func (r *IndexerRepository) UpdateCoin(ctx context.Context, coinAddress string, coin interface{}) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"coinAddress": coinAddress}, coin)
	return err
}

func (r *IndexerRepository) DeleteCoin(ctx context.Context, coinAddress string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"coinAddress": coinAddress})
	return err
}

func NewIndexerRepository(collectionName string) CoinRepository {
	collection, err := config.GetCollection(collectionName)
	if err != nil {
		logger.Errorf("Error while getting collection %v", err)
	}
	return &IndexerRepository{collection: collection}
}
