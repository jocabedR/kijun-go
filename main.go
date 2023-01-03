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

func dbConnection() *sql.DB {
	conn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return db
}

func getUserByUsername(c *gin.Context) {
	db := dbConnection()
	if db == nil {
		fmt.Println("Cannot connect to PostgreSQL!")
		db.Close()
	}
	defer db.Close()

	usernamePath := c.Param("username")
	rows, _ := db.Query("SELECT id, username, name, birth_date, registration_date FROM users WHERE username = $1", usernamePath)
	defer rows.Close()

	var id int
	var username, name, birth_date, registration_date string
	var user_unmarshal User
	for rows.Next() {
		rows.Scan(&id, &username, &name, &birth_date, &registration_date)
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
	router := setupRouter()
	router.Run("localhost:3000")
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/whoami/:username", getUserByUsername)

	return router
}
