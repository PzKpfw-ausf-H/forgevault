package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/PzKpfw-ausf-H/forgevault/internal/auth"
	"github.com/PzKpfw-ausf-H/forgevault/internal/httpapi"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo/postgres"
	"github.com/PzKpfw-ausf-H/forgevault/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	// TOKEN MANAGER CONFIG
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "forgevault"
	}

	ttlMin := 60
	if v := os.Getenv("JWT_TTL_MIN"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal("JWT_TTL_MIN must be a positive int")
		}
		ttlMin = n
	}

	tm := auth.NewTokenManager([]byte(secret), time.Duration(ttlMin)*time.Minute, issuer)

	//refresh token TTL
	refreshTTLDays := 7
	if v := os.Getenv("REFRESH_TTL_DAYS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("REFRESH_TTL_DAYS must be a positive int")
		}
		refreshTTLDays = n
	}

	//pgxpool
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("db url required")
	}
	pool, poolErr := pgxpool.New(context.Background(), dbURL)
	if poolErr != nil {
		log.Fatalf("error during creating pool: %v", poolErr)
	}
	defer pool.Close()

	// assets config
	assetsRepo := postgres.NewAssetsSQLRepo(pool)
	svc := service.NewAssetService(assetsRepo)
	h := httpapi.NewAssetsHandler(svc)

	//users config
	umemrepo := postgres.NewUsersSQLRepo(pool)
	refreshSessions := postgres.NewRefreshSessionSQLRepo(pool)
	usvc := service.NewUserService(umemrepo, tm, refreshSessions, time.Duration(refreshTTLDays)*24*time.Hour)
	uh := httpapi.NewUsersHandler(usvc)

	router := httpapi.NewRouter(h, tm, uh)

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
