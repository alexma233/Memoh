package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/memohai/memoh/internal/chat"
	"github.com/memohai/memoh/internal/config"
	ctr "github.com/memohai/memoh/internal/containerd"
	"github.com/memohai/memoh/internal/db"
	dbsqlc "github.com/memohai/memoh/internal/db/sqlc"
	"github.com/memohai/memoh/internal/embeddings"
	"github.com/memohai/memoh/internal/handlers"
	"github.com/memohai/memoh/internal/mcp"
	"github.com/memohai/memoh/internal/memory"
	"github.com/memohai/memoh/internal/models"
	"github.com/memohai/memoh/internal/providers"
	"github.com/memohai/memoh/internal/server"
)

func main() {
	ctx := context.Background()
	cfgPath := os.Getenv("CONFIG_PATH")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if strings.TrimSpace(cfg.Auth.JWTSecret) == "" {
		log.Fatalf("jwt secret is required")
	}
	jwtExpiresIn, err := time.ParseDuration(cfg.Auth.JWTExpiresIn)
	if err != nil {
		log.Fatalf("invalid jwt expires in: %v", err)
	}

	addr := cfg.Server.Addr
	if value := os.Getenv("HTTP_ADDR"); value != "" {
		addr = value
	}

	socketPath := cfg.Containerd.SocketPath
	if value := os.Getenv("CONTAINERD_SOCKET"); value != "" {
		socketPath = value
	}
	factory := ctr.DefaultClientFactory{SocketPath: socketPath}
	client, err := factory.New(ctx)
	if err != nil {
		log.Fatalf("connect containerd: %v", err)
	}
	defer client.Close()

	service := ctr.NewDefaultService(client, cfg.Containerd.Namespace)
	manager := mcp.NewManager(service, cfg.MCP)

	pingHandler := handlers.NewPingHandler()
	containerdHandler := handlers.NewContainerdHandler(service, cfg.MCP, cfg.Containerd.Namespace)

	conn, err := db.Open(ctx, cfg.Postgres)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer conn.Close()
	manager.WithDB(conn)
	queries := dbsqlc.New(conn)
	modelsService := models.NewService(queries)

	authHandler := handlers.NewAuthHandler(conn, cfg.Auth.JWTSecret, jwtExpiresIn)

	// Initialize chat resolver for both chat and memory operations
	chatResolver := chat.NewResolver(modelsService, queries, 30*time.Second)

	// Create LLM client for memory operations using chat provider
	var llmClient memory.LLM
	memoryModel, memoryProvider, err := models.SelectMemoryModel(ctx, modelsService, queries)
	if err != nil {
		log.Fatalf("select memory model: %v\nPlease configure at least one chat model in the database.", err)
	}

	log.Printf("Using memory model: %s (provider: %s)", memoryModel.ModelID, memoryProvider.ClientType)
	provider, err := chat.CreateProvider(memoryProvider, 30*time.Second)
	if err != nil {
		log.Fatalf("create memory provider: %v", err)
	}
	llmClient = memory.NewProviderLLMClient(provider, memoryModel.ModelID)

	resolver := embeddings.NewResolver(modelsService, queries, 10*time.Second)
	vectors, textModel, multimodalModel, hasModels, err := embeddings.CollectEmbeddingVectors(ctx, modelsService)
	if err != nil {
		log.Fatalf("embedding models: %v", err)
	}

	var memoryService *memory.Service
	var memoryHandler *handlers.MemoryHandler

	if !hasModels {
		log.Println("WARNING: No embedding models configured. Memory service will not be available.")
		log.Println("You can add embedding models via the /models API endpoint.")
		memoryHandler = handlers.NewMemoryHandler(nil)
	} else {
		if textModel.ModelID == "" {
			log.Println("WARNING: No text embedding model configured. Text embedding features will be limited.")
		}
		if multimodalModel.ModelID == "" {
			log.Println("WARNING: No multimodal embedding model configured. Multimodal embedding features will be limited.")
		}

		var textEmbedder embeddings.Embedder
		var store *memory.QdrantStore

		if textModel.ModelID != "" && textModel.Dimensions > 0 {
			textEmbedder = &embeddings.ResolverTextEmbedder{
				Resolver: resolver,
				ModelID:  textModel.ModelID,
				Dims:     textModel.Dimensions,
			}

			if len(vectors) > 0 {
				store, err = memory.NewQdrantStoreWithVectors(
					cfg.Qdrant.BaseURL,
					cfg.Qdrant.APIKey,
					cfg.Qdrant.Collection,
					vectors,
					time.Duration(cfg.Qdrant.TimeoutSeconds)*time.Second,
				)
				if err != nil {
					log.Fatalf("qdrant named vectors init: %v", err)
				}
			} else {
				store, err = memory.NewQdrantStore(
					cfg.Qdrant.BaseURL,
					cfg.Qdrant.APIKey,
					cfg.Qdrant.Collection,
					textModel.Dimensions,
					time.Duration(cfg.Qdrant.TimeoutSeconds)*time.Second,
				)
				if err != nil {
					log.Fatalf("qdrant init: %v", err)
				}
			}
		}

		memoryService = memory.NewService(llmClient, textEmbedder, store, resolver, textModel.ModelID, multimodalModel.ModelID)
		memoryHandler = handlers.NewMemoryHandler(memoryService)
	}
	embeddingsHandler := handlers.NewEmbeddingsHandler(modelsService, queries)
	swaggerHandler := handlers.NewSwaggerHandler()
	chatHandler := handlers.NewChatHandler(chatResolver)

	// Initialize providers and models handlers
	providersService := providers.NewService(queries)
	providersHandler := handlers.NewProvidersHandler(providersService)
	modelsHandler := handlers.NewModelsHandler(modelsService)
	srv := server.NewServer(addr, cfg.Auth.JWTSecret, pingHandler, authHandler, memoryHandler, embeddingsHandler, swaggerHandler, chatHandler, providersHandler, modelsHandler, containerdHandler)

	if err := srv.Start(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
