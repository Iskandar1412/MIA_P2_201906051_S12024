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

func Values_CHOWN(instructions []string) (string, string, bool, bool) {
	var path, user string
	var r bool

	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "path") {
			var value = utils.TienePathFilePermitions("CHOWN", valor)
			path = value
		} else if strings.HasPrefix(strings.ToLower(valor), "user") {
			var value = utils.TieneUser("CHOWN", valor)
			user = value
		} else if strings.HasPrefix(strings.ToLower(valor), "r") {
			var value = utils.TieneRPermitionsFile("CHOWN", valor)
			r = value
		} else {
			color.Yellow("[CHOWN]: Atributo no reconocido")
			return "", "", false, false
		}
	}
	if path == "" || len(path) == 0 {
		color.Red("[CHOWN]: No hay path")
		return "", "", false, false
	} else if user == "" || len(user) == 0 {
		color.Red("[CHOWN]: No tiene user")
		return "", "", false, false
	}
	return path, user, r, true
}

func CHOWN_EXECUTE(path string, user string, r bool) {
	if !global.UsuarioLogeado.Logged_in {
		color.Red("[CHOWN]: No hay usuario logeado")
		return
	}

	nodo := global.UsuarioLogeado.Mounted
	file, err := os.OpenFile(nodo.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[CHOWN]: Error al abrir archivo")
		return
	}
	defer file.Close()

	var start int32 = 0
	if nodo.Es_Particion_L {
		start = nodo.Particion_L.Part_start + size.SizeEBR()
		if nodo.Particion_L.Part_mount != 1 {
			color.Red("[CHOWN]: El disco (logico) no se ha formateado -> " + utils.ToString(nodo.Particion_L.Name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
			return
		}
	} else if nodo.Es_Particion_P {
		start = nodo.Particion_P.Part_start
		if nodo.Particion_P.Part_status != 1 {
			color.Red("[CHOWN]: El disco (primario) no se ha formateado -> " + utils.ToString(nodo.Particion_P.Part_name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
			return
		}
	}
	if _, err := file.Seek(int64(start), 0); err != nil {
		color.Red("[CHOWN]: Error en mover puntero")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &utils.Sb_System); err != nil {
		color.Red("[CHOWN]: Error en la lectura del archivo")
		return
	}

	content := utils.GetContentAdminUsers(utils.Sb_System.S_inode_start + size.SizeTablaInodo())
	split_content := strings.Split(content, "\n")
	if !utils.UsrExist(split_content, user) {
		color.Red("[CHOWN]: Usuario <" + user + "> no existente")
		return
	}

	id_usr := GetUsrSystem(split_content, user)
	id_grp := GetGrpSystem(split_content, user)
	if id_grp == -1 && id_usr == -1 {
		return
	}

	var posInodoF int32
	if path != "/" {
		rutaS := utils.SplitRuta(path)
		if len(rutaS) == 0 {
			color.Red("[CHOWN]: Ruta invalida")
			return
		}
		posInodoF = utils.GetInodoFSystem(rutaS, 0, int32(len(rutaS)-1), utils.Sb_System.S_inode_start, nodo.Path)
		if posInodoF == -1 {
			color.Red("[CHOWN]: Archivo no encontrado")
			return
		}
	} else {
		posInodoF = utils.Sb_System.S_inode_start
	}

	utils.ChownSystem(posInodoF, id_usr, id_grp, r)
	if utils.Sb_System.S_filesistem_type == 3 {
		utils.EscribirJournalSystem("chown", '1', path, user, nodo)
	}
	color.Green("[CHOWN]: Se actualizo propietario de la carpeta -> " + path)
}

func GetUsrSystem(content []string, usr string) int32 {
	for _, xy := range content {
		if strings.Contains(xy, ",U,") {
			if strings.Contains(xy, usr) {
				split := strings.Split(xy, ",")
				idusr, _ := strconv.Atoi(split[0])
				return int32(idusr)
			}
		}
	}
	return -1
}

func GetGrpSystem(content []string, usr string) int32 {
	var grp_usr string
	for _, xy := range content {
		if strings.Contains(xy, ",U,") {
			if strings.Contains(xy, usr) {
				split := strings.Split(xy, ",")
				// idusr, _ := strconv.Atoi(split[0])
				grp_usr = split[2]
				goto g0
			}
		}
	}
	return -1

g0:
	for _, xy := range content {
		if strings.Contains(xy, ",G,") {
			if strings.Contains(xy, grp_usr) {
				split := strings.Split(xy, ",")
				idgrp, _ := strconv.Atoi(split[0])
				return int32(idgrp)
			}
		}
	}
	return -1
}
