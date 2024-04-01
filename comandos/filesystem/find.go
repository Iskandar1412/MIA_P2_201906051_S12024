package filesystem

import (
	"proyecto/comandos/utils"
	"strings"

	"github.com/fatih/color"
)

func Values_FIND(instructions []string) (string, string, bool) {
	var path, name string

	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "path") {
			var value = utils.TienePathFilePermitions("FIND", valor)
			path = value
		} else if strings.HasPrefix(strings.ToLower(valor), "name") {
			var value = utils.TieneNameRename("FIND", valor)
			name = value
		} else {
			color.Yellow("[FIND]: Atributo no reconocido")
			return "", "", false
		}
	}
	if path == "" || len(path) == 0 {
		color.Red("[FIND]: No hay path")
		return "", "", false
	}
	if name == "" || len(name) == 0 {
		color.Red("[FIND]: No hay name")
		return "", "", false
	}
	return path, name, true
}

func FIND_EXECUTE(path string, name string) {

}
