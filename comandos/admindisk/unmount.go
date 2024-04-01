package admindisk

import (
	"encoding/binary"
	"os"
	"proyecto/comandos/global"
	"proyecto/comandos/utils"
	"proyecto/estructuras/size"
	"proyecto/estructuras/structures"
	"reflect"
	"strings"

	"github.com/fatih/color"
)

func Values_Unmount(instructions []string) (string, bool) {
	var _id = ""
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "id") {
			var value = utils.TieneID("UNMOUNT", valor)
			_id = value
			break
		} else {
			color.Yellow("[UNMOUNT]: Atributo no reconocido")
			_id = ""
			break
		}
	}
	if _id == "" || len(_id) == 0 || len(_id) > 4 {
		return "", false
	} else {
		return _id, true
	}
}

func UNMOUNT_EXECUTE(comando string, nombre string) {
	bandera := false
	for _, disco := range global.Mounted_Partitions {
		if disco.ID_Particion == utils.IDParticionByte(nombre) {
			mbr, embr := utils.Obtener_FULL_MBR_FDISK(disco.Path)
			if !embr {
				return
			}

			if disco.DriveLetter != nombre[0] {
				bandera = false
				continue
			}

			bandera = true

			var nombre_particion [16]byte
			if disco.Es_Particion_L {
				nombre_particion = disco.Particion_L.Name
			} else if disco.Es_Particion_P {
				nombre_particion = disco.Particion_P.Part_name
			}

			conjunto, econ := utils.BuscarParticion(mbr, nombre_particion[:], disco.Path)
			if !econ {
				return
			}

			file, err := os.OpenFile(disco.Path, os.O_RDWR, 0666)
			if err != nil {
				color.Red("[Unmount]: Error al abrir archivo")
				return
			}

			// Inicio SB
			inicio := int32(0)

			// Verificar Logica
			logica := structures.EBR{}
			if temp, ok := conjunto[0].(structures.EBR); ok {
				v := reflect.ValueOf(temp)
				reflect.ValueOf(&logica).Elem().Set(v)

				inicio = logica.Part_start + size.SizeEBR()

				if _, err := file.Seek(int64(logica.Part_start), 0); err != nil {
					color.Red("[Unmount]: Error en mover puntero")
					return
				}
				if err := binary.Write(file, binary.LittleEndian, &logica); err != nil {
					color.Red("[Unmount]: Error en la escritura del EBR")
					return
				}

				color.Green("[Unmount]: Particion (logica) «" + utils.ToString(logica.Name[:]) + "» desmontada - (ID): -> " + nombre)
			}

			// Verificar Primaria
			primaria := structures.Partition{}
			if temp, ok := conjunto[0].(structures.Partition); ok {
				v := reflect.ValueOf(temp)
				reflect.ValueOf(&primaria).Elem().Set(v)

				inicio = primaria.Part_start
				primaria.Part_id = [4]byte{'\x00', '\x00', '\x00', '\x00'}

				count := 0
				for _, c := range mbr.Mbr_partitions {
					if utils.ToString(c.Part_name[:]) == utils.ToString(nombre_particion[:]) {
						break
					}
					count++
				}
				mbr.Mbr_partitions[count] = primaria

				if _, err := file.Seek(0, 0); err != nil {
					color.Red("[Unmount]: Error en mover puntero")
					return
				}
				if err := binary.Write(file, binary.LittleEndian, &mbr); err != nil {
					color.Red("[Unmount]: Error en la escritura del MBR")
					return
				}

				color.Green("[Unmount]: Particion (primaria) «" + utils.ToString(primaria.Part_name[:]) + "» desmontada - (ID): -> " + nombre)
			}

			var nueva_lista []global.ParticionesMontadas
			for _, c := range global.Mounted_Partitions {
				if c.ID_Particion == utils.IDParticionByte(nombre) {
					if c.DriveLetter == nombre[0] {
						continue
					}
					continue
				}
				nueva_lista = append(nueva_lista, c)
			}

			global.Mounted_Partitions = nueva_lista

			// SuperBloque
			superblock := structures.SuperBloque{}
			if _, err := file.Seek(int64(inicio), 0); err != nil {
				color.Red("[Unmount]: Error en mover puntero")
				return
			}
			if err := binary.Read(file, binary.LittleEndian, &superblock); err != nil {
				color.Red("[Unmount]: Error en la lectura del EBR")
				return
			}

			if superblock.S_mnt_count > 0 {
				superblock.S_mtime = utils.ObFechaInt()
				superblock.S_mnt_count += 1
				if _, err := file.Seek(int64(inicio), 0); err != nil {
					color.Red("[Unmount]: Error en mover puntero")
					return
				}
				if err := binary.Write(file, binary.LittleEndian, &superblock); err != nil {
					color.Red("[Unmount]: Error en la escritura del SuperBloque")
					return
				}
			}

		}
	}
	if !bandera {
		color.Red("[Unmount]: ID de Particion [" + nombre + "] no encontrada")
		return
	}
}
