package drivers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gmllt/clariti/logger"
	"github.com/gmllt/clariti/models/event"
)

// S3Storage implements EventStorage interface using AWS S3
type S3Storage struct {
	client  *s3.Client
	bucket  string
	prefix  string
	timeout time.Duration
}

// S3Config holds S3-specific configuration
type S3Config struct {
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string
	Prefix          string
}

// NewS3Storage creates a new S3 storage driver
func NewS3Storage(cfg S3Config) (*S3Storage, error) {
	log := logger.GetDefault().WithComponent("S3Storage")
	log.WithField("bucket", cfg.Bucket).WithField("region", cfg.Region).Info("Initializing S3 storage driver")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build AWS config
	var opts []func(*config.LoadOptions) error

	// Set region
	opts = append(opts, config.WithRegion(cfg.Region))

	// Set credentials if provided
	if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
		log.Debug("Using provided AWS credentials")
		opts = append(opts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		))
	} else {
		log.Debug("Using default AWS credential chain")
	}

	awsConfig, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		log.WithError(err).Error("Failed to load AWS config")
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client
	s3Client := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			log.WithField("endpoint", cfg.Endpoint).Debug("Using custom S3 endpoint")
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			// For MinIO and other S3-compatible services, use path-style addressing
			o.UsePathStyle = true
		}
	})

	// Test connection by listing objects (with limit 1)
	log.Debug("Testing S3 connection")
	_, err = s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:  aws.String(cfg.Bucket),
		MaxKeys: aws.Int32(1),
	})
	if err != nil {
		log.WithError(err).WithField("bucket", cfg.Bucket).Warn("Failed to test S3 connection, but continuing anyway")
		// Don't fail here - the bucket might exist but be empty, or there might be permission issues
		// We'll catch real issues when we try to use it
	}

	prefix := cfg.Prefix
	if prefix != "" && prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}

	log.WithField("prefix", prefix).Info("S3 storage driver initialized successfully")
	return &S3Storage{
		client:  s3Client,
		bucket:  cfg.Bucket,
		prefix:  prefix,
		timeout: 30 * time.Second,
	}, nil
}

// Helper methods for S3 operations

func (s *S3Storage) getKey(category, id string) string {
	return fmt.Sprintf("%s%s/%s.json", s.prefix, category, id)
}

// createBucketIfNotExists creates the bucket if it doesn't exist
func (s *S3Storage) createBucketIfNotExists() error {
	log := logger.GetDefault().WithComponent("S3Storage")
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	log.WithField("bucket", s.bucket).Info("Creating bucket")
	_, err := s.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		// If bucket already exists, that's fine
		if strings.Contains(err.Error(), "BucketAlreadyExists") || strings.Contains(err.Error(), "BucketAlreadyOwnedByYou") {
			log.WithField("bucket", s.bucket).Debug("Bucket already exists")
			return nil
		}
		log.WithError(err).WithField("bucket", s.bucket).Error("Failed to create bucket")
		return err
	}
	log.WithField("bucket", s.bucket).Info("Bucket created successfully")
	return nil
}

func (s *S3Storage) putObject(key string, data interface{}) error {
	log := logger.GetDefault().WithComponent("S3Storage")
	log.WithField("key", key).Debug("Putting object to S3")

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.WithError(err).WithField("key", key).Error("Failed to marshal data")
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	log.WithField("key", key).WithField("size_bytes", len(jsonData)).Debug("Uploading object to S3")
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(jsonData),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		log.WithError(err).WithField("key", key).Error("Failed to put object to S3")
		return fmt.Errorf("failed to put object %s: %w", key, err)
	}

	log.WithField("key", key).Info("Object uploaded successfully")
	return nil
}

func (s *S3Storage) getObject(key string, dest interface{}) error {
	log := logger.GetDefault().WithComponent("S3Storage")
	log.WithField("key", key).Debug("Getting object from S3")

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		// Check if it's a NoSuchKey error (object not found)
		if strings.Contains(err.Error(), "NoSuchKey") || strings.Contains(err.Error(), "NotFound") {
			log.WithField("key", key).Debug("Object not found in S3")
			return ErrNotFound
		}
		log.WithError(err).WithField("key", key).Error("Failed to get object from S3")
		return fmt.Errorf("failed to get object %s: %w", key, err)
	}
	defer result.Body.Close()

	log.WithField("key", key).Debug("Decoding object from S3")
	if err := json.NewDecoder(result.Body).Decode(dest); err != nil {
		log.WithError(err).WithField("key", key).Error("Failed to decode object")
		return fmt.Errorf("failed to decode object %s: %w", key, err)
	}

	log.WithField("key", key).Debug("Object retrieved successfully")
	return nil
}

func (s *S3Storage) deleteObject(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object %s: %w", key, err)
	}

	return nil
}

func (s *S3Storage) listObjects(prefix string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	var keys []string
	fullPrefix := s.prefix + prefix

	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(fullPrefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects with prefix %s: %w", fullPrefix, err)
		}

		for _, obj := range page.Contents {
			if obj.Key != nil {
				keys = append(keys, *obj.Key)
			}
		}
	}

	return keys, nil
}

