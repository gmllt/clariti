package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/common/version"
)

var (
	// HTTP Metrics - navigation instruments for request monitoring
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "clariti_http_requests_total",
			Help: "Total number of HTTP requests processed by flight control",
		},
		[]string{"method", "path", "status_code"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "clariti_http_request_duration_seconds",
			Help:    "Duration of HTTP requests - flight time measurement",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Business Metrics - operational status instruments
	incidentsTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "clariti_incidents_total",
			Help: "Current number of incidents by severity - system health indicators",
		},
		[]string{"severity", "status"},
	)

	plannedMaintenancesTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "clariti_planned_maintenances_total",
			Help: "Current number of planned maintenances by status - scheduled operations",
		},
		[]string{"status"},
	)

	componentsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "clariti_components_total",
			Help: "Total number of monitored components - fleet size",
		},
	)

	// System Metrics - aircraft status instruments
	applicationInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "clariti_application_info",
			Help: "Application information - aircraft identification",
		},
		[]string{"version", "build_date", "git_commit"},
	)

	uptimeSeconds = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "clariti_uptime_seconds",
			Help: "Application uptime in seconds - flight duration",
		},
	)
)

var startTime = time.Now()

// RecordHTTPRequest logs HTTP request metrics - flight data recording
func RecordHTTPRequest(method, path string, statusCode int, duration time.Duration) {
	httpRequestsTotal.WithLabelValues(method, path, strconv.Itoa(statusCode)).Inc()
	httpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

// UpdateIncidentMetrics refreshes incident counters - emergency status update
func UpdateIncidentMetrics(severity, status string, count float64) {
	incidentsTotal.WithLabelValues(severity, status).Set(count)
}

// UpdatePlannedMaintenanceMetrics refreshes maintenance counters - scheduled operation status
func UpdatePlannedMaintenanceMetrics(status string, count float64) {
	plannedMaintenancesTotal.WithLabelValues(status).Set(count)
}

// UpdateComponentsCount refreshes component counter - fleet size update
func UpdateComponentsCount(count float64) {
	componentsTotal.Set(count)
}

// UpdateUptimeMetrics refreshes uptime counter - flight time update
func UpdateUptimeMetrics() {
	uptimeSeconds.Set(time.Since(startTime).Seconds())
}

// InitializeApplicationInfo sets application metadata - aircraft registration
func InitializeApplicationInfo() {
	applicationInfo.WithLabelValues(
		version.Version,
		version.BuildDate,
		version.Revision,
	).Set(1)
}

// GetRegistry returns the Prometheus registry - control tower access
func GetRegistry() *prometheus.Registry {
	return prometheus.DefaultRegisterer.(*prometheus.Registry)
}
