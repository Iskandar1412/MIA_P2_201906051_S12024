package report

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"proyecto/comandos/utils"
	"proyecto/estructuras/size"
	"proyecto/estructuras/structures"
	"strings"

	"github.com/fatih/color"
)

func Report_BM_Block(name string, path string, ruta string, id_disco string) {
	disco_buscado, edb := utils.ObtenerDiscoID(id_disco)
	if !edb {
		return
	}

	nombre_sin_extension := strings.Split(name, ".")
	_, esimagen := utils.EsImagen(strings.ToLower(nombre_sin_extension[1]), "jpg", "png", "svg")

	if esimagen {
		rutaB := path + "/" + nombre_sin_extension[0] + ".dot"
		dot, err := os.Create(rutaB)
		if err != nil {
			color.Red("Error al crear el archivo <" + name + ">")
			return
		}

		file, err := os.OpenFile(disco_buscado.Path, os.O_RDWR, 0666)
		if err != nil {
			color.Red("[REP]: Error al abrir archivo")
			return
		}
		defer file.Close()
		//----------
		inicioSB := int32(0)
		if disco_buscado.Es_Particion_L {
			if disco_buscado.Particion_L.Part_mount != 1 {
				color.Red("[REP]: El disco (logico) no se ha formateado -> " + utils.ToString(disco_buscado.Particion_L.Name[:]) + " - (ID): -> " + utils.ToString(disco_buscado.ID_Particion[:]))
				return
			}
			inicioSB = disco_buscado.Particion_L.Part_start + size.SizeEBR()
		} else if disco_buscado.Es_Particion_P {
			if disco_buscado.Particion_P.Part_status != 1 {
				color.Red("[REP]: El disco (primario) no se ha formateado -> " + utils.ToString(disco_buscado.Particion_P.Part_name[:]) + " - (ID): -> " + utils.ToString(disco_buscado.ID_Particion[:]))
				return
			}
			inicioSB = disco_buscado.Particion_P.Part_start
		}

		sb := structures.SuperBloque{}
		if _, err := file.Seek(int64(inicioSB), 0); err != nil {
			color.Red("[REP]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &sb); err != nil {
			color.Red("[REP]: Error en la lectura del SuperBloque")
			return
		}
		start := sb.S_bm_block_start
		end := start + sb.S_blocks_count
		var bit byte
		controlador := int32(0)

		fmt.Fprintln(dot, ""+`digraph G {`)
		fmt.Fprintln(dot, "\t"+`node[shape=none, lblstyle="align=left"];`)

		fmt.Fprint(dot, "\t"+`start[label="`)
		for i := start; i < end; i++ {
			if _, err := file.Seek(int64(i), 0); err != nil {
				color.Red("[REP]: Error en mover puntero")
				return
			}
			if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
				color.Red("[REP]: Error en la lectura del Bit")
				return
			}
			fmt.Fprint(dot, ""+``+string(bit)+" ")
			if controlador == 19 {
				fmt.Fprint(dot, "\\n")
				controlador = 0
			} else {
				controlador++
			}
		}
		fmt.Fprintln(dot, ""+`"];`)

		fmt.Fprintln(dot, ""+`}`)
		//----------
		dot.Close()
		file.Close()
		// Generacion del reporte
		imagePath := path + "/" + name

		cmd := exec.Command("/usr/bin/dot", "-T"+nombre_sin_extension[1], rutaB, "-o", imagePath)
		err = cmd.Run()
		if err != nil {
			color.Red("[REP]: Error al generar imagen")
			return
		}
	} else { //caso de ser txt

		rutaB := path + "/" + nombre_sin_extension[0] + "." + nombre_sin_extension[1]
		dot, err := os.Create(rutaB)
		if err != nil {
			color.Red("Error al crear el archivo <" + name + ">")
			return
		}

		file, err := os.OpenFile(disco_buscado.Path, os.O_RDWR, 0666)
		if err != nil {
			color.Red("[REP]: Error al abrir archivo")
			return
		}
		defer file.Close()
		//----------
		inicioSB := int32(0)
		if disco_buscado.Es_Particion_L {
			if disco_buscado.Particion_L.Part_mount != 1 {
				color.Red("[REP]: El disco (logico) no se ha formateado -> " + utils.ToString(disco_buscado.Particion_L.Name[:]) + " - (ID): -> " + utils.ToString(disco_buscado.ID_Particion[:]))
				return
			}
			inicioSB = disco_buscado.Particion_L.Part_start + size.SizeEBR()
		} else if disco_buscado.Es_Particion_P {
			if disco_buscado.Particion_P.Part_status != 1 {
				color.Red("[REP]: El disco (primario) no se ha formateado -> " + utils.ToString(disco_buscado.Particion_P.Part_name[:]) + " - (ID): -> " + utils.ToString(disco_buscado.ID_Particion[:]))
				return
			}
			inicioSB = disco_buscado.Particion_P.Part_start
		}

		sb := structures.SuperBloque{}
		if _, err := file.Seek(int64(inicioSB), 0); err != nil {
			color.Red("[REP]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &sb); err != nil {
			color.Red("[REP]: Error en la lectura del SuperBloque")
			return
		}
		start := sb.S_bm_block_start
		end := start + sb.S_blocks_count
		var bit byte
		controlador := int32(0)

		for i := start; i < end; i++ {
			if _, err := file.Seek(int64(i), 0); err != nil {
				color.Red("[REP]: Error en mover puntero")
				return
			}
			if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
				color.Red("[REP]: Error en la lectura del Bit")
				return
			}
			fmt.Fprint(dot, ""+``+string(bit)+" ")
			if controlador == 19 {
				fmt.Fprint(dot, "\n")
				controlador = 0
			} else {
				controlador++
			}
		}
		//----------
		file.Close()
		dot.Close()
	}

	color.Green("[REP]: Bitmap Block «" + name + "» generated Sucessfull")
}
