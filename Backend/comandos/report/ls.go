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

func Report_LS(name string, path string, ruta string, id_disco string) {
	disco_buscado, edb := utils.ObtenerDiscoID(id_disco)
	if !edb {
		return
	}

	file, err := os.OpenFile(disco_buscado.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[REP]: Error al abrir archivo")
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

	posInodoF := int32(0)
	var rutaS []string
	//----------
	if ruta == "/" {
		posInodoF = sb.S_inode_start
		rutaS = append(rutaS, "/")
	} else {
		rutaS = utils.SplitRuta(ruta)
		if len(rutaS) == 0 {
			color.Red("[REP]: Ruta invalida")
			return
		}
		posInodoF = utils.GetInodoF(rutaS, 0, int32(len(rutaS)-1), sb.S_inode_start, disco_buscado.Path)
	}

	// Obtencion del inodo
	inodo := structures.TablaInodo{}
	if posInodoF == -1 {
		color.Red("[REP]: Archivo no encontrado")
		return
	}

	if _, err := file.Seek(int64(posInodoF), 0); err != nil {
		color.Red("[REP]: Error en mover puntero")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[REP]: Error en la lectura del Inodo")
		return
	}

	inodo.I_atime = utils.ObFechaInt()
	if _, err := file.Seek(int64(posInodoF), 0); err != nil {
		color.Red("[REP]: Error en mover puntero")
		return
	}
	if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[REP]: Error en la escritura de la Tabla de inodos")
		return
	}

	nombre_sin_extension := strings.Split(name, ".")
	rutaB := "Rep/" + nombre_sin_extension[0] + ".dot"
	dot, err := os.Create(rutaB)
	if err != nil {
		color.Red("Error al crear el archivo <" + name + ">")
		return
	}

	file.Close()
	fmt.Fprintln(dot, ""+`digraph G {`)
	fmt.Fprintln(dot, "\t"+`node[shape=none];`)
	fmt.Fprintln(dot, "\t\t"+`start[label=<<table>`)
	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t"+`<td>Permisos</td>`)
	fmt.Fprintln(dot, "\t\t\t"+`<td>Usuarios</td>`)
	fmt.Fprintln(dot, "\t\t\t"+`<td>Grupo</td>`)
	fmt.Fprintln(dot, "\t\t\t"+`<td>Size</td>`)
	fmt.Fprintln(dot, "\t\t\t"+`<td>Fecha</td>`)
	fmt.Fprintln(dot, "\t\t\t"+`<td>Tipo</td>`)
	fmt.Fprintln(dot, "\t\t\t"+`<td>Nombre</td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)
	// Obtener inodo para graficar
	nodograph := utils.LsInodo(posInodoF, rutaS[len(rutaS)-1], disco_buscado.Path, sb)
	fmt.Fprintln(dot, nodograph)
	//fmt.Println(nodograph)
	//----------------------------
	fmt.Fprintln(dot, "\t\t"+`</table>>];`)
	fmt.Fprintln(dot, ""+`}`)

	dot.Close()

	// Generacion del reporte
	// imagePath := "Rep/" + name

	// cmd := exec.Command("/usr/bin/dot", "-T"+nombre_sin_extension[1], rutaB, "-o", imagePath)
	// err = cmd.Run()
	// if err != nil {
	// color.Red("[REP]: Error al generar imagen")
	// return
	// }

	color.Green("[REP]: LS «" + name + "» generated Sucessfull")
}
