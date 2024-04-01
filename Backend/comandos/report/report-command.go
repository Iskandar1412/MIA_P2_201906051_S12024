package report

import (
	"os"
	"proyecto/comandos/utils"
	"strings"

	"github.com/fatih/color"
)

func ReportCommandProps(command string, instructions []string) {
	// path := "MIA/P1/Disks/"
	var _name string
	var _path string
	var _id string
	var _ruta string
	var er bool
	/*
	 */
	if strings.ToUpper(command) == "REP" {
		//fmt.Println("GENERANDO REPORTE")
		_name, _path, _id, _ruta, er = Sub_Reports(instructions)
		if !er {
			color.Red("[REP]: Error to asign values")
		} else {
			REP_EXECUTE(_name, _path, _id, _ruta)
		}
	} else {
		color.Red("[ReportCommandProps]: Internal Error")
	}
}

func Sub_Reports(instructions []string) (string, string, string, string, bool) {
	var _name string
	var _path string
	var _id string
	var _ruta string
	var er bool
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "name") {
			_name, er = utils.TieneNameRep(valor)
			if !er {
				color.Red("[REP]: Error al asignar name")
				return "", "", "", "", false
			}
		} else if strings.HasPrefix(strings.ToLower(valor), "path") {
			_path, er = utils.TienePathRep(valor)
			if !er {
				color.Red("[REP]: Error al asignar path")
				return "", "", "", "", false
			}
		} else if strings.HasPrefix(strings.ToLower(valor), "id") {
			_id, er = utils.TieneIDRep(valor)
			if !er {
				color.Red("[REP]: Error al asignar id")
				return "", "", "", "", false
			}
		} else if strings.HasPrefix(strings.ToLower(valor), "ruta") {
			_ruta, er = utils.TieneRutaRep(valor)
			if !er {
				color.Red("[REP]: Error al asignar id")
				return "", "", "", "", false
			}
		}
	}
	if _id == "" || _path == "" || _name == "" {
		// color.Red("[REP]: Valores no Asignados")
		return "", "", "", "", false
	}
	return _name, _path, _id, _ruta, true
}

func REP_EXECUTE(_name string, _path string, _id string, _ruta string) {

	ruta_separada := strings.Split(_path, "/")
	cantidad_carpetas := len(ruta_separada)
	nombre_archivo := ruta_separada[cantidad_carpetas-1]
	ruta_sin_archivo := strings.ReplaceAll(_path, "/"+nombre_archivo, "")

	if _, ecarpeta := os.Stat(ruta_sin_archivo); os.IsNotExist(ecarpeta) {
		err := os.MkdirAll(ruta_sin_archivo, 0777)
		if err != nil {
			color.Red("Error al crear carpeta")
			return
		}

		color.Green("Carpeta: " + ruta_sin_archivo + " --> Creada")

	} else {
		color.Yellow("Carpeta: " + ruta_sin_archivo + " --> Existente")
	}

	if _name == "mbr" {
		Report_MBR(nombre_archivo, ruta_sin_archivo, _ruta, _id)
	} else if _name == "disk" {
		Report_DISK(nombre_archivo, ruta_sin_archivo, _ruta, _id)
	} else if _name == "inode" {
		Report_INODE(nombre_archivo, ruta_sin_archivo, _ruta, _id)
	} else if _name == "journaling" {
		Report_Journal(nombre_archivo, ruta_sin_archivo, _ruta, _id)
	} else if _name == "block" {
		Report_BLOCK(nombre_archivo, ruta_sin_archivo, _ruta, _id)
	} else if _name == "bm_inode" {
		Report_BM_Inode(nombre_archivo, ruta_sin_archivo, _ruta, _id)
	} else if _name == "bm_block" {
		Report_BM_Block(nombre_archivo, ruta_sin_archivo, _ruta, _id)
	} else if _name == "tree" {
		Report_TREE(nombre_archivo, ruta_sin_archivo, _ruta, _id)
	} else if _name == "sb" {
		Report_SUPERBLOCK(nombre_archivo, ruta_sin_archivo, _ruta, _id)
	} else if _name == "file" {
		Report_FILE(nombre_archivo, ruta_sin_archivo, _ruta, _id)
	} else if _name == "ls" {
		Report_LS(nombre_archivo, ruta_sin_archivo, _ruta, _id)
	} else {
		color.Red("[REP]: Internal Error")
		return
	}
}
