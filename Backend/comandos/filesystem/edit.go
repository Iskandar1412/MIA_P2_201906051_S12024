package filesystem

import (
	"proyecto/comandos/utils"
	"strings"

	"github.com/fatih/color"
)

func Values_EDIT(instructions []string) (string, string, bool) {
	var path, cont string

	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "path") {
			var value = utils.TienePathFilePermitions("EDIT", valor)
			path = value
		} else if strings.HasPrefix(strings.ToLower(valor), "cont") {
			var value = utils.TieneContFile("EDIT", valor)
			cont = value
		} else {
			color.Yellow("[EDIT]: Atributo no reconocido")
			return "", "", false
		}
	}
	if path == "" || len(path) == 0 {
		color.Red("[EDIT]: No hay path")
		return "", "", false
	}
	if cont == "" || len(cont) == 0 {
		color.Red("[EDIT]: No hay cont")
		return "", "", false
	}
	return path, cont, true
}

func EDIT_EXECUTE(path string, cont string) {

}
