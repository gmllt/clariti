package config

import (
	"testing"
)

func TestConfig_StorageDriverValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
	}{
		{
			name: "RAM driver valid",
			config: Config{
				Storage: StorageConfig{
					Driver: "ram",
				},
			},
			expectError: false,
		},
		{
			name: "S3 driver valid",
			config: Config{
				Storage: StorageConfig{
					Driver: "s3",
					S3: S3Config{
						Region: "us-east-1",
						Bucket: "test-bucket",
					},
				},
			},
			expectError: false,
		},
		{
			name: "S3 driver missing region",
			config: Config{
				Storage: StorageConfig{
					Driver: "s3",
					S3: S3Config{
						Bucket: "test-bucket",
					},
				},
			},
			expectError: true,
		},
		{
			name: "S3 driver missing bucket",
			config: Config{
				Storage: StorageConfig{
					Driver: "s3",
					S3: S3Config{
						Region: "us-east-1",
					},
				},
			},
			expectError: true,
		},
		{
			name: "Unsupported driver",
			config: Config{
				Storage: StorageConfig{
					Driver: "mongodb",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validateStorageConfig()
			if tt.expectError && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestConfig_GetStorageDriver(t *testing.T) {
	config := &Config{
		Storage: StorageConfig{
			Driver: "s3",
		},
	}

	if got := config.GetStorageDriver(); got != "s3" {
		t.Errorf("GetStorageDriver() = %v, want %v", got, "s3")
	}
}

func TestConfig_GetS3Config(t *testing.T) {
	config := &Config{
		Storage: StorageConfig{
			S3: S3Config{
				Region: "us-west-2",
				Bucket: "my-bucket",
				Prefix: "clariti/",
			},
		},
	}

	s3Config := config.GetS3Config()
	if s3Config.Region != "us-west-2" {
		t.Errorf("Expected region us-west-2, got %s", s3Config.Region)
	}
	if s3Config.Bucket != "my-bucket" {
		t.Errorf("Expected bucket my-bucket, got %s", s3Config.Bucket)
	}
	if s3Config.Prefix != "clariti/" {
		t.Errorf("Expected prefix clariti/, got %s", s3Config.Prefix)
	}
}
