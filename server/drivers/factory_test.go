package drivers

import (
	"testing"

	"github.com/gmllt/clariti/server/config"
)

func TestNewStorage_RAM(t *testing.T) {
	cfg := &config.Config{
		Storage: config.StorageConfig{
			Driver: "ram",
		},
	}

	storage, err := NewStorage(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, ok := storage.(*RAMStorage); !ok {
		t.Fatalf("expected RAMStorage, got %T", storage)
	}
}

func TestNewStorage_S3(t *testing.T) {
	cfg := &config.Config{
		Storage: config.StorageConfig{
			Driver: "s3",
			S3: config.S3Config{
				Region: "us-east-1",
				Bucket: "test-bucket",
			},
		},
	}

	// This test will fail without proper AWS credentials/bucket
	// but we can at least test the factory logic
	_, err := NewStorage(cfg)

	// We expect an error because we don't have valid AWS credentials
	// or the bucket doesn't exist, but we should not get an "unsupported driver" error
	if err != nil && err.Error() == "unsupported storage driver: s3" {
		t.Fatalf("factory should support s3 driver")
	}
}

func TestNewStorage_UnsupportedDriver(t *testing.T) {
	cfg := &config.Config{
		Storage: config.StorageConfig{
			Driver: "unknown",
		},
	}

	_, err := NewStorage(cfg)
	if err == nil {
		t.Fatalf("expected error for unsupported driver")
	}

	expected := "unsupported storage driver: unknown"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}
