package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/Grs2080w/worker-knoteq/packages/github"
	"github.com/Grs2080w/worker-knoteq/packages/google"
	"github.com/Grs2080w/worker-knoteq/packages/prometheus"
	"github.com/Grs2080w/worker-knoteq/packages/supa"
	"github.com/Grs2080w/worker-knoteq/packages/supa/get"
)

func Worker(supaPublic *supa.SupabasePublic, supaAuth *supa.SupabaseAuth) (*get.Job, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	job, err := supaPublic.GetJob(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	// if no job pending, return
	if job.Id == "" {
		return &job, nil
	}

	fmt.Println(job.Id)

	err = supaPublic.UpdateJobStatus(ctx, job.Id, "progress")
	if err != nil {
		supaPublic.UpdateJobFailed(ctx, job.Id, err.Error(), job.Attempts)
		return &get.Job{}, err
	}
	
	tok, err := supaAuth.GetTokenUser(ctx, job.Provider, job.User_id)
	if err != nil {
		supaPublic.UpdateJobFailed(ctx, job.Id, err.Error(), job.Attempts)
		return &get.Job{}, err
	}


	switch job.Provider {
		
		case "google":

			drive, err := google.New(ctx, tok)
			if err != nil {
				supaPublic.UpdateJobFailed(ctx, job.Id, err.Error(), job.Attempts)
				return &get.Job{}, err
			}
			
			err = drive.UpdateFile(ctx, job.Provider_id, []byte(job.Payload))
			if err != nil {
				supaPublic.UpdateJobFailed(ctx, job.Id, err.Error(), job.Attempts)
				return &get.Job{}, err
			}


		case "github":

			git := github.New(tok.Access_token, job.Owner, job.Repo)
	
			err = git.UpdateFile(ctx, job.NameFile, job.Payload)
			if err != nil {
				supaPublic.UpdateJobFailed(ctx, job.Id, err.Error(), job.Attempts)
				return &get.Job{}, err	
			}

	}
	
	
	err = supaPublic.UpdateJobStatus(ctx, job.Id, "done")
	if err != nil {
		supaPublic.UpdateJobFailed(ctx, job.Id, err.Error(), job.Attempts)
		return &get.Job{}, err
	}

	prometheus.JobsProcessed.Inc()

	
	return &job, nil
}