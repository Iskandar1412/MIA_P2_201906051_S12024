package obtencionpf

import (
	"bufio"
	"fmt"
	"os"
	"proyecto/comandos/global"
	"proyecto/comandos/utils"
	"proyecto/estructuras/structures"
	"strings"

	"github.com/fatih/color"
)

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
	}

	dir, err := os.Open(directorio)
	if err != nil {
		return []structures.MBR{}, false
	}
	defer dir.Close()

	filename, err := dir.Readdirnames(-1)
	if err != nil {
		return []structures.MBR{}, false
	}
	mbrsObtenidos = []structures.MBR{}
	for _, discos := range filename {
		fulpath := directorio + discos
		mbr, embr := utils.Obtener_FULL_MBR_FDISK(fulpath)
		if !embr {
			return []structures.MBR{}, false
		}

		mbrsObtenidos = append(mbrsObtenidos, mbr)
	}

	return mbrsObtenidos, true

}

func MostrarParticionesMontadas() {
	if len(global.Mounted_Partitions) > 0 {
		for _, i := range global.Mounted_Partitions {
			fmt.Println("Disco:", string(i.DriveLetter), " - Particion:", utils.ToString(i.Particion_P.Part_name[:]), " (id): --> "+utils.ToString(i.Particion_P.Part_id[:]))
		}
	}
}

func ObtenerReportes() ([]byte, error) {
	path := "MIA/P1/Reports/"

	isEmpty, err := isDirEmpty(path)
	if err != nil {
		color.Red("Error al verificar directorio")
		return []byte{}, fmt.Errorf("Error")
	}

	if isEmpty {
		return []byte{}, fmt.Errorf("Error")
	}

	dir, err := os.Open(path)
	if err != nil {
		return []byte{}, fmt.Errorf("Error")
	}
	defer dir.Close()

	filename, err := dir.Readdirnames(-1)
	if err != nil {
		return []byte{}, fmt.Errorf("Error")
	}

	cont := "["
	for _, dir := range filename {
		dir2 := strings.Split(dir, ".")

		cont += "{"
		cont += "\"dot\":\"" + dir + "\","
		cont += "\"path\":\"" + path + dir + "\","
		cont += "\"extension\":\"" + dir2[1] + "\""
		cont += "}"
	}
	cont += "]"
	cont = strings.ReplaceAll(cont, "}{", "},{")
	// fmt.Println(cont)
	return []byte(cont), nil
}

func ObtenerDot(path string) string {
	dot := ""
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		dot += line
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return ""
	}
	// fmt.Println(dot)
	return dot
}
