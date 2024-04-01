package permitions

import (
	"encoding/binary"
	"os"
	"proyecto/comandos/global"
	"proyecto/comandos/utils"
	"proyecto/estructuras/size"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func Values_CHMOD(instructions []string) (string, string, bool, bool) {
	var path, ugo string
	var r bool

	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "path") {
			var value = utils.TienePathFilePermitions("CHMOD", valor)
			path = value
		} else if strings.HasPrefix(strings.ToLower(valor), "ugo") {
			var value = utils.TieneUGO("CHMOD", valor)
			ugo = value
		} else if strings.HasPrefix(strings.ToLower(valor), "r") {
			var value = utils.TieneRPermitionsFile("CHMOD", valor)
			r = value
		} else {
			color.Yellow("[CHMOD]: Atributo no reconocido")
			return "", "", false, false
		}
	}
	if path == "" || len(path) == 0 {
		color.Red("[CHMOD]: No hay path")
		return "", "", false, false
	} else if ugo == "" || len(ugo) == 0 || len(ugo) > 3 {
		color.Red("[CHMOD]: No tiene ugo")
		return "", "", false, false
	}
	return path, ugo, r, true
}

func CHMOD_EXECUTE(path string, ugo string, r bool) {
	if !global.UsuarioLogeado.Logged_in {
		color.Red("[CHMOD]: No hay usuario logeado")
		return
	}
	if global.UsuarioLogeado.UID != 1 && global.UsuarioLogeado.GID != 1 {
		color.Red("[CHMOD]: Solo usuario root puede usar el comando")
		return
	}

	if len(ugo) == 1 {
		ugo = "00" + ugo
	} else if len(ugo) == 2 {
		ugo = "0" + ugo
	}

	nodo := global.UsuarioLogeado.Mounted
	file, err := os.OpenFile(nodo.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[CHMOD]: Error al abrir archivo")
		return
	}
	defer file.Close()

	var start int32 = 0
	if nodo.Es_Particion_L {
		start = nodo.Particion_L.Part_start + size.SizeEBR()
		if nodo.Particion_L.Part_mount != 1 {
			color.Red("[CHMOD]: El disco (logico) no se ha formateado -> " + utils.ToString(nodo.Particion_L.Name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
			return
		}
	} else if nodo.Es_Particion_P {
		start = nodo.Particion_P.Part_start
		if nodo.Particion_P.Part_status != 1 {
			color.Red("[CHMOD]: El disco (primario) no se ha formateado -> " + utils.ToString(nodo.Particion_P.Part_name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
			return
		}
	}
	if _, err := file.Seek(int64(start), 0); err != nil {
		color.Red("[CHMOD]: Error en mover puntero")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &utils.Sb_System); err != nil {
		color.Red("[CHMOD]: Error en la lectura del archivo")
		return
	}

	var posInodoF int32
	if path != "/" {
		rutaS := utils.SplitRuta(path)
		if len(rutaS) == 0 {
			color.Red("[CHMOD]: Ruta invalida")
			return
		}
		posInodoF = utils.GetInodoFSystem(rutaS, 0, int32(len(rutaS)-1), utils.Sb_System.S_inode_start, nodo.Path)
		if posInodoF == -1 {
			color.Red("[CHMOD]: Archivo no encontrado")
			return
		}
	} else {
		posInodoF = utils.Sb_System.S_inode_start
	}

	ugo_int, _ := strconv.Atoi(ugo)
	utils.ChmodR(posInodoF, int32(ugo_int), r)

	if utils.Sb_System.S_filesistem_type == 3 {
		utils.EscribirJournalSystem("chmod", '1', path, ugo, nodo)
	}

	color.Green("Se actualizaron los permisos de la ruta -> " + path)
}
