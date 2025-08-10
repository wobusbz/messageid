package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/wobusbz/messageid/cache"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var (
		serverAddr = flag.String("server", "localhost:8080", "cache server address")
	)
	flag.Parse()

	// Connect to the cache server
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := cache.NewCacheServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Example 1: Set player data
	playerKey := "/players/12345"
	playerData := []byte(`{"name": "TestPlayer", "level": 50, "gold": 1000}`)

	log.Printf("Setting player data for key: %s", playerKey)
	setResp, err := client.Set(ctx, &cache.SetRequest{
		Key:  playerKey,
		Data: playerData,
	})
	if err != nil {
		log.Fatalf("Failed to set data: %v", err)
	}
	if !setResp.Success {
		log.Printf("Set failed: %s", setResp.Error)
	} else {
		log.Printf("Successfully set player data")
	}

	// Example 2: Get player data
	log.Printf("Getting player data for key: %s", playerKey)
	getResp, err := client.Get(ctx, &cache.GetRequest{
		Key: playerKey,
	})
	if err != nil {
		log.Fatalf("Failed to get data: %v", err)
	}
	if !getResp.Found {
		log.Printf("Data not found: %s", getResp.Error)
	} else {
		log.Printf("Retrieved player data: %s", string(getResp.Data))
	}

	// Example 3: Set camp wall data (slave table)
	campWallKey := "/players/12345/campwall/1"
	campWallData := []byte(`{"wall_type": "stone", "durability": 100}`)

	log.Printf("Setting camp wall data for key: %s", campWallKey)
	setResp, err = client.Set(ctx, &cache.SetRequest{
		Key:  campWallKey,
		Data: campWallData,
	})
	if err != nil {
		log.Fatalf("Failed to set camp wall data: %v", err)
	}
	if !setResp.Success {
		log.Printf("Set failed: %s", setResp.Error)
	} else {
		log.Printf("Successfully set camp wall data")
	}

	// Example 4: Try to set data with invalid key
	invalidKey := "/invalid/key/path"
	log.Printf("Trying to set data with invalid key: %s", invalidKey)
	setResp, err = client.Set(ctx, &cache.SetRequest{
		Key:  invalidKey,
		Data: []byte("test"),
	})
	if err != nil {
		log.Fatalf("Failed to call set: %v", err)
	}
	if !setResp.Success {
		log.Printf("Expected failure - Set failed: %s", setResp.Error)
	}

	// Example 5: Delete player data
	log.Printf("Deleting player data for key: %s", playerKey)
	deleteResp, err := client.Delete(ctx, &cache.DeleteRequest{
		Key: playerKey,
	})
	if err != nil {
		log.Fatalf("Failed to delete data: %v", err)
	}
	if !deleteResp.Success {
		log.Printf("Delete failed: %s", deleteResp.Error)
	} else {
		log.Printf("Successfully deleted player data")
	}

	// Example 6: Verify deletion
	log.Printf("Verifying deletion by getting player data for key: %s", playerKey)
	getResp, err = client.Get(ctx, &cache.GetRequest{
		Key: playerKey,
	})
	if err != nil {
		log.Fatalf("Failed to get data: %v", err)
	}
	if !getResp.Found {
		log.Printf("Confirmed: player data not found after deletion")
	} else {
		log.Printf("Unexpected: still found data: %s", string(getResp.Data))
	}

	log.Printf("Cache client example completed")
}