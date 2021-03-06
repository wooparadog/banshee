// Copyright 2015 Eleme Inc. All rights reserved.

package metricdb

import (
	"fmt"
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util"
	"strconv"
)

// Stamp
const (
	// Futher timestamps will be stored as the diff to this horizon for storage
	// cost reason.
	horizon uint32 = 1450322633
	// Timestamps will be converted to 36-hex string format before they are put
	// to db, also for storage cost reason.
	convBase = 36
	// A 36-hex string format timestamp with this length is enough to use for
	// later 90 years.
	stampLen = 7
)

// Horizon returns the timestamp horizon
func Horizon() uint32 {
	return horizon
}

// encodeKey encodes db key from metric.
func encodeKey(m *models.Metric) []byte {
	stamp := m.Stamp
	// As horizon if stamp is too small.
	if m.Stamp < horizon {
		stamp = horizon
	}
	// Key format is Name+Stamp.
	t := stamp - horizon
	v := strconv.FormatUint(uint64(t), convBase)
	s := fmt.Sprintf("%s%0*s", m.Name, stampLen, v)
	return []byte(s)
}

// encodeValue encodes db value from metric.
func encodeValue(m *models.Metric) []byte {
	// Value format is Value:Score:Average.
	value := util.ToFixed(m.Value, 5)
	score := util.ToFixed(m.Score, 5)
	average := util.ToFixed(m.Average, 5)
	s := fmt.Sprintf("%s:%s:%s", value, score, average)
	return []byte(s)
}

// decodeKey decodes db key into metric, this will fill metric name and metric
// stamp.
func decodeKey(key []byte, m *models.Metric) error {
	s := string(key)
	if len(s) <= stampLen {
		return ErrCorrupted
	}
	// First substring is Name.
	idx := len(s) - stampLen
	m.Name = s[:idx]
	// Last substring is Stamp.
	str := s[idx:]
	n, err := strconv.ParseUint(str, convBase, 32)
	if err != nil {
		return err
	}
	m.Stamp = horizon + uint32(n)
	return nil
}

// decodeValue decodes db value into metric, this will fill metric value,
// average and stddev.
func decodeValue(value []byte, m *models.Metric) error {
	n, err := fmt.Sscanf(string(value), "%f:%f:%f", &m.Value, &m.Score, &m.Average)
	if err != nil {
		return err
	}
	if n != 3 {
		return ErrCorrupted
	}
	return nil
}
