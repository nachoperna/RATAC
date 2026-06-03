package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Creamos un mock de os.DirEntry falso para poder pasarlo a la función sin depender de archivos reales
type mockDirEntry struct {
	name string
}

func (m mockDirEntry) Name() string               { return m.name }
func (m mockDirEntry) IsDir() bool                { return false }
func (m mockDirEntry) Type() os.FileMode          { return 0 }
func (m mockDirEntry) Info() (os.FileInfo, error) { return nil, nil }

// --- Tests ---

func TestProcesarJson_archivoInexistente_retornaError(t *testing.T) {
	// Le pasamos el nombre de un archivo que sabemos que no existe
	entradaMock := mockDirEntry{name: "fantasma.json"}
	ctx := context.Background()

	// Pasamos nil en qtx porque el código debería fallar antes de llegar a usar la base de datos
	err := procesarJson(entradaMock, "./ruta_falsa", nil, ctx)

	assert.Error(t, err)
	assert.Equal(t, "No se pudo abrir archivo", err.Error())
}

func TestProcesarJson_jsonInvalido_retornaError(t *testing.T) {
	// Setup: Usamos t.TempDir() que es una magia de Go para crear una carpeta temporal
	// que se borra sola cuando termina el test, así no ensuciamos tu compu.
	dirTemp := t.TempDir()
	nombreArchivo := "roto.json"
	rutaCompleta := filepath.Join(dirTemp, nombreArchivo)

	// Escribimos un JSON totalmente mal formado a propósito
	err := os.WriteFile(rutaCompleta, []byte(`{ "Protocolo": "CAN-2026", "Edad": faltan_comillas }`), 0644)
	assert.NoError(t, err)

	entradaMock := mockDirEntry{name: nombreArchivo}
	ctx := context.Background()

	err = procesarJson(entradaMock, dirTemp, nil, ctx)

	assert.Error(t, err)
	assert.Equal(t, "No se pudo decodificar JSON", err.Error())
}
