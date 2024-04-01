package utils

import (
	"encoding/binary"
	"os"
	"proyecto/comandos/global"
	"proyecto/estructuras/structures"

	"github.com/fatih/color"
)

func ChmodR(posI int32, u int32, r bool) {
	var inodo structures.TablaInodo
	var carpeta structures.BloqueCarpeta
	var apuntador1, apuntador2, apuntador3 structures.BloqueApuntador
	path := global.UsuarioLogeado.Mounted.Path
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return
	}
	defer file.Close()
	if _, err := file.Seek(int64(posI), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[/]: Error en la lectura del archivo")
		return
	}
	inodo.I_mtime = ObFechaInt()
	inodo.I_perm = u
	if _, err := file.Seek(int64(posI), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return
	}
	if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[/]: Error en la escritura del archivo")
		return
	}

	if (inodo.I_type == 0) && r {
		for i := 0; i < 15; i++ {
			if inodo.I_block[i] != -1 {
				if i < 12 {
					if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return
					}
					if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
						color.Red("[/]: Error en la lectura del archivo")
						return
					}
					for c := 0; c < 4; c++ {
						nombre := ToString(carpeta.B_content[c].B_name[:])
						if (carpeta.B_content[c].B_inodo != -1) && (nombre != "." && nombre != "..") {
							ChmodR(carpeta.B_content[c].B_inodo, u, r)
						}
					}
				} else if i == 12 {
					if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return
					}
					if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
						color.Red("[/]: Error en la lectura del archivo")
						return
					}
					for j := 0; j < 16; j++ {
						if apuntador1.B_pointers[j] != -1 {
							if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
								color.Red("[/]: Error en mover puntero")
								return
							}
							if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
								color.Red("[/]: Error en la lectura del archivo")
								return
							}
							for c := 0; c < 4; c++ {
								nombre := ToString(carpeta.B_content[c].B_name[:])
								if (carpeta.B_content[c].B_inodo != -1) && (nombre != "." && nombre != "..") {
									ChmodR(carpeta.B_content[c].B_inodo, u, r)
								}
							}
						}
					}
				} else if i == 13 {
					if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return
					}
					if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
						color.Red("[/]: Error en la lectura del archivo")
						return
					}
					for j := 0; j < 16; j++ {
						if apuntador1.B_pointers[j] != -1 {
							if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
								color.Red("[/]: Error en mover puntero")
								return
							}
							if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
								color.Red("[/]: Error en la lectura del archivo")
								return
							}
							for k := 0; k < 16; k++ {
								if apuntador2.B_pointers[k] != -1 {
									if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
										color.Red("[/]: Error en mover puntero")
										return
									}
									if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
										color.Red("[/]: Error en la lectura del archivo")
										return
									}
									for c := 0; c < 4; c++ {
										nombre := ToString(carpeta.B_content[c].B_name[:])
										if (carpeta.B_content[c].B_inodo != -1) && (nombre != "." && nombre != "..") {
											ChmodR(carpeta.B_content[c].B_inodo, u, r)
										}
									}
								}
							}
						}
					}
				} else if i == 14 {
					if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
						color.Red("[/]: Error en mover puntero")
						return
					}
					if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
						color.Red("[/]: Error en la lectura del archivo")
						return
					}
					for j := 0; j < 16; j++ {
						if apuntador1.B_pointers[j] != -1 {
							if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
								color.Red("[/]: Error en mover puntero")
								return
							}
							if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
								color.Red("[/]: Error en la lectura del archivo")
								return
							}
							for k := 0; k < 16; k++ {
								if apuntador2.B_pointers[k] != -1 {
									if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
										color.Red("[/]: Error en mover puntero")
										return
									}
									if err := binary.Read(file, binary.LittleEndian, &apuntador3); err != nil {
										color.Red("[/]: Error en la lectura del archivo")
										return
									}
									for l := 0; l < 16; l++ {
										if apuntador3.B_pointers[l] != -1 {
											if _, err := file.Seek(int64(apuntador3.B_pointers[l]), 0); err != nil {
												color.Red("[/]: Error en mover puntero")
												return
											}
											if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
												color.Red("[/]: Error en la lectura del archivo")
												return
											}
											for c := 0; c < 4; c++ {
												nombre := ToString(carpeta.B_content[c].B_name[:])
												if (carpeta.B_content[c].B_inodo != -1) && (nombre != "." && nombre != "..") {
													ChmodR(carpeta.B_content[c].B_inodo, u, r)
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
	}
}
