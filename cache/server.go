package cache

import (
	"context"
	"log"
)

// Server implements the CacheService gRPC interface
type Server struct {
	UnimplementedCacheServiceServer
	cache *Cache
}

// NewServer creates a new gRPC server with the given cache instance
func NewServer(cache *Cache) *Server {
	return &Server{
		cache: cache,
	}
}

// Get implements the Get RPC method
func (s *Server) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	if req.Key == "" {
		return &GetResponse{
			Found: false,
			Error: "key cannot be empty",
		}, nil
	}

	data, found, err := s.cache.Get(req.Key)
	if err != nil {
		log.Printf("Error getting key %s: %v", req.Key, err)
		return &GetResponse{
			Found: false,
			Error: err.Error(),
		}, nil
	}

	return &GetResponse{
		Found: found,
		Data:  data,
	}, nil
}

// Set implements the Set RPC method
func (s *Server) Set(ctx context.Context, req *SetRequest) (*SetResponse, error) {
	if req.Key == "" {
		return &SetResponse{
			Success: false,
			Error:   "key cannot be empty",
		}, nil
	}

	err := s.cache.Set(req.Key, req.Data)
	if err != nil {
		log.Printf("Error setting key %s: %v", req.Key, err)
		return &SetResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &SetResponse{
		Success: true,
	}, nil
}

// Delete implements the Delete RPC method
func (s *Server) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	if req.Key == "" {
		return &DeleteResponse{
			Success: false,
			Error:   "key cannot be empty",
		}, nil
	}

	err := s.cache.Delete(req.Key)
	if err != nil {
		log.Printf("Error deleting key %s: %v", req.Key, err)
		return &DeleteResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &DeleteResponse{
		Success: true,
	}, nil
}