package permitions

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

func Values_CHGRP(instructions []string) (string, string, bool) {
	var grp, user string

	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "grp") {
			var value = utils.TieneGRP("CHGRP", valor)
			grp = value
		} else if strings.HasPrefix(strings.ToLower(valor), "user") {
			var value = utils.TieneUser("CHGRP", valor)
			user = value
		} else {
			color.Yellow("[CHGRP]: Atributo no reconocido")
			return "", "", false
		}
	}
	if grp == "" || len(grp) == 0 {
		color.Red("[CHGRP]: No hay path")
		return "", "", false
	} else if user == "" || len(user) == 0 {
		color.Red("[CHGRP]: No tiene grp")
		return "", "", false
	}
	return user, grp, true
}

func CHGRP_EXECUTE(user string, grp string) {
	if !global.UsuarioLogeado.Logged_in {
		color.Red("[CHGRP]: No hay usuario loggeado")
		return
	}
	if global.UsuarioLogeado.UID != 1 && global.UsuarioLogeado.GID != 1 {
		color.Red("[CHGRP]: No tienes permisos para ejecutar este comando")
		return
	}
	if user == utils.ToString(global.UsuarioLogeado.User[:]) {
		color.Red("[CHGRP]: No puedes cambiar el grupo al usuario Root")
		return
	}
	mount := global.UsuarioLogeado.Mounted
	file, err := os.OpenFile(mount.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[CHGRP]: Error al abrir archivo")
		return
	}
	defer file.Close()

	if mount.Es_Particion_L {
		inicio := mount.Particion_L.Part_start + size.SizeEBR()
		if _, err := file.Seek(int64(inicio), 0); err != nil {
			color.Red("[CHGRP]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &utils.Sb_AdminUsr); err != nil {
			color.Red("[CHGRP]: Error en la lectura del SuperBloque")
			return
		}
	} else if mount.Es_Particion_P {
		inicio := mount.Particion_P.Part_start
		if _, err := file.Seek(int64(inicio), 0); err != nil {
			color.Red("[CHGRP]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &utils.Sb_AdminUsr); err != nil {
			color.Red("[CHGRP]: Error en la lectura del SuperBloque")
			return
		}
	}

	content := utils.GetContentAdminUsers(utils.Sb_AdminUsr.S_inode_start + size.SizeTablaInodo())
	split_content := strings.Split(content, "\n")
	// cantBlockAnt := int32(len(utils.SplitContent(content)))

	if !utils.GrupoExist(split_content, grp) {
		color.Red("[CHGRP]: Grupo <" + grp + "> no existente")
		return
	}

	if utils.UsrExist(split_content, user) {
		// var newContent []string
		var content_string string
		for _, econ := range split_content {
			if econ != "" {
				if strings.Contains(econ, ",U,") || strings.Contains(econ, ",G,") {
					if strings.Contains(econ, ",U,") {
						if strings.Contains(econ, user) {
							split_user := strings.Split(econ, ",")
							new_grupo_user := split_user[0] + "," + split_user[1] + "," + grp + "," + split_user[3] + "," + split_user[4]
							content_string += new_grupo_user + "\n"
							// newContent = append(newContent, new_grupo_user)
							continue
						}
					}
					content_string += econ + "\n"
					// newContent = append(newContent, econ)
				}
			}
		}

		// if len(newContent) == 0 {
		// 	return
		// }

		usersTxt := utils.SplitContent(content_string)

		inodo := structures.TablaInodo{}
		seekInodo := utils.Sb_AdminUsr.S_inode_start + size.SizeTablaInodo()

		if _, err := file.Seek(int64(seekInodo), 0); err != nil {
			color.Red("[CHGRP]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
			color.Red("[CHGRP]: Error en la lectura del Inodo")
			return
		}

		tamanio := int32(0)
		for tm := range usersTxt {
			tamanio += int32(len(usersTxt[tm]))
		}
		inodo.I_s = tamanio
		inodo.I_mtime = utils.ObFechaInt()
		inodo.I_atime = utils.ObFechaInt()

		var j = 0
		for j < len(usersTxt) {
			inodo = utils.AgregarArchivo(usersTxt[j], inodo, int32(j), -1)
			j++
		}

		if _, err := file.Seek(int64(seekInodo), 0); err != nil {
			color.Red("[CHGRP]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
			color.Red("[CHGRP]: Error en la escritura del Inodo")
			return
		}

		if mount.Es_Particion_P {
			if _, err := file.Seek(int64(mount.Particion_P.Part_start), 0); err != nil {
				color.Red("[MKGRP]: Error en mover puntero")
				return
			}
			if err := binary.Write(file, binary.LittleEndian, &utils.Sb_AdminUsr); err != nil {
				color.Red("[MKGRP]: Error en la escritura del SuperBloque")
				return
			}
		} else if mount.Es_Particion_L {
			if _, err := file.Seek(int64(mount.Particion_L.Part_start+size.SizeEBR()), 0); err != nil {
				color.Red("[MKGRP]: Error en mover puntero")
				return
			}
			if err := binary.Write(file, binary.LittleEndian, &utils.Sb_AdminUsr); err != nil {
				color.Red("[MKGRP]: Error en la escritura del SuperBloque")
				return
			}
		}

		if utils.Sb_AdminUsr.S_filesistem_type == 3 {
			//escribimos journal
			utils.EscribirJournal("chgrp", '1', "users.txt", user)
		}
	} else {
		color.Red("[CHGRP]: Usuario <" + user + "> no existente")
		return
	}
	color.Green("[CHGRP]: Grupo «" + grp + "» cambiado correctamente al usuario -> " + user)
}
