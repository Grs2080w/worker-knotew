package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Grs2080w/worker-knoteq/packages/memory"
	"github.com/Grs2080w/worker-knoteq/packages/prometheus"
	"github.com/Grs2080w/worker-knoteq/packages/redis"
	"github.com/Grs2080w/worker-knoteq/packages/supa"
	"github.com/Grs2080w/worker-knoteq/packages/supa/get"
	"github.com/Grs2080w/worker-knoteq/packages/worker"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func jobFetcher(ctx context.Context, supaPublic *supa.SupabasePublic, out chan<- get.Job) {
    for {
        jobs, err := supaPublic.GetJobs(ctx)
        if err != nil {
			prometheus.IdleIterations.Inc()
            time.Sleep(10 * time.Second)
            continue
        }

		for _, job := range jobs {
			
			if job.Id == "" {
				continue
			}
			
			out <- job
			
		}
		
		time.Sleep(25*time.Second)
		
    }
}



func main() {
	
	ctx := context.Background()

	prometheus.InitPrometheus()
	go memory.MonitorResources()

	supaPublic, err := supa.NewPublic()
	if err != nil { panic(err) }
	
	supaAuth, err := supa.NewAuth()
	if err != nil { panic(err) }

	redisClient, err := redis.New()
	if err != nil { panic(err) }
	
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Worker Metrics exposed.")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	numWorkers := 3

	jobsChan := make(chan get.Job, 100)
	go jobFetcher(ctx, supaPublic, jobsChan)
	
	log.Println("Worker started")

    for i := 0; i < numWorkers; i++ {
       go worker.Worker(supaPublic, supaAuth, redisClient, jobsChan)
    }

	select {}

}