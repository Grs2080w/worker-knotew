package client

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

func GetClientPublic() (*supabase.Client, error) {

	godotenv.Load()

	SUPABASE_URL := os.Getenv("SUPABASE_URL")
	SUPABASE_PUBLISHABLE_KEY := os.Getenv("SUPABASE_PUBLISHABLE_KEY")
	
	client, err := supabase.NewClient(SUPABASE_URL, SUPABASE_PUBLISHABLE_KEY, &supabase.ClientOptions{})
	
	if err != nil {
		log.Println("Failed to initalize the client: ", err)
		return nil, err
	}

	return client, nil
}

func GetClientAuth() (*supabase.Client, error) {

	godotenv.Load()

	SUPABASE_URL := os.Getenv("SUPABASE_URL")
	SUPABASE_PUBLISHABLE_KEY := os.Getenv("SUPABASE_PUBLISHABLE_KEY")

	options := &supabase.ClientOptions{
		Schema: "next_auth",
	}

	client, err := supabase.NewClient(SUPABASE_URL, SUPABASE_PUBLISHABLE_KEY, options)
	
	if err != nil {
		log.Println("Failed to initalize the client: ", err)
		return nil, err
	}

	return client, nil
}