package typesense_prometheus_exporter

import "github.com/prometheus/client_golang/prometheus"

func getMetricsDesc() map[string]*prometheus.Desc {
	return map[string]*prometheus.Desc{
		"system_cpu1_active_percentage": prometheus.NewDesc(
			"typesense_system_cpu1_active_percentage",
			"System CPU core 1 active percentage",
			metricLabels, nil,
		),
		"system_cpu2_active_percentage": prometheus.NewDesc(
			"typesense_system_cpu2_active_percentage",
			"System CPU core 2 active percentage",
			metricLabels, nil,
		),
		"system_cpu3_active_percentage": prometheus.NewDesc(
			"typesense_system_cpu3_active_percentage",
			"System CPU core 3 active percentage",
			metricLabels, nil,
		),
		"system_cpu4_active_percentage": prometheus.NewDesc(
			"typesense_system_cpu4_active_percentage",
			"System CPU core 4 active percentage",
			metricLabels, nil,
		),
		"system_cpu_active_percentage": prometheus.NewDesc(
			"typesense_system_cpu_active_percentage",
			"System overall CPU active percentage",
			metricLabels, nil,
		),
		"system_disk_total_bytes": prometheus.NewDesc(
			"typesense_system_disk_total_bytes",
			"Total system disk space in bytes",
			metricLabels, nil,
		),
		"system_disk_used_bytes": prometheus.NewDesc(
			"typesense_system_disk_used_bytes",
			"Used system disk space in bytes",
			metricLabels, nil,
		),
		"system_memory_total_bytes": prometheus.NewDesc(
			"typesense_system_memory_total_bytes",
			"Total system memory in bytes",
			metricLabels, nil,
		),
		"system_memory_used_bytes": prometheus.NewDesc(
			"typesense_system_memory_used_bytes",
			"Used system memory in bytes",
			metricLabels, nil,
		),
		"system_memory_total_swap_bytes": prometheus.NewDesc(
			"typesense_system_memory_total_swap_bytes",
			"Total system swap memory in bytes",
			metricLabels, nil,
		),
		"system_memory_used_swap_bytes": prometheus.NewDesc(
			"typesense_system_memory_used_swap_bytes",
			"Used system swap memory in bytes",
			metricLabels, nil,
		),
		"system_network_received_bytes": prometheus.NewDesc(
			"typesense_system_network_received_bytes",
			"Total network received bytes",
			metricLabels, nil,
		),
		"system_network_sent_bytes": prometheus.NewDesc(
			"typesense_system_network_sent_bytes",
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
			"typesense_uptime_seconds",
			"Uptime in seconds",
			metricLabels,
			nil,
		),
		"memory.used": prometheus.NewDesc(
			"typesense_memory_used_bytes",
			"Memory used in bytes",
			metricLabels,
			nil,
		),
		"memory.total": prometheus.NewDesc(
			"typesense_memory_total_bytes",
			"Total memory in bytes",
			metricLabels,
			nil,
		),
		"memory.available": prometheus.NewDesc(
			"typesense_memory_available_bytes",
			"Available memory in bytes",
			metricLabels,
			nil,
		),
		"memory.resident": prometheus.NewDesc(
			"typesense_memory_resident_bytes",
			"Resident memory in bytes",
			metricLabels,
			nil,
		),
		"delete_latency_ms": prometheus.NewDesc(
			"typesense_delete_latency_ms",
			"Latency of delete operations",
			metricLabels,
			nil,
		),
		"delete_requests_per_second": prometheus.NewDesc(
			"typesense_delete_requests_per_second",
			"Delete requests per second",
			metricLabels,
			nil,
		),
		"import_latency_ms": prometheus.NewDesc(
			"typesense_import_latency_ms",
			"Latency of import operations",
			metricLabels,
			nil,
		),
		"import_requests_per_second": prometheus.NewDesc(
			"typesense_import_requests_per_second",
			"Import requests per second",
			metricLabels,
			nil,
		),
		"overloaded_requests_per_second": prometheus.NewDesc(
			"typesense_overloaded_requests_per_second",
			"Overloaded requests per second",
			metricLabels,
			nil,
		),
		"pending_write_batches": prometheus.NewDesc(
			"typesense_pending_write_batches",
			"Pending write batches",
			metricLabels,
			nil,
		),
		"search_latency_ms": prometheus.NewDesc(
			"typesense_search_latency_ms",
			"Latency of search operations",
			metricLabels,
			nil,
		),
		"search_requests_per_second": prometheus.NewDesc(
			"typesense_search_requests_per_second",
			"Search requests per second",
			metricLabels,
			nil,
		),
		"total_requests_per_second": prometheus.NewDesc(
			"typesense_total_requests_per_second",
			"Total requests per second",
			metricLabels,
			nil,
		),
		"write_latency_ms": prometheus.NewDesc(
			"typesense_write_latency_ms",
			"Latency of write operations",
			metricLabels,
			nil,
		),
		"write_requests_per_second": prometheus.NewDesc(
			"typesense_write_requests_per_second",
			"Write requests per second",
			metricLabels,
			nil,
		),
		"latency_ms": prometheus.NewDesc(
			"typesense_latency_ms",
			"Latency for specific endpoints",
			endpointLabels,
			nil,
		),
		"requests_per_second": prometheus.NewDesc(
			"typesense_requests_per_second",
			"Requests per second for specific endpoints",
			endpointLabels,
			nil,
		),
	}
}
