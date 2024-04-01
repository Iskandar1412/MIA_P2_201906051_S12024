package partition

import (
	"encoding/binary"
	"os"
	"proyecto/comandos/utils"
	"proyecto/estructuras/size"
	"proyecto/estructuras/structures"
	"strconv"

	"github.com/fatih/color"
)

func ParticionLogica(_size int32, _driveletter byte, _name []byte, _unit byte, _type byte, _fit byte, _delete string, _add int32) {
	path := "MIA/P1/Disks/" + string(_driveletter) + ".dsk"
	if !utils.ExisteArchivo("FDISK", path) {
		color.Yellow("[FDISK] Cancel the operation because not yet a file")
		return
	}

	if utils.ExisteParticionExt(path) {
		if !utils.ExisteNombreP(path, utils.ToString(_name)) {
			pos := -1
			mbr, embr := utils.Obtener_FULL_MBR_FDISK(path)
			if !embr {
				return
			}
			for i := range mbr.Mbr_partitions {
				if mbr.Mbr_partitions[i].Part_type == 'E' {
					pos = i
					break
				}
			}
			if pos != -1 {
				ebrAux := structures.EBR{}
				espacioOcupado := int32(0)
				file, err := os.OpenFile(path, os.O_RDWR, 0666)
				if err != nil {
					color.Red("[*]: Error al abrir archivo")
					return
				}
				defer file.Close()
				if _, err := file.Seek(int64(mbr.Mbr_partitions[pos].Part_start), 0); err != nil {
					color.Red("[*]: Error en mover puntero")
					return
				}
				if err := binary.Read(file, binary.LittleEndian, &ebrAux); err != nil {
					color.Red("[*]: Error en la lectura del MBR")
					return
				}
				//exista por lo menos una particion logica
				if (ebrAux.Part_next != -1) || (ebrAux.Part_s != -1) {
					espacioOcupado += size.SizeEBR() + ebrAux.Part_s
					for (ebrAux.Part_next != -1) && (utils.Ftell(file) < (mbr.Mbr_partitions[pos].Part_start + mbr.Mbr_partitions[pos].Part_s)) {
						if _, err := file.Seek(int64(ebrAux.Part_next), 0); err != nil {
							color.Red("[FDISK]: Error en mover puntero")
							return
						}
						if err := binary.Read(file, binary.LittleEndian, &ebrAux); err != nil {
							color.Red("[FDISK]: Error en la lectura del EBR")
							return
						}
						espacioOcupado += size.SizeEBR() + ebrAux.Part_s
					}
					newE := structures.EBR{}
					newE.Part_fit = _fit
					newE.Part_start = ebrAux.Part_start + size.SizeEBR() + ebrAux.Part_s
					newE.Part_mount = -1
					newE.Part_next = -1
					newE.Name = utils.DevolverNombreByte(utils.ToString(_name))
					newE.Part_s = utils.Tamano(_size, _unit)
					espacioD := mbr.Mbr_partitions[pos].Part_s - espacioOcupado
					espacioNewE := size.SizeEBR() + newE.Part_s
					ebrAux.Part_next = newE.Part_start
					if espacioD > espacioNewE {
						if _, err := file.Seek(int64(ebrAux.Part_start), 0); err != nil {
							color.Red("[FDISK]: Error en mover puntero")
							return
						}
						if err := binary.Write(file, binary.LittleEndian, &ebrAux); err != nil {
							color.Red("[FDISK]: Error en la escritura del EBR")
							return
						}

						if _, err := file.Seek(int64(newE.Part_start), 0); err != nil {
							color.Red("[FDISK]: Error en mover puntero")
							return
						}
						if err := binary.Write(file, binary.LittleEndian, &newE); err != nil {
							color.Red("[FDISK]: Error en la escritura del EBR")
							return
						}
						file.Close()

						color.Green("-----------------------------------------------------------")
						color.Blue("Se creo la particion")
						color.Blue("Particion: " + utils.ToString(newE.Name[:]))
						color.Blue("Tipo Logica")
						color.Blue("Inicio: " + strconv.Itoa(int(newE.Part_start)))
						color.Blue("Tamaño: " + strconv.Itoa(int(newE.Part_s)))
						color.Green("-----------------------------------------------------------")
					} else {
						color.Red("[FDISK]: No hay espacio para crear particion logica")
					}
				} else { //primera particion logica
					ebrAux.Part_fit = _fit
					ebrAux.Part_start = mbr.Mbr_partitions[pos].Part_start
					ebrAux.Part_mount = -1
					ebrAux.Part_s = utils.Tamano(_size, _unit)
					ebrAux.Part_next = -1
					ebrAux.Name = utils.DevolverNombreByte(utils.ToString(_name))

					if mbr.Mbr_partitions[pos].Part_s >= (ebrAux.Part_s + size.SizeEBR()) {
						if _, err := file.Seek(0, 0); err != nil {
							color.Red("[FDISK]: Error en mover puntero")
							return
						}
						if err := binary.Write(file, binary.LittleEndian, &mbr); err != nil {
							color.Red("[FDISK]: Error en la escritura del MBR")
							return
						}

						if _, err := file.Seek(int64(ebrAux.Part_start), 0); err != nil {
							color.Red("[FDISK]: Error en mover puntero")
							return
						}
						if err := binary.Write(file, binary.LittleEndian, &ebrAux); err != nil {
							color.Red("[FDISK]: Error en la escritura del EBR")
							return
						}

						file.Close()

						color.Green("-----------------------------------------------------------")
						color.Blue("Se creo la particion")
						color.Blue("Particion: " + utils.ToString(ebrAux.Name[:]))
						color.Blue("Tipo Logica")
						color.Blue("Inicio: " + strconv.Itoa(int(ebrAux.Part_start)))
						color.Blue("Tamaño: " + strconv.Itoa(int(ebrAux.Part_s)))
						color.Green("-----------------------------------------------------------")
					} else {
						color.Red("[FDISK]: No hay espacio para crear particion logica -> " + utils.ToString(_name))
						return
					}
				}
			}
		} else {
			color.Red("[FDISK]: La particion ya existe -> " + utils.ToString(_name))
			return
		}
	} else {
		color.Red("[FDISK]: No existe particion extendida para almacenar particion logica")
		return
	}
}
