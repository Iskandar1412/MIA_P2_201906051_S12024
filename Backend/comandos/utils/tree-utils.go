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

func Conexiones(inicio int32, final int32) string {
	dot := "n" + fmt.Sprint(inicio) + ":" + fmt.Sprint(final) + "->n" + fmt.Sprint(final) + ";\n"
	return dot
}

func TreeBlock(pos int32, _type int32, path string) string {
	dot := ``
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[+]: Error al abrir archivo")
		return ""
	}
	defer file.Close()

	if _type == 0 {
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
		dot += "\t\t\t" + `<td colspan="2" bgcolor="#f34037">Bloque Carpeta ` + fmt.Sprint(pos) + `</td>` + "\n"
		dot += "\t\t" + `</tr>` + "\n"
		//--/--/--/--/--/--/--/--/--/
		//--/--/--/--/--/--/--/--/--/
		for i := 0; i < 4; i++ {
			dot += "\t\t" + `<tr>` + "\n"
			dot += "\t\t\t" + `<td>` + Returnstring(ToString(carpeta.B_content[i].B_name[:])) + `</td>` + "\n"
			dot += "\t\t\t" + `<td port="` + fmt.Sprint(carpeta.B_content[i].B_inodo) + `">` + fmt.Sprint(carpeta.B_content[i].B_inodo) + `</td>` + "\n"
			dot += "\t\t" + `</tr>` + "\n"
		}
		//--/--/--/--/--/--/--/--/--/
		//--/--/--/--/--/--/--/--/--/
		dot += "\t" + `</table>>];` + "\n"

		for i := 0; i < 4; i++ {
			nom := ToString(carpeta.B_content[i].B_name[:])
			if carpeta.B_content[i].B_inodo != -1 && (nom != "." && nom != "..") {
				dot += Conexiones(pos, carpeta.B_content[i].B_inodo)
			}
		}
	} else if _type == 1 {
		// contenido := ""
		archivo := structures.BloqueArchivo{}
		if _, err := file.Seek(int64(pos), 0); err != nil {
			color.Red("[+]: Error en mover puntero")
			return ""
		}
		if err := binary.Read(file, binary.LittleEndian, &archivo); err != nil {
			color.Red("[+]: Error en la lectura del Bloque de Archivo")
			return ""
		}

		contenido := ToString(archivo.B_content[:])
		contenido = strings.ReplaceAll(contenido, "\n", "<br/>")
		dot += "\tn" + strconv.Itoa(int(pos)) + `[label=<<table>` + "\n"
		dot += "\t\t" + `<tr>` + "\n"
		dot += "\t\t\t" + `<td bgcolor="#c3f8b6">Bloque Archivo</td>` + "\n"
		dot += "\t\t" + `</tr>` + "\n"
		dot += "\t\t" + `<tr>` + "\n"
		dot += "\t\t\t" + `<td>` + Returnstring(contenido) + `</td>` + "\n"
		dot += "\t\t" + `</tr>` + "\n"
		dot += "\t" + `</table>>];` + "\n"
	} else {
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
		for i := 0; i < 16; i++ {
			dot += "\t\t" + `<tr>` + "\n"
			dot += "\t\t\t" + `<td>` + "b_pointer " + fmt.Sprint(i) + `</td>` + "\n"
			dot += "\t\t\t" + `<td>` + fmt.Sprint(apuntador.B_pointers[i]) + `</td>` + "\n"
			dot += "\t\t" + `</tr>` + "\n"
		}
		dot += "\t" + `</table>>];` + "\n"
	}

	return dot
}

