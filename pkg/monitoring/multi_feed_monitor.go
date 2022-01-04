package monitoring

import (
	"context"
	"sync"

	"github.com/smartcontractkit/chainlink-solana/pkg/monitoring/config"
	"github.com/smartcontractkit/chainlink/core/logger"
)

type MultiFeedMonitor interface {
	Start(ctx context.Context, wg *sync.WaitGroup)
}

func NewMultiFeedMonitor(
	solanaConfig config.Solana,
	feeds []config.Feed,

	log logger.Logger,
	transmissionReader, stateReader AccountReader,
	producer Producer,
	metrics Metrics,

	configSetTopic string,
	configSetSimplifiedTopic string,
	transmissionTopic string,

	configSetSchema Schema,
	configSetSimplifiedSchema Schema,
	transmissionSchema Schema,
) MultiFeedMonitor {
	return &multiFeedMonitor{
		solanaConfig,
		feeds,

		log,
		transmissionReader, stateReader,
		producer,
		metrics,

		configSetTopic, configSetSimplifiedTopic, transmissionTopic,
		configSetSchema, configSetSimplifiedSchema, transmissionSchema,
	}
}

type multiFeedMonitor struct {
	solanaConfig config.Solana
	feeds        []config.Feed

	log                logger.Logger
	transmissionReader AccountReader
	stateReader        AccountReader
	producer           Producer
	metrics            Metrics

	configSetTopic           string
	configSetSimplifiedTopic string
	transmissionTopic        string

	configSetSchema           Schema
	configSetSimplifiedSchema Schema
	transmissionSchema        Schema
}

const bufferCapacity = 100

// Start should be executed as a goroutine.
func (m *multiFeedMonitor) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(len(m.feeds))
	for _, feedConfig := range m.feeds {
		go func(feedConfig config.Feed) {
			defer wg.Done()

			feedLogger := m.log.With(
				"feed", feedConfig.FeedName,
				"network", m.solanaConfig.NetworkName,
			)

			transmissionPoller := NewPoller(
				feedLogger.With("component", "transmissions-poller", "address", feedConfig.TransmissionsAccount.String()),
				feedConfig.TransmissionsAccount,
				m.transmissionReader,
				m.solanaConfig.PollInterval,
				m.solanaConfig.ReadTimeout,
				bufferCapacity,
			)
			statePoller := NewPoller(
				feedLogger.With("component", "state-poller", "address", feedConfig.StateAccount.String()),
				feedConfig.StateAccount,
				m.stateReader,
				m.solanaConfig.PollInterval,
				m.solanaConfig.ReadTimeout,
				bufferCapacity,
			)

			wg.Add(2)
			go func() {
				defer wg.Done()
				transmissionPoller.Start(ctx)
			}()
			go func() {
				defer wg.Done()
				statePoller.Start(ctx)
			}()

			exporters := []Exporter{
				NewPrometheusExporter(
					m.solanaConfig,
					feedConfig,
					feedLogger.With("component", "prometheus-exporter"),
					m.metrics,
				),
				NewKafkaExporter(
					m.solanaConfig,
					feedConfig,
					feedLogger.With("component", "kafka-exporter"),
					m.producer,

					m.configSetSchema,
					m.configSetSimplifiedSchema,
					m.transmissionSchema,

					m.configSetTopic,
					m.configSetSimplifiedTopic,
					m.transmissionTopic,
				),
			}

			feedMonitor := NewFeedMonitor(
				feedLogger,
				transmissionPoller, statePoller,
				exporters,
			)
			feedMonitor.Start(ctx, wg)
		}(feedConfig)
	}
}
