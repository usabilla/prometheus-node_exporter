// Copyright 2019 Usabilla
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

package collector

import (
	"fmt"

	"github.com/prometheus/procfs"

	"github.com/prometheus/client_golang/prometheus"
)

type swapCollector struct {
	fs       procfs.FS
	Size     *prometheus.Desc
	Used     *prometheus.Desc
	Priority *prometheus.Desc
}

func init() {
	registerCollector("swap", defaultEnabled, NewSwapCollector)
}

func NewSwapCollector() (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %v", err)
	}
	return &swapCollector{
		fs: fs,
		Size: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "swap", "size_bytes"),
			"The size of the swap file.",
			[]string{"filename", "type"},
			nil,
		),
		Used: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "swap", "used_bytes"),
			"The amount of used space in the swap file.",
			[]string{"filename", "type"},
			nil,
		),
		Priority: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "swap", "priority"),
			"The priority of the swap file.",
			[]string{"filename", "type"},
			nil,
		),
	}, nil
}

func (c *swapCollector) Update(ch chan<- prometheus.Metric) error {
	swaps, err := c.fs.Swaps()
	if err != nil {
		return err
	}

	for _, s := range swaps {
		ch <- prometheus.MustNewConstMetric(c.Size, prometheus.GaugeValue, float64(s.Size*1024), s.Filename, s.Type)
		ch <- prometheus.MustNewConstMetric(c.Used, prometheus.GaugeValue, float64(s.Used*1024), s.Filename, s.Type)
		ch <- prometheus.MustNewConstMetric(c.Priority, prometheus.GaugeValue, float64(s.Priority), s.Filename, s.Type)
	}

	return nil
}
