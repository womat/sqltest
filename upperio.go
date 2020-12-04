package main

import (
	"fmt"
	"log"
	"upper.io/db.v3/mssql"
)

func storedProcedure() {
	var settings = mssql.ConnectionURL{
		Database: `MM-Karton.SignIT`, // Database name
		Host:     `dev-sql-05`,       // Server IP or name
		User:     `signit`,           // Username
		Password: `signit`,           // Password
	}

	// The database connection is attempted.
	sess, err := mssql.Open(settings)
	if err != nil {
		log.Fatalf("db.Open(): %q\n", err)
	}
	defer sess.Close() // Closing the session is a good practice.

	result, err := sess.Exec(fmt.Sprintf("sppurge_signit_documents @d=%v", 5))

	i, _ := result.RowsAffected()
	fmt.Println(result, i, err)
}
