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

func Report_SUPERBLOCK(name string, path string, ruta string, id_disco string) {
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
	//------------

	fmt.Fprintln(dot, ""+`digraph G {`)
	fmt.Fprintln(dot, "\t"+`node[shape=none];`)
	fmt.Fprintln(dot, "\t"+`start[label=<`)
	fmt.Fprintln(dot, "\t\t"+`<table>`)
	//------------
	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td colspan="2" bgcolor="#927d55"><font point-size="20" color="white"><b>REPORTE DE SUPERBLOQUE</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#555092" width="300"><font color="gray"><b>sb_nombre_hd</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#555092" width="300"><font color="gray"><b>`+string(disco_buscado.DriveLetter)+`</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#2B228C" width="300"><font color="gray"><b>s_filesystem_type</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#2B228C" width="300"><font color="gray"><b>`+fmt.Sprint(sb.S_filesistem_type)+`</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#555092" width="300"><font color="gray"><b>s_inodes_count</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#555092" width="300"><font color="gray"><b>`+fmt.Sprint(sb.S_inodes_count)+`</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#2B228C" width="300"><font color="gray"><b>s_blocks_count</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#2B228C" width="300"><font color="gray"><b>`+fmt.Sprint(sb.S_blocks_count)+`</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#555092" width="300"><font color="gray"><b>s_free_blocks_count</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#555092" width="300"><font color="gray"><b>`+fmt.Sprint(sb.S_free_blocks_count)+`</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#2B228C" width="300"><font color="gray"><b>s_free_inodes_count</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#2B228C" width="300"><font color="gray"><b>`+fmt.Sprint(sb.S_free_inodes_count)+`</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#555092" width="300"><font color="gray"><b>s_mtime</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#555092" width="300"><font color="gray"><b>`+utils.IntFechaToStr(sb.S_mtime)+`</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#2B228C" width="300"><font color="gray"><b>s_umtime</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#2B228C" width="300"><font color="gray"><b>`+utils.IntFechaToStr(sb.S_umtime)+`</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#555092" width="300"><font color="gray"><b>s_mnt_count</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#555092" width="300"><font color="gray"><b>`+fmt.Sprint(sb.S_mnt_count)+`</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#2B228C" width="300"><font color="gray"><b>s_magic</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#2B228C" width="300"><font color="gray"><b>`+fmt.Sprint(sb.S_magic)+`</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#555092" width="300"><font color="gray"><b>s_inode_s</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#555092" width="300"><font color="gray"><b>`+fmt.Sprint(sb.S_inode_s)+`</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#2B228C" width="300"><font color="gray"><b>s_block_s</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#2B228C" width="300"><font color="gray"><b>`+fmt.Sprint(sb.S_block_s)+`</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#555092" width="300"><font color="gray"><b>s_first_ino</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#555092" width="300"><font color="gray"><b>`+fmt.Sprint(sb.S_first_ino)+`</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#2B228C" width="300"><font color="gray"><b>s_first_blo</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#2B228C" width="300"><font color="gray"><b>`+fmt.Sprint(sb.S_first_blo)+`</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#555092" width="300"><font color="gray"><b>s_bm_inode_start</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#555092" width="300"><font color="gray"><b>`+fmt.Sprint(sb.S_bm_inode_start)+`</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#2B228C" width="300"><font color="gray"><b>s_bm_block_start</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#2B228C" width="300"><font color="gray"><b>`+fmt.Sprint(sb.S_bm_block_start)+`</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#555092" width="300"><font color="gray"><b>s_inode_start</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#555092" width="300"><font color="gray"><b>`+fmt.Sprint(sb.S_inode_start)+`</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)

	fmt.Fprintln(dot, "\t\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#2B228C" width="300"><font color="gray"><b>s_block_start</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t\t"+`<td bgcolor="#2B228C" width="300"><font color="gray"><b>`+fmt.Sprint(sb.S_block_start)+`</b></font></td>`)
	fmt.Fprintln(dot, "\t\t\t"+`</tr>`)
	//------------
	fmt.Fprintln(dot, "\t\t"+`</table>`)
	fmt.Fprintln(dot, "\t"+`>];`)
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

	color.Green("[REP]: SuperBlock «" + name + "» generated Sucessfull")
}
