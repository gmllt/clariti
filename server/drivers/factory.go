package drivers

import (
	"fmt"

	"github.com/gmllt/clariti/server/config"
)

// NewStorage creates a new storage driver based on the configuration
func NewStorage(cfg *config.Config) (EventStorage, error) {
	switch cfg.GetStorageDriver() {
	case "ram":
		return NewRAMStorage(), nil
	case "s3":
		s3Config := cfg.GetS3Config()
		s3StorageConfig := S3Config{
			Region:          s3Config.Region,
			Bucket:          s3Config.Bucket,
			AccessKeyID:     s3Config.AccessKeyID,
			SecretAccessKey: s3Config.SecretAccessKey,
			Endpoint:        s3Config.Endpoint,
			Prefix:          s3Config.Prefix,
		}
		return NewS3Storage(s3StorageConfig)
	default:
		return nil, fmt.Errorf("unsupported storage driver: %s", cfg.GetStorageDriver())
	}
}
