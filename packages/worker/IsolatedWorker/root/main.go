package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Grs2080w/worker-knoteq/packages/memory"
	"github.com/Grs2080w/worker-knoteq/packages/prometheus"
	"github.com/Grs2080w/worker-knoteq/packages/supa"
	"github.com/Grs2080w/worker-knoteq/packages/worker/IsolatedWorker/worker"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*

	This version is for isolated worker, here just run one worker

*/

func main() {
	
	ctx := context.Background()

	prometheus.InitPrometheus()
	go memory.MonitorResources()

	supaPublic, err := supa.NewPublic()
	if err != nil { panic(err) }
	
	supaAuth, err := supa.NewAuth()
	if err != nil { panic(err) }
	
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Prometheus logs exposed!")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()
	
	log.Println("Worker started")

	idle := 0
	
	for {
		start := time.Now()
	
		job, err := worker.Worker(supaPublic, supaAuth)

		duration := time.Since(start).Seconds()
		
		// if no job pending, continue with idle
		if job == nil || job.Id == "" {
			time.Sleep(time.Duration(idle) * time.Second)
			if idle < 120 {
				idle += 10
			}
			
			prometheus.IdleIterations.Inc()
			continue
		} 
			
		if job.Id != "" {
			prometheus.JobExecutionDuration.Observe(duration)
		}
	
		if err != nil {
			log.Printf("worker error: job=%s err=%v", job.Id, err)
			prometheus.JobsFailed.Inc()
			
			if job.Attempts >= 3 {
				err := supaPublic.UpdateJobError(ctx, job.Id, err.Error())
				if err != nil {
					log.Printf("worker error attempts more than 3: job=%s err=%v", job.Id, err)
				}
			} else {
				err := supaPublic.UpdateJobFailed(ctx, job.Id, err.Error(), job.Attempts)
				if err != nil {
					log.Printf("worker error attempts less than 3: job=%s err=%v", job.Id, err)
				}
			}
		} else {
			prometheus.JobsProcessed.Inc()
		}
	
		time.Sleep(1500 * time.Millisecond)
		idle = 0
	}

	} 