package typesense_prometheus_exporter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type TypesenseCollector struct {
	ctx        context.Context
	logger     *slog.Logger
	endPoint   string
	apiKey     string
	namespace  string
	cluster    string
	httpClient *http.Client
	metrics    map[string]*prometheus.Desc
}

var (
	labels = []string{"namespace", "typesense_cluster"}
)

func NewTypesenseCollector(ctx context.Context, logger *slog.Logger, config Config) *TypesenseCollector {
	collector := &TypesenseCollector{
		ctx:       ctx,
		logger:    logger,
		endPoint:  fmt.Sprintf("%s://%s:%d", config.Protocol, config.Host, config.ApiPort),
		apiKey:    config.ApiKey,
		namespace: config.Namespace,
		cluster:   config.Cluster,
		httpClient: &http.Client{
			Timeout: 500 * time.Millisecond,
		},
		metrics: map[string]*prometheus.Desc{
			"system_cpu1_active_percentage": prometheus.NewDesc(
				"system_cpu1_active_percentage",
				"System CPU core 1 active percentage",
				labels, nil,
			),
			"system_cpu2_active_percentage": prometheus.NewDesc(
				"system_cpu2_active_percentage",
				"System CPU core 2 active percentage",
				labels, nil,
			),
			"system_cpu3_active_percentage": prometheus.NewDesc(
				"system_cpu3_active_percentage",
				"System CPU core 3 active percentage",
				labels, nil,
			),
			"system_cpu4_active_percentage": prometheus.NewDesc(
				"system_cpu4_active_percentage",
				"System CPU core 4 active percentage",
				labels, nil,
			),
			"system_cpu_active_percentage": prometheus.NewDesc(
				"system_cpu_active_percentage",
				"System overall CPU active percentage",
				labels, nil,
			),
			"system_disk_total_bytes": prometheus.NewDesc(
				"system_disk_total_bytes",
				"Total system disk space in bytes",
				labels, nil,
			),
			"system_disk_used_bytes": prometheus.NewDesc(
				"system_disk_used_bytes",
				"Used system disk space in bytes",
				labels, nil,
			),
			"system_memory_total_bytes": prometheus.NewDesc(
				"system_memory_total_bytes",
				"Total system memory in bytes",
				labels, nil,
			),
			"system_memory_used_bytes": prometheus.NewDesc(
				"system_memory_used_bytes",
				"Used system memory in bytes",
				labels, nil,
			),
			"system_memory_total_swap_bytes": prometheus.NewDesc(
				"system_memory_total_swap_bytes",
				"Total system swap memory in bytes",
				labels, nil,
			),
			"system_memory_used_swap_bytes": prometheus.NewDesc(
				"system_memory_used_swap_bytes",
				"Used system swap memory in bytes",
				labels, nil,
			),
			"system_network_received_bytes": prometheus.NewDesc(
				"system_network_received_bytes",
				"Total network received bytes",
				labels, nil,
			),
			"system_network_sent_bytes": prometheus.NewDesc(
				"system_network_sent_bytes",
				"Total network sent bytes",
				labels, nil,
			),
			"typesense_memory_active_bytes": prometheus.NewDesc(
				"typesense_memory_active_bytes",
				"Typesense active memory usage in bytes",
				labels, nil,
			),
			"typesense_memory_allocated_bytes": prometheus.NewDesc(
				"typesense_memory_allocated_bytes",
				"Typesense allocated memory in bytes",
				labels, nil,
			),
			"typesense_memory_fragmentation_ratio": prometheus.NewDesc(
				"typesense_memory_fragmentation_ratio",
				"Typesense memory fragmentation ratio",
				labels, nil,
			),
			"typesense_memory_mapped_bytes": prometheus.NewDesc(
				"typesense_memory_mapped_bytes",
				"Typesense memory mapped in bytes",
				labels, nil,
			),
			"typesense_memory_metadata_bytes": prometheus.NewDesc(
				"typesense_memory_metadata_bytes",
				"Typesense memory metadata size in bytes",
				labels, nil,
			),
			"typesense_memory_resident_bytes": prometheus.NewDesc(
				"typesense_memory_resident_bytes",
				"Typesense resident memory usage in bytes",
				labels, nil,
			),
			"typesense_memory_retained_bytes": prometheus.NewDesc(
				"typesense_memory_retained_bytes",
				"Typesense retained memory in bytes",
				labels, nil,
			),
		},
	}

	return collector
}

// Describe sends the metric descriptors to the Prometheus channel
func (c *TypesenseCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric
	}
}

// Collect fetches the metrics from the Typesense endpoint and sends them to the Prometheus channel
func (c *TypesenseCollector) Collect(ch chan<- prometheus.Metric) {
	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, fmt.Sprintf("%s/metrics.json", c.endPoint), nil)
	if err != nil {
		c.logger.Error(fmt.Sprintf("error creating request: %v", err))
		return
	}

	req.Header.Set("x-typesense-api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error(fmt.Sprintf("error fetching metrics: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error(fmt.Sprintf("error fetching metrics: %v", resp.Status))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error(fmt.Sprintf("error reading response body: %v", err))
		return
	}

	var data map[string]string
	if err := json.Unmarshal(body, &data); err != nil {
		c.logger.Error(fmt.Sprintf("error unmarshalling metrics.json body: %v", err))
		return
	}

	for key, value := range data {
		select {
		case <-c.ctx.Done():
			c.logger.Error(fmt.Sprintf("context canceled, stopping collection"))
			return
		default:
		}

		if desc, ok := c.metrics[key]; ok {
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				c.logger.Error(fmt.Sprintf("error converting value for %s: %v", key, err))
				continue
			}

			metric := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, val, c.namespace, c.cluster)
			c.logger.Debug("collected metric", "fqName", key, "value", val)

			ch <- metric
		}
	}
}
