package utils

import (
	"encoding/binary"
	"io"
	"os"
	"proyecto/estructuras/size"
	"proyecto/estructuras/structures"

	"github.com/fatih/color"
)

func Obtener_FULL_MBR_FDISK(path string) (structures.MBR, bool) {
	mbr := structures.MBR{}
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[FDISK]: Error al abrir archivo")
		return structures.MBR{}, false
	}
	defer file.Close()
	if _, err := file.Seek(0, 0); err != nil {
		color.Red("[FDISK]: Error en mover puntero")
		return structures.MBR{}, false
	}
	if err := binary.Read(file, binary.LittleEndian, &mbr); err != nil {
		color.Red("[FDISK]: Error en la lectura del MBR")
		return structures.MBR{}, false
	}
	return mbr, true
}

func EspacioDisponible(s int32, path string, u byte, pos int32) bool {
	mbr, embr := Obtener_FULL_MBR_FDISK(path)
	if !embr {
		return false
	}

	if pos > -1 {
		if Tamano(s, u) > 0 {
			espacioRestante := 0
			if pos == 0 {
				espacioRestante = int(mbr.Mbr_tamano) - int(size.SizeMBR())
			} else {
				espacioRestante = int(mbr.Mbr_tamano) - int(mbr.Mbr_partitions[pos-1].Part_start) - int(mbr.Mbr_partitions[pos-1].Part_s)
			}
			return espacioRestante >= int(Tamano(s, u))
		}
	}
	return false
}

func ExisteNombreP(path string, name string) bool {
	mbr, embr := Obtener_FULL_MBR_FDISK(path)
	if !embr {
		return true
	}

	for i := range mbr.Mbr_partitions {
		if ToString(mbr.Mbr_partitions[i].Part_name[:]) == name {
			return true
		}
		if mbr.Mbr_partitions[i].Part_type == 'E' {
			EBR := structures.EBR{}
			file, err := os.OpenFile(path, os.O_RDWR, 0666)
			if err != nil {
				color.Red("[*]: Error al abrir archivo")
				return true
			}
			defer file.Close()
			if _, err := file.Seek(int64(mbr.Mbr_partitions[i].Part_start), 0); err != nil {
				color.Red("[*]: Error en mover puntero")
				return true
			}
			if err := binary.Read(file, binary.LittleEndian, &EBR); err != nil {
				color.Red("[*]: Error en la lectura del MBR")
				return true
			}
			if EBR.Part_next != -1 || EBR.Part_s != -1 {
				if ToString(EBR.Name[:]) == name {
					return true
				}
				for EBR.Part_next != -1 {
					if ToString(EBR.Name[:]) == name {
						return true
					}
					if _, err := file.Seek(int64(EBR.Part_next), 0); err != nil {
						color.Red("[*]: Error en mover puntero")
						return true
					}
					if err := binary.Read(file, binary.LittleEndian, &EBR); err != nil {
						color.Red("[*]: Error en la lectura del MBR")
						return true
					}
					//si la particion que le sigue
					if ToString(EBR.Name[:]) == ToString([]byte(name)) {
						return true
					}
				}
			}
		}
	}
	return false
}

func ExisteParticionExt(path string) bool {
	mbr, embr := Obtener_FULL_MBR_FDISK(path)
	if !embr {
		return false
	}

	for i := range mbr.Mbr_partitions {
		if mbr.Mbr_partitions[i].Part_type == 'E' {
			return true
		}
	}
	return false
}

func Ftell(file *os.File) int32 {
	pos, _ := file.Seek(0, io.SeekCurrent)
	return int32(pos)
}
