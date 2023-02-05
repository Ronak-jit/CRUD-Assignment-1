package store

import (
	"context"
	"fmt"
	"log"
	"regexp"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

func SetUpSpanner(dbURI string) {
	ctx := context.Background()
	dbExist, _ := CheckExistingDb(ctx, dbURI)
	if !dbExist {
		err := CreateDatabase(ctx, dbURI)
		if err != nil {
			log.Fatalf("error while creating database: %v", err)
		} else {
			fmt.Println("Database created")
		}
	}
}

func CreateDatabase(ctx context.Context, dbURI string) error {
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(dbURI)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("invalid database id %s", dbURI)
	}

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	op, err := adminClient.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
		Parent:          matches[1],
		CreateStatement: "CREATE DATABASE `" + matches[2] + "`",
		ExtraStatements: []string{
			`CREATE TABLE Users (
                                Id   STRING(1024) NOT NULL,
                                FirstName  STRING(1024),
                                LastName   STRING(1024),
                        ) PRIMARY KEY (Id)`,
		},
	})
	if err != nil {
		return err
	}
	if _, err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Printf("Created database [%s]\n", dbURI)
	return nil
}
