package main

import (
	
	"database/sql"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	_"github.com/go-sql-driver/mysql"
)

var database = "backoffice"

type Employee struct {
	Id       int    `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
}

func defineEndpointsAndRunApp() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/employees", GetEmployees)
	r.Run() // listen and serve on 0.0.0.0:8080

}

func GetEmployees(c *gin.Context){
	// Read employees
	var employees []Employee
	employees, err := ReadEmpleados()	
	c.JSON(http.StatusOK, employees)
	fmt.Printf("Read %d row(s) successfully.\n", err)

}

func obtenerBaseDeDatos() (db *sql.DB, e error) {
	usuario := "root"
	pass := "1234"
	host := "tcp(35.226.5.182:3306)"
	nombreBaseDeDatos := "backoffice"
	// Debe tener la forma usuario:contraseña@host/nombreBaseDeDatos
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", usuario, pass, host, nombreBaseDeDatos))
	if err != nil {
		return nil, err
	}
	return db, nil
}


func ReadEmpleados() ([]Employee, error) {
	employees := []Employee{}
	db, err := obtenerBaseDeDatos()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	filas, err := db.Query("Select * from employees")

	if err != nil {
		return nil, err
	}
	// Si llegamos aquí, significa que no ocurrió ningún error
	defer filas.Close()

	// Aquí vamos a "mapear" lo que traiga la consulta en el while de más abajo
	var emp Employee

	// Recorrer todas las filas, en un "while"
	for filas.Next() {
		err = filas.Scan(&emp.Id, &emp.Name, &emp.Location)
		// Al escanear puede haber un error
		if err != nil {
			return nil, err
		}
		// Y si no, entonces agregamos lo leído al arreglo
		employees = append(employees, emp)
	}
	fmt.Printf("Conculta hecha correctamente")
	// Vacío o no, regresamos el arreglo de contactos

	
	return employees, nil
}

func main() {
	db, err := obtenerBaseDeDatos()
	if err != nil {
		fmt.Printf("Error obteniendo base de datos: %v", err)
		return
	}
	// Terminar conexión al terminar función
	defer db.Close()

	// Ahora vemos si tenemos conexión
	err = db.Ping()
	if err != nil {
		fmt.Printf("Error conectando: %v", err)
		return
	}
	// Listo, aquí ya podemos usar a db!
	fmt.Printf("Conectado correctamente")
	defineEndpointsAndRunApp()
	
}