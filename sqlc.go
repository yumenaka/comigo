package main

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"reflect"

	"github.com/yumenaka/comigo/sqlc"
	_ "modernc.org/sqlite"
)

// Controller/Handler → Service → Repository → sqlc.Queries
//    ↑                    ↑           ↑           ↑
//  HTTP层              业务逻辑层     数据访问层    数据库层

//go:embed sqlc/schema.sql
var ddl string

func run() error {
	ctx := context.Background()

	db, err := sql.Open("sqlite", ":memory_repository:")
	if err != nil {
		return err
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	queries := sqlc.New(db)

	// list all authors
	authors, err := queries.ListStores(ctx)
	if err != nil {
		return err
	}
	log.Println(authors)

	// create a Store
	insertedStore, err := queries.CreateStore(ctx, sqlc.CreateStoreParams{
		Name:          "Brian Kernighan",
		Description:   sql.NullString{String: "The C Programming Language and The Go Programming Language", Valid: true},
		FileBackendID: 0,
	})
	if err != nil {
		return err
	}
	log.Println(insertedStore)

	// get the Store we just inserted
	fetchedStore, err := queries.GetStoreByID(ctx, insertedStore.ID)
	if err != nil {
		return err
	}

	// prints true
	log.Println(reflect.DeepEqual(insertedStore, fetchedStore))
	return nil
}

// func main() {
// 	if err := run(); err != nil {
// 		log.Fatal(err)
// 	}
// }
