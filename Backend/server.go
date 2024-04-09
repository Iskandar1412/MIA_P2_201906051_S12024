package main

//para reiniciar servidor automaticamente
//go get -u github.com/gin-gonic/gin
//gin run <archivo>.go

import (
	"encoding/json"
	"fmt"
	"net/http"
	"proyecto/comandos/comandos"
	"proyecto/comandos/obtencionpf"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()

	// Configurar CORS
	c := cors.AllowAll()

	// Manejar las rutas
	mux.HandleFunc("/command", handleCommand)
	mux.HandleFunc("/login", handleLogin)
	mux.HandleFunc("/logout", handleLogout)
	mux.HandleFunc("/obtainmbr", handleObtainMBR)
	mux.HandleFunc("/reportesobtener", handleReportsObtener)
	mux.HandleFunc("/graphs", handleGraph)
	mux.HandleFunc("/obtain-carpetas-archivos", handleObtainCarpetasArchivos)

	handler := c.Handler(mux)

	fmt.Println("Backend server is on 8080")
	comandos.CrearCarpeta()
	obtencionpf.ObtenerMBR_Mounted()
	obtencionpf.MostrarParticionesMontadas()
	http.ListenAndServe(":8080", handler)
}

func handleCommand(w http.ResponseWriter, r *http.Request) {
	// command
	if r.Method == http.MethodOptions {
		// Establecer encabezados CORS para las solicitudes OPTIONS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Solo permitir solicitudes POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Configurar encabezados CORS para las solicitudes POST
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Decodificar el cuerpo JSON de la solicitud
	var requestBody struct {
		Comando string `json:"comando"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	var ejecutar []string
	ejecutar = append(ejecutar, requestBody.Comando)
	comandos.GlobalCom(ejecutar)
	obtencionpf.ObtenerMBR_Mounted()
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(requestBody.Comando)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	// login
	if r.Method == http.MethodOptions {
		// Establecer encabezados CORS para las solicitudes OPTIONS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Solo permitir solicitudes POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Configurar encabezados CORS para las solicitudes POST
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var requestBody struct {
		Comando string `json:"comando"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	var ejecutar []string
	ejecutar = append(ejecutar, requestBody.Comando)
	if !comandos.GlobalCom(ejecutar) {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	obtencionpf.ObtenerMBR_Mounted()
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(requestBody.Comando)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	// logout
	if r.Method == http.MethodOptions {
		// Establecer encabezados CORS para las solicitudes OPTIONS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Solo permitir solicitudes POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Configurar encabezados CORS para las solicitudes POST
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var requestBody struct {
		Comando string `json:"comando"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	var ejecutar []string
	ejecutar = append(ejecutar, requestBody.Comando)
	// fmt.Println(ejecutar)
	if !comandos.GlobalCom(ejecutar) {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	obtencionpf.ObtenerMBR_Mounted()
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(requestBody.Comando)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
}

func handleObtainMBR(w http.ResponseWriter, r *http.Request) {
	// obtainmbr
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	mbr, emb := obtencionpf.Retorno_MBR()
	if emb != nil {
		return
	}

	//fmt.Println(string(mbr))
	data := map[string]interface{}{
		"datos": string(mbr),
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
}

func handleGraph(w http.ResponseWriter, r *http.Request) {
	// Graphs
	queryValues := r.URL.Query()

	// verificacion de id
	id := queryValues.Get("id")
	if id == "" {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	path := "MIA/P1/Reports/" + id
	dot_Obtenido := obtencionpf.ObtenerDot(path)
	if dot_Obtenido == "" {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"datos": dot_Obtenido,
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
}

func handleReportsObtener(w http.ResponseWriter, r *http.Request) {
	// reportesobtener
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	rep, emb := obtencionpf.ObtenerReportes()
	if emb != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	//fmt.Println(string(rep))
	data := map[string]interface{}{
		"datos": string(rep),
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
}

func handleObtainCarpetasArchivos(w http.ResponseWriter, r *http.Request) {
	// obtain-carpetas-archivos
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	path, error := obtencionpf.Retorno_Paths()
	if error != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// fmt.Println(string(path))
	data := map[string]interface{}{
		"datos": string(path),
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
}

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
//-----------------------------FUNCIONES PROYECTO-----------------------------
//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

// *****************************************************************************
// *****************************************************************************
// *****************************************************************************
// ----------------------------------OBTENCION----------------------------------
// *****************************************************************************
// *****************************************************************************
// *****************************************************************************

// *****************************************************************************
// *****************************************************************************
// *****************************************************************************
// ----------------------------------OBTENCION----------------------------------
// *****************************************************************************
// *****************************************************************************
// *****************************************************************************

//##################################################################################
//##################################################################################
//##################################################################################
//---------------------------------OBTENER CARPETAS---------------------------------
//##################################################################################
//##################################################################################
//##################################################################################

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
//-----------------------------FUNCIONES PROYECTO-----------------------------
//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
