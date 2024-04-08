package obtencionpf

import (
	"encoding/binary"
	"fmt"
	"os"
	"proyecto/comandos/utils"
	"proyecto/estructuras/structures"
	"strings"
)

func ObtenerPaths(id_disco string) (string, error) {
	disco_buscado, edb := utils.ObtenerDiscoID(id_disco)
	if !edb {
		return "", fmt.Errorf("Error")
	}

	inicioSB := disco_buscado.Particion_P.Part_start
	sb := structures.SuperBloque{}
	file, err := os.OpenFile(disco_buscado.Path, os.O_RDWR, 0666)
	if err != nil {
		//color.Red("[Mount]: Error al abrir archivo")
		return "", fmt.Errorf("Error")
	}
	defer file.Close()

	if _, err := file.Seek(int64(inicioSB), 0); err != nil {
		return "", fmt.Errorf("Error")
	}
	if err := binary.Read(file, binary.LittleEndian, &sb); err != nil {
		return "", fmt.Errorf("Error")
	}

	inodo := structures.TablaInodo{}
	if _, err := file.Seek(int64(sb.S_inode_start), 0); err != nil {
		return "", fmt.Errorf("Error")
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		return "", fmt.Errorf("Error")
	}

	// fmt.Println(inodo)

	con := "["

	for i := 0; i < 15; i++ {
		if inodo.I_block[i] != -1 {
			if i < 12 {
				if inodo.I_type == 0 {
					//carpeta
					carpeta := structures.BloqueCarpeta{}
					if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
						return "", fmt.Errorf("Error")
					}
					if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
						return "", fmt.Errorf("Error")
					}
					// fmt.Println("carpeta", carpeta)
					for c := 0; c < 4; c++ {
						if carpeta.B_content[c].B_inodo != -1 {
							name := utils.ToString(carpeta.B_content[c].B_name[:])
							if (carpeta.B_content[c].B_inodo != -1) && (name != "." && name != "..") {
								inodo2 := structures.TablaInodo{}
								if _, err := file.Seek(int64(carpeta.B_content[c].B_inodo), 0); err != nil {
									return "", fmt.Errorf("Error")
								}
								if err := binary.Read(file, binary.LittleEndian, &inodo2); err != nil {
									return "", fmt.Errorf("Error")
								}
								if inodo2.I_perm != 0 {
									if inodo2.I_type == 0 {
										con += "{"
										con += "\"nombre\": \"" + utils.ToString(carpeta.B_content[c].B_name[:]) + "\","
										con += "\"tipo\": \"carpeta\","
										con += "\"contenido\": ["
										con += CarpetaArchivo(disco_buscado.Path, carpeta.B_content[c].B_inodo)
										con += "]"
										con += "}"

										// fmt.Println(inodo2)
									} else if inodo2.I_type == 1 {
										con += "{"
										con += "\"nombre\": \"" + utils.ToString(carpeta.B_content[c].B_name[:]) + "\","
										con += "\"tipo\": \"archivo\","
										con += "\"contenido\": \""
										con += Archivo(disco_buscado.Path, carpeta.B_content[c].B_inodo)
										con += "\""
										con += "}"
										// fmt.Println(inodo2)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	con += "]"

	con = strings.ReplaceAll(con, "}{", "},{")
	return con, nil
}
