package runner

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/mdonahue-godaddy/go-http-server/pkg/config"
	"github.com/mdonahue-godaddy/go-http-server/pkg/http/server"
	"github.com/mdonahue-godaddy/go-http-server/pkg/metrics"
	"github.com/mdonahue-godaddy/go-http-server/pkg/shared"
)

const (
	ServiceName string = "go-http-server" // do not change this value, it is used in multiple locations
)

func GetAllEnvironmentVariables() map[string]string {
	envMap := map[string]string{}

	for _, element := range os.Environ() {
		parts := strings.Split(element, "=")
		if len(parts) == 2 && len(parts[0]) > 0 {
			envMap[parts[0]] = parts[1]
		}
	}

	return envMap
}

func SetupDefaultLogrusConfig() {
	//nolint
	// Log as JSON with UTC times instead of the default formatter.
	jsonFormatter := log.JSONFormatter{
		TimestampFormat: shared.TimestampFormat,
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "@timestamp",
			log.FieldKeyLevel: "level",
			log.FieldKeyMsg:   "message",
			log.FieldKeyFunc:  "caller",
		},
	}

	utcFormatter := shared.UTCFormatter{
		Formatter: &jsonFormatter,
	}

	log.SetFormatter(utcFormatter)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// set default log level
	log.SetLevel(log.InfoLevel)
}

// Run - do the work
func Run(version, date, commit string) {
	//nolint
	method := fmt.Sprintf("runner.Run('%s', '%s', '%s')", version, date, commit)
	ctx := shared.CreateContext(context.Background(), ServiceName, shared.ActionTypeService)

	// Setup universal context values
	shared.Init(ServiceName)

	SetupDefaultLogrusConfig()

	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Infof("%s starting...", method)
	//log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false, shared.KeyProcessConfigEnv, GetAllEnvironmentVariables())).Infof("%s Environment Variables", method)

	jsonFileName := fmt.Sprintf("./config/%s.json", ServiceName)
	cfg, jerr := config.LoadSettings(jsonFileName)
	if jerr != nil {
		//To Do
	}

	value, err := shared.Struct2JSONString(cfg)
	if err == nil {
		log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false, shared.KeyAppConfig, value)).Infof("%s current config", method)
	}

	// Set log leve from config
	lvl, perr := log.ParseLevel(cfg.Logging.Level)
	if perr == nil {
		log.SetLevel(lvl)
	}

	// start pprof & metrics services
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Infof("%s setup metrics pprof end point", method)
	dialAddress := net.JoinHostPort(cfg.Metrics.HTTP.Server.IPv4Address, strconv.FormatUint(uint64(cfg.Metrics.HTTP.Server.Port), 10))
	go metrics.StartServer(ctx, dialAddress, cfg.Metrics.HealthcheckEnabled, cfg.Metrics.PPRofEnabled)

	// setup forwarding service
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Infof("%s setup forwarding end point", method)
	server := server.NewServer(ServiceName, cfg)
	server.Init()

	exitCode := 1
	everlastingGobstopper := make(chan bool)
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-osSignals
		log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Infof("%s signal received.  Shutting down service.", method)
		server.Shutdown()
		log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Infof("%s Deactivating the Everlasting Gobstopper.", method)
		everlastingGobstopper <- true
	}()

	server.Run()

	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Infof("%s returned from server.Run(), waiting for Everlasting Gobstopper", method)

	<-everlastingGobstopper

	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Infof("%s exiting", method)

	os.Exit(exitCode)
}
