package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Depado/ginprom"
	"github.com/aurowora/compress"
	"github.com/earthboundkid/versioninfo/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/kofalt/go-memoize"
	"github.com/map-services/street-manager-relay/internal"
	"github.com/map-services/street-manager-relay/internal/middleware"
	"github.com/map-services/street-manager-relay/internal/promoter"
	"github.com/map-services/street-manager-relay/internal/routes"
	"github.com/rm-hull/godx"
	"github.com/tavsec/gin-healthcheck/checks"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"

	healthcheck "github.com/tavsec/gin-healthcheck"
	hc_config "github.com/tavsec/gin-healthcheck/config"
)

func ApiServer(dbPath string, port int, debug bool, excludeActivityTypes []string) {
	logger := internal.SetupLogger()
	godx.Diagnostics(logger)

	organisations, err := promoter.GetPromoterOrgsMap()
	if err != nil {
		slog.Error("failed to initialize promoter organisations", "error", err)
		os.Exit(1)
	}

	repo, err := internal.NewDbRepository(dbPath)
	if err != nil {
		slog.Error("Failed to initialize db repository", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := repo.Close(); err != nil {
			slog.Error("Error closing database", "error", err)
		}
	}()

	err = sentry.Init(sentry.ClientOptions{
		Dsn:         os.Getenv("SENTRY_DSN"),
		Debug:       debug,
		Release:     versioninfo.Revision[:7],
		Environment: os.Getenv("MODE"),
	})
	if err != nil {
		slog.Error("sentry.Init", "error", err)
		os.Exit(1)
	}
	defer sentry.Flush(2 * time.Second)

	r := gin.New()

	prometheus := ginprom.New(
		ginprom.Engine(r),
		ginprom.Path("/metrics"),
		ginprom.Ignore("/healthz"),
	)

	r.Use(
		sentrygin.New(sentrygin.Options{
			Repanic:         true,
			WaitForDelivery: false,
			Timeout:         5 * time.Second,
		}),
		gin.Recovery(),
		middleware.RequestLogger(logger, "/healthz", "/metrics"),
		prometheus.Instrument(),
		compress.Compress(),
		cors.Default(),
		sentryErrorHandler(),
	)

	if debug {
		slog.Warn("pprof endpoints are enabled and exposed. Do not run with this flag in production.")
		pprof.Register(r)
	}

	err = healthcheck.New(r, hc_config.DefaultConfig(), []checks.Check{
		repo.HealthCheck(),
	})
	if err != nil {
		slog.Error("failed to initialize healthcheck", "error", err)
		os.Exit(1)
	}

	if len(excludeActivityTypes) > 0 {
		slog.Info("Excluing events","sctivityTypes", excludeActivityTypes)
	}

	certManager := internal.NewCertManager(memoize.NewMemoizer(24*time.Hour, 1*time.Hour))

	r.POST("/v1/street-manager-relay/sns", routes.HandleSNSMessage(repo, certManager, excludeActivityTypes))
	r.GET("/v1/street-manager-relay/search", routes.HandleSearch(repo, organisations))
	r.GET("/v1/street-manager-relay/refdata", routes.HandleRefData(repo, memoize.NewMemoizer(10*time.Minute, 1*time.Hour)))

	addr := fmt.Sprintf(":%d", port)
	slog.Info("Starting HTTP API Server", "port", port)
	err = r.Run(addr)
	slog.Error("HTTP API Server failed to start", "port", port, "error", err)
	os.Exit(1)
}

func sentryErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			hub := sentrygin.GetHubFromContext(c)
			for _, e := range c.Errors {
				if hub != nil {
					hub.CaptureException(e.Err)
				} else {
					sentry.CaptureException(e.Err)
				}
			}
		}
	}
}
