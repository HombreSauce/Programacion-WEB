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

	//CREAMOS OTRO USUARIO
	usuarioCreado, err := queries.CrearUsuario(ctx,
		sqlc.CrearUsuarioParams{
			Dni:             "12345678",
			Nombre:          "Roman",
			Apellido:        "Gonzalez",
			Sexo:            "Masculino",
			FechaNacimiento: time.Date(1997, 6, 30, 0, 0, 0, 0, time.UTC),
			Email:           "romanGonzalez@ejemplo.com",
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

	//CREAMOS UNA RELACION MEDICO OBRA
	relacionCreada, err := queries.CrearRelacionMedicoObra(ctx, sqlc.CrearRelacionMedicoObraParams{
		IDMedico:         usuarioCreado.ID,
		ObraSocialNombre: "IOMA",
	})
	if err != nil {
		log.Fatalf("No se pudo crear la relación medico obra social: %v", err)
	}
	fmt.Printf("Relacion creada: %+v\n", relacionCreada)

	//INTENTAMOS CREAR UNA RELACION MEDICO OBRA CON UNA OBRA SOCIAL QUE NO EXISTE
	relacionCreada, err = queries.CrearRelacionMedicoObra(ctx, sqlc.CrearRelacionMedicoObraParams{
		IDMedico:         usuarioCreado.ID,
		ObraSocialNombre: "IOMA SQL",
	})
	if err != nil {
		log.Fatalf("No se pudo crear la relación medico obra social: %v", err)
	}
	fmt.Printf("Relacion creada: %+v\n", relacionCreada)

	//VEMOS LAS OBRAS SOCIALES POR MEDICO
	obrasSociales, err := queries.ListarObrasPorMedico(ctx, usuarioCreado.ID)
	if err != nil {
		log.Fatalf("No se pueden listar las obras sociales: %v", err)
	}
	fmt.Printf("Las obras sociales por las que atiende el medico son: %+v\n", obrasSociales)

	//VEMOS LOS MEDICOS POR OBRA SOCIAL
	medicos, err := queries.ListarMedicosPorObra(ctx, "IOMA")
	if err != nil {
		log.Fatalf("No se pueden listar los medicos: %v", err)
	}
	fmt.Printf("Los medicos que atiendes por esa obra social son: %+v\n", medicos)

	medicos, err = queries.ListarMedicosPorObra(ctx, "PAMI - INSSJP")
	if err != nil {
		log.Fatalf("No se pueden listar los medicos: %v", err)
	}
	fmt.Printf("Los medicos que atiendes por esa obra social son: %+v\n", medicos)

	//ELIMINAMOS RELACION MEDICO OBRA SOCIAL
	filas, err := queries.EliminarRelacionMedicoObra(ctx, sqlc.EliminarRelacionMedicoObraParams{
		IDMedico:         usuarioCreado.ID,
		ObraSocialNombre: "IOMA",
	})
	if err != nil {
		log.Fatalf("No se pudo borrar la relación: %v", err)
	}
	if filas == 0 {
		log.Fatalf("No se encontró la relación especificada")
	}
	fmt.Println("La relación se borró satisfactoriamente")

	filas, err = queries.EliminarRelacionMedicoObra(ctx, sqlc.EliminarRelacionMedicoObraParams{
		IDMedico:         usuarioCreado.ID,
		ObraSocialNombre: "PAMI",
	})
	if err != nil {
		log.Fatalf("No se pudo borrar la relación: %v", err)
	}
	if filas == 0 {
		log.Fatalf("No se encontró la relación especificada")
	}
	fmt.Println("La relación se borró satisfactoriamente")
}

// func serveIndex(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/" || r.Method != http.MethodGet {
// 		http.NotFound(w, r)
// 		return
// 	}
// 	w.Header().Set("Content-type", "text/html; charset=utf-8")
// 	http.ServeFile(w, r, "./static/index.html")
// }
