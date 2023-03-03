package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Password := "1098741116"
	NombreBD := "empleadosapigo"

	conexion, err := sql.Open(Driver, Usuario+":"+Password+"@tcp(127.0.0.1)/"+NombreBD)

	if err != nil {
		panic(err.Error())
	}

	return conexion
}

var templates = template.Must(template.ParseGlob("templates/*"))

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/createEmployee", createEmployee)
	http.HandleFunc("/insertar", insertar)
	http.HandleFunc("/deleteEmployee", delete)
	http.HandleFunc("/updateEmployee", update)
	http.HandleFunc("/actualizar", actualizar)
	fmt.Println("server on")
	http.ListenAndServe(":3000", nil)
}

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

func index(w http.ResponseWriter, r *http.Request) {
	conexionEstablecida := conexionBD()
	registros, err := conexionEstablecida.Query("SELECT * FROM empleados")

	if err != nil {
		panic(err.Error())
	}
	empleado := Empleado{}
	arregloEmpleado := []Empleado{}

	for registros.Next() {
		var id int
		var nombre, correo string
		err = registros.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado)
	}

	templates.ExecuteTemplate(w, "index", arregloEmpleado)
}

func createEmployee(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "create", nil)
}

func insertar(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre,correo) VALUES(?,?)")

		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(nombre, correo)

		http.Redirect(w, r, "/", 301)
	}
}

func delete(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")

	conexionEstablecida := conexionBD()
	borrarRegistro, err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE id=?")

	if err != nil {
		panic(err.Error())
	}
	borrarRegistro.Exec(idEmpleado)

	http.Redirect(w, r, "/", 301)
}

func update(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	conexionEstablecida := conexionBD()
	registro, err := conexionEstablecida.Query("SELECT * FROM empleados WHERE id=?", idEmpleado)

	if err != nil {
		panic(err.Error())
	}

	empleado := Empleado{}
	for registro.Next() {
		var id int
		var nombre, correo string
		err = registro.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

	}

	templates.ExecuteTemplate(w, "update", empleado)
}

func actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		idEmpleado := r.URL.Query().Get("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()
		actualizarRegistros, err := conexionEstablecida.Prepare("UPDATE empleados SET nombre=?,correo=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}
		actualizarRegistros.Exec(nombre, correo, idEmpleado)

		http.Redirect(w, r, "/", 301)
	}
}
