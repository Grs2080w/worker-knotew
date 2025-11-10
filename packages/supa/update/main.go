package update

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

type Jobs struct {
	Id string `json:"id"`
	Provider_id string `json:"provider_id"`
	Provider string `json:"provider"`
	Payload string `json:"payload"`
	Status string `json:"status"`
	Attempts int `json:"attempts"`
	Last_error string `json:"last_error"`
	Updated_at string `json:"updated_at"`
}

type JobUpdate struct {
    Status     *string `json:"status,omitempty"`
    Attempts   *int    `json:"attempts,omitempty"`
    Last_error *string `json:"last_error,omitempty"`
    Updated_at *string `json:"updated_at,omitempty"`
}

func UpdateJobStatus(ctx context.Context, client *supabase.Client, id string, status string) error {

	godotenv.Load()

	if err := ctx.Err(); err != nil {
        return err
    }

	if status != "done" && status != "progress" {
		return errors.New("status must be done or progress")
	}

	if id == "" {
		return errors.New("id is required")
	}

	update := time.Now().UTC().Format(time.RFC3339)

	_, _, err := client.From(os.Getenv("SUPABASE_TABLE_JOBS")).Update(JobUpdate{
		Status: &status,
		Updated_at: &update,
	}, "", "exact").Eq("id", id).Execute()

	if err != nil {
		return err
	}

	return nil

}

func UpdateJobFailed(ctx context.Context, client *supabase.Client, id string, errr string, attempts int) error {

	godotenv.Load()

	if err := ctx.Err(); err != nil {
        return err
    }

	if id == "" {
		return errors.New("id is required")
	}

	status := "pending"
	totalAttempts := attempts + 1
	update := time.Now().UTC().Format(time.RFC3339)

	_, _, err := client.From(os.Getenv("SUPABASE_TABLE_JOBS")).Update(JobUpdate{
		Status: &status,
		Last_error: &errr,
		Attempts: &totalAttempts,
		Updated_at: &update,
	}, "", "exact").Eq("id", id).Execute()

	if err != nil {
		return err
	}

	return nil

}

func UpdateJobError(ctx context.Context, client *supabase.Client, id string, errr string) error {

	godotenv.Load()

	if err := ctx.Err(); err != nil {
        return err
    }

	if id == "" {
		return errors.New("id is required")
	}

	status := "failed"
	update := time.Now().UTC().Format(time.RFC3339)

	_, _, err := client.From(os.Getenv("SUPABASE_TABLE_JOBS")).Update(JobUpdate{
		Status: &status,
		Last_error: &errr,
		Updated_at: &update,
	}, "", "exact").Eq("id", id).Execute()

	if err != nil {
		return err
	}

	return nil

}

func DeleteJobDone(ctx context.Context, client *supabase.Client, id string) error {

	godotenv.Load()

	if err := ctx.Err(); err != nil {
		return err
	}

	if id == "" {
		return errors.New("id is required")
	}

	_, _, err := client.From(os.Getenv("SUPABASE_TABLE_JOBS")).Delete("", "").Eq("id", id).Execute()

	if err != nil {
		return err
	}

	return nil

}
