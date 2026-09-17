package minio

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestStorageMinIO_ObjectLifecycle(t *testing.T) {
	client := openTestStorage(t)
	idKey := uuid.NewString()
	key := "tests/" + idKey + "/test.txt"
	ctx := context.Background()
	ttl := time.Minute
	content := "hello forgevault"

	presigned, err := client.PresignUpload(ctx, key, ttl)
	if err != nil {
		t.Fatalf("presign upload: %v", err)
	}

	t.Cleanup(func() {
		err := client.Delete(context.Background(), key)
		if err != nil {
			t.Errorf("cleanup object: %v", err)
		}
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, presigned, strings.NewReader(content))
	if err != nil {
		t.Fatalf("upload request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("default client do: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d but got %d", http.StatusOK, resp.StatusCode)
	}

	exists, err := client.Exists(ctx, key)
	if err != nil {
		t.Fatalf("exists: %v", err)
	}

	if !exists {
		t.Fatalf("expected exists to be true but got false")
	}

	download, err := client.PresignDownload(ctx, key, ttl)
	if err != nil {
		t.Fatalf("presign download: %v", err)
	}

	downloadReq, err := http.NewRequestWithContext(ctx, http.MethodGet, download, nil)
	if err != nil {
		t.Fatalf("download request: %v", err)
	}

	downloadResp, err := http.DefaultClient.Do(downloadReq)
	if err != nil {
		t.Fatalf("download response: %v", err)
	}
	defer downloadResp.Body.Close()

	if downloadResp.StatusCode != http.StatusOK {
		t.Fatalf("download status: expected status %d but got %d", http.StatusOK, downloadResp.StatusCode)
	}

	body, err := io.ReadAll(downloadResp.Body)
	if err != nil {
		t.Fatalf("io read all: %v", err)
	}

	if string(body) != content {
		t.Fatalf("expected content %q but got %q", content, string(body))
	}

	err = client.Delete(ctx, key)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}

	exists, err = client.Exists(ctx, key)
	if err != nil {
		t.Fatalf("exists after delete: %v", err)
	}

	if exists {
		t.Fatalf("exists after delete: expected to be false")
	}
}
