package store

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func CreateUser(dbURI string, user UserInfo) (UserInfo, error) {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, dbURI)
	if err != nil {
		return UserInfo{}, err
	}
	defer client.Close()

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT Users (Id, FirstName, LastName) VALUES
                                (@id, @firstName, @lastName)`,
			Params: map[string]interface{}{
				"id":        user.Id,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
			},
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Printf("%d record(s) inserted.\n", rowCount)
		return err
	})
	if err != nil {
		return UserInfo{}, err
	} else {
		return user, nil
	}
}

func GetUser(userId, dbURI string) (UserInfo, error) {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, dbURI)
	if err != nil {
		return UserInfo{}, err
	}
	defer client.Close()

	stmt := spanner.Statement{SQL: `SELECT Id, FirstName, LastName FROM Users 
									WHERE Id=@id`,
		Params: map[string]interface{}{
			"id": userId,
		}}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	user := UserInfo{}
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return user, nil
		}
		if err != nil {
			return user, err
		}
		var Id, FirstName, LastName string

		if err := row.Columns(&Id, &FirstName, &LastName); err != nil {
			return user, err
		}
		user = UserInfo{
			Id,
			FirstName,
			LastName,
		}
	}
}

func UpdateUser(dbURI string, user UserInfo) error {
	ctx := context.Background()

	client, err := spanner.NewClient(ctx, dbURI)
	if err != nil {
		return err
	}
	defer client.Close()

	cols := []string{"Id", "FirstName", "LastName"}
	_, err = client.Apply(ctx, []*spanner.Mutation{
		spanner.Update("Users", cols, []interface{}{user.Id, user.FirstName, user.LastName}),
	})
	return err
}

func DeleteUser(dbURI string, userId string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, dbURI)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{SQL: `DELETE FROM Users WHERE Id = @id`,
			Params: map[string]interface{}{
				"id": userId,
			},
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Printf("%d record(s) deleted.\n", rowCount)
		return nil
	})
	return err
}
