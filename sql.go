package main

import (
	"database/sql"
	"fmt"
	"net/url"
)

func querySelect() {
	query := url.Values{}
	query.Add("app name", "MyAppName")

	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword("signit", "signit"),
		Host:   fmt.Sprintf("%s:%d", "dev-sql-05", 1433),
		// Path:  instance, // if connecting to an instance instead of a port
		RawQuery: query.Encode(),
	}
	db, _ := sql.Open("sqlserver", u.String())

	//engine, err := xorm.NewEngine(db, dataSourceName)

	rows, err := db.Query("SELECT *  FROM [dbo].[identity_user]")
	fmt.Println(err)
	x, _ := rows.Columns()
	fmt.Println(x)
}
