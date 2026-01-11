package connection

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var (
	clientInstance *mongo.Client
	mongoOnce      sync.Once // Ensures thread-safe initialization
)

// GetDatabase returns a pointer to the specific database instance.
func GetDatabase() *mongo.Database {
	if clientInstance == nil {
		mongoOnce.Do(func() {
			// If ConnectMongo hasn't been called yet, initialize it
			client, _ := ConnectMongo()
			clientInstance = client
		})
	}

	dbName := os.Getenv("MY_MONGODB_NAME")
	if dbName == "" {
		dbName = "goplaygrounds_default" // Default fallback
	}
	return clientInstance.Database(dbName)
}

// GetClient provides access to the raw client for Disconnect operations
func GetClient() *mongo.Client {
	return clientInstance
}

func ConnectMongo() (*mongo.Client, context.CancelFunc) {
	uri := os.Getenv("MY_MONGODB_URL")
	if uri == "" {
		log.Fatal("MY_MONGODB_URL not set in .env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// v2 Driver: Connect starts background monitoring immediately
	client, err := mongo.Connect(options.Client().ApplyURI(uri).
		SetBSONOptions(&options.BSONOptions{
			UseJSONStructTags: true,
			NilSliceAsEmpty:   true,
		}))
	if err != nil {
		cancel()
		log.Fatal("Failed to create MongoDB client:", err)
	}

	// Verify server is reachable
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		cancel()
		_ = client.Disconnect(context.Background())
		log.Fatal("Could not connect to MongoDB (Ping failed):", err)
	}

	clientInstance = client // Global assignment
	//log.Println("Successfully connected to MongoDB")
	return client, cancel
}
