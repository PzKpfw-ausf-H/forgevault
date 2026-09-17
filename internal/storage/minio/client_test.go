package minio

import (
	"testing"

	"github.com/PzKpfw-ausf-H/forgevault/internal/config"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		wantErr  bool
	}{
		{
			name:     "http endpoint",
			endpoint: "http://localhost:9000",
			wantErr:  false,
		},
		{
			name:     "https endpoint",
			endpoint: "https://storage.test.com",
			wantErr:  false,
		},
		{
			name:     "invalid endpoint",
			endpoint: "ftp://localhost:9000",
			wantErr:  true,
		},
		{
			name:     "empty host",
			endpoint: "http://",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			cfg := config.S3Config{
				Endpoint:  tt.endpoint,
				AccessKey: "test-access",
				SecretKey: "test-secret",
			}

			client, err := NewClient(cfg)
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tt.wantErr && client == nil {
				t.Fatalf("expected client to be non-nil")
			}

			if tt.wantErr && err == nil {
				t.Fatalf("expected error")
			}
		})

	}
}
