package db

import (
	sqlc "RATAC/DB/sqlc"
	"RATAC/domain"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
)

type AdminRepository struct {
	queries *sqlc.Queries
}

func NewAdminRepository(queries *sqlc.Queries) *AdminRepository {
	return &AdminRepository { queries: queries }
}

func (r *AdminRepository) MapeoDocumento(contenido bytes.Buffer) (*domain.Paciente, error) {
	var paciente domain.Paciente
	err := json.NewDecoder(&contenido).Decode(&paciente)
	if err != nil {
		return nil, errors.New("No se pudo decodificar JSON")
	}
	return &paciente, nil
}

func (r *AdminRepository) PacienteYaRegistrado(ctx context.Context, protocolo string) bool {
	_, err := r.queries.GetPaciente(ctx, protocolo)
	if err == sql.ErrNoRows{
		return false
	}
	return true
}

func (r *AdminRepository) GetUltimosDiagnosticosCargados(ctx context.Context, offset int8, token string) ([]domain.Paciente, int16, error) {
	pacientes_rows, err := r.queries.ListPacientesByLab(ctx, sqlc.ListPacientesByLabParams{
		Offset: int32(offset),
		Token: token,
	})
	if err != nil {
		return nil, 0, err
	}
	var pacientes []domain.Paciente
	var total int16 = 0
	if len(pacientes_rows) > 0 {
		total = int16(pacientes_rows[0].Total)
	}
	for _, p := range pacientes_rows {
		paciente := pacienteByLabORM(p)
		pacientes = append(pacientes, paciente)
	}
	return pacientes, total, nil
}

func (r *AdminRepository) GetSolicitudes(ctx context.Context) ([]domain.Solicitud, error) {
	var arr_solicitudes []domain.Solicitud
	solicitudes, err := r.queries.GetSolicitudes(ctx)
	if err != nil {
		return nil, err
	}
	for _, solicitud := range solicitudes {
		veterinarios, err := r.queries.GetVeterinarios(ctx, solicitud.Email)
		if err != nil {
			return nil, err
		}
		arr_solicitudes = append(arr_solicitudes, solicitudORM(solicitud, veterinarios))
	}
	return arr_solicitudes, nil
}

func solicitudORM(soli sqlc.GetSolicitudesRow, vet []sqlc.GetVeterinariosRow) domain.Solicitud {
	var vetes []domain.Veterinarios
	for _, v := range vet {
		vetes = append(vetes, domain.Veterinarios{
			Matricula: v.Matricula,
			Nombre: v.Nombre,
		})
	}
	return domain.Solicitud{
		Email: soli.Email,
		Nombre: soli.NombreLab,
		Ciudad: soli.CiudadOrigen,
		Veterinarios: vetes,
	}
}

func (r *AdminRepository) SolicitudAprobada (ctx context.Context, email string) (string, string, error) {
	err := r.queries.CambiarRol(ctx, sqlc.CambiarRolParams{
		Email: email,
		Column2: sqlc.RolesLaboratorio,
	})
	if err != nil {
		return "", "", err
	}
	usuario, err := r.queries.GetUsuario(ctx, email)
	if err != nil {
		return "", "", err
	}
	return usuario.NombreLab, usuario.ContraseñaHash, nil
}

func (r *AdminRepository) SolicitudRechazada (ctx context.Context, email string) (string, error) {
	return r.queries.EliminarUsuario(ctx, email)
}
