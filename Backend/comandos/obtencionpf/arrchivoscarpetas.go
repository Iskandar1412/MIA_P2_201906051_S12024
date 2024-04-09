package obtencionpf

import (
	"encoding/binary"
	"os"
	"proyecto/comandos/utils"
	"proyecto/estructuras/structures"
	"strings"
)

func CarpetaArchivo(path string, pos int32) string {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return ""
	}
	defer file.Close()
	inodo := structures.TablaInodo{}
	if _, err := file.Seek(int64(pos), 0); err != nil {
		return ""
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		return ""
	}
	// fmt.Println("carpetaarchivo", inodo)
	con := ""

	for i := 0; i < 15; i++ {
		if inodo.I_block[i] != -1 {
			if i < 12 {
				if inodo.I_type == 0 {
					//carpeta
					carpeta := structures.BloqueCarpeta{}
					if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
						return ""
					}
					if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
						return ""
					}
					for c := 0; c < 4; c++ {
						if carpeta.B_content[c].B_inodo != -1 {
							name := utils.ToString(carpeta.B_content[c].B_name[:])
							if (carpeta.B_content[c].B_inodo != -1) && (name != "." && name != "..") {
								// fmt.Println(utils.ToString(carpeta.B_content[c].B_name[:]))
								inodo2 := structures.TablaInodo{}
								if _, err := file.Seek(int64(carpeta.B_content[c].B_inodo), 0); err != nil {
									return ""
								}
								if err := binary.Read(file, binary.LittleEndian, &inodo2); err != nil {
									return ""
								}
								// fmt.Println("consoel", inodo2)
								if inodo2.I_perm != 0 {
									if inodo2.I_type == 0 {
										con += "{"
										con += "\"nombre\": \"" + utils.ToString(carpeta.B_content[c].B_name[:]) + "\","
										con += "\"tipo\": \"carpeta\","
										con += "\"contenido\": ["
										con += CarpetaArchivo(path, carpeta.B_content[c].B_inodo)
										con += "]"
										con += "}"
									} else if inodo2.I_type == 1 {
										con += "{"
										con += "\"nombre\": \"" + utils.ToString(carpeta.B_content[c].B_name[:]) + "\","
										con += "\"tipo\": \"archivo\","
										con += "\"contenido\": \""
										con += Archivo(path, carpeta.B_content[c].B_inodo)
										con += "\""
										con += "}"
									}
								}
							}
						}
					}
				}
			}
		}
	}

	con += ""
	return con
}

func Archivo(path string, pos int32) string {
	dot := ""
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return ""
	}
	defer file.Close()
	inodo := structures.TablaInodo{}
	archivo := structures.BloqueArchivo{}
	if _, err := file.Seek(int64(pos), 0); err != nil {
		return ""
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		return ""
	}

	for i := 0; i < 15; i++ {
		if inodo.I_block[i] != -1 {
			if i < 12 {
				content := ""
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					return ""
				}
				if err := binary.Read(file, binary.LittleEndian, &archivo); err != nil {
					return ""
				}
				content = utils.ToString(archivo.B_content[:])
				dot += strings.ReplaceAll(content, "\n", "\\n")
			}
		}
	}
	dot += ""
	// fmt.Println("archivo", dot)
	return dot
}
