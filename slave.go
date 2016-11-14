package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

func newSlaveCollector(httpClient *httpClient) prometheus.Collector {
	metrics := map[prometheus.Collector]func(metricMap, prometheus.Collector) error{
		// CPU/Disk/Mem resources in free/used
		gauge("slave", "cpus", "Current CPU resources in cluster.", "type"): func(m metricMap, c prometheus.Collector) error {
			total, ok := m["slave/cpus_total"]
			used, ok := m["slave/cpus_used"]
			if !ok {
				return notFoundInMap
			}
			c.(*prometheus.GaugeVec).WithLabelValues("free").Add(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used").Add(used)
			return nil
		},
		gauge("slave", "cpus_revocable", "Current revocable CPU resources in cluster.", "type"): func(m metricMap, c prometheus.Collector) error {
			total, ok := m["slave/cpus_revocable_total"]
			used, ok := m["slave/cpus_revocable_used"]
			if !ok {
				return notFoundInMap
			}
			c.(*prometheus.GaugeVec).WithLabelValues("free").Add(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used").Add(used)
			return nil
		},
		gauge("slave", "mem", "Current memory resources in cluster.", "type"): func(m metricMap, c prometheus.Collector) error {
			total, ok := m["slave/mem_total"]
			used, ok := m["slave/mem_used"]
			if !ok {
				return notFoundInMap
			}
			c.(*prometheus.GaugeVec).WithLabelValues("free").Add(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used").Add(used)
			return nil
		},
		gauge("slave", "mem_revocable", "Current revocable memory resources in cluster.", "type"): func(m metricMap, c prometheus.Collector) error {
			total, ok := m["slave/mem_revocable_total"]
			used, ok := m["slave/mem_revocable_used"]
			if !ok {
				return notFoundInMap
			}
			c.(*prometheus.GaugeVec).WithLabelValues("free").Add(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used").Add(used)
			return nil
		},
		gauge("slave", "disk", "Current disk resources in cluster.", "type"): func(m metricMap, c prometheus.Collector) error {
			total, ok := m["slave/disk_total"]
			used, ok := m["slave/disk_used"]
			if !ok {
				return notFoundInMap
			}
			c.(*prometheus.GaugeVec).WithLabelValues("free").Add(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used").Add(used)
			return nil
		},
		gauge("slave", "disk_revocable", "Current disk resources in cluster.", "type"): func(m metricMap, c prometheus.Collector) error {
			total, ok := m["slave/disk_revocable_total"]
			used, ok := m["slave/disk_revocable_used"]
			if !ok {
				return notFoundInMap
			}
			c.(*prometheus.GaugeVec).WithLabelValues("free").Add(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used").Add(used)
			return nil
		},

		// Slave stats about uptime and connectivity
		prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "slave",
			Name:      "registered",
			Help:      "1 if slave is registered with master, 0 if not.",
		}): func(m metricMap, c prometheus.Collector) error {
			registered, ok := m["slave/registered"]
			if !ok {
				return notFoundInMap
			}
			c.(prometheus.Gauge).Add(registered)
			return nil
		},
		prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "slave",
			Name:      "uptime_seconds",
			Help:      "Number of seconds the master process is running.",
		}): func(m metricMap, c prometheus.Collector) error {
			uptime, ok := m["slave/uptime_secs"]
			if !ok {
				return notFoundInMap
			}
			c.(prometheus.Gauge).Add(uptime)
			return nil
		},

		// Slave stats about frameworks and executors
		gauge("slave", "executor_state", "Current number of executors by state.", "state"): func(m metricMap, c prometheus.Collector) error {
			registering, ok := m["slave/executors_registering"]
			running, ok := m["slave/executors_running"]
			terminating, ok := m["slave/executors_terminating"]
			if !ok {
				return notFoundInMap
			}
			c.(*prometheus.GaugeVec).WithLabelValues("registering").Add(registering)
			c.(*prometheus.GaugeVec).WithLabelValues("running").Add(running)
			c.(*prometheus.GaugeVec).WithLabelValues("terminating").Add(terminating)
			return nil
		},
		prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "slave",
			Name:      "frameworks_active",
			Help:      "Current number of active frameworks",
		}): func(m metricMap, c prometheus.Collector) error {
			active, ok := m["slave/frameworks_active"]
			if !ok {
				return notFoundInMap
			}
			c.(prometheus.Gauge).Set(active)
			return nil
		},
		prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "mesos",
			Subsystem: "slave",
			Name:      "executors_terminated",
			Help:      "Total number of executor terminations.",
		}): func(m metricMap, c prometheus.Collector) error {
			terminated, ok := m["slave/executors_terminated"]
			if !ok {
				return notFoundInMap
			}
			c.(prometheus.Counter).Set(terminated)
			return nil
		},
		prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "mesos",
			Subsystem: "slave",
			Name:      "executors_preempted",
			Help:      "Total number of executor preemptions.",
		}): func(m metricMap, c prometheus.Collector) error {
			preempted, ok := m["slave/executors_preempted"]
			if !ok {
				return notFoundInMap
			}
			c.(prometheus.Counter).Set(preempted)
			return nil
		},

		// Slave stats about tasks
		counter("slave", "task_states_exit_total", "Total number of tasks processed by exit state.", "state"): func(m metricMap, c prometheus.Collector) error {
			errored, ok := m["slave/tasks_error"]
			failed, ok := m["slave/tasks_failed"]
			finished, ok := m["slave/tasks_finished"]
			killed, ok := m["slave/tasks_killed"]
			lost, ok := m["slave/tasks_lost"]
			if !ok {
				return notFoundInMap
			}
			c.(*prometheus.CounterVec).WithLabelValues("errored").Add(errored)
			c.(*prometheus.CounterVec).WithLabelValues("failed").Add(failed)
			c.(*prometheus.CounterVec).WithLabelValues("finished").Add(finished)
			c.(*prometheus.CounterVec).WithLabelValues("killed").Add(killed)
			c.(*prometheus.CounterVec).WithLabelValues("lost").Add(lost)
			return nil
		},
		counter("slave", "task_states_current", "Current number of tasks by state.", "state"): func(m metricMap, c prometheus.Collector) error {
			running, ok := m["slave/tasks_running"]
			staging, ok := m["slave/tasks_staging"]
			starting, ok := m["slave/tasks_starting"]
			if !ok {
				return notFoundInMap
			}
			c.(*prometheus.CounterVec).WithLabelValues("running").Add(running)
			c.(*prometheus.CounterVec).WithLabelValues("staging").Add(staging)
			c.(*prometheus.CounterVec).WithLabelValues("starting").Add(starting)
			return nil
		},

		// Slave stats about messages
		counter("slave", "messages_outcomes_total",
			"Total number of messages by outcome of operation",
			"type", "outcome"): func(m metricMap, c prometheus.Collector) error {

			frameworkMessagesValid, ok := m["slave/valid_framework_messages"]
			frameworkMessagesInvalid, ok := m["slave/invalid_framework_messages"]
			statusUpdateValid, ok := m["slave/valid_status_updates"]
			statusUpdateInvalid, ok := m["slave/invalid_status_updates"]

			if !ok {
				return notFoundInMap
			}
			c.(*prometheus.CounterVec).WithLabelValues("framework", "valid").Add(frameworkMessagesValid)
			c.(*prometheus.CounterVec).WithLabelValues("framework", "invalid").Add(frameworkMessagesInvalid)
			c.(*prometheus.CounterVec).WithLabelValues("status", "valid").Add(statusUpdateValid)
			c.(*prometheus.CounterVec).WithLabelValues("status", "invalid").Add(statusUpdateInvalid)

			return nil
		},
	}
	return newMetricCollector(httpClient, metrics)
}
