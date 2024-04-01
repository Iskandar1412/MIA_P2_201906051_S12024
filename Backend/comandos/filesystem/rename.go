package filesystem

import (
	"encoding/binary"
	"os"
	"proyecto/comandos/global"
	"proyecto/comandos/utils"
	"proyecto/estructuras/size"
	"proyecto/estructuras/structures"
	"strings"

	"github.com/fatih/color"
)

func Values_RENAME(instructions []string) (string, string, bool) {
	var path, name string

	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "path") {
			var value = utils.TienePathFilePermitions("RENAME", valor)
			path = value
		} else if strings.HasPrefix(strings.ToLower(valor), "name") {
			var value = utils.TieneNameRename("RENAME", valor)
			name = value
		} else {
			color.Yellow("[RENAME]: Atributo no reconocido")
			return "", "", false
		}
	}
	if path == "" || len(path) == 0 {
		color.Red("[RENAME]: No hay path")
		return "", "", false
	}
	if name == "" || len(name) == 0 {
		color.Red("[RENAME]: No hay name")
		return "", "", false
	}
	return path, name, true
}

func RENAME_EXECUTE(path string, name string) {
	if !global.UsuarioLogeado.Logged_in {
		color.Red("[RENAME]: No hay usuario logeado")
		return
	}

	nodo := global.UsuarioLogeado.Mounted
	file, err := os.OpenFile(nodo.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[RENAME]: Error al abrir archivo")
		return
	}
	defer file.Close()

	var start int32
	if nodo.Es_Particion_L {
		if nodo.Particion_L.Part_mount != 1 {
			color.Red("[RENAME]: El disco (logico) no se ha formateado -> " + utils.ToString(nodo.Particion_L.Name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
			return
		} else {
			start = nodo.Particion_L.Part_start + size.SizeEBR()
		}
	} else if nodo.Es_Particion_P {
		if nodo.Particion_P.Part_status != 1 {
			color.Red("[RENAME]: El disco (primario) no se ha formateado -> " + utils.ToString(nodo.Particion_P.Part_name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
			return
		} else {
			start = nodo.Particion_P.Part_start
		}
	}

	if _, err := file.Seek(int64(start), 0); err != nil {
		color.Red("[RENAME]: Error en mover puntero")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &utils.Sb_System); err != nil {
		color.Red("[RENAME]: Error en la lectura del archivo")
		return
	}

	rutaS := utils.SplitRuta(path)
	if len(rutaS) == 0 {
		color.Red("[RENAME]: Ruta invalida")
		return
	}

	posInodoI := utils.Sb_System.S_inode_start
	posInodoF := posInodoI
	var existP = true
	var inodo structures.TablaInodo

	for i := 0; i < len(rutaS); i++ {
		if existP {
			aux := posInodoI
			posInodoI = utils.ExistPathSystem(rutaS, int32(i), posInodoI, nodo.Path)
			if i != (len(rutaS) - 1) {
				posInodoF = posInodoI
			}
			if posInodoI == aux {
				existP = false
			}
		}
		if !existP {
			color.Red("[RENAME]: no se encontro el directorio -> " + path)
			return
		}
	}

	rutaS = append(rutaS, name)
	if posInodoF != utils.ExistPathSystem(rutaS, int32(len(rutaS)-1), posInodoF, nodo.Path) {
		color.Red("[RENAME]: El nombre esta en uso -> " + name)
		return
	}

	if _, err := file.Seek(int64(posInodoF), 0); err != nil {
		color.Red("[RENAME]: Error en mover puntero")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[RENAME]: Error en la lectura del archivo")
		return
	}

	if !utils.ValidarPermisoWSystem(posInodoI, nodo.Path) {
		color.Red("[RENAME]: No se puede reescribir el archivo -> " + path + " - por falta de permisos")
		return
	}

	utils.CambiarNombre(inodo, rutaS[len(rutaS)-2], rutaS[len(rutaS)-1])
	inodo.I_mtime = utils.ObFechaInt()

	if _, err := file.Seek(int64(posInodoF), 0); err != nil {
		color.Red("[RENAME]: Error en mover puntero")
		return
	}

	if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[RENAME]: Error en la escritura del archivo")
		return
	}

	if utils.Sb_System.S_filesistem_type == 3 {
		utils.EscribirJournalSystem("rename", '1', path, name, nodo)
	}

	color.Green("[RENAME]: Se cambio nombre del archivo -> " + path + " -->> " + name)
}
