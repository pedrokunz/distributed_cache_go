package main

import (
	httpHandler "github.com/pedrokunz/distributed_cache_go/api/http"
	"github.com/pedrokunz/distributed_cache_go/internal/cache_manager"
	"log"
	"net/http"
)

func main() {
	cacheManager := cache_manager.New()
	handler := httpHandler.New(cacheManager)

	log.Println("Cache Manager is running on port 8080")

	log.Fatalln(http.ListenAndServe(":8080", handler))
}
