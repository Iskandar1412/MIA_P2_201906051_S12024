package adminusers

import (
	"encoding/binary"
	"fmt"
	"os"
	"proyecto/comandos/global"
	"proyecto/comandos/utils"
	"proyecto/estructuras/size"
	"proyecto/estructuras/structures"
	"strings"

	"github.com/fatih/color"
)

func Values_MKUSR(instructions []string) (string, string, string, bool) {
	var _nam string
	var _pas string
	var _grp string
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "user") {
			var value = utils.TieneUser("MKUSR", valor)
			if len(value) > 10 {
				color.Red("No se puede ser un nombre mayor a 10 caracteres")
				_nam = ""
				break
			} else {
				_nam = value
				continue
			}
		} else if strings.HasPrefix(strings.ToLower(valor), "pass") {
			var value = utils.TienePassword("MKUSR", valor)
			if len(value) > 10 {
				color.Red("No se puede ser un password mayor a 10 caracteres")
				_pas = ""
				break
			} else {
				_pas = value
				continue
			}
		} else if strings.HasPrefix(strings.ToLower(valor), "grp") {
			var value = utils.TieneGRP("MKUSR", valor)
			if len(value) > 10 {
				color.Red("No se puede ser un grp mayor a 10 caracteres")
				_grp = ""
				break
			} else {
				_grp = value
				continue
			}
		} else {
			color.Yellow("[MKUSR]: Atrubuto no reconocido")
		}
	}
	if _nam == "" || len(_nam) == 0 || len(_nam) > 10 {
		return "", "", "", false
	} else if _pas == "" || len(_pas) == 0 || len(_pas) > 10 {
		return "", "", "", false
	} else if _grp == "" || len(_grp) == 0 || len(_grp) > 10 {
		return "", "", "", false
	}

	return _nam, _pas, _grp, true
}

