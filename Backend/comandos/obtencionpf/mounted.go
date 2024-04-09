package obtencionpf

import (
	"os"
	"proyecto/comandos/global"
	"proyecto/comandos/utils"
	"proyecto/estructuras/structures"

	"github.com/fatih/color"
)

func ObtenerMBR_Mounted() ([]structures.MBR_Obtener, bool) {
	directorio := "MIA/P1/Disks/"
	// var mbrsObtenidos = []structures.MBR{}
	isEmpty, err := isDirEmpty(directorio)
	if err != nil {
		color.Red("Error al verificar directorio")
		return []structures.MBR_Obtener{}, false
	}

	if isEmpty {
		//caso en el que esta vacio

		return []structures.MBR_Obtener{}, false
	}

	dir, err := os.Open(directorio)
	if err != nil {
		return []structures.MBR_Obtener{}, false
	}
	defer dir.Close()

	filename, err := dir.Readdirnames(-1)
	if err != nil {
		return []structures.MBR_Obtener{}, false
	}
	global.Mounted_Partitions = []global.ParticionesMontadas{}
	ObtenerEstructuras = []structures.MBR_Obtener{}
	for _, discos := range filename {
		fulpath := directorio + discos
		mbr, embr := utils.Obtener_FULL_MBR_FDISK(fulpath)
		if !embr {
			return []structures.MBR_Obtener{}, false
		}

		for _, _content := range mbr.Mbr_partitions {
			if utils.ToString(_content.Part_id[:]) != "" {

				nuevoMBR := global.ParticionesMontadas{}
				nuevoMBR.DriveLetter = discos[0]
				nuevoMBR.Es_Particion_P = true
				nuevoMBR.Es_Particion_L = false
				nuevoMBR.Particion_P = _content
				nuevoMBR.Path = fulpath
				nuevoMBR.ID_Particion = _content.Part_id
				nuevoMBR.Type = 'P'

				global.Mounted_Partitions = append(global.Mounted_Partitions, nuevoMBR)
			}
			continue
		}
		temp_ext := structures.MBR_Obtener{}

		temp_ext.Disco = string(discos[0]) + ".dsk"
		temp_ext.Disco_Path = fulpath
		temp_ext.Mbr_partitions[0].Id_mounted = utils.ToString(mbr.Mbr_partitions[0].Part_id[:])
		temp_ext.Mbr_partitions[1].Id_mounted = utils.ToString(mbr.Mbr_partitions[1].Part_id[:])
		temp_ext.Mbr_partitions[2].Id_mounted = utils.ToString(mbr.Mbr_partitions[2].Part_id[:])
		temp_ext.Mbr_partitions[3].Id_mounted = utils.ToString(mbr.Mbr_partitions[3].Part_id[:])
		temp_ext.Mbr_partitions[0].Particion = utils.ToString(mbr.Mbr_partitions[0].Part_name[:])
		temp_ext.Mbr_partitions[1].Particion = utils.ToString(mbr.Mbr_partitions[1].Part_name[:])
		temp_ext.Mbr_partitions[2].Particion = utils.ToString(mbr.Mbr_partitions[2].Part_name[:])
		temp_ext.Mbr_partitions[3].Particion = utils.ToString(mbr.Mbr_partitions[3].Part_name[:])
		temp_ext.Mbr_partitions[0].Status = mbr.Mbr_partitions[0].Part_status
		temp_ext.Mbr_partitions[1].Status = mbr.Mbr_partitions[1].Part_status
		temp_ext.Mbr_partitions[2].Status = mbr.Mbr_partitions[2].Part_status
		temp_ext.Mbr_partitions[3].Status = mbr.Mbr_partitions[3].Part_status
		temp_ext.Mbr_partitions[0].Type = string(mbr.Mbr_partitions[0].Part_type)
		temp_ext.Mbr_partitions[1].Type = string(mbr.Mbr_partitions[1].Part_type)
		temp_ext.Mbr_partitions[2].Type = string(mbr.Mbr_partitions[2].Part_type)
		temp_ext.Mbr_partitions[3].Type = string(mbr.Mbr_partitions[3].Part_type)
		ObtenerEstructuras = append(ObtenerEstructuras, temp_ext)
		continue

	}

	return ObtenerEstructuras, true
}
