package lists

import (
	"context"
	"fmt"
	"strconv"

	"github.com/elastic/beats/v7/metricbeat/mb"
	"github.com/elastic/beats/v7/metricbeat/module/pgbouncer"
)

// init registers the MetricSet with the central registry.
func init() {
	mb.Registry.MustAddMetricSet("pgbouncer", "lists", New,
		mb.WithHostParser(pgbouncer.ParseURL),
		mb.DefaultMetricSet(),
	)
}

// MetricSet type defines all fields of the MetricSet
type MetricSet struct {
	*pgbouncer.MetricSet
}

// New creates a new instance of the MetricSet.
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	ms, err := pgbouncer.NewMetricSet(base)
	if err != nil {
		return nil, err
	}
	return &MetricSet{MetricSet: ms}, nil
}

// Fetch methods implements the data gathering and data conversion to the right format
// It publishes the event which is then forwarded to the output. In case of an error, an error is reported.
func (m *MetricSet) Fetch(reporter mb.ReporterV2) error {
	ctx := context.Background()
	results, err := m.QueryStats(ctx, "SHOW LISTS;")
	fmt.Printf("\nPOOLS PRINTING results: %v\n\n", results)
	if err != nil {
		return fmt.Errorf("error in QueryStats: %w", err)
	}
	resultMap := make(map[string]interface{})
	for _, s := range results {
		key := s["list"].(string)
		// Convert value to int if necessary
		var value int
		switch v := s["items"].(type) {
		case string:
			value, err = strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("error converting value to int: %w", err)
			}
		case int:
			value = v
		default:
			return fmt.Errorf("unexpected type for value: %T", v)
		}

		// Assign the integer value
		resultMap[key] = value
	}
	event, err := MapResult(resultMap)
	if err != nil {
		return fmt.Errorf("error mapping result: %w", err)
	}
	reporter.Event(mb.Event{
		MetricSetFields: event,
	})
	return nil
}
