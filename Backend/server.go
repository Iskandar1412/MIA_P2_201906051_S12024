package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func server() {
	http.HandleFunc("/command", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			// Establecer encabezados CORS para las solicitudes OPTIONS
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Manejar las solicitudes POST
		if r.Method != http.MethodPost {
			http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Obtener estructura
		var requestBody struct {
			Comando string `json:"comando"`
		}

		//fmt.Println((r.Body))
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, "Error al leer cuerpo de la solicitud", http.StatusBadRequest)
			return
		}

		fmt.Println(requestBody.Comando)

	})

	http.HandleFunc("/obtain-mbr", func(w http.ResponseWriter, r *http.Request) {
		//leerr mbr
		mbr, emb := Retorno_MBR()
		if emb != nil {
			fmt.Println("Error")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mbr)
	})

	fmt.Println("Backend server is on 8080")
	http.ListenAndServe(":8080", nil)
}
