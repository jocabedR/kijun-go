package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "kijun"
)

type User struct {
	Id                int
	Name              string
	Birth_date        string
	Registration_date string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	username := "jduesterrr"

	rows, err := db.Query("SELECT id, name, birth_date, registration_date FROM users WHERE username = $1", username)

	defer rows.Close()
	for rows.Next() {
		var id int
		var name, birth_date, registration_date string

		err = rows.Scan(&id, &name, &birth_date, &registration_date)

		user_unmarshal := User{
			Id:                id,
			Name:              name,
			Birth_date:        birth_date,
			Registration_date: registration_date,
		}

		fmt.Println(user_unmarshal)

		b, err := json.Marshal(user_unmarshal)
		if err != nil {
			fmt.Println("error:", err)
		}

		os.Stdout.Write(b)
	}
}
