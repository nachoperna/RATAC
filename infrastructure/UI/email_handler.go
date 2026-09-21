package ui

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func enviarMail (nombre, email, contraseña string) error {
	m := gomail.NewMessage()

	setHeader(m, os.Getenv("FROM_MAIL"), email)
	incluirImagenes(m)

	var html bytes.Buffer
	tmp, err := template.ParseFiles("./infrastructure/UI/static/email_template.html")
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

func setHeader (m *gomail.Message, from, to string) {
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "RATAC - Solicitud Aprobada")
}

func incluirImagenes (m *gomail.Message) {
	m.Embed("./infrastructure/UI/static/img/RATAC LOGO FULL COLOR HORIZONTAL-high.png", 
		gomail.SetHeader(map[string][]string{
			"Content-ID": {"<ratac_logo>"}, // en el html se usa el alias ratac_logo
		}),
	)
}
