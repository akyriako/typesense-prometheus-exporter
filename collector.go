package typesense_prometheus_exporter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type TypesenseCollector struct {
	ctx        context.Context
	logger     *slog.Logger
	endPoint   string
	apiKey     string
	cluster    string
	httpClient *http.Client
	metrics    map[string]*prometheus.Desc
	stats      map[string]*prometheus.Desc
	mutex      sync.Mutex
}

var (
	metricLabels   = []string{"typesense_cluster"}
	endpointLabels = []string{"typesense_cluster", "endpoint"}
)

func NewTypesenseCollector(ctx context.Context, logger *slog.Logger, config Config) *TypesenseCollector {
	collector := &TypesenseCollector{
		ctx:      ctx,
		logger:   logger,
		endPoint: fmt.Sprintf("%s://%s:%d", config.Protocol, config.Host, config.ApiPort),
		apiKey:   config.ApiKey,
		cluster:  config.Cluster,
		httpClient: &http.Client{
			Timeout: 500 * time.Millisecond,
		},
		metrics: getMetricsDesc(),
		stats:   getStatsDesc(),
	}

	return collector
}

// Describe sends the metric descriptors to the Prometheus channel
func (c *TypesenseCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric
	}

	for _, stat := range c.stats {
		ch <- stat
	}
}

// Collect fetches the metrics from the Typesense endpoint and sends them to the Prometheus channel
func (c *TypesenseCollector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	targets := []string{"metrics", "stats"}
	for _, target := range targets {
		data, err := c.fetch(target)
		if err != nil {
			return
		}

		c.collect(target, data, ch)
	}
}

func (c *TypesenseCollector) fetch(target string) (map[string]interface{}, error) {
	start := time.Now()
	url := fmt.Sprintf("%s/%s.json", c.endPoint, target)
	c.logger.Info(fmt.Sprintf("collecting %s...", target), "cluster", c.cluster, "url", url)

	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, url, nil)
	if err != nil {
		c.logger.Error(fmt.Sprintf("error creating request: %v", err))
		return nil, err
	}

	req.Header.Set("x-typesense-api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error(fmt.Sprintf("error fetching %s: %v", target, err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error(fmt.Sprintf("error fetching %s: %v", target, resp.Status))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error(fmt.Sprintf("error reading response body from %s: %v", url, err))
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		c.logger.Error(fmt.Sprintf("error unmarshalling %s.json body: %v", target, err))
		return nil, err
	}

	defer func(count int) {
		elapsed := time.Since(start)
		c.logger.Info(fmt.Sprintf("collecting %s completed", target), "count", count, "duration", elapsed)
	}(len(data))

	return data, nil
}

func (c *TypesenseCollector) collect(target string, data map[string]interface{}, ch chan<- prometheus.Metric) {
	for key, value := range data {
		select {
		case <-c.ctx.Done():
			c.logger.Error(fmt.Sprintf("context canceled, stopping collection"))
			return
		default:
		}

		switch target {
		case "metrics":
			if desc, ok := c.metrics[key]; ok {
				if sval, ok := value.(string); ok {
					val, err := strconv.ParseFloat(sval, 64)
					if err != nil {
						c.logger.Error(fmt.Sprintf("error converting value for %s: %v", key, err))
						continue
					}
					metric := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, val, c.cluster)
					c.logger.Debug(fmt.Sprintf("collected %s", target), "fqName", key, "value", val)

					ch <- metric
				}
			}
		case "stats":
			if nestedData, ok := data[key]; ok {
				if endpoints, ok := nestedData.(map[string]interface{}); ok {
					for endpoint, endpointVal := range endpoints {
						if desc, ok := c.stats[key]; ok {
							if val, ok := endpointVal.(float64); ok {
								stat := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, val, c.cluster, endpoint)
								c.logger.Debug(fmt.Sprintf("collected %s", target), "fqName", key, "endpoint", endpoint, "value", val)

								ch <- stat
							}
						}
					}
				} else {
					if desc, ok := c.stats[key]; ok {
						if val, ok := value.(float64); ok {
							stat := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, val, c.cluster)
							c.logger.Debug(fmt.Sprintf("collected %s", target), "fqName", key, "value", val)

							ch <- stat
						}
					}
				}
			}
		}
	}
}

