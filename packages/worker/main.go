package worker

import (
	"context"
	"log"
	"time"

	"github.com/Grs2080w/worker-knoteq/packages/github"
	"github.com/Grs2080w/worker-knoteq/packages/google"
	"github.com/Grs2080w/worker-knoteq/packages/prometheus"
	"github.com/Grs2080w/worker-knoteq/packages/redis"
	"github.com/Grs2080w/worker-knoteq/packages/supa"
	"github.com/Grs2080w/worker-knoteq/packages/supa/get"
	"github.com/Grs2080w/worker-knoteq/packages/supa/token"
)

func Worker(supaPublic *supa.SupabasePublic, supaAuth *supa.SupabaseAuth, redisClient *redis.Redis , in <-chan get.Job) {

	for job := range in {
			
		// try 3 times 
		for range 3 {
			
			ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
			defer cancel()

			start := time.Now()
			
			tok, err := redisClient.GetToken(ctx, job.User_id+job.Provider, func() ( token.Token, error) {
				return supaAuth.GetTokenUser(ctx, job.Provider, job.User_id)
			})
			
			if err != nil {
				redisClient.InvalidateKey(ctx, job.User_id+job.Provider)
				supaPublic.UpdateJobFailedError(ctx, err, job)
				job.Attempts++
				continue
			}
			
			err = supaPublic.UpdateJobStatus(ctx, job.Id, "progress")
			if err != nil {
				supaPublic.UpdateJobFailedError(ctx, err, job)
				job.Attempts++
				continue
			}

			switch job.Provider {
				
				case "google":

					drive, err := google.New(ctx, tok)
					if err != nil {
						supaPublic.UpdateJobFailedError(ctx, err, job)
						job.Attempts++
						continue
					}
					
					err = drive.UpdateFile(ctx, job.Provider_id, []byte(job.Payload))
					if err != nil {
						supaPublic.UpdateJobFailedError(ctx, err, job)
						job.Attempts++
						continue
					}


				case "github":

					git := github.New(tok.Access_token, job.Owner, job.Repo)
			
					err = git.UpdateFile(ctx, job.NameFile, job.Payload)
					if err != nil {
						supaPublic.UpdateJobFailedError(ctx, err, job)
						job.Attempts++
						continue
					}

			}
		
			err = supaPublic.DeleteJobDone(ctx, job.Id)
			if err != nil {
				log.Println(err)
				job.Attempts++
				continue
			}

			// Metric for duration job
			elapsed := time.Since(start).Seconds()

			prometheus.JobExecutionDuration.Observe(elapsed)
			prometheus.JobsProcessed.Inc()

			break

		}
		
	}
	
}