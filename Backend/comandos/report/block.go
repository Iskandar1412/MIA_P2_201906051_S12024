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

func Report_BLOCK(name string, path string, ruta string, id_disco string) {
	disco_buscado, edb := utils.ObtenerDiscoID(id_disco)
	if !edb {
		return
	}

	file, err := os.OpenFile(disco_buscado.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[Mount]: Error al abrir archivo")
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

	start := sb.S_bm_inode_start
	end := start + sb.S_inodes_count
	var bit byte
	cont := int32(0)
	inodo := structures.TablaInodo{}
	apuntador1, apuntador2, apuntador3 := structures.BloqueApuntador{}, structures.BloqueApuntador{}, structures.BloqueApuntador{}

	///-------------------
	nombre_sin_extension := strings.Split(name, ".")
	rutaB := path + "/" + nombre_sin_extension[0] + ".dot"
	dot, err := os.Create(rutaB)
	if err != nil {
		color.Red("Error al crear el archivo <" + name + ">")
		return
	}
	///-------------------
	//-//-//-//-//-//-//-//
	fmt.Fprintln(dot, ""+`digraph G {`)
	fmt.Fprintln(dot, "\t"+`node[shape=none];`)
	//-//-//-//-//-//-//-//
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
			if _, err := file.Seek(int64(sb.S_inode_start+(cont*size.SizeTablaInodo())), 0); err != nil {
				color.Red("[REP]: Error en mover puntero")
				return
			}
			if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
				color.Red("[REP]: Error en la lectura del byte")
				return
			}
			for i := 0; i < 15; i++ {
				if inodo.I_block[i] != -1 {
					if i < 12 {
						if inodo.I_type == 0 { //carpeta
							fmt.Fprintln(dot, utils.GraphBlockCarpeta(inodo.I_block[i], disco_buscado.Path))
						} else if inodo.I_type == 1 { //archivo
							fmt.Fprintln(dot, utils.GraphBlockArchivo(inodo.I_block[i], disco_buscado.Path))
						}
					} else if i == 12 {
						if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
							color.Red("[REP]: Error en mover puntero")
							return
						}
						if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
							color.Red("[REP]: Error en la lectura del Bloque Apuntador 1")
							return
						}
						fmt.Fprintln(dot, utils.GraphBlockApuntador(inodo.I_block[i], disco_buscado.Path))
						for j := 0; j < 16; j++ {
							if apuntador1.B_pointers[j] != -1 {
								if inodo.I_type == 0 {
									fmt.Fprintln(dot, utils.GraphBlockCarpeta(apuntador1.B_pointers[j], disco_buscado.Path))
								} else if inodo.I_type == 1 {
									fmt.Fprintln(dot, utils.GraphBlockArchivo(apuntador1.B_pointers[j], disco_buscado.Path))
								}
							}
						}
					} else if i == 13 {
						if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
							color.Red("[REP]: Error en mover puntero")
							return
						}
						if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
							color.Red("[REP]: Error en la lectura del Bloque Apuntador 1")
							return
						}
						fmt.Fprintln(dot, utils.GraphBlockApuntador(inodo.I_block[i], disco_buscado.Path))
						for j := 0; j < 16; j++ {
							if apuntador1.B_pointers[j] != -1 {
								if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
									color.Red("[REP]: Error en mover puntero")
									return
								}
								if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
									color.Red("[REP]: Error en la lectura del Bloque Apuntador 2")
									return
								}
								fmt.Fprintln(dot, utils.GraphBlockApuntador(apuntador1.B_pointers[j], disco_buscado.Path))
								for k := 0; k < 16; k++ {
									if apuntador2.B_pointers[k] != -1 {
										if inodo.I_type == 0 {
											fmt.Fprintln(dot, utils.GraphBlockCarpeta(apuntador2.B_pointers[k], disco_buscado.Path))
										} else if inodo.I_type == 1 {
											fmt.Fprintln(dot, utils.GraphBlockArchivo(apuntador2.B_pointers[k], disco_buscado.Path))
										}
									}
								}

							}
						}
					} else if i == 14 {
						if _, err := file.Seek(int64(inodo.I_block[i]), 0); err != nil {
							color.Red("[REP]: Error en mover puntero")
							return
						}
						if err := binary.Read(file, binary.LittleEndian, &apuntador1); err != nil {
							color.Red("[REP]: Error en la lectura del Bloque Apuntador 1")
							return
						}
						fmt.Fprintln(dot, utils.GraphBlockApuntador(inodo.I_block[i], disco_buscado.Path))
						for j := 0; j < 16; j++ {
							if apuntador1.B_pointers[j] != -1 {
								if _, err := file.Seek(int64(apuntador1.B_pointers[j]), 0); err != nil {
									color.Red("[REP]: Error en mover puntero")
									return
								}
								if err := binary.Read(file, binary.LittleEndian, &apuntador2); err != nil {
									color.Red("[REP]: Error en la lectura del Bloque Apuntador 2")
									return
								}
								fmt.Fprintln(dot, utils.GraphBlockApuntador(apuntador1.B_pointers[j], disco_buscado.Path))
								for k := 0; k < 16; k++ {
									if apuntador2.B_pointers[k] != -1 {
										if _, err := file.Seek(int64(apuntador2.B_pointers[k]), 0); err != nil {
											color.Red("[REP]: Error en mover puntero")
											return
										}
										if err := binary.Read(file, binary.LittleEndian, &apuntador3); err != nil {
											color.Red("[REP]: Error en la lectura del Bloque Apuntador 3")
											return
										}
										fmt.Fprintln(dot, utils.GraphBlockApuntador(apuntador2.B_pointers[k], disco_buscado.Path))
										for l := 0; l < 16; l++ {
											if apuntador3.B_pointers[l] != -1 {
												if inodo.I_type == 0 {
													fmt.Fprintln(dot, utils.GraphBlockCarpeta(apuntador3.B_pointers[l], disco_buscado.Path))
												} else if inodo.I_type == 1 {
													fmt.Fprintln(dot, utils.GraphBlockArchivo(apuntador3.B_pointers[l], disco_buscado.Path))
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
		cont++
	}
	//-//-//-//-//-//-//-//
	fmt.Fprintln(dot, ""+`}`)
	//-//-//-//-//-//-//-//
	///-------------------
	dot.Close()
	// Generacion del reporte
	// imagePath := path + "/" + name

	// cmd := exec.Command("/usr/bin/dot", "-T"+nombre_sin_extension[1], rutaB, "-o", imagePath)
	// err = cmd.Run()
	// if err != nil {
	// color.Red("[REP]: Error al generar imagen")
	// return
	// }

	color.Green("[REP]: Block «" + name + "» generated Sucessfull")
}
