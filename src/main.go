package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

var db *sql.DB
var server = "backofficedb.database.windows.net"
var port = 1433
var user = "backoffice"
var password = "D4t4c3nt3rCl0ud"
var database = "backofficeDB"

type Employee struct {
	Id       int    `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
}

func main() {
	//Abrir conexion de base de datos
	openConnection()
	defineEndpointsAndRunApp()
}
func openConnection() {
	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	var err error

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")
}

func defineEndpointsAndRunApp() {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"greet": "hello, world!",
		})
	})

	r.GET("/echo/:echo", func(c *gin.Context) {
		echo := c.Param("echo")
		c.JSON(http.StatusOK, gin.H{
			"echo": echo,
		})
	})

	r.POST("/upload", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)

			// Upload the file to specific dst.
			// c.SaveUploadedFile(file, dst)
		}
		c.JSON(http.StatusOK, gin.H{
			"uploaded": len(files),
		})
	})

	r.GET("/employees", GetEmployees)
	r.Run() // listen and serve on 0.0.0.0:8080

}

func GetEmployees(c *gin.Context) {
	//Modificado para conectar con base de datos Microsoft SQL Server

	// Read employees
	var employees []Employee
	count, err, employees := ReadEmployees()
	if err != nil {
		log.Fatal("Error reading Employees: ", err.Error())
	}
	c.JSON(http.StatusOK, employees)
	fmt.Printf("Read %d row(s) successfully.\n", count)

}

func ReadEmployees() (int, error, []Employee) {
	ctx := context.Background()
	var employees []Employee
	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return -1, err, employees
	}

	tsql := fmt.Sprintf("SELECT Id, Name, Location FROM SalesLT.Employees;")

	// Execute query
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		return -1, err, employees
	}

	defer rows.Close()

	var count int
	// Iterate through the result set.
	for rows.Next() {
		var emp Employee
		// Get values from row.
		err := rows.Scan(&emp.Id, &emp.Name, &emp.Location)
		if err != nil {
			return -1, err, employees
		}
		employees = append(employees, emp)
		count++
	}
	return count, nil, employees
}