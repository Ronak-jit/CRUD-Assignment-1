package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

type UserInfo struct {
	Id        string
	FirstName string
	LastName  string
}

// CheckExistingDb checks whether the database with dbURI exists or not.
// If API call doesn't respond then user is informed after every 5 minutes on command line.
func CheckExistingDb(ctx context.Context, dbURI string) (bool, error) {
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return false, err
	}
	defer adminClient.Close()
	gotResponse := make(chan bool)
	go func() {
		_, err = adminClient.GetDatabase(ctx, &adminpb.GetDatabaseRequest{Name: dbURI})
		gotResponse <- true
	}()
	for {
		select {
		case <-time.After(5 * time.Minute):
			fmt.Println("WARNING! API call not responding: make sure that spanner api endpoint is configured properly")
		case <-gotResponse:
			if err != nil {
				if containsAny(strings.ToLower(err.Error()), []string{"database not found"}) {
					return false, nil
				}
				return false, fmt.Errorf("can't get database info: %s", err)
			}
			return true, nil
		}
	}
}

func containsAny(s string, l []string) bool {
	for _, a := range l {
		if strings.Contains(s, a) {
			return true
		}
	}
	return false
}
