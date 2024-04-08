package main

//go mod init <paquete>
//go mod tidy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"proyecto/comandos/comandos"
	"proyecto/comandos/global"
	obtencionpf "proyecto/comandos/obtencion-pf"
	"proyecto/comandos/utils"
	"proyecto/estructuras/structures"

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
	ObtenerMBR_Mounted()
	MostrarParticionesMontadas()
	Retorno_MBR()

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

func ObtenerMBR_s() ([]structures.MBR, bool) {
	directorio := "MIA/P1/Disks/"
	var mbrsObtenidos = []structures.MBR{}
	isEmpty, err := isDirEmpty(directorio)
	if err != nil {
		color.Red("Error al verificar directorio")
		return []structures.MBR{}, false
	}

	if isEmpty {
		//caso en el que esta vacio
		return []structures.MBR{}, false
	} else {
		for i := 0; i < 26; i++ {
			letra := fmt.Sprintf("%c.dsk", 'A'+i)
			archivo := directorio + letra
			if _, err := os.Stat(archivo); os.IsNotExist(err) {
				break
			} else {
				mbr, embr := utils.Obtener_FULL_MBR_FDISK(archivo)
				if !embr {
					return []structures.MBR{}, false
				}
				mbrsObtenidos = append(mbrsObtenidos, mbr)
				continue
				//color.Yellow("[MKDISK]: Disco Existente")
			}
		}
		return mbrsObtenidos, true
	}
}

func ObtenerMBR_Mounted() ([]structures.MBR, bool) {
	directorio := "MIA/P1/Disks/"
	var mbrsObtenidos = []structures.MBR{}
	isEmpty, err := isDirEmpty(directorio)
	if err != nil {
		color.Red("Error al verificar directorio")
		return []structures.MBR{}, false
	}

	if isEmpty {
		//caso en el que esta vacio
		fmt.Println("No hay discos")
		return []structures.MBR{}, false
	} else {
		for i := 0; i < 26; i++ {
			letra := fmt.Sprintf("%c.dsk", 'A'+i)
			archivo := directorio + letra
			if _, err := os.Stat(archivo); os.IsNotExist(err) {
				break
			} else {
				mbr, embr := utils.Obtener_FULL_MBR_FDISK(archivo)
				if !embr {
					return []structures.MBR{}, false
				}
				for _, _content := range mbr.Mbr_partitions {
					if utils.ToString(_content.Part_id[:]) != "" {

						nuevoMBR := global.ParticionesMontadas{}
						nuevoMBR.DriveLetter = letra[0]
						nuevoMBR.Es_Particion_P = true
						nuevoMBR.Es_Particion_L = false
						nuevoMBR.Particion_P = _content
						nuevoMBR.Path = archivo
						nuevoMBR.ID_Particion = _content.Part_id
						nuevoMBR.Type = 'P'

						global.Mounted_Partitions = append(global.Mounted_Partitions, nuevoMBR)
					}
					continue
				}

				temp_ext := structures.MBR_Obtener{}

				temp_ext.Disco = string(letra[0]) + ".dsk"
				temp_ext.Disco_Path = archivo
				temp_ext.Mbr_partitions[0].Id_mounted = utils.ToString(mbr.Mbr_partitions[0].Part_id[:])
				temp_ext.Mbr_partitions[1].Id_mounted = utils.ToString(mbr.Mbr_partitions[1].Part_id[:])
				temp_ext.Mbr_partitions[2].Id_mounted = utils.ToString(mbr.Mbr_partitions[2].Part_id[:])
				temp_ext.Mbr_partitions[3].Id_mounted = utils.ToString(mbr.Mbr_partitions[3].Part_id[:])
				temp_ext.Mbr_partitions[0].Particion = utils.ToString(mbr.Mbr_partitions[0].Part_name[:])
				temp_ext.Mbr_partitions[1].Particion = utils.ToString(mbr.Mbr_partitions[1].Part_name[:])
				temp_ext.Mbr_partitions[2].Particion = utils.ToString(mbr.Mbr_partitions[2].Part_name[:])
				temp_ext.Mbr_partitions[3].Particion = utils.ToString(mbr.Mbr_partitions[3].Part_name[:])
				temp_ext.Mbr_partitions[0].Status = mbr.Mbr_partitions[0].Part_status
				temp_ext.Mbr_partitions[1].Status = mbr.Mbr_partitions[1].Part_status
				temp_ext.Mbr_partitions[2].Status = mbr.Mbr_partitions[2].Part_status
				temp_ext.Mbr_partitions[3].Status = mbr.Mbr_partitions[3].Part_status
				temp_ext.Mbr_partitions[0].Type = string(mbr.Mbr_partitions[0].Part_type)
				temp_ext.Mbr_partitions[1].Type = string(mbr.Mbr_partitions[1].Part_type)
				temp_ext.Mbr_partitions[2].Type = string(mbr.Mbr_partitions[2].Part_type)
				temp_ext.Mbr_partitions[3].Type = string(mbr.Mbr_partitions[3].Part_type)
				ObtenerEstructuras = append(ObtenerEstructuras, temp_ext)
				continue
				//color.Yellow("[MKDISK]: Disco Existente")
			}
		}
		return mbrsObtenidos, true
	}
}

func MostrarParticionesMontadas() {
	if len(global.Mounted_Partitions) > 0 {
		for _, i := range global.Mounted_Partitions {
			fmt.Println("Disco:", string(i.DriveLetter), " - Particion:", utils.ToString(i.Particion_P.Part_name[:]), " (id): --> "+utils.ToString(i.Particion_P.Part_id[:]))
		}
	}
}

func isDirEmpty(dirname string) (bool, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // trata de leer algun nombre de directorio
	if err == nil {
		// existe algun archivo en el directorio
		return false, nil
	}

	// directorio vacio
	if err == io.EOF {
		return true, nil
	}

	return false, err
}

var ObtenerEstructuras []structures.MBR_Obtener

func Retorno_MBR() ([]byte, error) {
	//fmt.Println(ObtenerEstructuras)

	jsonData, err := json.Marshal(ObtenerEstructuras)
	if err != nil {
		return []byte{}, fmt.Errorf("error")
	}

	fmt.Println(string(jsonData))
	return jsonData, nil
	//return "", fmt.Errorf("error")
}

func Retorno_Paths() ([]byte, error) {
	log := utils.ToString(global.UsuarioLogeado.Mounted.ID_Particion[:])
	fmt.Println(log)
	ino, efno := obtencionpf.ObtenerPaths("A151")
	if efno != nil {
		return []byte{}, fmt.Errorf("Error")
	}
	fmt.Println(ino)
	return []byte(ino), nil
}
