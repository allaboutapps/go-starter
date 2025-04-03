package users

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricsCollector interface {
	GetTotalUsersCount(ctx context.Context) float64
}

const (
	MetricNameTotalUserCount = "total_user_count"
)

func Metrics(ctx context.Context, collector MetricsCollector) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name: MetricNameTotalUserCount,
			},
			func() float64 { return collector.GetTotalUsersCount(ctx) },
		),
	}
}
