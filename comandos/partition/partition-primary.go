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

func ParticionPrimaria(_size int32, _driveletter byte, _name []byte, _unit byte, _type byte, _fit byte, _delete string, _add int32) {
	path := "MIA/P1/Disks/" + string(_driveletter) + ".dsk"
	if !utils.ExisteArchivo("FDISK", path) {
		color.Yellow("[FDISK] Disco «" + string(_driveletter) + "» no existente")
		return
	}

	particion := utils.PartitionVacia()
	mbr, embr := utils.Obtener_FULL_MBR_FDISK(path)
	if !embr {
		return
	}

	pos := -1
	for i := range mbr.Mbr_partitions {
		if mbr.Mbr_partitions[i].Part_start == -1 {
			pos = i
			break
		}
	}

	if utils.ExisteNombreP(path, utils.ToString(_name)) {
		color.Red("[FDISK]: La particion «" + utils.ToString(_name) + "» ya existe")
		return
	}

	if utils.EspacioDisponible(_size, path, _unit, int32(pos)) {
		if !utils.ExisteNombreP(path, utils.ToString(_name)) {
			particion.Part_fit = _fit
			particion.Part_type = _type
			particion.Part_name = utils.DevolverNombreByte(utils.ToString(_name))
			particion.Part_status = -1
			particion.Part_correlative = int32(pos + 1)
			particion.Part_s = utils.Tamano(_size, _unit)
			if pos == 0 {
				particion.Part_start = size.SizeMBR()
			} else {
				particion.Part_start = mbr.Mbr_partitions[pos-1].Part_start + mbr.Mbr_partitions[pos-1].Part_s
			}
			mbr.Mbr_partitions[pos] = particion
			file, err := os.OpenFile(path, os.O_RDWR, 0666)
			if err != nil {
				color.Red("[FDISK]: Error al abrir archivo")
				return
			}
			defer file.Close()
			if _, err := file.Seek(0, 0); err != nil {
				color.Red("[FDISK]: Error en mover puntero")
				return
			}
			if err := binary.Write(file, binary.LittleEndian, &mbr); err != nil {
				color.Red("[FDISK]: Error en la lectura del MBR")
				return
			}
			file.Close()
			//comprobación
			comprobacion := structures.MBR{}
			file, err = os.OpenFile(path, os.O_RDWR, 0666)
			if err != nil {
				color.Red("[FDISK]: Error al abrir archivo")
				return
			}
			defer file.Close()
			if _, err := file.Seek(0, 0); err != nil {
				color.Red("[FDISK]: Error en mover puntero")
				return
			}
			if err := binary.Read(file, binary.LittleEndian, &comprobacion); err != nil {
				color.Red("[FDISK]: Error en la lectura del MBR")
				return
			}
			file.Close()
			color.Green("-----------------------------------------------------------")
			color.Blue("Se creo la particion #" + strconv.Itoa(int(comprobacion.Mbr_partitions[pos].Part_correlative)))
			color.Blue("Particion: " + utils.ToString(comprobacion.Mbr_partitions[pos].Part_name[:]))
			color.Blue("Tipo Primaria")
			color.Blue("Inicio: " + strconv.Itoa(int(comprobacion.Mbr_partitions[pos].Part_start)))
			color.Blue("Tamaño: " + strconv.Itoa(int(comprobacion.Mbr_partitions[pos].Part_s)))
			color.Green("-----------------------------------------------------------")
		} else {
			color.Yellow("[FDISK]: Particion <" + utils.ToString(_name) + "> existente")
			return
		}
	} else {
		color.Red("[FDISK]: Espacio Insuficiente")
		return
	}
}
