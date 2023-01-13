// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build darwin linux openbsd
// +build !nomeminfo

package collector

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector("cputemperature", defaultEnabled, NewCPUTemperatureCollector)
}

type CPUTemperatureCollector struct {
	entries *prometheus.Desc
	logger  log.Logger
}

func NewCPUTemperatureCollector(logger log.Logger) (Collector, error) {
	return &CPUTemperatureCollector{
		entries: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "cpu_temperature"),
			"The cpu Temperature,information from /sys/class/hwmon/hwmon0/temp1_input",
			[]string{"virt"}, nil),
		logger: logger,
	}, nil
}

func (c *CPUTemperatureCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.entries
}

// Update implements Collector
func (c *CPUTemperatureCollector) Update(ch chan<- prometheus.Metric) error {
	entries, err := getCPUTemperature()
	if err != nil {
		level.Error(c.logger).Log("Get CPU temperature failed, %v", err)
	}
	fmt.Println(entries)
	ch <- prometheus.MustNewConstMetric(
		c.entries, prometheus.GaugeValue, entries, *instanceUUID)
	return nil
}

func getCPUTemperature() (float64, error) {
	content, err := ioutil.ReadFile("/sys/class/hwmon/hwmon0/temp1_input")
	if err != nil {
		return 0, err
	}
	result := strings.TrimSpace(string(content))
	entries, _ := strconv.ParseFloat(result, 64)

	return entries / 1000, nil
}
