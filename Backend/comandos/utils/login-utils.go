package utils

import (
	"encoding/binary"
	"os"
	"proyecto/estructuras/structures"

	"github.com/fatih/color"
)

func GetContent(inodoStart int32, path string) (string, bool) {
	inodo := structures.TablaInodo{}
	archivo := structures.BloqueArchivo{}
	apuntador1, apuntador2, apuntador3 := structures.BloqueApuntador{}, structures.BloqueApuntador{}, structures.BloqueApuntador{}
	content := ""

	// Obtener tabla de inodos
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[*]: Error al abrir archivo")
		return "", false
	}
	if _, err := file.Seek(int64(inodoStart), 0); err != nil {
		color.Red("[*]: Error en mover puntero")
		return "", false
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[*]: Error en la lectura de la Tabla de Inodos")
		return "", false
	}
	for i := 0; i < 15; i++ {
		if inodo.I_block[i] != -1 {
			if i < 12 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[*]: Error en mover puntero")
					return "", false
				}
				if err := binary.Read(file, binary.LittleEndian, &archivo); err != nil {
					color.Red("[*]: Error en la lectura de la Tabla de Bloque de Archivos")
					return "", false
				}
				content += ToString(archivo.B_content[:])
			} else if i == 12 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[*]: Error en mover puntero")
					return "", false
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[*]: Error en la lectura de la Tabla de Bloque de Archivos")
					return "", false
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[*]: Error en mover puntero")
							return "", false
						}
						if err := binary.Read(file, binary.LittleEndian, &archivo); err != nil {
							color.Red("[*]: Error en la lectura de la Tabla de Bloque de Archivos")
							return "", false
						}
						content += ToString(archivo.B_content[:])
					}
				}
			} else if i == 13 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[*]: Error en mover puntero")
					return "", false
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[*]: Error en la lectura de la Tabla de Bloque de Apuntadores 1")
					return "", false
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[*]: Error en mover puntero")
							return "", false
						}
						if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
							color.Red("[*]: Error en la lectura de la Tabla de Bloque de Bloque de Apuntadores 2")
							return "", false
						}
						for k := 0; k < 16; k++ {
							if apuntador2.B_pointers[k] != -1 {
								if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
									color.Red("[*]: Error en mover puntero")
									return "", false
								}
								if err := binary.Read(file, binary.LittleEndian, &archivo); err != nil {
									color.Red("[*]: Error en la lectura de la Tabla de Bloque de Archivos")
									return "", false
								}
								content += ToString(archivo.B_content[:])
							}
						}
					}
				}
			} else if i == 14 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[*]: Error en mover puntero")
					return "", false
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[*]: Error en la lectura de la Tabla de Bloque de Apuntadores 1")
					return "", false
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[*]: Error en mover puntero")
							return "", false
						}
						if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
							color.Red("[*]: Error en la lectura de la Tabla de Bloque de Bloque de Apuntadores 2")
							return "", false
						}
						for k := 0; k < 16; k++ {
							if apuntador2.B_pointers[k] != -1 {
								if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
									color.Red("[*]: Error en mover puntero")
									return "", false
								}
								if err := binary.Read(file, binary.LittleEndian, &apuntador3); err != nil {
									color.Red("[*]: Error en la lectura de la Tabla de Bloque de Bloque de Apuntadores 3")
									return "", false
								}
								for l := 0; l < 16; l++ {
									if _, err := file.Seek(int64(apuntador3.B_pointers[l]), 0); err != nil {
										color.Red("[*]: Error en mover puntero")
										return "", false
									}
									if err := binary.Read(file, binary.LittleEndian, &archivo); err != nil {
										color.Red("[*]: Error en la lectura de la Tabla de Bloque de Archivos")
										return "", false
									}
									content += ToString(archivo.B_content[:])
								}
							}
						}
					}
				}
			}
		}
	}
	return content, true
}
