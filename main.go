package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	sqlc "tp_especial.com/servidor-go/db/sqlc"

	_ "github.com/lib/pq"
)

func main() {
	//INICIALIZAR SERVIDOR PARA PAGINA WEB
	// http.HandleFunc("/", serveIndex)
	// port := ":8080"
	// fmt.Printf("Servidor escuchando en el puerto http://localhost%s\n", port)
	// err := http.ListenAndServe(port, nil) //inicia el servidor
	// if err != nil {
	// 	fmt.Printf("Error: %s\n", err)
	// }

	//BASE DE DATOS
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=base_turnero sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()
	queries := sqlc.New(db)
	ctx := context.Background()

	//CREAMOS UN USUARIO
	usuarioCreado, err := queries.CrearUsuario(ctx,
		sqlc.CrearUsuarioParams{
			Dni:             "12345678",
			Nombre:          "Juan",
			Apellido:        "Caballo",
			Sexo:            "Masculino",
			FechaNacimiento: time.Date(2000, 5, 20, 0, 0, 0, 0, time.UTC),
			Email:           "juan.caballo@ejemplo.com",
			Telefono:        "2494505050",
			Rol:             "P",
		})

	if err != nil {
		log.Fatalf("No se pudo crear el usuario: %v", err)
	}
	fmt.Printf("Usuario creado: %+v\n", usuarioCreado)

	//OBTENEMOS SU INFORMACIÓN
	usuario, err := queries.ObtenerUsuario(ctx, usuarioCreado.ID)
	if err != nil {
		log.Fatalf("No se pudo obtener el usuario: %v", err)
	}
	fmt.Printf("Usuario obtenido: %+v\n", usuario)

	//LO GUARDAMOS COMO PACIENTE
	pacienteCreado, err := queries.CrearPaciente(ctx,
		sqlc.CrearPacienteParams{
			IDPaciente: usuarioCreado.ID,
		})
	if err != nil {
		log.Fatalf("No se pudo crear el paciente: %v", err)
	}
	fmt.Printf("Paciente creado: %+v\n", pacienteCreado)

	//CREAMOS OTRO USUARIO
	usuarioCreado, err = queries.CrearUsuario(ctx,
		sqlc.CrearUsuarioParams{
			Dni:             "12345678",
			Nombre:          "Juana",
			Apellido:        "Maria",
			Sexo:            "Femenino",
			FechaNacimiento: time.Date(1997, 6, 30, 0, 0, 0, 0, time.UTC),
			Email:           "juanita.maria@ejemplo.com",
			Telefono:        "2494505045",
			Rol:             "M",
		})
	if err != nil {
		log.Fatalf("No se pudo crear el usuario: %v", err)
	}
	fmt.Printf("Usuario creado: %+v\n", usuarioCreado)

	//LO GUARDAMOS COMO MEDICO
	medicoCreado, err := queries.CrearMedico(ctx,
		sqlc.CrearMedicoParams{
			IDMedico:     usuarioCreado.ID,
			NroMatricula: 100,
			Especialidad: "Cardiología",
		})

	if err != nil {
		log.Fatalf("No se pudo crear el medico: %v", err)
	}
	fmt.Printf("Medico creado: %+v\n", medicoCreado)

	//OBTENEMOS EL MEDICO
	medico, err := queries.ObtenerMedico(ctx, usuarioCreado.ID)
	if err != nil {
		log.Fatalf("No se pudo obtener el médico: %v", err)
	}
	fmt.Printf("Médico obtenido: %+v\n", medico)

	//Actualicemos la información del primer usuario
	err = queries.ActualizarUsuario(ctx,
		sqlc.ActualizarUsuarioParams{
			ID:              1,
			Dni:             "12345678",
			Nombre:          "Juan",
			Apellido:        "Caballo",
			Sexo:            "Masculino",
			FechaNacimiento: time.Date(2000, 5, 20, 0, 0, 0, 0, time.UTC),
			Email:           "juan.caballo123@ejemplo.com",
			Telefono:        "2494505050",
			Rol:             "P",
		})
	if err != nil {
		log.Fatalf("No se pudo actualizar el usuario: %v", err)
	}
	fmt.Printf("El usuario se actualizó correctamente \n")

	//VEMOS TODOS LOS USUARIOS QUE TENEMOS CREADOS
	usuarios, err := queries.ListarUsuarios(ctx)
	if err != nil {
		log.Fatalf("No se pueden listar los usuarios: %v", err)
	}
	fmt.Printf("Todos los usuarios: %+v\n", usuarios)

	//BORRAMOS UN USUARIO QUE ES PACIENTE
	err = queries.EliminarUsuario(ctx, 1)
	if err != nil {
		log.Fatalf("No se pudo borrar el usuario: %v", err)
	}
	fmt.Println("El usuario ID 1 se borró satisfactoriamente")

	//VEMOS TODOS LOS MEDICOS
	medicos, err := queries.ListarMedicos(ctx)
	if err != nil {
		log.Fatalf("No se pueden listar los medicos: %v", err)
	}
	fmt.Printf("Todos los medicos: %+v\n", medicos)

	//VEMOS TODOS LOS PACIENTES
	pacientes, err := queries.ListarPacientes(ctx)
	if err != nil {
		log.Fatalf("No se pueden listar los pacientes: %v", err)
	}
	fmt.Printf("Todos los pacientes: %+v\n", pacientes)

	//ELIMINAMOS MEDICO
	err = queries.EliminarMedico(ctx, 2)
	if err != nil {
		log.Fatalf("No se pudo borrar el medico: %v", err)
	}
	fmt.Println("El medico ID 2 se borró satisfactoriamente")

	//AHORA VEMOS QUE SE BORRO DE MEDICO, PERO NO DE USUARIOS
	medicos, err = queries.ListarMedicos(ctx)
	if err != nil {
		log.Fatalf("No se pueden listar los medicos: %v", err)
	}
	fmt.Printf("Todos los medicos: %+v\n", medicos)

	usuarios, err = queries.ListarUsuarios(ctx)
	if err != nil {
		log.Fatalf("No se pueden listar los usuarios: %v", err)
	}
	fmt.Printf("Todos los usuarios: %+v\n", usuarios)
}

// func serveIndex(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/" || r.Method != http.MethodGet {
// 		http.NotFound(w, r)
// 		return
// 	}
// 	w.Header().Set("Content-type", "text/html; charset=utf-8")
// 	http.ServeFile(w, r, "./static/index.html")
// }
