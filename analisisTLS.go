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
	http.HandleFunc("/", mostrarFormulario)
	http.HandleFunc("/analizar", analizarDominio)

	fmt.Println("Servidor iniciado en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func mostrarFormulario(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
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

	fmt.Fprint(w, "<h2>Resultado</h2>")
	fmt.Fprintf(w, "<p><b>Dominio:</b> %s</p>", dominio)
	fmt.Fprintf(w, "<p><b>Estado:</b> %s</p>", respuesta.Estado)

	fmt.Fprint(w, "<table border='1'>")
	fmt.Fprint(w, "<tr><th>IP</th><th>Nota TLS</th></tr>")

	for _, servidor := range respuesta.Servidores {
		fmt.Fprintf(w,
			"<tr><td>%s</td><td>%s</td></tr>",
			servidor.IP, servidor.Nota)
	}

	fmt.Fprint(w, "</table>")
	fmt.Fprint(w, "<br><a href='/'>Volver</a>")
}
