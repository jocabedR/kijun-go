package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
	Username          string
	Name              string
	Birth_date        string
	Registration_date string
}

func getUserByUsername(c *gin.Context) {

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

	usernamePath := c.Param("username")
	rows, err := db.Query("SELECT id, username, name, birth_date, registration_date FROM users WHERE username = $1", usernamePath)

	defer rows.Close()

	var id int
	var username, name, birth_date, registration_date string
	var user_unmarshal User

	for rows.Next() {
		err = rows.Scan(&id, &username, &name, &birth_date, &registration_date)
		user_unmarshal.Id = id
		user_unmarshal.Username = username
		user_unmarshal.Name = name
		user_unmarshal.Birth_date = birth_date
		user_unmarshal.Registration_date = registration_date
	}

	if user_unmarshal.Username == usernamePath {
		c.IndentedJSON(http.StatusOK, user_unmarshal)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

func main() {
	router := gin.Default()
	router.GET("/whoiam/:username", getUserByUsername)

	router.Run("localhost:3000")
}
