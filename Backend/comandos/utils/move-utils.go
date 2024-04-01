package utils

import (
	"encoding/binary"
	"os"
	"proyecto/comandos/global"
	"proyecto/estructuras/structures"

	"github.com/fatih/color"
)

func GetPosInodo(inodo structures.TablaInodo, nombreBuscar string) int32 {
	var carpeta structures.BloqueCarpeta
	var apuntador1, apuntador2, apuntador3 structures.BloqueApuntador
	nodo := global.UsuarioLogeado.Mounted
	file, err := os.OpenFile(nodo.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[RENAME]: Error al abrir archivo")
		return -1
	}
	defer file.Close()

	for i := 0; i < 15; i++ {
		if inodo.I_block[i] != -1 {
			if i < 12 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[RENAME]: Error en mover puntero")
					return -1
				}
				if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
					color.Red("[RENAME]: Error en la lectura del archivo")
					return -1
				}
				for c := 0; c < 15; c++ {
					if ToString(carpeta.B_content[c].B_name[:]) == nombreBuscar {
						aux := carpeta.B_content[c].B_inodo
						carpeta.B_content[c].B_inodo = -1
						if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
							color.Red("[REMOVE]: Error en mover puntero")
							return -1
						}
						if err := binary.Write(file, binary.LittleEndian, &carpeta); err != nil {
							color.Red("[/]: Error en la escritura del archivo")
							return -1
						}
						return aux
					}
				}
			} else if i == 12 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[RENAME]: Error en mover puntero")
					return -1
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[RENAME]: Error en la lectura del archivo")
					return -1
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[RENAME]: Error en mover puntero")
							return -1
						}
						if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
							color.Red("[RENAME]: Error en la lectura del archivo")
							return -1
						}
						for c := 0; c < 15; c++ {
							if ToString(carpeta.B_content[c].B_name[:]) == nombreBuscar {
								aux := carpeta.B_content[c].B_inodo
								carpeta.B_content[c].B_inodo = -1
								if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
									color.Red("[REMOVE]: Error en mover puntero")
									return -1
								}
								if err := binary.Write(file, binary.LittleEndian, &carpeta); err != nil {
									color.Red("[/]: Error en la escritura del archivo")
									return -1
								}
								return aux
							}
						}
					}
				}
			} else if i == 13 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[RENAME]: Error en mover puntero")
					return -1
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[RENAME]: Error en la lectura del archivo")
					return -1
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[RENAME]: Error en mover puntero")
							return -1
						}
						if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
							color.Red("[RENAME]: Error en la lectura del archivo")
							return -1
						}
						for k := 0; k < 16; k++ {
							if apuntador2.B_pointers[k] != -1 {
								if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
									color.Red("[RENAME]: Error en mover puntero")
									return -1
								}
								if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
									color.Red("[RENAME]: Error en la lectura del archivo")
									return -1
								}
								for c := 0; c < 15; c++ {
									if ToString(carpeta.B_content[c].B_name[:]) == nombreBuscar {
										aux := carpeta.B_content[c].B_inodo
										carpeta.B_content[c].B_inodo = -1
										if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
											color.Red("[REMOVE]: Error en mover puntero")
											return -1
										}
										if err := binary.Write(file, binary.LittleEndian, &carpeta); err != nil {
											color.Red("[/]: Error en la escritura del archivo")
											return -1
										}
										return aux
									}
								}
							}
						}
					}
				}
			} else if i == 14 {
				if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
					color.Red("[RENAME]: Error en mover puntero")
					return -1
				}
				if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
					color.Red("[RENAME]: Error en la lectura del archivo")
					return -1
				}
				for j := 0; j < 16; j++ {
					if apuntador1.B_pointers[j] != -1 {
						if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
							color.Red("[RENAME]: Error en mover puntero")
							return -1
						}
						if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
							color.Red("[RENAME]: Error en la lectura del archivo")
							return -1
						}
						for k := 0; k < 16; k++ {
							if apuntador2.B_pointers[k] != -1 {
								if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
									color.Red("[RENAME]: Error en mover puntero")
									return -1
								}
								if err := binary.Read(file, binary.LittleEndian, &apuntador3); err != nil {
									color.Red("[RENAME]: Error en la lectura del archivo")
									return -1
								}
								for l := 0; l < 16; l++ {
									if apuntador3.B_pointers[l] != -1 {
										if _, err := file.Seek(int64(apuntador3.B_pointers[l]), 0); err != nil {
											color.Red("[RENAME]: Error en mover puntero")
											return -1
										}
										if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
											color.Red("[RENAME]: Error en la lectura del archivo")
											return -1
										}
										for c := 0; c < 15; c++ {
											if ToString(carpeta.B_content[c].B_name[:]) == nombreBuscar {
												aux := carpeta.B_content[c].B_inodo
												carpeta.B_content[c].B_inodo = -1
												if _, err := file.Seek(int64(apuntador3.B_pointers[l]), 0); err != nil {
													color.Red("[REMOVE]: Error en mover puntero")
													return -1
												}
												if err := binary.Write(file, binary.LittleEndian, &carpeta); err != nil {
													color.Red("[/]: Error en la escritura del archivo")
													return -1
												}
												return aux
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
	return -1
}