// Incidents implementation

func (s *S3Storage) CreateIncident(incident *event.Incident) error {
	log := logger.GetDefault().WithComponent("S3Storage")
	log.WithField("incident_id", incident.GUID).Info("Creating incident in S3")

	// Check if incident already exists
	key := s.getKey("incidents", incident.GUID)
	log.WithField("incident_id", incident.GUID).WithField("key", key).Debug("Checking if incident exists")

	var existing event.Incident
	if err := s.getObject(key, &existing); err == nil {
		log.WithField("incident_id", incident.GUID).Warn("Incident already exists")
		return ErrExists
	} else if err != ErrNotFound {
		log.WithError(err).WithField("incident_id", incident.GUID).Error("Error checking incident existence")
		return err
	}

	log.WithField("incident_id", incident.GUID).Debug("Incident does not exist, proceeding with creation")
	if err := s.putObject(key, incident); err != nil {
		log.WithError(err).WithField("incident_id", incident.GUID).Error("Failed to create incident")
		return err
	}

	log.WithField("incident_id", incident.GUID).Info("Incident created successfully in S3")
	return nil
}

func (s *S3Storage) GetIncident(id string) (*event.Incident, error) {
	key := s.getKey("incidents", id)
	var incident event.Incident
	if err := s.getObject(key, &incident); err != nil {
		return nil, err
	}
	return &incident, nil
}

func (s *S3Storage) GetAllIncidents() ([]*event.Incident, error) {
	log := logger.GetDefault().WithComponent("S3Storage")
	log.Info("Retrieving all incidents from S3")

	keys, err := s.listObjects("incidents/")
	if err != nil {
		log.WithError(err).Error("Failed to list incident objects")
		return nil, err
	}

	log.WithField("count", len(keys)).Info("Found incident objects, loading data")
	var incidents []*event.Incident
	errorCount := 0

	for _, key := range keys {
		var incident event.Incident
		if err := s.getObject(key, &incident); err != nil {
			log.WithError(err).WithField("key", key).Warn("Failed to load incident, skipping")
			errorCount++
			continue
		}
		incidents = append(incidents, &incident)
	}

	if errorCount > 0 {
		log.WithField("errors", errorCount).WithField("loaded", len(incidents)).Warn("Some incidents failed to load")
	}

	log.WithField("count", len(incidents)).Info("All incidents retrieved successfully")
	return incidents, nil
}

func (s *S3Storage) UpdateIncident(incident *event.Incident) error {
	// Check if incident exists
	key := s.getKey("incidents", incident.GUID)
	var existing event.Incident
	if err := s.getObject(key, &existing); err != nil {
		if err == ErrNotFound {
			return ErrNotFound
		}
		return err
	}

	return s.putObject(key, incident)
}

func (s *S3Storage) DeleteIncident(id string) error {
	// Check if incident exists
	key := s.getKey("incidents", id)
	var existing event.Incident
	if err := s.getObject(key, &existing); err != nil {
		if err == ErrNotFound {
			return ErrNotFound
		}
		return err
	}

	return s.deleteObject(key)
}

// Planned Maintenances implementation

func (s *S3Storage) CreatePlannedMaintenance(pm *event.PlannedMaintenance) error {
	// Check if planned maintenance already exists
	key := s.getKey("planned_maintenances", pm.GUID)
	var existing event.PlannedMaintenance
	if err := s.getObject(key, &existing); err == nil {
		return ErrExists
	} else if err != ErrNotFound {
		return err
	}

	return s.putObject(key, pm)
}

func (s *S3Storage) GetPlannedMaintenance(id string) (*event.PlannedMaintenance, error) {
	key := s.getKey("planned_maintenances", id)
	var pm event.PlannedMaintenance
	if err := s.getObject(key, &pm); err != nil {
		return nil, err
	}
	return &pm, nil
}

func (s *S3Storage) GetAllPlannedMaintenances() ([]*event.PlannedMaintenance, error) {
	keys, err := s.listObjects("planned_maintenances/")
	if err != nil {
		return nil, err
	}

	var pms []*event.PlannedMaintenance
	for _, key := range keys {
		var pm event.PlannedMaintenance
		if err := s.getObject(key, &pm); err != nil {
			// Log error but continue with other planned maintenances
			continue
		}
		pms = append(pms, &pm)
	}

	return pms, nil
}

func (s *S3Storage) UpdatePlannedMaintenance(pm *event.PlannedMaintenance) error {
	// Check if planned maintenance exists
	key := s.getKey("planned_maintenances", pm.GUID)
	var existing event.PlannedMaintenance
	if err := s.getObject(key, &existing); err != nil {
		if err == ErrNotFound {
			return ErrNotFound
		}
		return err
	}

	return s.putObject(key, pm)
}

func (s *S3Storage) DeletePlannedMaintenance(id string) error {
	// Check if planned maintenance exists
	key := s.getKey("planned_maintenances", id)
	var existing event.PlannedMaintenance
	if err := s.getObject(key, &existing); err != nil {
		if err == ErrNotFound {
			return ErrNotFound
		}
		return err
	}

	return s.deleteObject(key)
}
