package admindisk

import (
	"encoding/binary"
	"fmt"
	"os"
	"proyecto/comandos/utils"
	"proyecto/estructuras/structures"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func Values_MKDISK(instructions []string) (int32, byte, byte) {
	var _size int32
	var _fit byte = 'F'
	var _unit byte = 'M'
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "size") {
			var value = utils.TieneSize("MKDISK", valor)
			_size = value
		} else if strings.HasPrefix(strings.ToLower(valor), "fit") {
			var value = utils.TieneFit("MKDISK", valor)
			_fit = value
		} else if strings.HasPrefix(strings.ToLower(valor), "unit") {
			var value = utils.TieneUnit("mkdisk", valor)
			_unit = value
		} else {
			color.Yellow("[MKDISK]: Atributo no reconocido")
			return -1, '0', '0'
		}
	}
	return _size, _fit, _unit
}

func MKDISK_Create(_size int32, _fit byte, _unit byte) {
	directorio := "MIA/P1/Disks/"
	for i := 0; i < 26; i++ {
		letra := fmt.Sprintf("%c.dsk", 'A'+i)
		archivo := directorio + letra
		if _, err := os.Stat(archivo); os.IsNotExist(err) {
			CreateFile(archivo, _size, _fit, _unit)
			color.Green("[MKDISK]: Disco '" + letra + "' Creado -> " + strconv.Itoa(int(_size)) + string(_unit))
			break
		} else {
			continue
			//color.Yellow("[MKDISK]: Disco Existente")
		}
	}
}

func CreateFile(archivo string, _size int32, _fit byte, _unit byte) {
	file, err := os.Create(archivo)
	if err != nil {
		color.Red("Error al crear el archivo")
		return
	}
	defer file.Close()
	//Escribir datos en archivo
	var estructura structures.MBR
	tamanio := utils.Tamano(_size, _unit)
	estructura.Mbr_tamano = tamanio
	estructura.Mbr_fecha_creacion = utils.ObFechaInt()
	estructura.Mbr_disk_signature = utils.ObDiskSignature()
	estructura.Dsk_fit = _fit
	for i := 0; i < len(estructura.Mbr_partitions); i++ {
		estructura.Mbr_partitions[i] = utils.PartitionVacia()
	}
	//Llenar el archivo
	bytes_llenar := make([]byte, int(tamanio))
	if _, err := file.Write(bytes_llenar); err != nil {
		color.Red("Error al escribir bytes en el archivo")
		return
	}

	//mover de posicion
	if _, err := file.Seek(0, 0); err != nil {
		color.Red("Error al mover puntero del archivo")
		return
	}

	if err := binary.Write(file, binary.LittleEndian, &estructura); err != nil {
		color.Red("Error al escribir datos del MBR")
		return
	}
}