func getMetricsDesc() map[string]*prometheus.Desc {
	return map[string]*prometheus.Desc{
		"system_cpu1_active_percentage": prometheus.NewDesc(
			"system_cpu1_active_percentage",
			"System CPU core 1 active percentage",
			metricLabels, nil,
		),
		"system_cpu2_active_percentage": prometheus.NewDesc(
			"system_cpu2_active_percentage",
			"System CPU core 2 active percentage",
			metricLabels, nil,
		),
		"system_cpu3_active_percentage": prometheus.NewDesc(
			"system_cpu3_active_percentage",
			"System CPU core 3 active percentage",
			metricLabels, nil,
		),
		"system_cpu4_active_percentage": prometheus.NewDesc(
			"system_cpu4_active_percentage",
			"System CPU core 4 active percentage",
			metricLabels, nil,
		),
		"system_cpu_active_percentage": prometheus.NewDesc(
			"system_cpu_active_percentage",
			"System overall CPU active percentage",
			metricLabels, nil,
		),
		"system_disk_total_bytes": prometheus.NewDesc(
			"system_disk_total_bytes",
			"Total system disk space in bytes",
			metricLabels, nil,
		),
		"system_disk_used_bytes": prometheus.NewDesc(
			"system_disk_used_bytes",
			"Used system disk space in bytes",
			metricLabels, nil,
		),
		"system_memory_total_bytes": prometheus.NewDesc(
			"system_memory_total_bytes",
			"Total system memory in bytes",
			metricLabels, nil,
		),
		"system_memory_used_bytes": prometheus.NewDesc(
			"system_memory_used_bytes",
			"Used system memory in bytes",
			metricLabels, nil,
		),
		"system_memory_total_swap_bytes": prometheus.NewDesc(
			"system_memory_total_swap_bytes",
			"Total system swap memory in bytes",
			metricLabels, nil,
		),
		"system_memory_used_swap_bytes": prometheus.NewDesc(
			"system_memory_used_swap_bytes",
			"Used system swap memory in bytes",
			metricLabels, nil,
		),
		"system_network_received_bytes": prometheus.NewDesc(
			"system_network_received_bytes",
			"Total network received bytes",
			metricLabels, nil,
		),
		"system_network_sent_bytes": prometheus.NewDesc(
			"system_network_sent_bytes",
			"Total network sent bytes",
			metricLabels, nil,
		),
		"typesense_memory_active_bytes": prometheus.NewDesc(
			"typesense_memory_active_bytes",
			"Typesense active memory usage in bytes",
			metricLabels, nil,
		),
		"typesense_memory_allocated_bytes": prometheus.NewDesc(
			"typesense_memory_allocated_bytes",
			"Typesense allocated memory in bytes",
			metricLabels, nil,
		),
		"typesense_memory_fragmentation_ratio": prometheus.NewDesc(
			"typesense_memory_fragmentation_ratio",
			"Typesense memory fragmentation ratio",
			metricLabels, nil,
		),
		"typesense_memory_mapped_bytes": prometheus.NewDesc(
			"typesense_memory_mapped_bytes",
			"Typesense memory mapped in bytes",
			metricLabels, nil,
		),
		"typesense_memory_metadata_bytes": prometheus.NewDesc(
			"typesense_memory_metadata_bytes",
			"Typesense memory metadata size in bytes",
			metricLabels, nil,
		),
		"typesense_memory_resident_bytes": prometheus.NewDesc(
			"typesense_memory_resident_bytes",
			"Typesense resident memory usage in bytes",
			metricLabels, nil,
		),
		"typesense_memory_retained_bytes": prometheus.NewDesc(
			"typesense_memory_retained_bytes",
			"Typesense retained memory in bytes",
			metricLabels, nil,
		),
	}
}

func getStatsDesc() map[string]*prometheus.Desc {
	return map[string]*prometheus.Desc{
		"uptime": prometheus.NewDesc(
			"uptime_seconds",
			"Uptime in seconds",
			metricLabels,
			nil,
		),
		"memory.used": prometheus.NewDesc(
			"memory_used_bytes",
			"Memory used in bytes",
			metricLabels,
			nil,
		),
		"memory.total": prometheus.NewDesc(
			"memory_total_bytes",
			"Total memory in bytes",
			metricLabels,
			nil,
		),
		"memory.available": prometheus.NewDesc(
			"memory_available_bytes",
			"Available memory in bytes",
			metricLabels,
			nil,
		),
		"memory.resident": prometheus.NewDesc(
			"memory_resident_bytes",
			"Resident memory in bytes",
			metricLabels,
			nil,
		),
		"delete_latency_ms": prometheus.NewDesc(
			"delete_latency_ms",
			"Latency of delete operations",
			metricLabels,
			nil,
		),
		"delete_requests_per_second": prometheus.NewDesc(
			"delete_requests_per_second",
			"Delete requests per second",
			metricLabels,
			nil,
		),
		"import_latency_ms": prometheus.NewDesc(
			"import_latency_ms",
			"Latency of import operations",
			metricLabels,
			nil,
		),
		"import_requests_per_second": prometheus.NewDesc(
			"import_requests_per_second",
			"Import requests per second",
			metricLabels,
			nil,
		),
		"overloaded_requests_per_second": prometheus.NewDesc(
			"overloaded_requests_per_second",
			"Overloaded requests per second",
			metricLabels,
			nil,
		),
		"pending_write_batches": prometheus.NewDesc(
			"pending_write_batches",
			"Pending write batches",
			metricLabels,
			nil,
		),
		"search_latency_ms": prometheus.NewDesc(
			"search_latency_ms",
			"Latency of search operations",
			metricLabels,
			nil,
		),
		"search_requests_per_second": prometheus.NewDesc(
			"search_requests_per_second",
			"Search requests per second",
			metricLabels,
			nil,
		),
		"total_requests_per_second": prometheus.NewDesc(
			"total_requests_per_second",
			"Total requests per second",
			metricLabels,
			nil,
		),
		"write_latency_ms": prometheus.NewDesc(
			"write_latency_ms",
			"Latency of write operations",
			metricLabels,
			nil,
		),
		"write_requests_per_second": prometheus.NewDesc(
			"write_requests_per_second",
			"Write requests per second",
			metricLabels,
			nil,
		),
		"latency_ms": prometheus.NewDesc(
			"latency_ms",
			"Latency for specific endpoints",
			endpointLabels,
			nil,
		),
		"requests_per_second": prometheus.NewDesc(
			"requests_per_second",
			"Requests per second for specific endpoints",
			endpointLabels,
			nil,
		),
	}
}
