package users

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricsCollector interface {
	GetTotalUsersCount(ctx context.Context) float64
}

const (
	MetricNameTotalUsers = "total_users"
)

func Metrics(ctx context.Context, collector MetricsCollector) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name: MetricNameTotalUsers,
				Help: "Total users",
			},
			func() float64 { return collector.GetTotalUsersCount(ctx) },
		),
	}
}
