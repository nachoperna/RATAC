package ui

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func enviarMailSolicitudAprobada (nombre, email, contraseña string) error {
	m := gomail.NewMessage()

	setHeader(m, os.Getenv("FROM_MAIL"), email, "RATAC - Solicitud Aprobada")
	incluirImagenes(m)

	var html bytes.Buffer
	tmp, err := template.ParseFiles("./infrastructure/UI/static/email_solicitud_aprobada.html")
	if err != nil {
		return fmt.Errorf("error al parsear template de mail: %w", err)
	}
	tmp.Execute(&html, struct{
		Nombre string
		Email string
		Contraseña string
	}{
		Nombre: nombre,
		Email: email,
		Contraseña: contraseña,
	})

	m.SetBody("text/html", html.String())

	port, _ := strconv.Atoi(os.Getenv("PUERTO_MAIL"))
	d := gomail.NewDialer(os.Getenv("HOST_MAIL"), port, os.Getenv("FROM_MAIL"), os.Getenv("CLAVE_ACCESO"))
	err = d.DialAndSend(m) 
	return err
}

func enviarMailSolicitudRechazada (nombre, email string) error {
	m := gomail.NewMessage()

	setHeader(m, os.Getenv("FROM_MAIL"), email, "RATAC - Solicitud Rechazada")
	incluirImagenes(m)

	var html bytes.Buffer
	tmp, err := template.ParseFiles("./infrastructure/UI/static/email_solicitud_rechazada.html")
	if err != nil {
		return fmt.Errorf("error al parsear template de mail: %w", err)
	}
	tmp.Execute(&html, struct{
		Nombre string
	}{
		Nombre: nombre,
	})

	m.SetBody("text/html", html.String())

	port, _ := strconv.Atoi(os.Getenv("PUERTO_MAIL"))
	d := gomail.NewDialer(os.Getenv("HOST_MAIL"), port, os.Getenv("FROM_MAIL"), os.Getenv("CLAVE_ACCESO"))
	err = d.DialAndSend(m) 
	return err
}

func enviarMailSolicitud (nombre, email, ciudad_origen string, veterinarios, matriculas []string) error {
	m := gomail.NewMessage()

	setHeader(m, os.Getenv("FROM_MAIL"), os.Getenv("FROM_MAIL"), "Nueva Solicitud de Colaboración")
	incluirImagenes(m)

	var html bytes.Buffer
	tmp, err := template.ParseFiles("./infrastructure/UI/static/email_solicitud.html")
	if err != nil {
		return fmt.Errorf("error al parsear template de mail: %w", err)
	}
	tmp.Execute(&html, struct{
		Nombre string
		Email string
		Ciudad string
		Veterinarios []string
		Matriculas []string
	}{
		Nombre: nombre,
		Email: email,
		Ciudad: ciudad_origen,
		Veterinarios: veterinarios,
		Matriculas: matriculas,
	})

	m.SetBody("text/html", html.String())

	port, _ := strconv.Atoi(os.Getenv("PUERTO_MAIL"))
	d := gomail.NewDialer(os.Getenv("HOST_MAIL"), port, os.Getenv("FROM_MAIL"), os.Getenv("CLAVE_ACCESO"))
	err = d.DialAndSend(m) 
	return err
}

func setHeader (m *gomail.Message, from, to, subject string) {
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
}

func incluirImagenes (m *gomail.Message) {
	m.Embed("./infrastructure/UI/static/img/RATAC LOGO FULL COLOR HORIZONTAL-high.png", 
		gomail.SetHeader(map[string][]string{
			"Content-ID": {"<ratac_logo>"}, // en el html se usa el alias ratac_logo
		}),
	)
}
