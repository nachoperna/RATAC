# ===== FLUJO 1: Registro =====
# Hurl captura y reenvía la cookie automaticamente entre peticiones del mismo archivo

# 1. Registrar un laboratorio
POST http://localhost:8080/registrarse
[FormParams]
email: lab@mail.com
contraseña: lab1password
ciudad_origen: Buenos Aires
HTTP 200
[Asserts]
header "Set-Cookie" contains "token_sesion"
[Captures]
error: body

# 2. La sesión quedó activa
GET http://localhost:8080/activo
HTTP 200
[Asserts]
body contains "SESION ACTIVA"

# 3. Cerrar sesión
GET http://localhost:8080/logout
HTTP 200
[Asserts]
body contains "DESLOGUEADO CON EXITO"


# ===== FLUJO 2: Login =====

# 4. Login con las mismas credenciales
POST http://localhost:8080/login
[FormParams]
email: lab@mail.com
contraseña: lab1password
HTTP 200
[Asserts]
header "Set-Cookie" contains "token_sesion"

# 5. Sesión activa tras el login
GET http://localhost:8080/activo
HTTP 200
[Asserts]
body contains "SESION ACTIVA"

# 6. Error: faltan credenciales
POST http://localhost:8080/login
[FormParams]
email: lab@mail.com
HTTP 400

