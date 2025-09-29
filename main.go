package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	db "tp_especial.com/servidor-go/db/sqlc"
	sqlc "tp_especial.com/servidor-go/db/sqlc"

	_ "github.com/lib/pq"
)

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Falta %s (ej.: postgres://user:pass@host:port/db?sslmode=disable)", key)
	}
	return v
}

func ensureSchema(ctx context.Context, sqlDB *sql.DB) {
	const ddl = `
CREATE TABLE IF NOT EXISTS turnos (
  id_turno   SERIAL PRIMARY KEY,
  id_medico  INTEGER NOT NULL,
  id_paciente INTEGER NOT NULL,
  estado     TEXT NOT NULL DEFAULT 'programado',
  fecha      DATE NOT NULL,
  hora       TIME NOT NULL
);
`
	if _, err := sqlDB.ExecContext(ctx, ddl); err != nil {
		log.Fatalf("Error creando schema: %v", err)
	}
}


func ProbarTurnos() {
    fmt.Println("Prueba de turnos iniciada...")

		ctx := context.Background()
	dsn := mustEnv("DATABASE_URL")

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("open DB: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatalf("ping DB: %v", err)
	}

	ensureSchema(ctx, sqlDB)

	q := db.New(sqlDB)

	// === 1) Crear un turno futuro programado ===
	maniana := time.Now().Add(24 * time.Hour)
	fecha := time.Date(maniana.Year(), maniana.Month(), maniana.Day(), 0, 0, 0, 0, time.Local)
	hora := time.Date(1, 1, 1, 10, 30, 0, 0, time.Local) // 10:30

	t1, err := q.CrearTurno(ctx, db.CrearTurnoParams{
		IDMedico:   101,
		IDPaciente: 201,
		Fecha:      fecha,
		Hora:       hora,
	})
	if err != nil {
		log.Fatalf("CrearTurno: %v", err)
	}
	fmt.Printf("CrearTurno -> ID=%d Estado=%s\n", t1.IDTurno, t1.Estado)

	// Esperado: estado "programado"
	if t1.Estado != "programado" {
		log.Fatalf("Esperaba estado programado, obtuve %q", t1.Estado)
	}

	// === 2) ObtenerTurnoPorId ===
	got, err := q.ObtenerTurnoPorId(ctx, t1.IDTurno)
	if err != nil {
		log.Fatalf("ObtenerTurnoPorId: %v", err)
	}
	fmt.Printf("ObtenerTurnoPorId -> ID=%d Medico=%d Paciente=%d Estado=%s Hora=%s\n",
		got.IDTurno, got.IDMedico, got.IDPaciente, got.Estado, got.Hora.Format("15:04"))

	// === 3) ActualizarTurnoDatos (cambio médico y hora) ===
	nuevaHora := time.Date(1, 1, 1, 11, 0, 0, 0, time.Local) // 11:00
	if err := q.ActualizarTurnoDatos(ctx, db.ActualizarTurnoDatosParams{
		IDTurno:    t1.IDTurno,
		IDMedico:   102,
		IDPaciente: 201,
		Fecha:      fecha,
		Hora:       nuevaHora,
	}); err != nil {
		log.Fatalf("ActualizarTurnoDatos: %v", err)
	}
	got, _ = q.ObtenerTurnoPorId(ctx, t1.IDTurno)
	fmt.Printf("Tras ActualizarTurnoDatos -> Medico=%d Hora=%s\n", got.IDMedico, got.Hora.Format("15:04"))

	// Esperado: id_medico=102 y hora=11:00
	if got.IDMedico != 102 || got.Hora.Format("15:04") != "11:00" {
		log.Fatalf("ActualizarTurnoDatos no aplicó como esperaba")
	}

	// === 4) CambiarEstadoTurno a "reprogramado" ===
	if err := q.CambiarEstadoTurno(ctx, db.CambiarEstadoTurnoParams{
		IDTurno: t1.IDTurno,
		Estado:  "reprogramado",
	}); err != nil {
		log.Fatalf("CambiarEstadoTurno: %v", err)
	}
	got, _ = q.ObtenerTurnoPorId(ctx, t1.IDTurno)
	fmt.Printf("Tras CambiarEstadoTurno -> Estado=%s\n", got.Estado)
	if got.Estado != "reprogramado" {
		log.Fatalf("Esperaba estado reprogramado")
	}

	// === 5) AtenderTurno requiere estado 'programado' ===
	// Lo vuelvo a programado y luego lo atiendo
	if err := q.CambiarEstadoTurno(ctx, db.CambiarEstadoTurnoParams{
		IDTurno: t1.IDTurno,
		Estado:  "programado",
	}); err != nil {
		log.Fatalf("Reponer a programado: %v", err)
	}
	if err := q.AtenderTurno(ctx, t1.IDTurno); err != nil {
		log.Fatalf("AtenderTurno: %v", err)
	}
	got, _ = q.ObtenerTurnoPorId(ctx, t1.IDTurno)
	fmt.Printf("Tras AtenderTurno -> Estado=%s\n", got.Estado)
	// Esperado: atendido
	if got.Estado != "atendido" {
		log.Fatalf("Esperaba estado atendido")
	}

	// === 6) CancelarTurno (debería pasar a 'cancelado') ===
	if err := q.CancelarTurno(ctx, t1.IDTurno); err != nil {
		log.Fatalf("CancelarTurno: %v", err)
	}
	got, _ = q.ObtenerTurnoPorId(ctx, t1.IDTurno)
	fmt.Printf("Tras CancelarTurno -> Estado=%s\n", got.Estado)
	// Esperado: cancelado
	if got.Estado != "cancelado" {
		log.Fatalf("Esperaba estado cancelado")
	}

	// === 7) Crear otro turno programado para probar listados ===
	t2, err := q.CrearTurno(ctx, db.CrearTurnoParams{
		IDMedico:   102,
		IDPaciente: 201,
		Fecha:      fecha,
		Hora:       time.Date(1, 1, 1, 12, 0, 0, 0, time.Local), // 12:00
	})
	if err != nil {
		log.Fatalf("CrearTurno t2: %v", err)
	}
	fmt.Printf("CrearTurno t2 -> ID=%d Estado=%s\n", t2.IDTurno, t2.Estado)

	// === 8) Obtener listas ===
	listPac, err := q.ObtenerListaTurnosDePaciente(ctx, 201)
	if err != nil {
		log.Fatalf("ObtenerListaTurnosDePaciente: %v", err)
	}
	fmt.Printf("ObtenerListaTurnosDePaciente(201) -> %d turno(s) programado(s)\n", len(listPac))

	listMed, err := q.ObtenerListaTurnosPorMedico(ctx, 102)
	if err != nil {
		log.Fatalf("ObtenerListaTurnosPorMedico: %v", err)
	}
	fmt.Printf("ObtenerListaTurnosPorMedico(102) -> %d turno(s) programado(s)\n", len(listMed))

	histPac, err := q.ObtenerHistorialTurnosDePaciente(ctx, 201)
	if err != nil {
		log.Fatalf("ObtenerHistorialTurnosDePaciente: %v", err)
	}
	fmt.Printf("Historial paciente 201 -> %d turno(s)\n", len(histPac))

	histMed, err := q.ObtenerHistorialTurnosPorMedico(ctx, 102)
	if err != nil {
		log.Fatalf("ObtenerHistorialTurnosPorMedico: %v", err)
	}
	fmt.Printf("Historial medico 102 -> %d turno(s)\n", len(histMed))

	// === Validaciones esperadas de los listados ===
	// Después de cancelar t1 y crear t2 (programado),
	// - Listas "programado" deberían contener al menos t2.
	if len(listPac) == 0 || len(listMed) == 0 {
		log.Fatalf("Listas de programados vacías; esperaba al menos 1 (t2)")
	}

	fmt.Println("Prueba de turnos terminada.")
	os.Exit(0)
}

func main() {
	ProbarTurnos()
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
