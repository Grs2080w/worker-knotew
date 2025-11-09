package supa

import (
	"context"
	"log"

	"github.com/Grs2080w/worker-knoteq/packages/prometheus"
	"github.com/Grs2080w/worker-knoteq/packages/supa/client"
	"github.com/Grs2080w/worker-knoteq/packages/supa/get"
	"github.com/Grs2080w/worker-knoteq/packages/supa/token"
	"github.com/Grs2080w/worker-knoteq/packages/supa/update"
	"github.com/supabase-community/supabase-go"
)

type SupabasePublic struct {
	Client *supabase.Client
}

type SupabaseAuth struct {
	Client *supabase.Client
}

// New: Create new Supabase client to have access to database
func NewPublic() (*SupabasePublic, error) {
	client, err := client.GetClientPublic()
	
	if err != nil {
		return &SupabasePublic{}, err
	}

	return &SupabasePublic{Client: client}, nil
}

func NewAuth() (*SupabaseAuth, error) {
	client, err := client.GetClientAuth()
	
	if err != nil {
		return &SupabaseAuth{}, err
	}

	return &SupabaseAuth{Client: client}, nil
}

// GetJobs: Get all jobs on database and return it
func (s *SupabasePublic) GetJobs(ctx context.Context) ([]get.Job, error) {
	return get.GetJobs(ctx, s.Client)
}

// GetJob: Get job from database order by updated_at desc
func (s *SupabasePublic) GetJob(ctx context.Context) (get.Job, error) {
	return get.GetJob(ctx, s.Client)
}

// UpdateJobStatus: Update job status on database
func (s *SupabasePublic) UpdateJobStatus(ctx context.Context, id string, status string) error {
	return update.UpdateJobStatus(ctx, s.Client, id, status)
}

// UpdateJobError: Update job status on database when job get error after 3 attempts
func (s *SupabasePublic) UpdateJobError(ctx context.Context, id string, err string) error {
	return update.UpdateJobError(ctx, s.Client, id, err)
}

// UpdateJobFailed: Update job status on database when job failed 
func (s *SupabasePublic) UpdateJobFailed(ctx context.Context, id string, err string, attemptsCurrent int) error {
	return update.UpdateJobFailed(ctx, s.Client, id, err, attemptsCurrent)
}

// UpdateJobFailedError: Update job status on database when job failed, and set error if attempts more than 3
func (s *SupabasePublic) UpdateJobFailedError(ctx context.Context, err error, job get.Job) {
	log.Printf("worker error: job=%s err=%v", job.Id, err)

	if job.Attempts >= 3 {
		err := s.UpdateJobError(ctx, job.Id, err.Error())
		if err != nil {}
		prometheus.JobsError.Inc()
	} else {
		prometheus.JobsFailed.Inc()
		err := s.UpdateJobFailed(ctx, job.Id, err.Error(), job.Attempts)
		if err != nil {}
	}
}

// DeleteJobDone: Delete job from database
func (s *SupabasePublic) DeleteJobDone(ctx context.Context, id string) error {
	return update.DeleteJobDone(ctx, s.Client, id)
}

// GetTokensUser: Get tokens of a user from database, get access token and refresh token
func (s *SupabaseAuth) GetTokenUser(ctx context.Context, provider string, user_id string) (token.Token, error) {
	return token.GetTokenUser(ctx, s.Client, provider, user_id)
}