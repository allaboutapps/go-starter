package metrics

import (
	"context"
	"database/sql"

	"allaboutapps.dev/aw/go-starter/internal/metrics/users"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/prometheus/client_golang/prometheus"
)

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) (*Service, error) {
	return &Service{
		db: db,
	}, nil
}

func (s *Service) RegisterMetrics(ctx context.Context) error {
	log := util.LogFromContext(ctx)

	var metrics []prometheus.Collector
	metrics = append(metrics, users.Metrics(ctx, users.NewDatabaseMetricsCollector(s.db))...)

	for _, metric := range metrics {
		if err := prometheus.Register(metric); err != nil {
			log.Error().Err(err).Msg("Failed to register metric")
			return err
		}
	}

	return nil
}
