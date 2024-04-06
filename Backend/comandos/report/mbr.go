package report

import (
	"encoding/binary"
	"fmt"
	"os"
	"proyecto/comandos/utils"
	"proyecto/estructuras/structures"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func Report_MBR(name string, path string, ruta string, id_disco string) {
	disco_buscado, edb := utils.ObtenerDiscoID(id_disco)
	if !edb {
		return
	}

	nombre_sin_extension := strings.Split(name, ".")
	rutaB := path + "/" + nombre_sin_extension[0] + ".dot"
	dot, err := os.Create(rutaB)
	if err != nil {
		color.Red("Error al crear el archivo <" + name + ">")
		return
	}

	file, err := os.OpenFile(disco_buscado.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[Mount]: Error al abrir archivo")
		return
	}
	defer file.Close()

	mbr, _ := utils.Obtener_FULL_MBR_FDISK(disco_buscado.Path)

	fmt.Fprintln(dot, "digraph G{")
	fmt.Fprintln(dot, "\tnode[shape=none];")
	fmt.Fprintln(dot, "\tstart[label=<<table>")
	fmt.Fprintln(dot, "\t\t"+`<tr><td colspan="2" bgcolor="#6308d8"><font color="white">REPORTE DE MBR</font></td></tr>`)
	fmt.Fprintln(dot, "\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t"+`<td>mbr_tamano</td>`)
	fmt.Fprintln(dot, "\t\t\t"+`<td>`+strconv.Itoa(int(mbr.Mbr_tamano))+`</td>`)
	fmt.Fprintln(dot, "\t\t"+`</tr>`)
	fmt.Fprintln(dot, "\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t"+`<td bgcolor="#b48be8">mbr_fecha_creacion</td>`)
	fmt.Fprintln(dot, "\t\t\t"+`<td bgcolor="#b48be8">`+utils.IntFechaToStr(mbr.Mbr_fecha_creacion)+`</td>`)
	fmt.Fprintln(dot, "\t\t"+`</tr>`)
	fmt.Fprintln(dot, "\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t"+`<td>mbr_disk_signature</td>`)
	fmt.Fprintln(dot, "\t\t\t"+`<td>`+strconv.Itoa(int(mbr.Mbr_disk_signature))+`</td>`)
	fmt.Fprintln(dot, "\t\t"+`</tr>`)

	for _, particion := range mbr.Mbr_partitions {
		fmt.Fprintln(dot, "\t\t"+`<tr><td colspan="2" bgcolor="#6308d8"><font color="white">Particion</font></td></tr>`)
		fmt.Fprintln(dot, "\t\t"+`<tr>`)
		fmt.Fprintln(dot, "\t\t\t"+`<td>part_status</td>`)
		fmt.Fprintln(dot, "\t\t\t"+`<td>`+strconv.Itoa(int(particion.Part_status))+`</td>`)
		fmt.Fprintln(dot, "\t\t"+`</tr>`)
		fmt.Fprintln(dot, "\t\t"+`<tr>`)
		fmt.Fprintln(dot, "\t\t\t"+`<td brcolor="#b48be8">part_type</td>`)
		fmt.Fprintln(dot, "\t\t\t"+`<td brcolor="#b48be8">`+utils.Returnstring(string(particion.Part_type))+`</td>`)
		fmt.Fprintln(dot, "\t\t"+`</tr>`)
		fmt.Fprintln(dot, "\t\t"+`<tr>`)
		fmt.Fprintln(dot, "\t\t\t"+`<td>part_fit</td>`)
		fmt.Fprintln(dot, "\t\t\t"+`<td>`+utils.Returnstring(string(particion.Part_fit))+`</td>`)
		fmt.Fprintln(dot, "\t\t"+`</tr>`)
		fmt.Fprintln(dot, "\t\t"+`<tr>`)
		fmt.Fprintln(dot, "\t\t\t"+`<td brcolor="#b48be8">part_start</td>`)
		fmt.Fprintln(dot, "\t\t\t"+`<td brcolor="#b48be8">`+strconv.Itoa(int(particion.Part_start))+`</td>`)
		fmt.Fprintln(dot, "\t\t"+`</tr>`)
		fmt.Fprintln(dot, "\t\t"+`<tr>`)
		fmt.Fprintln(dot, "\t\t\t"+`<td>part_size</td>`)
		fmt.Fprintln(dot, "\t\t\t"+`<td>`+strconv.Itoa(int(particion.Part_s))+`</td>`)
		fmt.Fprintln(dot, "\t\t"+`</tr>`)
		fmt.Fprintln(dot, "\t\t"+`<tr>`)
		fmt.Fprintln(dot, "\t\t\t"+`<td brcolor="#b48be8">part_name</td>`)
		fmt.Fprintln(dot, "\t\t\t"+`<td brcolor="#b48be8">`+utils.Returnstring(utils.ToString(particion.Part_name[:]))+`</td>`)
		fmt.Fprintln(dot, "\t\t"+`</tr>`)
		fmt.Fprintln(dot, "\t\t"+`<tr>`)
		fmt.Fprintln(dot, "\t\t\t"+`<td>part_id</td>`)
		fmt.Fprintln(dot, "\t\t\t"+`<td>`+utils.Returnstring(utils.ToString(particion.Part_id[:]))+`</td>`)
		fmt.Fprintln(dot, "\t\t"+`</tr>`)
		fmt.Fprintln(dot, "\t\t"+`<tr>`)
		fmt.Fprintln(dot, "\t\t\t"+`<td brcolor="#b48be8">part_correlative</td>`)
		fmt.Fprintln(dot, "\t\t\t"+`<td brcolor="#b48be8">`+strconv.Itoa(int(particion.Part_correlative))+`</td>`)
		fmt.Fprintln(dot, "\t\t"+`</tr>`)
		if particion.Part_type == 'E' {
			ebr := structures.EBR{}
			if _, err := file.Seek(int64(particion.Part_start), 0); err != nil {
				color.Red("[REP]: Error en mover puntero")
				return
			}
			if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
				color.Red("[REP]: Error en la lectura del EBR")
				return
			}
			if !(ebr.Part_s == -1 && ebr.Part_next == -1) {
				for { //while true
					fmt.Fprintln(dot, "\t\t"+`<tr><td colspan="2" bgcolor="#ff738c"><font color="white">Particion Logica</font></td></tr>`)
					fmt.Fprintln(dot, "\t\t"+`<tr>`)
					fmt.Fprintln(dot, "\t\t\t"+`<td>part_status</td>`)
					fmt.Fprintln(dot, "\t\t\t"+`<td>`+strconv.Itoa(int(ebr.Part_mount))+`</td>`)
					fmt.Fprintln(dot, "\t\t"+`</tr>`)
					fmt.Fprintln(dot, "\t\t"+`<tr>`)
					fmt.Fprintln(dot, "\t\t\t"+`<td bgbolor="#ffb2c0">part_next</td>`)
					fmt.Fprintln(dot, "\t\t\t"+`<td bgbolor="#ffb2c0">`+strconv.Itoa(int(ebr.Part_next))+`</td>`)
					fmt.Fprintln(dot, "\t\t"+`</tr>`)
					fmt.Fprintln(dot, "\t\t"+`<tr>`)
					fmt.Fprintln(dot, "\t\t\t"+`<td>part_fit</td>`)
					fmt.Fprintln(dot, "\t\t\t"+`<td>`+utils.Returnstring(string(ebr.Part_fit))+`</td>`)
					fmt.Fprintln(dot, "\t\t"+`</tr>`)
					fmt.Fprintln(dot, "\t\t"+`<tr>`)
					fmt.Fprintln(dot, "\t\t\t"+`<td bgbolor="#ffb2c0">part_start</td>`)
					fmt.Fprintln(dot, "\t\t\t"+`<td bgbolor="#ffb2c0">`+strconv.Itoa(int(ebr.Part_start))+`</td>`)
					fmt.Fprintln(dot, "\t\t"+`</tr>`)
					fmt.Fprintln(dot, "\t\t"+`<tr>`)
					fmt.Fprintln(dot, "\t\t\t"+`<td>part_size</td>`)
					fmt.Fprintln(dot, "\t\t\t"+`<td>`+strconv.Itoa(int(ebr.Part_s))+`</td>`)
					fmt.Fprintln(dot, "\t\t"+`</tr>`)
					fmt.Fprintln(dot, "\t\t"+`<tr>`)
					fmt.Fprintln(dot, "\t\t\t"+`<td bgbolor="#ffb2c0">part_name</td>`)
					fmt.Fprintln(dot, "\t\t\t"+`<td bgbolor="#ffb2c0">`+utils.Returnstring(utils.ToString(ebr.Name[:]))+`</td>`)
					fmt.Fprintln(dot, "\t\t"+`</tr>`)
					if ebr.Part_next == -1 {
						break
					}
					if _, err := file.Seek(int64(ebr.Part_next), 0); err != nil {
						color.Red("[REP]: Error en mover puntero")
						return
					}
					if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
						color.Red("[REP]: Error en la lectura del EBR")
						return
					}
				}
			}
		}
	}
	fmt.Fprintln(dot, "\t"+`</table>>];`)
	fmt.Fprintln(dot, `}`)
	dot.Close()

	// Generacion del reporte
	// imagePath := path + "/" + name

	// cmd := exec.Command("/usr/bin/dot", "-T"+nombre_sin_extension[1], rutaB, "-o", imagePath)
	// err = cmd.Run()
	// if err != nil {
	// color.Red("[REP]: Error al generar imagen")
	// return
	// }

	color.Green("[REP]: MBR «" + name + "» generated Sucessfull")
}
