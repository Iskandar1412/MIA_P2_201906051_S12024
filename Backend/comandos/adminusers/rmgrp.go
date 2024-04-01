package adminusers

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

func Values_RMGRP(instructions []string) (string, bool) {
	var _name string
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "name") {
			var value = utils.TieneNombre("RMGRP", valor)
			if len(value) > 10 {
				color.Red("No puede ser nun nombre mayor a 10")
				_name = ""
				break
			} else {
				_name = value
			}
		} else {
			color.Yellow("[RMGRP]: Atributo no reconocido")
		}
	}
	if _name == "" || len(_name) == 0 {
		return "", false
	}
	return _name, true
}

func RMGRP_EXECUTE(_name string) {
	if !global.UsuarioLogeado.Logged_in {
		color.Red("[RMGRP]: No hay usuario loggeado")
		return
	}
	if global.UsuarioLogeado.UID != 1 && global.UsuarioLogeado.GID != 1 {
		color.Red("[RMGRP]: No tienes permisos para ejecutar este comando")
		return
	}

	mount := global.UsuarioLogeado.Mounted
	file, err := os.OpenFile(mount.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[RMGRP]: Error al abrir archivo")
		return
	}
	defer file.Close()

	if mount.Es_Particion_L {
		inicio := mount.Particion_L.Part_start + size.SizeEBR()
		if _, err := file.Seek(int64(inicio), 0); err != nil {
			color.Red("[RMGRP]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &utils.Sb_AdminUsr); err != nil {
			color.Red("[RMGRP]: Error en la lectura del SuperBloque")
			return
		}
	} else if mount.Es_Particion_P {
		inicio := mount.Particion_P.Part_start
		if _, err := file.Seek(int64(inicio), 0); err != nil {
			color.Red("[RMGRP]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &utils.Sb_AdminUsr); err != nil {
			color.Red("[RMGRP]: Error en la lectura del SuperBloque")
			return
		}
	}

	content := utils.GetContentAdminUsers(utils.Sb_AdminUsr.S_inode_start + size.SizeTablaInodo())
	usersTxt_Old := utils.SplitContent(content)
	split_content := strings.Split(content, "\n")
	if utils.GrupoExist(split_content, _name) {
		var newContentSplit []string
		for _, econ := range split_content {
			if econ != "" {
				if strings.Contains(econ, ",G,") || strings.Contains(econ, ",U,") {
					if strings.Contains(econ, ",G,") {
						if strings.Contains(econ, _name) {
							continue
						}
					}
					newContentSplit = append(newContentSplit, econ)
				}
			}
		}
		var newContent string
		for _, econ := range newContentSplit {
			newContent += econ + "\n"
		}
		content = newContent

		usersTxt := utils.SplitContent(content)
		var inodo structures.TablaInodo
		seekInodo := utils.Sb_AdminUsr.S_inode_start + size.SizeTablaInodo()
		if _, err := file.Seek(int64(seekInodo), 0); err != nil {
			color.Red("[RMGRP]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
			color.Red("[RMGRP]: Error en la lectura del Inodo")
			return
		}

		tamanio := int32(0)
		for tm := range usersTxt {
			tamanio += int32(len(usersTxt[tm]))
		}

		inodo.I_s = tamanio
		inodo.I_atime = utils.ObFechaInt()
		inodo.I_mtime = utils.ObFechaInt()

		var j = 0
		for j < len(usersTxt) {
			inodo = utils.AgregarArchivo(usersTxt[j], inodo, int32(j), -1)
			j++
		}

		//-----------------------
		if len(usersTxt_Old) > len(usersTxt) {
			conInRes := 0
			for iz, z := range inodo.I_block {
				if z != -1 {
					conInRes = iz
				} else {
					break
				}
			}
			inodo.I_block[conInRes] = -1

			var bit byte
			start := utils.Sb_AdminUsr.S_bm_block_start
			end := start + utils.Sb_AdminUsr.S_block_start
			cantContiguos := int32(0)
			inicioBM := int32(-1)
			// inicioB := int32(-1)
			contadorA := int32(0)
			for i := start; i < end; i++ {
				if _, err := file.Seek(int64(i), 0); err != nil {
					color.Red("[RMGRP]: Error en mover puntero")
					return
				}
				if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
					color.Red("[RMGRP]: Error en la lectura del Archivo")
					return
				}
				if bit == '1' { //ocupado
					cantContiguos = 0
					inicioBM = -1
					// inicioB = -1
				} else { // desocupados
					if cantContiguos == 0 {
						inicioBM = int32(i)
						// inicioB = contadorA
					}
					cantContiguos++
				}
				if cantContiguos >= int32(len(usersTxt_Old)-len(usersTxt)) {
					break
				}

				contadorA++
			}
			var zero byte = '0'
			if _, err := file.Seek(int64(inicioBM-1), 0); err != nil {
				color.Red("[RMGRP]: Error en mover puntero")
				return
			}
			if err := binary.Write(file, binary.LittleEndian, &zero); err != nil {
				color.Red("[RMGRP]: Error en la escritura del Archivo")
				return
			}
			utils.Sb_AdminUsr.S_free_blocks_count += int32(len(usersTxt_Old) - len(usersTxt))
			bit2 := int32(0)
			for k := start; k < end; k++ {
				if _, err := file.Seek(int64(k), 0); err != nil {
					color.Red("[RMGRP]: Error en mover puntero")
					return
				}
				if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
					color.Red("[RMGRP]: Error en la lectura del Archivo")
					return
				}
				if bit == '0' {
					break
				}
				bit2++
			}
			utils.Sb_AdminUsr.S_first_blo = bit2
		}
		//-----------------------

		if _, err := file.Seek(int64(seekInodo), 0); err != nil {
			color.Red("[RMGRP]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
			color.Red("[RMGRP]: Error en la escritura del Inodo")
			return
		}

		if mount.Es_Particion_P {
			if _, err := file.Seek(int64(mount.Particion_P.Part_start), 0); err != nil {
				color.Red("[RMGRP]: Error en mover puntero")
				return
			}
			if err := binary.Write(file, binary.LittleEndian, &utils.Sb_AdminUsr); err != nil {
				color.Red("[RMGRP]: Error en la escritura del SuperBloque")
				return
			}
		} else if mount.Es_Particion_L {
			if _, err := file.Seek(int64(mount.Particion_L.Part_start+size.SizeEBR()), 0); err != nil {
				color.Red("[RMGRP]: Error en mover puntero")
				return
			}
			if err := binary.Write(file, binary.LittleEndian, &utils.Sb_AdminUsr); err != nil {
				color.Red("[RMGRP]: Error en la escritura del SuperBloque")
				return
			}
		}

		if utils.Sb_AdminUsr.S_filesistem_type == 3 {
			//escribimos journal
			utils.EscribirJournal("rmgrp", '1', "users.txt", _name)
		}

	} else {
		color.Red("[RMGRP]: El grupo «" + _name + "» no existe")
		return
	}
	color.Blue("[RMGRP]: Grupo «" + _name + "» eliminado correctamente")
}
