package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/wobusbz/messageid/cache"
	"google.golang.org/grpc"
)

func main() {
	var (
		port       = flag.String("port", "8080", "gRPC server port")
		configFile = flag.String("config", "config/tables.json", "path to table configuration file")
	)
	flag.Parse()

	// Load configuration
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create cache instance
	cacheInstance := cache.NewCache(config)
	log.Printf("Cache server initialized with %d master tables and %d slave tables", 
		len(config.MasterTables), len(config.SlaveTables))

	// Create gRPC server
	grpcServer := grpc.NewServer()
	cacheServer := cache.NewServer(cacheInstance)
	cache.RegisterCacheServiceServer(grpcServer, cacheServer)

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Cache server starting on port %s", *port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func loadConfig(filename string) (*cache.TableConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config cache.TableConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return &config, nil
}