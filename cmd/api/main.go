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
	"github.com/PzKpfw-ausf-H/forgevault/internal/storage/minio"
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
		if err != nil || n <= 0 {
			log.Fatal("JWT_TTL_MIN must be a positive int")
		}
		ttlMin = n
	}

	tm := auth.NewTokenManager([]byte(secret), time.Duration(ttlMin)*time.Minute, issuer)

	// S3 config env
	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		log.Fatal("S3_BUCKET is required")
	}

	s3TtlMin := 15
	if v := os.Getenv("S3_PRESIGN_TTL_MIN"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal("S3_PRESIGN_TTL_MIN must be a positive int")
		}
		s3TtlMin = n
	}

	//refresh token TTL
	refreshTTLDays := 7
	if v := os.Getenv("REFRESH_TTL_DAYS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
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

	//S3 config
	S3, err := minio.NewFromEnv()
	if err != nil {
		log.Fatal("error creating S3 storage from env")
	}
	if err := S3.EnsureBucket(context.Background(), bucket); err != nil {
		log.Fatalf("ensure bucket: %v", err)
	}

	// files config
	filesRepo := postgres.NewAssetFilesSQLRepo(pool)
	fileSvc := service.NewFileService(filesRepo, bucket, S3, time.Duration(s3TtlMin)*time.Minute, assetsRepo)
	fh := httpapi.NewFilesHandler(fileSvc)

	//users config
	usersRepo := postgres.NewUsersSQLRepo(pool)
	refreshSessions := postgres.NewRefreshSessionSQLRepo(pool)
	usvc := service.NewUserService(usersRepo, tm, refreshSessions, time.Duration(refreshTTLDays)*24*time.Hour)
	uh := httpapi.NewUsersHandler(usvc)

	router := httpapi.NewRouter(h, tm, uh, fh)

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
