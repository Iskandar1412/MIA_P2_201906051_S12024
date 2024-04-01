package admindisk

import (
	"os"
	"proyecto/comandos/utils"
	"strings"

	"github.com/fatih/color"
)

func Values_RMDISK(instructions []string) (byte, bool) {
	var _driveletter byte = '0'
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "driveletter") {
			var value = utils.TieneDriveLetter("RMDISK", valor)
			_driveletter = value
			break
		} else {
			color.Yellow("[RMDISK]: Atributo no reconocido")
			_driveletter = '0'
			break
		}
	}
	if _driveletter == '0' {
		return '0', false
	} else {
		return _driveletter, true
	}
}

func RMDISK_EXECUTE(_driveletter byte) {
	PATH := "MIA/P1/Disks/" + string(_driveletter) + ".dsk"
	if _, err := os.Stat(PATH); os.IsNotExist(err) {
		color.Red("[RMDISK]: No existe el disco")
		return
	}
	err := os.Remove(PATH)
	if err != nil {
		color.Red("[RMDISK]: Error al borrar el disco")
		return
	}
	color.Green("[RMDISK]: Disco '" + string(_driveletter) + ".dsk' Borrado")
}
