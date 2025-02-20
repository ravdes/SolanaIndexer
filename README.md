# Solana Indexer

## Overview

The **Solana Indexer** is a Go-based application designed to track the creation of new coins on **Pumpfun** and the migration of assets from **Pumpfun** to **Raydium**. By utilizing the [Yellowstone gRPC](https://github.com/rpcpool/yellowstone-grpc/tree/master) service application monitors transactions in real time and stores relevant data in MongoDB for efficient querying and analysis.

The entire application is dockerized, enabling seamless setup and deployment for developers and teams.

## Requirements

- **Docker** installed on your machine
- Access to a **gRPC Solana endpoint** for transaction monitoring

## Setup

To get started with the Solana Indexer, follow the configuration steps outlined below:

### Configuration

Before running the application, configure your environment by editing the `.env` file with the following settings:

```dotenv
MONGO_DB_NAME=solanaIndexer
MONGO_USERNAME=example1
MONGO_PASSWORD=example2
GRPC=grpcEndpoint
```
- MONGO_DB_NAME: The name of the MongoDB database to store indexed data.
- MONGO_USERNAME: The MongoDB username for database access.
- MONGO_PASSWORD: The MongoDB password for secure access.
- GRPC: The Solana gRPC endpoint for monitoring transactions in real time.

### Running the Application
Once the .env file is properly configured, you can start the Solana Indexer application using Docker. Simply run the following command to build and launch the container:

```env
docker-compose up
```
This will build the Docker image, pull any necessary dependencies, and start the application.

Upon startup, the application will simultaneously monitor both of the things:
- New coin creations on Pumpfun
- Migrations from Pumpfun to Raydium

### Database structure

**Pumpfun created coin**  structure in database
```go
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
```
**Raydium migrated coin** structure in database
```go
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
```


### Features

- Monitors new coin creation on Pumpfun.
- Tracks migration events from Pumpfun to Raydium.
- Stores transaction data securely in MongoDB.
- Fully Dockerized for easy setup and deployment.
- Real-time monitoring using the Yellowstone gRPC service.

## Demo



https://github.com/user-attachments/assets/12c80dc9-bca9-4449-82eb-2599198d239a



## Usage

Feel free to fork this repository and extend it to suit your needs.  
