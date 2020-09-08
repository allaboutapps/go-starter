package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/mailer"
	"allaboutapps.dev/aw/go-starter/internal/mailer/transport"
	"allaboutapps.dev/aw/go-starter/internal/push"
	"allaboutapps.dev/aw/go-starter/internal/push/provider"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	// Import postgres driver for database/sql package
	_ "github.com/lib/pq"
)

type Router struct {
	Routes     []*echo.Route
	Root       *echo.Group
	Management *echo.Group
	APIV1Auth  *echo.Group
	APIV1Push  *echo.Group
}

type Server struct {
	Config config.Server
	DB     *sql.DB
	Echo   *echo.Echo
	Router *Router
	Mailer *mailer.Mailer
	Push   *push.Service
}

func NewServer(config config.Server) *Server {
	s := &Server{
		Config: config,
		DB:     nil,
		Echo:   nil,
		Router: nil,
		Mailer: nil,
		Push:   nil,
	}

	return s
}

func (s *Server) Ready() bool {
	return s.DB != nil &&
		s.Echo != nil &&
		s.Router != nil &&
		s.Mailer != nil &&
		s.Push != nil
}

func (s *Server) InitDB(ctx context.Context) error {
	db, err := sql.Open("postgres", s.Config.Database.ConnectionString())
	if err != nil {
		return err
	}

	if s.Config.Database.MaxOpenConns > 0 {
		db.SetMaxOpenConns(s.Config.Database.MaxOpenConns)
	}
	if s.Config.Database.MaxIdleConns > 0 {
		db.SetMaxIdleConns(s.Config.Database.MaxIdleConns)
	}
	if s.Config.Database.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(s.Config.Database.ConnMaxLifetime)
	}

	if err := db.PingContext(ctx); err != nil {
		return err
	}

	s.DB = db

	return nil
}

func (s *Server) InitMailer() error {
	switch config.MailerTransporter(s.Config.Mailer.Transporter) {
	case config.MailerTransporterMock:
		log.Warn().Msg("Initializing mock mailer")
		s.Mailer = mailer.New(s.Config.Mailer, transport.NewMock())
	case config.MailerTransporterSMTP:
		s.Mailer = mailer.New(s.Config.Mailer, transport.NewSMTP(s.Config.SMTP))
	default:
		return fmt.Errorf("Unsupported mail transporter: %s", s.Config.Mailer.Transporter)
	}

	return s.Mailer.ParseTemplates()
}

func (s *Server) InitPush() error {
	s.Push = push.New(s.DB)

	if s.Config.Push.UseFCMProvider {
		fcmProvider, err := provider.NewFCM(s.Config.FCMConfig)
		if err != nil {
			return err
		}
		s.Push.RegisterProvider(fcmProvider)
	}

	if s.Config.Push.UseMockProvider {
		log.Warn().Msg("Initializing mock push provider")
		mockProvider := provider.NewMock(push.ProviderTypeFCM)
		s.Push.RegisterProvider(mockProvider)
	}

	if s.Push.GetProviderCount() < 1 {
		log.Warn().Msg("No providers registered for push service")
	}

	return nil
}

func (s *Server) Start() error {
	if !s.Ready() {
		return errors.New("server is not ready")
	}

	return s.Echo.Start(s.Config.Echo.ListenAddress)
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Warn().Msg("Shutting down server")

	if s.DB != nil {
		log.Debug().Msg("Closing database connection")

		if err := s.DB.Close(); err != nil && err != sql.ErrConnDone {
			log.Error().Err(err).Msg("Failed to close database connection")
		}
	}

	log.Debug().Msg("Shutting down echo server")

	return s.Echo.Shutdown(ctx)
}