func MKUSR_EXECUTE(_us string, _pas string, _grp string) {
	if !global.UsuarioLogeado.Logged_in {
		color.Red("[MKUSR]: No hay usuario loggeado")
		return
	}
	if global.UsuarioLogeado.UID != 1 && global.UsuarioLogeado.GID != 1 {
		color.Red("[MKUSR]: No tienes permisos para ejecutar este comando")
		return
	}

	mount := global.UsuarioLogeado.Mounted

	file, err := os.OpenFile(mount.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[MKUSR]: Error al abrir archivo")
		return
	}
	defer file.Close()

	if mount.Es_Particion_L {
		inicio := mount.Particion_L.Part_start + size.SizeEBR()
		if _, err := file.Seek(int64(inicio), 0); err != nil {
			color.Red("[MKUSR]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &utils.Sb_AdminUsr); err != nil {
			color.Red("[MKUSR]: Error en la lectura del SuperBloque")
			return
		}
	} else if mount.Es_Particion_P {
		inicio := mount.Particion_P.Part_start
		if _, err := file.Seek(int64(inicio), 0); err != nil {
			color.Red("[MKUSR]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &utils.Sb_AdminUsr); err != nil {
			color.Red("[MKUSR]: Error en la lectura del SuperBloque")
			return
		}
	}

	content := utils.GetContentAdminUsers(utils.Sb_AdminUsr.S_inode_start + size.SizeTablaInodo())

	cantBlockAnt := int32(len(utils.SplitContent(content)))

	nuevoUsr := ""
	split_content := strings.Split(content, "\n")

	if utils.UsrExist(split_content, _us) {
		color.Red("[MKUSR]: Usuario «" + _us + "» ya existente")
		return
	}

	if utils.GrupoExist(split_content, _grp) {
		nuevoUsr = fmt.Sprint(utils.GetUID(split_content)) + ",U," + _grp + "," + _us + "," + _pas + "\n"
		content += nuevoUsr
		usersTxt := utils.SplitContent(content)
		cantBlockAct := int32(len(usersTxt))
		if len(usersTxt) > 4380 {
			color.Red("[MKUSR]: No se pueden crear más grupos")
			return
		}

		if utils.Sb_AdminUsr.S_free_blocks_count < (cantBlockAct - cantBlockAnt) {
			color.Red("[MKUSR]: No hay bloques suficientes para crear archivo")
			return
		}

		//Buscar bitmap de bloques

		var bit byte
		start := utils.Sb_AdminUsr.S_bm_block_start
		end := start + utils.Sb_AdminUsr.S_block_start
		cantContiguos := int32(0)
		inicioBM := int32(-1)
		inicioB := int32(-1)
		contadorA := int32(0)
		if (cantBlockAct - cantBlockAnt) > 0 {
			for i := start; i < end; i++ {
				if _, err := file.Seek(int64(i), 0); err != nil {
					color.Red("[MKUSR]: Error en mover puntero")
					return
				}
				if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
					color.Red("[MKUSR]: Error en la lectura del Archivo")
					return
				}
				if bit == '1' { //ocupado
					cantContiguos = 0
					inicioBM = -1
					inicioB = -1
				} else { // desocupados
					if cantContiguos == 0 {
						inicioBM = int32(i)
						inicioB = contadorA
					}
					cantContiguos++
				}
				if cantContiguos >= (cantBlockAct - cantBlockAnt) {
					break
				}

				contadorA++
			}
			if (inicioBM == -1) || (cantContiguos != (cantBlockAct - cantBlockAnt)) {
				color.Red("[MKUSR]: No hay bloques suficientes para actualizar el archivo")
				return
			}

			for i := inicioBM; i < (inicioBM + (cantBlockAct - cantBlockAnt)); i++ {
				var uno byte = '1'
				if _, err := file.Seek(int64(i), 0); err != nil {
					color.Red("[MKUSR]: Error en mover puntero")
					return
				}

				if err := binary.Write(file, binary.LittleEndian, &uno); err != nil {
					color.Red("[MPGRP]: Error en la escritura del archivo")
					return
				}
			}

			utils.Sb_AdminUsr.S_free_blocks_count -= (cantBlockAct - cantBlockAnt)
			bit2 := int32(0)
			for k := start; k < end; k++ {
				if _, err := file.Seek(int64(k), 0); err != nil {
					color.Red("[MKUSR]: Error en mover puntero")
					return
				}
				if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
					color.Red("[MKUSR]: Error en la lectura del Archivo")
					return
				}
				if bit == '0' {
					break
				}
				bit2++
			}
			utils.Sb_AdminUsr.S_first_blo = bit2
		}

		inodo := structures.TablaInodo{}
		seekInodo := utils.Sb_AdminUsr.S_inode_start + size.SizeTablaInodo()
		if _, err := file.Seek(int64(seekInodo), 0); err != nil {
			color.Red("[MKUSR]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
			color.Red("[MKUSR]: Error en la lectura del Inodo")
			return
		}

		tamanio := int32(0)
		for tm := range usersTxt {
			tamanio += int32(len(usersTxt[tm]))
		}

		inodo.I_s = tamanio
		inodo.I_mtime = utils.ObFechaInt()

		var j, contador = 0, 0
		for j < len(usersTxt) {
			utils.CambioCont = false
			inodo = utils.AgregarArchivo(usersTxt[j], inodo, int32(j), (inicioB + int32(contador)))
			if utils.CambioCont {
				contador++
			}
			j++
		}

		if _, err := file.Seek(int64(seekInodo), 0); err != nil {
			color.Red("[MKUSR]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
			color.Red("[MKUSR]: Error en la escritura del Inodo")
			return
		}
		if mount.Es_Particion_P {
			if _, err := file.Seek(int64(mount.Particion_P.Part_start), 0); err != nil {
				color.Red("[MKUSR]: Error en mover puntero")
				return
			}
			if err := binary.Write(file, binary.LittleEndian, &utils.Sb_AdminUsr); err != nil {
				color.Red("[MKUSR]: Error en la escritura del SuperBloque")
				return
			}
		} else if mount.Es_Particion_L {
			if _, err := file.Seek(int64(mount.Particion_L.Part_start+size.SizeEBR()), 0); err != nil {
				color.Red("[MKUSR]: Error en mover puntero")
				return
			}
			if err := binary.Write(file, binary.LittleEndian, &utils.Sb_AdminUsr); err != nil {
				color.Red("[MKUSR]: Error en la escritura del SuperBloque")
				return
			}
		}

		if utils.Sb_AdminUsr.S_filesistem_type == 3 {
			//escribimos journal
			utils.EscribirJournal("mkusr", '1', "users.txt", nuevoUsr)
		}
	} else {
		color.Red("[MKUSR]: El grupo «" + _grp + "» no existe")
		return
	}

	color.Green("[MKUSR]: Usuario «" + _us + "» guardado correctamente")
}
