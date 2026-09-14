# 7. Error: sesión vencida o cookie ausente
GET http://localhost:8080/activo
[Cookies]
token_sesion: eyJhbGciOiJIUzI1NiJ9.eyJleHAiOjB9.sinvalida
HTTP 400
[Asserts]
body contains "SESION TERMINADA"
