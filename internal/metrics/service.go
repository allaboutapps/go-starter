package metrics

import (
	"context"
	"database/sql"
	"fmt"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/metrics/users"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/dlmiddlecote/sqlstats"
	"github.com/prometheus/client_golang/prometheus"
)

type Service struct {
	config config.Server
	db     *sql.DB
}

func New(config config.Server, db *sql.DB) (*Service, error) {
	return &Service{
		config: config,
		db:     db,
	}, nil
}

func (s *Service) RegisterMetrics(ctx context.Context) error {
	log := util.LogFromContext(ctx)

	var metrics []prometheus.Collector

	// custom metrics
	metrics = append(metrics, users.Metrics(ctx, users.NewDatabaseMetricsCollector(s.db))...)

	// sqlstats metrics, see https://github.com/dlmiddlecote/sqlstats?tab=readme-ov-file#exposed-metrics for the exposed metrics
	metrics = append(metrics, sqlstats.NewStatsCollector(s.config.Database.Database, s.db))

	for _, metric := range metrics {
		if err := prometheus.Register(metric); err != nil {
			log.Error().Err(err).Msg("Failed to register metric")
			return fmt.Errorf("failed to register metric: %w", err)
		}
	}

	return nil
}
