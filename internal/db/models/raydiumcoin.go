package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RaydiumCoin struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	MigratedAt  string             `bson:"migratedAt"`
	CoinAddress string             `bson:"coinAddress"`
	PoolId      string             `bson:"poolId"`
	Pool1       string             `bson:"pool1"`
	Pool2       string             `bson:"pool2"`
	Block       uint64             `bson:"block"`
	Signature   string             `bson:"signature"`
}
