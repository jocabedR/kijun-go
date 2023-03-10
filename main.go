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

var DB *sql.DB

func dbConnection() {
	conn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
	}

	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(1)
	db.SetConnMaxLifetime(5)

	DB = db
}

func getUserByUsername(c *gin.Context) {
	db := DB

	usernamePath := c.Param("username")
	rows, _ := db.Query("SELECT id, username, name, birth_date, registration_date FROM users WHERE username = $1", usernamePath)
	defer rows.Close()

	var id int
	var username, name, birth_date, registration_date string
	var user_unmarshal User
	for rows.Next() {
		if err := rows.Scan(&id, &username, &name, &birth_date, &registration_date); err != nil {
			fmt.Println(err)
		}
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
	dbConnection()
	if DB == nil {
		fmt.Println("Database connection failed.")
	} else {
		fmt.Println("Sucessful database connection!")
		router := setupRouter()
		err := router.Run("localhost:3000")
		if err != nil {
			fmt.Println(err)
		}
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/whoami/:username", getUserByUsername)

	return router
}
