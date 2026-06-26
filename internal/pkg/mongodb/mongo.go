package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"godan/internal/config"
)

var Client *mongo.Client
var DB *mongo.Database

func Init(cfg *config.MongoDBConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(cfg.URI).SetMaxPoolSize(uint64(cfg.PoolSize)))
	if err != nil {
		return fmt.Errorf("mongo connect: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("mongo ping: %w", err)
	}

	Client = client
	DB = client.Database(cfg.Database)
	return nil
}

func Collection(name string) *mongo.Collection {
	return DB.Collection(name)
}

func Close() {
	if Client != nil {
		Client.Disconnect(context.Background())
	}
}