func TreeInodo(pos int32, path string) string {
	dot := ``
	inodo := structures.TablaInodo{}
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[+]: Error al abrir archivo")
		return ""
	}
	defer file.Close()

	if _, err := file.Seek(int64(pos), 0); err != nil {
		color.Red("[+]: Error en mover puntero")
		return ""
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[+]: Error en la lectura del Bloque de Apuntadores")
		return ""
	}

	dot += "\tn" + strconv.Itoa(int(pos)) + `[label=<<table>` + "\n"
	dot += "\t\t" + `<tr>` + "\n"
	dot += "\t\t\t" + `<td bgcolor="#376ef3" colspan="2">Inodo  ` + fmt.Sprint(pos) + `</td>` + "\n"
	dot += "\t\t" + `</tr>` + "\n"

	dot += "\t\t" + `<tr>` + "\n"
	dot += "\t\t\t" + `<td>i_uid</td>` + "\n"
	dot += "\t\t\t" + `<td>` + fmt.Sprint(inodo.I_uid) + `</td>` + "\n"
	dot += "\t\t" + `</tr>` + "\n"

	dot += "\t\t" + `<tr>` + "\n"
	dot += "\t\t\t" + `<td>i_gid</td>` + "\n"
	dot += "\t\t\t" + `<td>` + fmt.Sprint(inodo.I_gid) + `</td>` + "\n"
	dot += "\t\t" + `</tr>` + "\n"

	dot += "\t\t" + `<tr>` + "\n"
	dot += "\t\t\t" + `<td>i_s</td>` + "\n"
	dot += "\t\t\t" + `<td>` + fmt.Sprint(inodo.I_s) + `</td>` + "\n"
	dot += "\t\t" + `</tr>` + "\n"

	dot += "\t\t" + `<tr>` + "\n"
	dot += "\t\t\t" + `<td>i_atime</td>` + "\n"
	dot += "\t\t\t" + `<td>` + IntFechaToStr(inodo.I_atime) + `</td>` + "\n"
	dot += "\t\t" + `</tr>` + "\n"

	dot += "\t\t" + `<tr>` + "\n"
	dot += "\t\t\t" + `<td>i_ctime</td>` + "\n"
	dot += "\t\t\t" + `<td>` + IntFechaToStr(inodo.I_ctime) + `</td>` + "\n"
	dot += "\t\t" + `</tr>` + "\n"

	dot += "\t\t" + `<tr>` + "\n"
	dot += "\t\t\t" + `<td>i_mtime</td>` + "\n"
	dot += "\t\t\t" + `<td>` + IntFechaToStr(inodo.I_mtime) + `</td>` + "\n"
	dot += "\t\t" + `</tr>` + "\n"

	for j := 0; j < 15; j++ {
		if inodo.I_block[j] != -1 {
			dot += "\t\t" + `<tr>` + "\n"
			dot += "\t\t\t" + `<td>i_block ` + fmt.Sprint(j) + `</td>` + "\n"
			dot += "\t\t\t" + `<td port="` + fmt.Sprint(inodo.I_block[j]) + `">` + fmt.Sprint(inodo.I_block[j]) + `</td>` + "\n"
			dot += "\t\t" + `</tr>` + "\n"
		} else {
			dot += "\t\t" + `<tr>` + "\n"
			dot += "\t\t\t" + `<td>i_block</td>` + "\n"
			dot += "\t\t\t" + `<td>-1</td>` + "\n"
			dot += "\t\t" + `</tr>` + "\n"
		}
	}

	dot += "\t\t" + `<tr>` + "\n"
	dot += "\t\t\t" + `<td>i_type</td>` + "\n"
	dot += "\t\t\t" + `<td>` + fmt.Sprint(inodo.I_type) + `</td>` + "\n"
	dot += "\t\t" + `</tr>` + "\n"

	dot += "\t\t" + `<tr>` + "\n"
	dot += "\t\t\t" + `<td>i_perm</td>` + "\n"
	dot += "\t\t\t" + `<td>` + fmt.Sprint(inodo.I_perm) + `</td>` + "\n"
	dot += "\t\t" + `</tr>` + "\n"

	dot += "\t" + `</table>>];` + "\n"

	return dot
}
