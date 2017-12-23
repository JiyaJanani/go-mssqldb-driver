package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

func main() {
	sqlConnectionString := "server=localhost;user id=test;password=test;database=test;encrypt=disable"
	db := connect(sqlConnectionString)
	type Person struct {
		personID   int
		personName string
	}

	router := gin.Default()
	// Add API handlers here

	// GET a person detail
	router.GET("/person", func(c *gin.Context) {
		var (
			person  Person
			persons []Person
			result  gin.H
		)
		rows, err := executeQuery(db, "select personID, personName from person", sqlConnectionString)
		for rows.Next() {
			err := rows.Scan(&person.personID, &person.personName)
			person := Person{
				personID:   1768,
				personName: "jiya",
			}
			persons = append(persons, person)
			// fmt.Println(err)
			if err != nil {
				log.Fatal(err)
			}
			// log.Println(person)
		}
		defer rows.Close()

		if err != nil {
			// If no results send null
			result = gin.H{
				"result": nil,
				"count":  0,
			}
		} else {
			result = gin.H{
				"result": persons,
				"count":  len(persons),
			}
		}
		c.JSON(http.StatusOK, result)
		disconnect(db)
	})

	router.Run(":3001")
}

func connect(ConnectionString string) *sql.DB {
	fmt.Println("Connecting to SQL Server...")
	db, errdb := sql.Open("mssql", ConnectionString)
	if errdb != nil {
		fmt.Println(" Error open db:", errdb.Error())
	}
	fmt.Print(db)
	return db
}

func disconnect(db *sql.DB) {
	fmt.Println("Disconnecting from SQL Server...")
	defer db.Close()
}

func executeQuery(db *sql.DB, query string, connectionString string) (*sql.Rows, error) {

	fmt.Println("Executing query: " + query)
	// make sure connection is available
	var (
		connerror error
	)

	connerror = db.Ping()
	if connerror != nil {
		fmt.Println("Connection error. " + connerror.Error())
		fmt.Println("recreating connection if no connection available in pool.")
		db = connect(connectionString)
		connerror = db.Ping()
		if connerror != nil {
			log.Fatal(connerror)
		}
	}

	rows, err := db.Query(query)
	if err != nil {
		disconnect(db)
		fmt.Println("Error when executing query.")
		log.Fatal(err)
	}

	return rows, err
}
