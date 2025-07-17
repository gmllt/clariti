# Storage Drivers Implementation Summary

## Overview
Added configurable storage drivers to support both in-memory (RAM) and persistent (S3) storage for event data.

## Files Added/Modified

### Configuration (`server/config/`)
- **config.go**: Added `StorageConfig` and `S3Config` structs
- **config_test.go**: Added validation tests for storage configuration

### Storage Drivers (`server/drivers/`)
- **s3.go**: Complete S3 storage driver implementation
- **factory.go**: Factory pattern for creating storage drivers
- **factory_test.go**: Tests for the storage factory

### Configuration Examples (`local/config/`)
- **config.ram.yaml**: RAM storage configuration example
- **config.s3.yaml**: S3 storage configuration example

## Storage Drivers

### RAM Driver (`ram`)
- **Type**: In-memory storage
- **Performance**: Very fast
- **Persistence**: Data lost on restart
- **Use Case**: Development, testing
- **Configuration**: No additional config required

### S3 Driver (`s3`)
- **Type**: AWS S3 or S3-compatible storage
- **Performance**: Network dependent
- **Persistence**: Data persists across restarts
- **Use Case**: Production environments
- **Features**:
  - AWS IAM role support
  - Access key authentication
  - S3-compatible endpoints (MinIO, etc.)
  - Configurable object key prefixes
  - Automatic bucket validation

## Configuration

```yaml
storage:
  driver: "s3"  # Options: "ram", "s3"
  s3:
    region: "us-east-1"
    bucket: "my-clariti-bucket"
    # Optional: credentials (prefer IAM roles)
    access_key_id: "AKIAIOSFODNN7EXAMPLE"
    secret_access_key: "secret"
    # Optional: S3-compatible endpoint
    endpoint: "http://localhost:9000"
    # Optional: object key prefix
    prefix: "clariti/"
```

## Usage

```go
// Load configuration
cfg, err := config.LoadConfig("config.yaml")
if err != nil {
    log.Fatal(err)
}

// Create storage driver
storage, err := drivers.NewStorage(cfg)
if err != nil {
    log.Fatal(err)
}

// Use storage driver (same interface for all drivers)
incident := &event.Incident{...}
err = storage.CreateIncident(incident)
```

## S3 Object Structure

```
bucket/
├── prefix/
│   ├── incidents/
│   │   ├── incident-guid-1.json
│   │   └── incident-guid-2.json
│   └── planned_maintenances/
│       ├── pm-guid-1.json
│       └── pm-guid-2.json
```

## Validation

- Configuration validation on startup
- S3 bucket connectivity check
- Required fields validation (region, bucket for S3)
- Unsupported driver error handling

## Dependencies Added

- `github.com/aws/aws-sdk-go-v2` - AWS SDK v2 core
- `github.com/aws/aws-sdk-go-v2/config` - AWS configuration
- `github.com/aws/aws-sdk-go-v2/credentials` - AWS credentials
- `github.com/aws/aws-sdk-go-v2/service/s3` - S3 service client

## Testing

- Factory tests for driver creation
- Configuration validation tests
- S3 connectivity testing (graceful failure without credentials)
- All existing tests continue to pass

## Security Considerations

- Prefer IAM roles over access keys in production
- Support for S3-compatible services with custom endpoints
- Object key prefixes for multi-tenancy
- Proper error handling for missing objects
- Connection timeout configuration
