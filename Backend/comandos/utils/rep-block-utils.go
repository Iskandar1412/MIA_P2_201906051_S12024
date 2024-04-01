package utils

import (
	"encoding/binary"
	"fmt"
	"os"
	"proyecto/estructuras/structures"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func GraphBlockCarpeta(pos int32, path string) string {
	dot := ``
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[+]: Error al abrir archivo")
		return ""
	}
	defer file.Close()

	carpeta := structures.BloqueCarpeta{}
	if _, err := file.Seek(int64(pos), 0); err != nil {
		color.Red("[+]: Error en mover puntero")
		return ""
	}
	if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
		color.Red("[+]: Error en la lectura del Bloque de Carpeta")
		return ""
	}

	dot += "\tn" + strconv.Itoa(int(pos)) + `[label=<<table>` + "\n"
	dot += "\t\t" + `<tr>` + "\n"
	dot += "\t\t\t" + `<td colspan="2" bgcolor="#c3f8b6">Bloque Carpeta</td>` + "\n"
	dot += "\t\t" + `</tr>` + "\n"
	//-------
	for i := 0; i < 4; i++ {
		dot += "\t\t" + `<tr>` + "\n"
		dot += "\t\t\t" + `<td colspan="2" bgcolor="#b6f8d3">b_content  ` + strconv.Itoa(i) + `</td>` + "\n"
		dot += "\t\t" + `</tr>` + "\n"

		dot += "\t\t" + `<tr>` + "\n"
		dot += "\t\t\t" + `<td>b_name</td>` + "\n"
		dot += "\t\t\t" + `<td>` + Returnstring(ToString(carpeta.B_content[i].B_name[:])) + `</td>` + "\n"
		dot += "\t\t" + `</tr>` + "\n"

		dot += "\t\t" + `<tr>` + "\n"
		dot += "\t\t\t" + `<td>b_inodo</td>` + "\n"
		dot += "\t\t\t" + `<td>` + fmt.Sprint(carpeta.B_content[i].B_inodo) + `</td>` + "\n"
		dot += "\t\t" + `</tr>` + "\n"
	}
	//-------
	dot += "\t" + `</table>>];` + "\n"
	return dot
}

func GraphBlockArchivo(pos int32, path string) string {
	dot := ``
	// content := ``
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[+]: Error al abrir archivo")
		return ""
	}
	defer file.Close()

	archivo := structures.BloqueArchivo{}
	if _, err := file.Seek(int64(pos), 0); err != nil {
		color.Red("[+]: Error en mover puntero")
		return ""
	}
	if err := binary.Read(file, binary.LittleEndian, &archivo); err != nil {
		color.Red("[+]: Error en la lectura del Bloque de Archivo")
		return ""
	}

	dot += "\tn" + strconv.Itoa(int(pos)) + `[label=<<table>` + "\n"
	dot += "\t\t" + `<tr>` + "\n"
	dot += "\t\t\t" + `<td bgcolor="#c3f8b6" width="700">Bloque Archivo</td>` + "\n"
	dot += "\t\t" + `</tr>` + "\n"
	//-------
	dot += "\t\t" + `<tr>` + "\n"
	str := ToString(archivo.B_content[:])
	dot += "\t\t\t" + `<td>` + Returnstring(strings.ReplaceAll(str, "\n", "<br/>")) + `</td>` + "\n"
	dot += "\t\t" + `</tr>` + "\n"
	//-------
	dot += "\t" + `</table>>];` + "\n"
	return dot
}

func GraphBlockApuntador(pos int32, path string) string {
	dot := ``
	// content := ``
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[+]: Error al abrir archivo")
		return ""
	}
	defer file.Close()

	apuntador := structures.BloqueApuntador{}
	if _, err := file.Seek(int64(pos), 0); err != nil {
		color.Red("[+]: Error en mover puntero")
		return ""
	}
	if err := binary.Read(file, binary.LittleEndian, &apuntador); err != nil {
		color.Red("[+]: Error en la lectura del Bloque de Apuntadores")
		return ""
	}

	dot += "\tn" + strconv.Itoa(int(pos)) + `[label=<<table>` + "\n"
	dot += "\t\t" + `<tr>` + "\n"
	dot += "\t\t\t" + `<td colspan="2" bgcolor="#c3f8b6">Bloque Apuntadores</td>` + "\n"
	dot += "\t\t" + `</tr>` + "\n"
	//-------
	for i := 0; i < 16; i++ {
		dot += "\t\t" + `<tr>` + "\n"
		dot += "\t\t\t" + `<td>` + "b_pointer " + fmt.Sprint(i) + `</td>` + "\n"
		dot += "\t\t\t" + `<td>` + fmt.Sprint(apuntador.B_pointers[i]) + `</td>` + "\n"
		dot += "\t\t" + `</tr>` + "\n"
	}
	//-------
	dot += "\t" + `</table>>];` + "\n"
	return dot
}
