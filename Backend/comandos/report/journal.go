package report

import (
	"encoding/binary"
	"fmt"
	"os"
	"proyecto/comandos/utils"
	"proyecto/estructuras/size"
	"proyecto/estructuras/structures"
	"strings"

	"github.com/fatih/color"
)

func Report_Journal(name string, path string, ruta string, id_disco string) {
	disco_buscado, edb := utils.ObtenerDiscoID(id_disco)
	if !edb {
		return
	}

	file, err := os.OpenFile(disco_buscado.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[Mount]: Error al abrir archivo")
		return
	}
	defer file.Close()

	//-----------
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

	if sb.S_filesistem_type != 3 {
		color.Red("[REP]: El sistema de archivos no es EXT3 -> " + utils.ToString([]byte(name)))
		return
	}

	journal := structures.Journal{}
	if err := binary.Read(file, binary.LittleEndian, &journal); err != nil {
		color.Red("[REP]: Error en la lectura del journal")
		return
	}

	nombre_sin_extension := strings.Split(name, ".")
	rutaB := path + "/" + nombre_sin_extension[0] + ".dot"
	dot, err := os.Create(rutaB)
	if err != nil {
		color.Red("Error al crear el archivo <" + name + ">")
		return
	}

	fmt.Fprintln(dot, ""+`digraph G {`)
	fmt.Fprintln(dot, "\t"+`node[shape=none];`)
	fmt.Fprintln(dot, "\t"+`start[label=<`)
	fmt.Fprintln(dot, "\t\t"+`<table>`)
	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td colspan="5" bgcolor="#6eff84" color="#298089"><font point-size="20"><b>Reporte Journaling</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#b7f4c1" color="#298089" width="100">Tipo Operacion</td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#b7f4c1" color="#298089" width="100">Tipo</td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#b7f4c1" color="#298089" width="600">Path</td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#b7f4c1" color="#298089" width="700">Contenido</td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#b7f4c1" color="#298089" width="150">Fecha</td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	for { //while true
		fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
		fmt.Fprintln(dot, "\t\t\t\t"+`<td>`+utils.Returnstring(utils.ToString(journal.J_Tipo_Operacion[:]))+`</td>`)
		fmt.Fprintln(dot, "\t\t\t\t"+`<td>`+string(journal.J_Tipo)+`</td>`)
		fmt.Fprintln(dot, "\t\t\t\t"+`<td>`+utils.Returnstring(utils.ToString(journal.J_Path[:]))+`</td>`)
		fmt.Fprintln(dot, "\t\t\t\t"+`<td>`+utils.Returnstring(utils.ToString(journal.J_Contenido[:]))+`</td>`)
		fmt.Fprintln(dot, "\t\t\t\t"+`<td>`+utils.IntFechaToStr(journal.J_Fecha)+`</td>`)
		fmt.Fprintln(dot, "\t\t\t"+`</tr>`)
		if journal.J_Sig != -1 {
			if journal.J_Sig != 0 {
				if _, err := file.Seek(int64(journal.J_Sig+size.SizeJournal()), 0); err != nil {
					color.Red("[REP]: Error en mover puntero")
					return
				}
				if err := binary.Read(file, binary.LittleEndian, &journal); err != nil {
					color.Red("[REP]: Error en la lectura del journal")
					return
				}
				if journal.J_Sig != -1 {
					if journal.J_Sig == 0 {
						break
					} else {
						continue
					}
				} else {
					break
				}
			} else {
				break
			}
		} else {
			break
		}
	}
	//----------------
	fmt.Fprintln(dot, "\t\t"+`</table>`)
	fmt.Fprintln(dot, "\t"+`>];`)

	fmt.Fprintln(dot, ""+`}`)
	//-------------
	dot.Close()
	// Generacion del reporte
	// imagePath := path + "/" + name

	// cmd := exec.Command("/usr/bin/dot", "-T"+nombre_sin_extension[1], rutaB, "-o", imagePath)
	// err = cmd.Run()
	// if err != nil {
	// color.Red("[REP]: Error al generar imagen")
	// return
	// }

	color.Green("[REP]: Journal «" + name + "» generated Sucessfull")
}
