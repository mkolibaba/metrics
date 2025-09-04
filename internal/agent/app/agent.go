package app

import (
	"context"
	"github.com/mkolibaba/metrics/internal/agent/collector"
	"github.com/mkolibaba/metrics/internal/agent/config"
	"github.com/mkolibaba/metrics/internal/agent/http/client"
	"github.com/mkolibaba/metrics/internal/agent/sender"
	"github.com/mkolibaba/metrics/internal/common/log"
	"github.com/mkolibaba/metrics/internal/common/rsa"
	"go.uber.org/zap"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"
)

// Run инициализирует и запускает агент сбора метрик.
// Функция настраивает конфигурацию, логгер, сборщик метрик и отправитель,
// а затем запускает основной цикл отправки метрик.
// Работает до прерывания.
func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadAgentConfig()
	if err != nil {
		stdlog.Fatalf("error creating config: %v", err)
	}

	logger := log.New()

	c := collector.NewMetricsCollector(cfg.PollInterval, logger)
	chGauges, chCounters := c.StartCollect(ctx)

	encryptor := rsa.NopEncryptor
	if cfg.CryptoKey != "" {
		encryptor, err = rsa.NewEncryptor(cfg.CryptoKey)
		if err != nil {
			logger.Fatalf("error creating rsa encryptor: %v", err)
		}
	}

	serverAPI := client.New(cfg.ServerAddress, cfg.Key, encryptor, logger)
	metricsSender := sender.NewMetricsSender(serverAPI, cfg.ReportInterval, cfg.RateLimit, logger)

	runAgent(ctx, metricsSender, chGauges, chCounters, logger)
}

func runAgent(
	ctx context.Context,
	metricsSender *sender.MetricsSender,
	chGauges <-chan map[string]float64,
	chCounters <-chan map[string]int64,
	logger *zap.SugaredLogger,
) {
	logger.Info("running agent")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		metricsSender.StartSend(ctx, chGauges, chCounters)
	}()

	<-interrupt

	logger.Info("agent stopped")
}
