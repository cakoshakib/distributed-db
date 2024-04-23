package metricsink

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/armon/go-metrics"
	log "github.com/cakoshakib/distributed-db/commons"
	"go.uber.org/zap"
)

// Create an in-memory sink with a 1-minute time window, data retained for 1 hour

type MetricHandler struct {
	Sink   *metrics.InmemSink
	Writer *csv.Writer
	File   *os.File
}

func NewMetricHandler(path string) (*MetricHandler, error) {
	handler := &MetricHandler{}

	sink, err := NewSink()
	if err != nil {
		return nil, err
	}
	handler.Sink = sink

	csv, writer, err := NewCSV(path)
	if err != nil {
		return nil, fmt.Errorf("metric csv file could not be created, err: %s", err)
	}
	handler.Writer = writer
	handler.File = csv

	return handler, nil
}

func NewSink() (*metrics.InmemSink, error) {
	sink := metrics.NewInmemSink(10*time.Second, time.Minute)
	config := metrics.DefaultConfig("raft")
	config.EnableHostname = false

	// Initialize the metrics system with the in-memory sink
	_, err := metrics.NewGlobal(config, sink)
	if err != nil {
		return nil, err
	}

	return sink, nil
}

func NewCSV(path string) (*os.File, *csv.Writer, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, nil, err
	}

	writer := csv.NewWriter(file)
	headers := []string{"Timestamp", "Metric", "Mean", "StdDev"}
	if err := writer.Write(headers); err != nil {
		file.Close()
		return nil, nil, err
	}
	writer.Flush()

	return file, writer, nil
}

func (m *MetricHandler) LogMetricsToCSV(ctx context.Context) {
	logger := log.LoggerFromContext(ctx)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping metrics logging due to cancellation.")
			m.Writer.Flush()
			m.File.Close()
			return
		case <-ticker.C:
			data := m.Sink.Data()
			if err := logMetrics(data, m.Writer); err != nil {
				logger.Warn("Error logging metrics, err: ", zap.Error(err))
				continue
			}
		}
	}
}

func logMetrics(data []*metrics.IntervalMetrics, writer *csv.Writer) error {
	for _, interval := range data {
		for key, sample := range interval.Samples {
			agg := sample.AggregateSample
			// Log all metrics for debugging
			//fmt.Printf("Metric key: %s, Mean: %f, StdDev: %f\n", key, agg.Mean(), agg.Stddev())

			if key == "raft.raft.fsm.apply" {
				record := []string{
					time.Now().Format(time.RFC3339),
					key,
					fmt.Sprintf("%f", agg.Mean()),
					fmt.Sprintf("%f", agg.Stddev()),
				}
				if err := writer.Write(record); err != nil {
					return fmt.Errorf("error writing to CSV: %w", err)
				}
				writer.Flush()
			}
		}
	}
	return nil
}

/*
func logMetrics(data []*metrics.IntervalMetrics, writer *csv.Writer) error {
	for _, interval := range data {
		for key, sample := range interval.Samples {
			if key == "raft.fsm.apply" {
				record := []string{
					time.Now().Format(time.RFC3339),
					key,
					fmt.Sprintf("%f", sample.Mean),
					fmt.Sprintf("%f", sample.Stddev),
				}
				if err := writer.Write(record); err != nil {
					return fmt.Errorf("error writing to CSV: %w", err)
				}
				writer.Flush()
				fmt.Println("Wrote to metrics.csv!!!")
			}
		}
	}
	return nil
}
*/
