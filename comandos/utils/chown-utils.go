package utils

import (
	"encoding/binary"
	"os"
	"proyecto/comandos/global"
	"proyecto/estructuras/structures"

	"github.com/fatih/color"
)

func ChownSystem(posI int32, idU int32, idG int32, r bool) {
	var inodo structures.TablaInodo
	var carpeta structures.BloqueCarpeta
	var apuntador1, apuntador2, apuntador3 structures.BloqueApuntador
	nodo := global.UsuarioLogeado.Mounted
	file, err := os.OpenFile(nodo.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[CHOWN]: Error al abrir archivo")
		return
	}
	defer file.Close()
	if _, err := file.Seek(int64(posI), 0); err != nil {
		color.Red("[CHOWN]: Error en mover puntero")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[CHOWN]: Error en la lectura del archivo")
		return
	}

	if (global.UsuarioLogeado.UID == 1) && (global.UsuarioLogeado.GID == 1) {
		// accion 1
		goto t0
	} else if (inodo.I_uid == global.UsuarioLogeado.UID) && (inodo.I_gid == global.UsuarioLogeado.GID) {
		//accion 1
		goto t0
	} else {
		//accion 2
		goto t1
	}

t0:
	inodo.I_mtime = ObFechaInt()
	inodo.I_uid = idU
	inodo.I_gid = idG
	if _, err := file.Seek(int64(posI), 0); err != nil {
		color.Red("[CHOWN]: Error en mover puntero")
		return
	}
	if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[CHOWN]: Error en la escritura del archivo")
		return
	}

t1:
	if (inodo.I_type == 0) && r {
		for i := 0; i < 15; i++ {
			if inodo.I_block[i] != -1 {
				if i < 12 {
					if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
						color.Red("[CHOWN]: Error en mover puntero")
						return
					}
					if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
						color.Red("[CHOWN]: Error en la lectura del archivo")
						return
					}
					for c := 0; c < 4; c++ {
						name := ToString(carpeta.B_content[c].B_name[:])
						if (carpeta.B_content[c].B_inodo != -1) && (name != "." && name != "..") {
							ChownSystem(carpeta.B_content[c].B_inodo, idU, idG, r)
						}
					}
				} else if i == 12 {
					if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
						color.Red("[CHOWN]: Error en mover puntero")
						return
					}
					if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
						color.Red("[CHOWN]: Error en la lectura del archivo")
						return
					}
					for j := 0; j < 16; j++ {
						if apuntador1.B_pointers[j] != -1 {
							if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
								color.Red("[CHOWN]: Error en mover puntero")
								return
							}
							if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
								color.Red("[CHOWN]: Error en la lectura del archivo")
								return
							}
							for c := 0; c < 4; c++ {
								name := ToString(carpeta.B_content[c].B_name[:])
								if (carpeta.B_content[c].B_inodo != -1) && (name != "." && name != "..") {
									ChownSystem(carpeta.B_content[c].B_inodo, idU, idG, r)
								}
							}
						}
					}
				} else if i == 13 {
					if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
						color.Red("[CHOWN]: Error en mover puntero")
						return
					}
					if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
						color.Red("[CHOWN]: Error en la lectura del archivo")
						return
					}
					for j := 0; j < 16; j++ {
						if apuntador1.B_pointers[j] != -1 {
							if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
								color.Red("[CHOWN]: Error en mover puntero")
								return
							}
							if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
								color.Red("[CHOWN]: Error en la lectura del archivo")
								return
							}
							for k := 0; k < 16; k++ {
								if apuntador2.B_pointers[k] != -1 {
									if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
										color.Red("[CHOWN]: Error en mover puntero")
										return
									}
									if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
										color.Red("[CHOWN]: Error en la lectura del archivo")
										return
									}
									for c := 0; c < 4; c++ {
										name := ToString(carpeta.B_content[c].B_name[:])
										if (carpeta.B_content[c].B_inodo != -1) && (name != "." && name != "..") {
											ChownSystem(carpeta.B_content[c].B_inodo, idU, idG, r)
										}
									}
								}
							}
						}
					}
				} else if i == 14 {
					if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
						color.Red("[CHOWN]: Error en mover puntero")
						return
					}
					if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
						color.Red("[CHOWN]: Error en la lectura del archivo")
						return
					}
					for j := 0; j < 16; j++ {
						if apuntador1.B_pointers[j] != -1 {
							if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
								color.Red("[CHOWN]: Error en mover puntero")
								return
							}
							if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
								color.Red("[CHOWN]: Error en la lectura del archivo")
								return
							}
							for k := 0; k < 16; k++ {
								if apuntador2.B_pointers[k] != -1 {
									if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
										color.Red("[CHOWN]: Error en mover puntero")
										return
									}
									if err := binary.Read(file, binary.LittleEndian, &apuntador3); err != nil {
										color.Red("[CHOWN]: Error en la lectura del archivo")
										return
									}
									for l := 0; l < 16; l++ {
										if apuntador3.B_pointers[l] != -1 {
											if _, err := file.Seek(int64(apuntador3.B_pointers[l]), 0); err != nil {
												color.Red("[CHOWN]: Error en mover puntero")
												return
											}
											if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
												color.Red("[CHOWN]: Error en la lectura del archivo")
												return
											}
											for c := 0; c < 4; c++ {
												name := ToString(carpeta.B_content[c].B_name[:])
												if (carpeta.B_content[c].B_inodo != -1) && (name != "." && name != "..") {
													ChownSystem(carpeta.B_content[c].B_inodo, idU, idG, r)
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
