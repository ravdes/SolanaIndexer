package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PumpfunCoin struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt              string             `bson:"createdAt"`
	CoinAddress            string             `bson:"coinAddress"`
	Creator                string             `bson:"creator"`
	BondingCurve           string             `bson:"bondingCurve"`
	AssociatedBondingCurve string             `bson:"associatedBondingCurve"`
	Block                  uint64             `bson:"block"`
	Signature              string             `bson:"signature"`
}
