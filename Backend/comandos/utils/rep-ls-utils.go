package utils

import (
	"encoding/binary"
	"fmt"
	"os"
	"proyecto/estructuras/size"
	"proyecto/estructuras/structures"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func LsInodo(pos int32, name1 string, path string, sb structures.SuperBloque) string {
	dot := ``
	inodo := structures.TablaInodo{}
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[-]: Error al abrir archivo")
		return ""
	}
	defer file.Close()

	if _, err := file.Seek(int64(pos), 0); err != nil {
		color.Red("[-]: Error en mover puntero")
		return ""
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[-]: Error en la lectura de la Tabla de Inodos")
		return ""
	}

	dot += "\t\t\t" + `<tr>`
	dot += "\t\t\t\t" + `<td>` + GetPermiso(inodo.I_perm) + `</td>`
	dot += "\t\t\t\t" + `<td>` + GetUsario(inodo.I_uid, path, sb) + `</td>`
	dot += "\t\t\t\t" + `<td>` + GetGrupo(inodo.I_gid, path, sb) + `</td>`
	dot += "\t\t\t\t" + `<td>` + fmt.Sprint(inodo.I_s) + `</td>`
	dot += "\t\t\t\t" + `<td>` + IntFechaToStr(inodo.I_mtime) + `</td>`
	if inodo.I_type == 1 {
		dot += "\t\t\t\t" + `<td>` + "Archivo" + `</td>`
		dot += "\t\t\t\t" + `<td>` + name1 + `</td>`
		dot += "\t\t\t" + `</tr>`
		return dot
	}
	dot += "\t\t\t\t" + `<td>` + "Carpeta" + `</td>`
	dot += "\t\t\t\t" + `<td>` + name1 + `</td>`
	dot += "\t\t\t" + `</tr>`

	apuntador1, apuntador2, apuntador3 := structures.BloqueApuntador{}, structures.BloqueApuntador{}, structures.BloqueApuntador{}
	carpeta := structures.BloqueCarpeta{}

	for i := 0; i < 15; i++ {
		if inodo.I_block[i] != -1 {
			if i < 12 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[-]: Error en mover puntero")
					return ""
				}
				if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
					color.Red("[-]: Error en la lectura de la Tabla de Bloque de Carpetas")
					return ""
				}
				for c := 0; c < 4; c++ {
					if carpeta.B_content[c].B_inodo != -1 {
						name1 := ToString(carpeta.B_content[c].B_name[:])
						if !(name1 == "." || name1 == "..") {
							dot += LsInodo(carpeta.B_content[c].B_inodo, name1, path, sb)
						}
					}
				}
			} else if i == 12 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[-]: Error en mover puntero")
					return ""
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[-]: Error en la lectura de la Tabla de Bloque Apuntador 1")
					return ""
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[-]: Error en mover puntero")
							return ""
						}
						if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
							color.Red("[-]: Error en la lectura de la Tabla de Bloque Apuntador 2")
							return ""
						}
						for c := 0; c < 4; c++ {
							if carpeta.B_content[c].B_inodo != -1 {
								name1 := ToString(carpeta.B_content[c].B_name[:])
								dot += LsInodo(carpeta.B_content[c].B_inodo, name1, path, sb)
							}
						}
					}
				}
			} else if i == 13 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[-]: Error en mover puntero")
					return ""
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[-]: Error en la lectura de la Tabla de Bloque Apuntador 1")
					return ""
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[i] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[-]: Error en mover puntero")
							return ""
						}
						if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
							color.Red("[-]: Error en la lectura de la Tabla de Bloque Apuntador 2")
							return ""
						}
						for k := 0; k < 16; k++ {
							if apuntador2.B_pointers[k] != -1 {
								if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
									color.Red("[-]: Error en mover puntero")
									return ""
								}
								if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
									color.Red("[-]: Error en la lectura de la Tabla de Bloque Carpeta")
									return ""
								}
								for c := 0; c < 4; c++ {
									if carpeta.B_content[c].B_inodo != -1 {
										name1 := ToString(carpeta.B_content[c].B_name[:])
										dot += LsInodo(carpeta.B_content[c].B_inodo, name1, path, sb)
									}
								}
							}
						}
					}
				}
			} else if i == 14 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[-]: Error en mover puntero")
					return ""
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[-]: Error en la lectura de la Tabla de Bloque Apuntador 1")
					return ""
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[-]: Error en mover puntero")
							return ""
						}
						if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
							color.Red("[-]: Error en la lectura de la Tabla de Bloque Apuntador 2")
							return ""
						}
						for k := 0; k < 16; k++ {
							if apuntador2.B_pointers[k] != -1 {
								if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
									color.Red("[-]: Error en mover puntero")
									return ""
								}
								if err := binary.Read(file, binary.LittleEndian, &apuntador3); err != nil {
									color.Red("[-]: Error en la lectura de la Tabla de Bloque Apuntador 3")
									return ""
								}
								for l := 0; l < 16; l++ {
									if apuntador3.B_pointers[l] != -1 {
										if _, err := file.Seek(int64(apuntador3.B_pointers[l]), 0); err != nil {
											color.Red("[-]: Error en mover puntero")
											return ""
										}
										if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
											color.Red("[-]: Error en la lectura de la Tabla de Bloque Carpeta")
											return ""
										}
										for c := 0; c < 4; c++ {
											if carpeta.B_content[c].B_inodo != -1 {
												name1 := ToString(carpeta.B_content[c].B_name[:])
												dot += LsInodo(carpeta.B_content[c].B_inodo, name1, path, sb)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return dot
}

func GetPermiso(permiso int32) string {
	perm := strconv.Itoa(int(permiso))
	dot := ""
	for i := 0; i < len(perm); i++ {
		dot += "-"
		if (perm[i] == '4') || (perm[i] == '5') || (perm[i] == '6') || (perm[i] == '7') {
			dot += "r"
		}
		if (perm[i] == '2') || (perm[i] == '3') || (perm[i] == '6') || (perm[i] == '7') {
			dot += "w"
		}
		if (perm[i] == '1') || (perm[i] == '3') || (perm[i] == '5') || (perm[i] == '7') {
			dot += "x"
		}
		dot += " "
	}
	return dot
}

func GetGrupo(id int32, path string, superbloque structures.SuperBloque) string {
	content := GetContentReport(superbloque.S_inode_start+size.SizeTablaInodo(), path)

	split_group := strings.Split(content, "\n")
	for _, con := range split_group {
		if strings.Contains(con, ",G,") {
			grupo := strings.Split(con, ",")
			if grupo[0] == strconv.Itoa(int(id)) {
				return grupo[2]
			}
		}
	}

	fmt.Println(content)
	return "NONE"
}

func GetUsario(id int32, path string, superbloque structures.SuperBloque) string {
	content := GetContentReport(superbloque.S_inode_start+size.SizeTablaInodo(), path)

	split_group := strings.Split(content, "\n")
	for _, con := range split_group {
		if strings.Contains(con, ",U,") {
			grupo := strings.Split(con, ",")
			if grupo[0] == strconv.Itoa(int(id)) {
				return grupo[3]
			}
		}
	}

	fmt.Println(content)
	return "NONE"
}
