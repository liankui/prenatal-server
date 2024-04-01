// Code generated by chaosmojo. DO NOT EDIT.
// Rerunning chaosmojo will overwrite this file.
// Version: 0.1.0
// Version Date:

package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"strings"
	"time"

	// 3d Party
	nclient "github.com/chaos-io/chaos/gokit/client"
	"github.com/chaos-io/chaos/gokit/kit"
	"github.com/chaos-io/chaos/gokit/metrics"
	"github.com/chaos-io/chaos/gokit/sd"
	nserver "github.com/chaos-io/chaos/gokit/server"
	"github.com/chaos-io/chaos/gokit/tracing"
	"github.com/chaos-io/chaos/gokit/utils/network"
	"github.com/chaos-io/chaos/logs"

	"github.com/etherlabsio/healthcheck/v2"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/gorilla/mux"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	_ "go.uber.org/automaxprocs"
	"google.golang.org/grpc"

	// This Service api
	pb "github.com/liankui/prenatal/go/pkg/prenatal/v1"
	"github.com/liankui/prenatal/service-go/pkg/quiz-service/handlers"
	"github.com/liankui/prenatal/service-go/pkg/quiz-service/svc"
)

var _ nclient.Config

func NewEndpoints(options map[string]interface{}) svc.Endpoints {
	// Business domain.
	var service pb.QuizServer
	{
		service = handlers.NewService()
		// Wrap Service with middlewares. See handlers/middlewares.go
		service = handlers.WrapService(service, options)
	}

	// Endpoint domain.
	var (
		createQuestionEndpoint = svc.MakeCreateQuestionEndpoint(service)
		getQuestionEndpoint    = svc.MakeGetQuestionEndpoint(service)
		updateQuestionEndpoint = svc.MakeUpdateQuestionEndpoint(service)
		deleteQuestionEndpoint = svc.MakeDeleteQuestionEndpoint(service)
		createAnswerEndpoint   = svc.MakeCreateAnswerEndpoint(service)
	)

	endpoints := svc.Endpoints{
		CreateQuestionEndpoint: createQuestionEndpoint,
		GetQuestionEndpoint:    getQuestionEndpoint,
		UpdateQuestionEndpoint: updateQuestionEndpoint,
		DeleteQuestionEndpoint: deleteQuestionEndpoint,
		CreateAnswerEndpoint:   createAnswerEndpoint,
	}

	// Wrap selected Endpoints with middlewares. See handlers/middlewares.go
	endpoints = handlers.WrapEndpoints(endpoints, options)

	return endpoints
}

func RegisterService(cfg nserver.Config, r *mux.Router, s *grpc.Server) svc.Endpoints {
	const FullServiceName = "prenatal.Quiz"

	// tracing init
	tracer, c := tracing.New(FullServiceName)
	if c != nil {
		defer c.Close()
	}

	// Create a single logger, which we'll use and give to other components.
	logger := kit.Logger()

	options := map[string]interface{}{
		"tracer": tracer,
		"logger": logger,
	}

	metricsConfig := metrics.NewConfig("metrics")
	if metricsConfig.Enabled() {
		fieldKeys := []string{"method", "access_key", "error"}
		count := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: metricsConfig.Department,
			Subsystem: metricsConfig.Project,
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys)

		latency := kitprometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Namespace: metricsConfig.Department,
			Subsystem: metricsConfig.Project,
			Name:      "request_latency_seconds",
			Help:      "Total duration of requests in seconds.",
		}, fieldKeys)

		options["count"] = count
		options["latency"] = latency
	}

	sdConfig := sd.NewConfig("sd")
	sdClient := sd.New(sdConfig, logger)

	if sdClient != nil {
		url := "etcd://" + network.GetHost() + ":" + getGrpcPort(cfg.GrpcAddr)
		err := sdClient.Register(url, FullServiceName, []string{})
		if err != nil {
			panic(err)
		}
		defer sdClient.Deregister()
	}

	// required service clients ...
	// xxClient := xx_client.New(nclient.NewConfig("xx"), sdClient.Instancer(FullServiceName), tracer, logger)
	// defer xxClient.Close()

	endpoints := NewEndpoints(options)

	svc.RegisterHttpHandler(r, endpoints, tracer, logger)
	pb.RegisterQuizServer(s, svc.MakeGRPCServer(endpoints, tracer, logger))

	return endpoints
}

// Run starts a new http server, gRPC server, and a debug server with the
// passed config and logger
func Run(cfg nserver.Config) {
	// Mechanical domain.
	errc := make(chan error)

	// Interrupt handler.
	go handlers.InterruptHandler(errc)

	// Debug listener.
	go func() {
		logs.Infow("begin debug server", "transport", "debug", "address", cfg.DebugAddr)

		m := http.NewServeMux()
		m.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		m.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		m.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		m.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		m.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

		m.Handle("/metrics", promhttp.Handler())

		m.Handle("/health", healthcheck.Handler(
			// WithTimeout allows you to set a max overall timeout.
			healthcheck.WithTimeout(5*time.Second),
			healthcheck.WithChecker("alive", healthcheck.CheckerFunc(func(ctx context.Context) error {
				conn, err := net.DialTimeout("tcp", cfg.HttpAddr, time.Second)
				if err != nil {
					return err
				}
				return conn.Close()
			})),
		))

		errc <- http.ListenAndServe(cfg.DebugAddr, m)
	}()

	s := grpc.NewServer(grpc.UnaryInterceptor(unaryServerFilter))
	r := mux.NewRouter()
	endpoints := RegisterService(cfg, r, s)

	// HTTP transport.
	go func() {
		logs.Infow("begin http server", "transport", "HTTP", "address", cfg.HttpAddr)
		h := cors.AllowAll().Handler(r)
		errc <- http.ListenAndServe(cfg.HttpAddr, h)
	}()

	// gRPC transport.
	go func() {
		logs.Infow("begin grpc server", "transport", "gRPC", "address", cfg.GrpcAddr)
		ln, err := net.Listen("tcp", cfg.GrpcAddr)
		if err != nil {
			errc <- err
			return
		}
		errc <- s.Serve(ln)
	}()

	// if watchObj, err := config.WatchFunc(level.ChangeLogLevel, level.LevelPath); err == nil {
	//    defer func() { _ = watchObj.Close() }()
	// } else {
	//    panic(err.Error())
	// }
	_ = endpoints

	// Run!
	logs.Info("prenatal.QuizServer", " started.")
	logs.Info("prenatal.QuizServer", <-errc)

	logs.Info("prenatal.QuizServer", " closed.")
}

func getGrpcPort(addr string) string {
	host := strings.Split(addr, ":")
	if len(host) < 2 {
		panic("host name is invalid (" + addr + ")")
	}
	return host[1]
}

func unaryServerFilter(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	// if err := middleware.Validate(req); err != nil {
	//	logs.Errorf("validate request failed, err: %s", err)
	//	return nil, core.NewError(http.StatusBadRequest, err.Error())
	// }

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	resp, err = handler(ctx, req)
	if err != nil {
		return resp, err
	}

	// var validatorCfg middleware.ValidatorConfig
	// _ = config.ScanKey("validator", &validatorCfg)
	// if !validatorCfg.CheckResponse {
	//	return
	// }
	// if err = middleware.Validate(resp); err != nil {
	//	logs.Errorf("validate response failed, err: %s", err)
	//	return nil, core.NewError(http.StatusInternalServerError, err.Error())
	// }
	return
}
