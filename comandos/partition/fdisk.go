package partition

import (
	"proyecto/comandos/utils"
	"strings"

	"github.com/fatih/color"
)

func Values_FDISK(instructions []string) (int32, byte, [16]byte, byte, byte, byte, string, int32) {
	var _size int32
	var _driveletter byte
	var _name [16]byte
	var _unit byte = 'K'
	var _type byte = 'P'
	var _fit byte = 'W'
	var _delete string = "None"
	var _add int32 = 0
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "size") {
			var value = utils.TieneSize("FDISK", valor)
			_size = value
		} else if strings.HasPrefix(strings.ToLower(valor), "driveletter") {
			var value = utils.TieneDriveLetter("FDISK", valor)
			_driveletter = value
		} else if strings.HasPrefix(strings.ToLower(valor), "name") {
			var value = utils.TieneNombre("FDISK", valor)
			if len(value) > 16 {
				color.Red("[FDISK]: El nombre no puede ser mayor a 16 caracteres")
				break
			} else {
				_name = utils.DevolverNombreByte(value)
			}
		} else if strings.HasPrefix(strings.ToLower(valor), "unit") {
			var value = utils.TieneUnit("FDISK", valor)
			_unit = value
		} else if strings.HasPrefix(strings.ToLower(valor), "type") {
			var value = utils.TieneTypeFDISK(valor)
			_type = value
		} else if strings.HasPrefix(strings.ToLower(valor), "fit") {
			var value = utils.TieneFit("FDISK", valor)
			_fit = value
		} else if strings.HasPrefix(strings.ToLower(valor), "delete") {
			var value = utils.TieneDelete(valor)
			_delete = value
		} else if strings.HasPrefix(strings.ToLower(valor), "add") {
			var value = utils.TieneAdd(valor)
			_add = value
		} else {
			color.Yellow("[FDISK]: Atributo no reconocido")
		}
	}
	return _size, _driveletter, _name, _unit, _type, _fit, _delete, _add
}

func FDISK_Create(_size int32, _driveletter byte, _name []byte, _unit byte, _type byte, _fit byte, _delete string, _add int32) {
	//fmt.Println(_name)
	path := "MIA/P1/Disks/" + string(_driveletter) + ".dsk"
	if !utils.ExisteArchivo("FDISK", path) {
		color.Yellow("[FDISK] Cancel the operation because not yet a file")
		return
	}

	// Delete
	if _delete != "None" {
		// Borrar particiones
		DeleteP(path, _name, _unit, _type, _fit)
		return
	}

	// Add
	if _add != 0 {
		if _add < 0 {
			RestE(path, _unit, _fit, _add, _name)
			return
		} else if _add > 0 {
			AddE(path, _unit, _fit, _add, _name)
			// fmt.Println("sumando")
			return
		}
	}

	if _type == 'P' {
		//primaria
		ParticionPrimaria(_size, _driveletter, _name, _unit, _type, _fit, _delete, _add)
	} else if _type == 'E' {
		//extended
		ParticionExtendida(_size, _driveletter, _name, _unit, _type, _fit, _delete, _add)
	} else if _type == 'L' {
		//logic
		ParticionLogica(_size, _driveletter, _name, _unit, _type, _fit, _delete, _add)
	} else {
		color.Red("[FDISK]: No reconocido Type")
		return
	}
}
