package main

import (
	"context"
	"flag"
	"os"
	"strings"

	"PIM_Server/config"
	"PIM_Server/log"
	"PIM_Server/plugins"
	"PIM_Server/service"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
)

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

func initTracer() func(context.Context) error {

	var secureOption otlptracegrpc.Option

	if strings.ToLower(insecure) == "false" || insecure == "0" || strings.ToLower(insecure) == "f" {
		secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)

	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Fatalf("Could not set resources: %v", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}

func main() {
	confFile := flag.String("f", "../etc/conf.yaml", "配置文件路径")
	flag.Parse()
	config.Init(*confFile)
	gin.SetMode(config.AppConfig().ServerInfo.Mode)

	log.InitLogger(config.AppConfig().LogInfo.Path,
		config.AppConfig().LogInfo.MaxSize,
		config.AppConfig().LogInfo.MaxBackUps,
		config.AppConfig().LogInfo.MaxAge,
		config.AppConfig().LogInfo.Level,
		config.AppConfig().LogInfo.CallerSkip)

	log.Infof("Start im_server, listen: %s", config.AppConfig().ServerInfo.Listen)

	cleanup := initTracer()
	defer cleanup(context.Background())

	// 初始化gin插件
	r := plugins.Init(serviceName)
	service.Init()

	if err := r.Run(config.AppConfig().ServerInfo.Listen); err != nil {
		log.Fatalf("startup service failed, err:%v", err)
	}
}
