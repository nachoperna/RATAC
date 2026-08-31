package domain

type Permiso string

const (
	PermVerPanel        Permiso = "ver_panel"
	PermCargarDiag      Permiso = "cargar_diagnostico"
	PermBorrarDiag      Permiso = "borrar_diagnostico"
	PermVerDiagnosticos Permiso = "ver_diagnosticos"
	PermGestionUsuarios Permiso = "gestionar_usuarios"
)

// Eje rol -> permiso. El otro eje (ruta -> permiso) vive en main.go.
// Para dar una capacidad nueva a un rol se agrega UNA linea aca.
var PermisosPorRol = map[Rol][]Permiso{
	RolPublico:     {},
	RolLaboratorio: {PermVerPanel, PermCargarDiag, PermBorrarDiag, PermVerDiagnosticos},
	RolAdmin:       {PermVerPanel, PermCargarDiag, PermBorrarDiag, PermVerDiagnosticos, PermGestionUsuarios},
}

func Puede(u *Usuario, p Permiso) bool {
	if u == nil {
		return false
	}
	for _, permiso := range PermisosPorRol[u.Rol] {
		if permiso == p {
			return true
		}
	}
	return false
}
