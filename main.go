package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_"github.com/go-sql-driver/mysql"
)

var plantillas = template.Must(template.ParseGlob("plantillas/*"))
type Usuario struct {
    Id int
    Nombre string
    Correo string
}

func main() {
    http.HandleFunc("/", Inicio)
    http.HandleFunc("/crear", Crear)
    http.HandleFunc("/insertar", Insertar)
    http.HandleFunc("/borrar", Borrar)
    http.HandleFunc("/editar", Editar)
    http.HandleFunc("/actualizar", Actualizar)
    log.Println("Servidor ejecutando...")
    http.ListenAndServe(":3000", nil)
}

func Actualizar(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        id := r.FormValue("id")
        nombre := r.FormValue("nombre")
        correo := r.FormValue("correo")

        conexionCreada := conexionDB()
        modificarRegistro, err := conexionCreada.Prepare("UPDATE usuarios SET nombre = ?, correo = ? WHERE id = ?")
        if err != nil {
            panic(err.Error())
        }

        modificarRegistro.Exec(nombre, correo, id)
        http.Redirect(w, r, "/", 301)
    }
}

func Editar(w http.ResponseWriter, r *http.Request) {
    idUsuario := r.URL.Query().Get("id")
    conexionCreada := conexionDB()
    registro, err := conexionCreada.Query("SELECT * FROM usuarios WHERE id = ?", idUsuario)
    usuario := Usuario{}
    for registro.Next() {
        var id  int
        var nombre, correo string
        err = registro.Scan(&id, &nombre, &correo)
        if err != nil {
            panic(err.Error())
        }
        usuario.Id = id
        usuario.Nombre = nombre
        usuario.Correo = correo
    }
    plantillas.ExecuteTemplate(w, "editar", usuario)
}

func Borrar(w http.ResponseWriter, r *http.Request) {
    idUsuario := r.URL.Query().Get("id")
    conexionCreada := conexionDB()
    borrarRegistro, err := conexionCreada.Prepare("DELETE FROM usuarios WHERE id = ?")
    if err != nil {
        panic(err.Error())
    }
    borrarRegistro.Exec(idUsuario)
    http.Redirect(w, r, "/", 301)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        nombre := r.FormValue("nombre")
        correo := r.FormValue("correo")

        conexionCreada := conexionDB()
        insertarRegistros,err := conexionCreada.Prepare("INSERT INTO usuarios(nombre, correo) VALUES(?, ?)")
        if err != nil {
            panic(err.Error())
        }
        insertarRegistros.Exec(nombre, correo)
        http.Redirect(w, r, "/", 301)
    }
}

func Inicio(w http.ResponseWriter, r *http.Request) {
    conexionCreada := conexionDB()
    registros, err := conexionCreada.Query("SELECT * FROM usuarios")
    if err != nil {
        panic(err.Error())
    }
    usuario := Usuario {}
    arregloUsuario := []Usuario{}
    for registros.Next() {
        var id int
        var nombre, correo string
        err = registros.Scan(&id, &nombre, &correo)
        if err != nil {
            panic(err.Error())
        }
        usuario.Id = id
        usuario.Nombre = nombre
        usuario.Correo = correo
        arregloUsuario = append(arregloUsuario, usuario)
    }
    plantillas.ExecuteTemplate(w, "inicio", arregloUsuario)
}

func Crear(w http.ResponseWriter, r *http.Request) {
    plantillas.ExecuteTemplate(w, "crear", nil)
}

func conexionDB()(conexion *sql.DB) {
    Driver := "mysql"
    Usuario := "laravel"
    Contra := "laravel"
    Nombre := "sistema_go"

    conexion, err := sql.Open(Driver, Usuario+":"+Contra+"@tcp(127.0.0.1:3306)/"+Nombre)
    if err != nil {
        panic(err.Error())
    }
    return conexion
}
