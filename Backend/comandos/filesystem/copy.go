package filesystem

import (
	"proyecto/comandos/utils"
	"strings"

	"github.com/fatih/color"
)

func Values_COPY(instructions []string) (string, string, bool) {
	var path, destino string

	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "path") {
			var value = utils.TienePathFilePermitions("COPY", valor)
			path = value
		} else if strings.HasPrefix(strings.ToLower(valor), "destino") {
			var value = utils.TieneDestinoFile("COPY", valor)
			destino = value
		} else {
			color.Yellow("[COPY]: Atributo no reconocido")
			return "", "", false
		}
	}
	if path == "" || len(path) == 0 {
		color.Red("[COPY]: No hay path")
		return "", "", false
	}
	if destino == "" || len(destino) == 0 {
		color.Red("[COPY]: No hay path")
		return "", "", false
	}
	return path, destino, true
}

func COPY_EXECUTE(path string, destino string) {

}
