package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PzKpfw-ausf-H/forgevault/internal/httpapi"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo/memory"
	"github.com/PzKpfw-ausf-H/forgevault/internal/service"
)

type healthResp struct {
	Status string `json:"status"`
}

func main() {
	memrepo := memory.NewMemRepo()
	svc := service.NewAssetService(memrepo)
	h := httpapi.NewAssetsHandler(svc)

	router := httpapi.NewRouter(h)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("Starting server on 8080...")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server listen and serve: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("error shutting down: %v", err)
	}

	log.Println("Server shut down")
}
