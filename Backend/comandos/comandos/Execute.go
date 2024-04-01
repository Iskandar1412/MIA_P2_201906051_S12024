package comandos

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

func Execute(x []string) []string {
	for _, y := range x {
		var path string
		if strings.HasPrefix(strings.ToLower(y), "path") {
			path = TienePath(y)
		} else {
			y := strings.Split(y, "=")
			color.Red("[EXECUTE] ( \"" + y[0] + "\" ): Comando no reconocido")
			break
		}
		if path == "nil" {
			return nil
		} else {
			return ExecuteFunc(path)
		}
	}
	return nil
}

func ExecuteFunc(x string) []string {
	file, err := os.Open(x)
	if err != nil {
		color.Red("Error al abrir archivo", err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lineas []string

	for scanner.Scan() {
		linea := strings.TrimSpace(scanner.Text())
		if len(linea) > 0 && !strings.HasPrefix(linea, "#") {
			lineas = append(lineas, linea)
		}
	}

	var exportar []string
	reg := regexp.MustCompile(`(.*?)\s*(?:#.*|$)`)
	for _, y := range lineas {
		match := reg.FindStringSubmatch(y)
		//fmt.Println(y, "asdf")
		if len(match) > 1 {
			exportar = append(exportar, match[1])
			//fmt.Println(match[0], "///", match[1])
		}
	}
	//fmt.Println(exportar)
	if err := scanner.Err(); err != nil {
		color.Red("Error en la lectura del archivo:", err)
		return nil
	}

	return exportar
}
