package main

//go mod init <paquete>
//go mod tidy

import (
	"bufio"
	"fmt"
	"os"
	"proyecto/comandos/comandos"

	"strings"

	"github.com/fatih/color"
)

// var Partitions_Mounted []interface{}

// MIA_P1_201906051/structures
// MIA_P1_201906051/size
// execute -path=/home/iskandar/Escritorio/prueba.asdj
// execute -path=/home/iskandar/Escritorio/Proyectos/Git/MIA_1S2024_201906051/Proyectos/MIA_P1_201906051/Pruebas/prueba.adsj
// execute -path=Pruebas/prueba.adsj
func main() {
	color.Blue("PROY1 - 201906051 - Juan Urbina")
	comandos.CrearCarpeta()
	for {
		// execute - path = Pruebas / avanzado.asdj
		// input := "execute -path=Pruebas/prueba.adsj"
		// input := "execute -path=Pruebas/basico.asdj"
		// input := "execute -path=Pruebas/avanzado.asdj"
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Ingresar comando EXECUTE (exit para salir): >")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if strings.ToLower(input) == "exit" {
			color.Cyan("Saliendo del programa")
			//color.Cyan(particiones_montadas)
			break
		} else {
			instrucciones := comandos.ObtenerComandos(input)
			if strings.HasPrefix(strings.ToLower(input), "execute") {
				//fmt.Println(instrucciones)
				ejecutar := comandos.Execute(instrucciones)
				comandos.GlobalCom(ejecutar)
			} else {
				//fmt.Println("comendo erroneo")
				var ejecutar []string
				ejecutar = append(ejecutar, input)
				comandos.GlobalCom(ejecutar)
			}
			//fmt.Println("instruciones", instrucciones)
			//fmt.Println(len(instrucciones))
		}
	}
}
