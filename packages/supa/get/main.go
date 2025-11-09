package get

import (
	"context"

	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type Job struct {
	Id string `json:"id"`
	Provider_id string `json:"provider_id"`
	User_id string `json:"user_id"`
	Provider string `json:"provider"`
	Payload string `json:"payload"`
	Status string `json:"status"`
	Attempts int `json:"attempts"`
	Last_error string `json:"last_error"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
	Owner string `json:"owner"`
	Repo string `json:"repo"`
	NameFile string `json:"nameFile"`
}


func GetJobs(ctx context.Context, client *supabase.Client) ([]Job, error) {

	if err := ctx.Err(); err != nil {
        return []Job{}, err
    }

	var res []Job	

	_, err := client.From("knoteq_sync_jobs").
	Select("*", "exact", false).
	Eq("status", "pending").
	Order("updated_at", &postgrest.OrderOpts{Ascending: true}).
	Order("status", &postgrest.OrderOpts{Ascending: true}).
	Limit(100, "").
	ExecuteTo(&res)

	if err != nil {
		return nil, err
	}

	return res, nil

}

func GetJob(ctx context.Context, client *supabase.Client) (Job, error) {

	if err := ctx.Err(); err != nil {
        return Job{}, err
    }

	var res []Job	

	_, err := client.From("knoteq_sync_jobs").
		Select("*", "exact", false).
		Eq("status", "pending").
		Order("updated_at", &postgrest.OrderOpts{Ascending: true}).
		Order("status", &postgrest.OrderOpts{Ascending: true}).
		Limit(1, "").
		ExecuteTo(&res)

	if err != nil {
		return Job{}, err
	}

	if len(res) == 0 {
    	return Job{}, nil
	}

	return res[0], nil

}