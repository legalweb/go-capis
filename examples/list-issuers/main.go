package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"lwebco.de/go-capis"
)

var (
	username = flag.String("username", "", "comparisonapis.com username")
	password = flag.String("password", "", "comparisonapis.com password")
	token    = flag.String("token", "", "comparisonapis.com token")
)

func main() {
	flag.Parse()

	ctx := context.Background()

	var auth capis.AuthProvider

	if *token != "" {
		auth = capis.StaticToken(*token)
	} else {
		auth = &capis.PasswordAuthentication{
			Username: *username,
			Password: *password,
		}
	}

	client, err := capis.New(capis.WithAuthProvider(auth))
	if err != nil {
		log.Fatalf("error initialising client %v\n", err)
	}

	issuers, err := client.ListIssuers(ctx, nil, 0, 10)
	if err != nil {
		log.Fatalf("error listing issuers %v\n", err)
	}

	for _, v := range issuers.Data {
		fmt.Printf("%s, ", v.Label)
	}
}
