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

func Report_TREE(name string, path string, ruta string, id_disco string) {

	disco_buscado, edb := utils.ObtenerDiscoID(id_disco)
	if !edb {
		return
	}

	nombre_sin_extension := strings.Split(name, ".")
	rutaB := "Rep/" + nombre_sin_extension[0] + ".dot"
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
	// apuntador1, apuntador2, apuntador3 := structures.BloqueApuntador{}, structures.BloqueApuntador{}, structures.BloqueApuntador{}
	var bit byte
	cont := int32(0)

	//*******************
	fmt.Fprintln(dot, ""+`digraph G {`)
	fmt.Fprintln(dot, "\t"+`rankdir=LR;`)
	fmt.Fprintln(dot, "\t"+`node[shape=none];`)
	//--/--/--/--/--/--/--/
	//*****//*****//*****//*****
	for i := start; i < end; i++ {
		if _, err := file.Seek(int64(i), 0); err != nil {
			color.Red("[REP]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &bit); err != nil {
			color.Red("[REP]: Error en la lectura del byte")
			return
		}
		if bit == '1' {
			posInodo := sb.S_inode_start + (cont * size.SizeTablaInodo())
			if _, err := file.Seek(int64(sb.S_inode_start+(cont*size.SizeTablaInodo())), 0); err != nil {
				color.Red("[REP]: Error en mover puntero")
				return
			}
			if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
				color.Red("[REP]: Error en la lectura del byte")
				return
			}
			fmt.Fprintln(dot, utils.TreeInodo(posInodo, disco_buscado.Path))
			for i := 0; i < 15; i++ {
				if inodo.I_block[i] != -1 {
					if i < 12 {
						if inodo.I_type == 0 {
							fmt.Fprintln(dot, utils.TreeBlock(inodo.I_block[i], 0, disco_buscado.Path))
						} else if inodo.I_type == 1 {
							fmt.Fprintln(dot, utils.TreeBlock(inodo.I_block[i], 1, disco_buscado.Path))
						}
						fmt.Fprintln(dot, utils.Conexiones(posInodo, inodo.I_block[i]))
					}
				}
			}
		}
		cont++
	}
	//*****//*****//*****//*****
	//--/--/--/--/--/--/--/
	fmt.Fprintln(dot, ""+`}`)
	//*******************
	dot.Close()

	// Generacion del reporte
	// imagePath := "Rep/" + name

	// cmd := exec.Command("/usr/bin/dot", "-T"+nombre_sin_extension[1], rutaB, "-o", imagePath)
	// err = cmd.Run()
	// if err != nil {
	// 	color.Red("[REP]: Error al generar imagen")
	// 	return
	// }

	color.Green("[REP]: Tree «" + name + "» generated Sucessfull")
}
