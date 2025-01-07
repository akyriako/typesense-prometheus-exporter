package typesense_prometheus_exporter

type Config struct {
	LogLevel    int    `env:"LOG_LEVEL" envDefault:"0"`
	ApiKey      string `env:"TYPESENSE_API_KEY,required"`
	Host        string `env:"TYPESENSE_HOST,required"`
	ApiPort     uint   `env:"TYPESENSE_PORT,required" envDefault:"8108"`
	MetricsPort uint   `env:"METRICS_PORT,required" envDefault:"9090"`
	Protocol    string `env:"TYPESENSE_PROTOCOL,required" envDefault:"http"`
	Namespace   string `env:"POD_NAMESPACE,required" envDefault:"~empty"`
	Cluster     string `env:"TYPESENSE_CLUSTER,required"`
}
