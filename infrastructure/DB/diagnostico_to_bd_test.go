package db

import (
	sqlc "RATAC/DB/sqlc"
	"RATAC/domain"
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var db *sql.DB
var queries *sqlc.Queries
var ctx context.Context
const email_lab = "mail@mail.com"

func new(s string) *string { return &s }

func TestTransformarDatos_ValoresValidos_ParseoExitoso(t *testing.T)  {
	fecha := "07-06-2001"
	edad := new("25")
	_, _, err := TransformarDatos(fecha, edad)
	assert.NoError(t, err)
}

func TestTransformarDatos_FechaInvalida_ParseoFallido(t *testing.T)  {
	fecha := "Esto-no-es-una-fecha"
	edad := new("25")
	_, _, err := TransformarDatos(fecha, edad)
	assert.EqualError(t, err, "Error al parsear fecha")
}

func TestTransformarDatos_EdadInvalida_ParseoFallido(t *testing.T)  {
	fecha := "07-06-2001"
	edad := new("EstoNoEsUnaEdad")
	_, _, err := TransformarDatos(fecha, edad)
	assert.EqualError(t, err, "Error al convertir edad")
}

func TestInsertarPaciente_DatosValidos_AltaPacienteDB(t *testing.T)  {
	fecha := "07-06-2001"
	edad := new("10")
	protocolo := "PROT-001"
	paciente := domain.Paciente{
		Protocolo: protocolo,
		Fecha: fecha,
		Solicitante: "Veterinaria",
		Email_Lab: email_lab,
		Tecnica: "HE",
		Familia: new("Perna"),
		Especie: new("Canino"),
		Raza: new("Doberman"),
		Edad: edad,
		NombrePaciente: *new("Nacho"),
		ReferenciasMastocitomas: false,
		Antecedentes: new("Test de antecedentes."),
		DescripcionMacroscopica: new("Test de descripcion macroscopica."),
		Descripciones_microscopicas: nil,
	}
	
	fecha_parseada, edad_parseada, err := TransformarDatos(fecha, edad)
	if err != nil {
		t.FailNow()
	}
	err = InsertarPaciente(paciente, fecha_parseada, edad_parseada, queries, ctx)
	assert.NoError(t, err)
	paciente_retorno, err := queries.GetPaciente(ctx, protocolo)
	assert.NoError(t, err)
	assert.EqualValues(t, paciente.Protocolo, paciente_retorno.Protocolo)
}

func TestInsertarPaciente_DatosInvalidos_SinAltaDB(t *testing.T)  {
	fecha := "07-06-2001"
	edad := new("10")
	paciente := domain.Paciente{
		Protocolo: "TestValidacionProtocoloMayorDeCincuentaCaracteresParaPruebas",
	}
	
	fecha_parseada, edad_parseada, err := TransformarDatos(fecha, edad)
	if err != nil {
		t.FailNow()
	}
	err = InsertarPaciente(paciente, fecha_parseada, edad_parseada, queries, ctx)
	assert.Error(t, err)
}

func TestDeletePaciente_PacienteTest_BajaPacienteDB(t *testing.T)  {
	assert.NoError(t, queries.DeletePaciente(ctx, "PROT-001"))
}

func getConnection() string {
	conn := os.Getenv("TEST_DB_URL")
	if conn != "" {
		return conn
	}
	return "host=localhost port=5433 user=admin password=password dbname=RATAC_DB_TEST sslmode=disable"
}

func TestMain(m *testing.M) {
	db, err := sql.Open("postgres", getConnection())
	if err != nil {
		log.Fatalf("Error al conectar con la Base de Datos: %v", err)
	}
	defer db.Close()
	queries = sqlc.New(db)
	ctx = context.Background()
	_, err = queries.CreateUsuario(ctx, sqlc.CreateUsuarioParams{
		Email: email_lab,
		NombreLab: "lab",
		ContraseñaHash: "no-importa",
		Rol: "laboratorio",
		CiudadOrigen: "testinglandia",
	})
	if err != nil {
		log.Fatalf("Error al insertar Usuario Lab a BD: %v", err)
	}
	m.Run() // ejecuto todos los test
}
