package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const urlAPI = "https://api.ssllabs.com/api/v3/analyze?host=%s"

type Respuesta struct {
	Estado     string     `json:"status"`
	Dominio    string     `json:"host"`
	Servidores []Servidor `json:"endpoints"`
}

type Servidor struct {
	IP   string `json:"ipAddress"`
	Nota string `json:"grade"`
}

func main() {
	http.HandleFunc("/", paginaPrincipal)
	http.HandleFunc("/analizar", analizarDominio)

	fmt.Println("Servidor iniciado en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func paginaPrincipal(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Analizador TLS</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f4f4f4;
				padding: 40px;
			}
			.contenedor {
				background: white;
				padding: 20px;
				max-width: 500px;
				margin: auto;
				border-radius: 5px;
			}
			input {
				width: 100%;
				padding: 8px;
				margin-bottom: 10px;
			}
			button {
				padding: 8px 15px;
			}
		</style>
	</head>
	<body>
		<div class="contenedor">
			<h2>Analizador TLS (SSL Labs)</h2>
			<form action="/analizar" method="post">
				<input type="text" name="dominio" placeholder="ejemplo.com" required>
				<button type="submit">Analizar</button>
			</form>
		</div>
	</body>
	</html>
	`
	fmt.Fprint(w, html)
}

func analizarDominio(w http.ResponseWriter, r *http.Request) {
	dominio := r.FormValue("dominio")
	url := fmt.Sprintf(urlAPI, dominio)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprint(w, "Error al consultar la API")
		return
	}
	defer resp.Body.Close()

	var respuesta Respuesta
	json.NewDecoder(resp.Body).Decode(&respuesta)

	fmt.Fprint(w, `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Resultado TLS</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f4f4f4;
				padding: 40px;
			}
			.contenedor {
				background: white;
				padding: 20px;
				max-width: 700px;
				margin: auto;
				border-radius: 5px;
			}
			table {
				border-collapse: collapse;
				width: 100%;
				margin-top: 15px;
			}
			th, td {
				border: 1px solid #ccc;
				padding: 8px;
				text-align: left;
			}
			th {
				background-color: #eee;
			}
			a {
				display: inline-block;
				margin-top: 15px;
			}
		</style>
	</head>
	<body>
		<div class="contenedor">
	`)

	fmt.Fprintf(w, "<h2>Resultado para %s</h2>", dominio)
	fmt.Fprintf(w, "<p><b>Estado:</b> %s</p>", respuesta.Estado)

	fmt.Fprint(w, `
	<table>
		<tr>
			<th>IP</th>
			<th>Nota TLS</th>
		</tr>
	`)

	for _, servidor := range respuesta.Servidores {
		fmt.Fprintf(w, `
		<tr>
			<td>%s</td>
			<td>%s</td>
		</tr>
		`, servidor.IP, servidor.Nota)
	}

	fmt.Fprint(w, `
	</table>
	<a href="/">Nueva consulta</a>
	</div>
	</body>
	</html>
	`)
}
