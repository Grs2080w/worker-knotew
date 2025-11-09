package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
    JobExecutionDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
        Name: "worker_job_execution_seconds",
        Help: "Tempo de execução de cada job",
        //Buckets: prometheus.LinearBuckets(0.1, 0.2, 15),
        Buckets: prometheus.LinearBuckets(2, 3, 15),
    })

    JobsProcessed = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "worker_jobs_processed_total",
        Help: "Total de jobs processados",
    })

    JobsError = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "worker_jobs_error_total",
        Help: "Total de jobs que deram erro",
    })
    
    JobsFailed = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "worker_jobs_failed_total",
        Help: "Total de jobs que deram erro",
    })

	IdleIterations = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "worker_idle_iterations_total",
		Help: "Total de idle iterations",
	})

	MemoryUsageBytes = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "worker_memory_usage_bytes",
            Help: "Quantidade de memória usada pelo processo.",
        },
    )

	Goroutines = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "worker_goroutines_total",
            Help: "Número total de goroutines.",
        },
    )

    GCCount = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "worker_gc_cycles_total",
            Help: "Total de ciclos de Garbage Collector.",
        },
    )
)

func InitPrometheus() {
    prometheus.MustRegister(JobExecutionDuration, JobsProcessed, JobsError, JobsFailed, IdleIterations, MemoryUsageBytes)
}
