package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/auth"
	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/data/dto"
	"allaboutapps.dev/aw/go-starter/internal/data/local"
	"allaboutapps.dev/aw/go-starter/internal/i18n"
	"allaboutapps.dev/aw/go-starter/internal/mailer"
	"allaboutapps.dev/aw/go-starter/internal/mailer/transport"
	"allaboutapps.dev/aw/go-starter/internal/metrics"
	"allaboutapps.dev/aw/go-starter/internal/push"
	"allaboutapps.dev/aw/go-starter/internal/push/provider"
	"github.com/dropbox/godropbox/time2"
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
	Config  config.Server
	DB      *sql.DB
	Echo    *echo.Echo
	Router  *Router
	Mailer  *mailer.Mailer
	Push    *push.Service
	I18n    *i18n.Service
	Clock   time2.Clock
	Auth    AuthService
	Local   *local.Service
	Metrics *metrics.Service
}

type AuthService interface {
	GetAppUserProfileIfExists(context.Context, string) (*dto.AppUserProfile, error)
	InitPasswordReset(context.Context, dto.InitPasswordResetRequest) (dto.InitPasswordResetResult, error)
	Login(context.Context, dto.LoginRequest) (dto.LoginResult, error)
	Logout(context.Context, dto.LogoutRequest) error
	Refresh(context.Context, dto.RefreshRequest) (dto.LoginResult, error)
	Register(context.Context, dto.RegisterRequest) (dto.RegisterResult, error)
	CompleteRegister(context.Context, dto.CompleteRegisterRequest) (dto.LoginResult, error)
	DeleteUserAccount(context.Context, dto.DeleteUserAccountRequest) error
	ResetPassword(context.Context, dto.ResetPasswordRequest) (dto.LoginResult, error)
	UpdatePassword(context.Context, dto.UpdatePasswordRequest) (dto.LoginResult, error)
}

func NewServer(config config.Server) *Server {
	s := &Server{
		Config: config,
		DB:     nil,
		Echo:   nil,
		Router: nil,
		Mailer: nil,
		Push:   nil,
		I18n:   nil,
		Clock:  nil,
		Auth:   nil,
		Local:  nil,
	}

	return s
}

func (s *Server) Ready() bool {
	return s.DB != nil &&
		s.Echo != nil &&
		s.Router != nil &&
		s.Mailer != nil &&
		s.Push != nil &&
		s.I18n != nil &&
		s.Clock != nil &&
		s.Auth != nil &&
		s.Local != nil
}

func (s *Server) InitCmd() *Server {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := s.InitDB(ctx); err != nil {
		cancel()
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	cancel()

	if err := s.InitClock(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize clock")
	}

	if err := s.InitMailer(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize mailer")
	}

	if err := s.InitPush(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize push service")
	}

	if err := s.InitI18n(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize i18n service")
	}

	if err := s.InitAuthService(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize auth service")
	}

	if err := s.InitLocalService(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize local service")
	}

	if err := s.InitMetricsService(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize metrics service")
	}

	return s
}

func (s *Server) InitAuthService() error {
	s.Auth = auth.NewService(s.Config, s.DB, s.Clock)

	return nil
}

func (s *Server) InitLocalService() error {
	s.Local = local.NewService(s.Config, s.DB, s.Clock)

	return nil
}

func (s *Server) InitMetricsService() error {
	var err error
	s.Metrics, err = metrics.New(s.Config, s.DB)
	if err != nil {
		return err
	}

	return nil
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

func (s *Server) InitClock() error {
	s.Clock = time2.DefaultClock
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

func (s *Server) InitI18n() error {
	i18nService, err := i18n.New(s.Config.I18n)

	if err != nil {
		return err
	}

	s.I18n = i18nService

	return nil
}

func (s *Server) Start() error {
	if !s.Ready() {
		return errors.New("server is not ready")
	}

	return s.Echo.Start(s.Config.Echo.ListenAddress)
}

func (s *Server) Shutdown(ctx context.Context) []error {
	log.Warn().Msg("Shutting down server")

	var errs []error

	if s.DB != nil {
		log.Debug().Msg("Closing database connection")

		if err := s.DB.Close(); err != nil && !errors.Is(err, sql.ErrConnDone) {
			log.Error().Err(err).Msg("Failed to close database connection")
			errs = append(errs, err)
		}
	}

	if s.Echo != nil {
		log.Debug().Msg("Shutting down echo server")

		if err := s.Echo.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("Failed to shutdown echo server")
			errs = append(errs, err)
		}

	}

	return errs
}
