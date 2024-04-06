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

func Report_INODE(name string, path string, ruta string, id_disco string) {
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
	//-----------

	InicioSB := int32(0)
	if disco_buscado.Es_Particion_L {
		InicioSB = disco_buscado.Particion_L.Part_start + size.SizeEBR()
	} else if disco_buscado.Es_Particion_P {
		InicioSB = disco_buscado.Particion_P.Part_start
	}

	sb := structures.SuperBloque{}
	if _, err := file.Seek(int64(InicioSB), 0); err != nil {
		color.Red("[REP]: Error en mover puntero")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &sb); err != nil {
		color.Red("[REP]: Error en la lectura del SuperBloque")
		return
	}
	start := sb.S_bm_inode_start
	end := start + sb.S_inodes_count
	inodo := structures.TablaInodo{}
	var bit byte
	cont := int32(0)
	// fecha := int32(0)

	fmt.Fprintln(dot, ""+`digraph G {`)
	fmt.Fprintln(dot, "\t"+`node[shape=none];`)
	//-----------
	for i := start; i < end; i++ {
		if _, err := file.Seek(int64(i), 0); err != nil {
			color.Red("[REP]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
			color.Red("[REP]: Error en la lectura del Inodo")
			return
		}
		// fmt.Println(fecha)
		if bit == '1' {
			if _, err := file.Seek(int64(sb.S_inode_start+(cont*size.SizeTablaInodo())), 0); err != nil {
				color.Red("[REP]: Error en mover puntero")
				return
			}
			if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
				color.Red("[REP]: Error en la lectura de la Tabla de Inodos")
				return
			}
			fmt.Fprintln(dot, "\t"+`n`+fmt.Sprint(cont)+`[label=<<table>`)
			//-
			fmt.Fprintln(dot, "\t\t"+`<tr>`)
			fmt.Fprintln(dot, "\t\t\t"+`<td colspan="2">Inodo `+fmt.Sprint(cont)+`</td>`)
			fmt.Fprintln(dot, "\t\t"+`</tr>`)

			fmt.Fprintln(dot, "\t\t"+`<tr>`)
			fmt.Fprintln(dot, "\t\t\t"+`<td width="100">i_uid</td>`)
			fmt.Fprintln(dot, "\t\t\t"+`<td width="200">`+fmt.Sprint(inodo.I_uid)+`</td>`)
			fmt.Fprintln(dot, "\t\t"+`</tr>`)

			fmt.Fprintln(dot, "\t\t"+`<tr>`)
			fmt.Fprintln(dot, "\t\t\t"+`<td width="100">i_gid</td>`)
			fmt.Fprintln(dot, "\t\t\t"+`<td width="200">`+fmt.Sprint(inodo.I_gid)+`</td>`)
			fmt.Fprintln(dot, "\t\t"+`</tr>`)

			fmt.Fprintln(dot, "\t\t"+`<tr>`)
			fmt.Fprintln(dot, "\t\t\t"+`<td width="100">i_s</td>`)
			fmt.Fprintln(dot, "\t\t\t"+`<td width="200">`+fmt.Sprint(inodo.I_s)+`</td>`)
			fmt.Fprintln(dot, "\t\t"+`</tr>`)

			fmt.Fprintln(dot, "\t\t"+`<tr>`)
			fmt.Fprintln(dot, "\t\t\t"+`<td width="100">i_atime</td>`)
			fmt.Fprintln(dot, "\t\t\t"+`<td width="200">`+utils.IntFechaToStr(inodo.I_atime)+`</td>`)
			fmt.Fprintln(dot, "\t\t"+`</tr>`)

			fmt.Fprintln(dot, "\t\t"+`<tr>`)
			fmt.Fprintln(dot, "\t\t\t"+`<td width="100">i_ctime</td>`)
			fmt.Fprintln(dot, "\t\t\t"+`<td width="200">`+utils.IntFechaToStr(inodo.I_ctime)+`</td>`)
			fmt.Fprintln(dot, "\t\t"+`</tr>`)

			fmt.Fprintln(dot, "\t\t"+`<tr>`)
			fmt.Fprintln(dot, "\t\t\t"+`<td width="100">i_mtime</td>`)
			fmt.Fprintln(dot, "\t\t\t"+`<td width="200">`+utils.IntFechaToStr(inodo.I_mtime)+`</td>`)
			fmt.Fprintln(dot, "\t\t"+`</tr>`)

			for k, iblo := range inodo.I_block {
				fmt.Fprintln(dot, "\t\t"+`<tr>`)
				fmt.Fprintln(dot, "\t\t\t"+`<td width="100">i_block `+fmt.Sprint(k)+`</td>`)
				fmt.Fprintln(dot, "\t\t\t"+`<td width="200">`+fmt.Sprint(iblo)+`</td>`)
				fmt.Fprintln(dot, "\t\t"+`</tr>`)
			}

			fmt.Fprintln(dot, "\t\t"+`<tr>`)
			fmt.Fprintln(dot, "\t\t\t"+`<td width="100">i_type</td>`)
			fmt.Fprintln(dot, "\t\t\t"+`<td width="200">`+fmt.Sprint(inodo.I_type)+`</td>`)
			fmt.Fprintln(dot, "\t\t"+`</tr>`)

			fmt.Fprintln(dot, "\t\t"+`<tr>`)
			fmt.Fprintln(dot, "\t\t\t"+`<td width="100">i_perm</td>`)
			fmt.Fprintln(dot, "\t\t\t"+`<td width="200">`+fmt.Sprint(inodo.I_perm)+`</td>`)
			fmt.Fprintln(dot, "\t\t"+`</tr>`)
			//-
			fmt.Fprintln(dot, "\t"+`n`+fmt.Sprint(cont)+`</table>>];`)
			// fecha++
		}
		cont++
	}
	//-----------
	fmt.Fprintln(dot, ""+`}`)
	//-----------
	dot.Close()

	// Generacion del reporte
	// imagePath := path + "/" + name

	// cmd := exec.Command("/usr/bin/dot", "-T"+nombre_sin_extension[1], rutaB, "-o", imagePath)
	// err = cmd.Run()
	// if err != nil {
	// 	color.Red("[REP]: Error al generar imagen")
	// 	return
	// }

	color.Green("[REP]: Inodos Table «" + name + "» generated Sucessfull")
}
