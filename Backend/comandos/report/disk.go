package report

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"proyecto/comandos/utils"
	"proyecto/estructuras/size"
	"proyecto/estructuras/structures"
	"strings"

	"github.com/fatih/color"
)

func Report_DISK(name string, path string, ruta string, id_disco string) {
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

	mbr, _ := utils.Obtener_FULL_MBR_FDISK(disco_buscado.Path)

	fmt.Fprintln(dot, "digraph G{")
	fmt.Fprintln(dot, "\tnode[shape=none];")
	fmt.Fprintln(dot, "\tstart[label=<<table border=\"1\" cellspacing=\"0\" cellpadding=\"5\" color=\"#000\">")
	fmt.Fprintln(dot, "\t\t"+`<tr>`)
	fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" bgcolor="#000"><font color="white"><b>MBR</b></font></td>`)

	// i := int32(0)
	// inicio := size.SizeMBR()
	//Empezaremos con las primeras particiones
	i := 0
	for i < 4 {
		disco := mbr.Mbr_partitions[i]
		if disco.Part_type == 'P' {
			if i == 3 {
				if disco.Part_s != -1 && disco.Part_start != -1 {
					if (disco.Part_s + disco.Part_start) < mbr.Mbr_tamano {
						espacio := (float32(disco.Part_s) / float32(mbr.Mbr_tamano)) * 100
						espacio_libre := ((float32(mbr.Mbr_tamano) - (float32(disco.Part_start + disco.Part_s))) / float32(mbr.Mbr_tamano)) * 100
						fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" width="300" bgcolor="#74719B"><b>Primaria `+utils.Returnstring(utils.ToString(disco.Part_name[:]))+`</b><br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
						fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" bgcolor="#CCCCCC"><b>Libre</b><br/>`+fmt.Sprint(utils.Redondeo(espacio_libre))+"%"+` del disco</td>`)
						i++
						break
						// goto tfinalPrimeraParte
					} else if (disco.Part_s + disco.Part_start) == mbr.Mbr_tamano {
						espacio := (float32(disco.Part_s) / float32(mbr.Mbr_tamano)) * 100
						fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" width="300" bgcolor="#74719B"><b>Primaria `+utils.Returnstring(utils.ToString(disco.Part_name[:]))+`</b><br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
						i++
						break
						// goto tfinalPrimeraParte
					}
				} else {
					espacio := (float32(mbr.Mbr_tamano-(mbr.Mbr_partitions[i-1].Part_start+mbr.Mbr_partitions[i-1].Part_s)) / float32(mbr.Mbr_tamano)) * 100
					fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" bgcolor="#CCCCCC"><b>Libre</b><br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
					i++
					break
					// goto tfinalPrimeraParte
				}
			} else {
				if disco.Part_s != -1 && disco.Part_start != -1 {
					// espacio := (float32(disco.Part_start+disco.Part_s) / float32(mbr.Mbr_tamano)) * 100
					// fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" width="300"><b>Primaria `+utils.Returnstring(utils.ToString(disco.Part_name[:]))+`</b><br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
					if (disco.Part_start + disco.Part_s) < (mbr.Mbr_partitions[i+1].Part_start) {
						espacio := (float32(disco.Part_s) / float32(mbr.Mbr_tamano)) * 100
						espacio_libre := ((float32(mbr.Mbr_partitions[i+1].Part_start) - float32(disco.Part_s+disco.Part_start)) / float32(mbr.Mbr_tamano)) * 100
						fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" width="300" bgcolor="#74719B"><b>Primaria `+utils.Returnstring(utils.ToString(disco.Part_name[:]))+`</b><br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
						fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" bgcolor="#CCCCCC"><b>Libre</b><br/>`+fmt.Sprint(utils.Redondeo(espacio_libre))+"%"+` del disco</td>`)
						i++
						continue
					} else if (disco.Part_start + disco.Part_s) == (mbr.Mbr_partitions[i+1].Part_start) {
						espacio := (float32(disco.Part_s) / float32(mbr.Mbr_tamano)) * 100
						fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" width="300" bgcolor="#74719B"><b>Primaria `+utils.Returnstring(utils.ToString(disco.Part_name[:]))+`</b><br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
						i++
						continue
					} else if mbr.Mbr_partitions[i+1].Part_start == -1 {
						espacio := (float32(disco.Part_s) / float32(mbr.Mbr_tamano)) * 100
						fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" width="300" bgcolor="#74719B"><b>Primaria `+utils.Returnstring(utils.ToString(disco.Part_name[:]))+`</b><br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
						i++
						continue
					}
				} else {
					if i == 0 {
						if mbr.Mbr_partitions[0].Part_s == -1 && mbr.Mbr_partitions[0].Part_start == -1 {
							if mbr.Mbr_partitions[1].Part_s == -1 && mbr.Mbr_partitions[1].Part_start == -1 {
								if mbr.Mbr_partitions[2].Part_s == -1 && mbr.Mbr_partitions[2].Part_start == -1 {
									if mbr.Mbr_partitions[3].Part_s == -1 && mbr.Mbr_partitions[3].Part_start == -1 {
										espacio := (float32(mbr.Mbr_tamano-size.SizeSuperBloque()) / float32(mbr.Mbr_tamano)) * 100
										fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" bgcolor="#CCCCCC"><b>Libre</b><br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
										// goto tfinalPrimeraParte
										break
									} else {
										espacio := ((float32(mbr.Mbr_partitions[i+3].Part_start - size.SizeMBR())) / float32(mbr.Mbr_tamano)) * 100
										fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" bgcolor="#CCCCCC"><b>Libre</b><br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
										i = 3
										continue
									}
								} else {
									espacio := ((float32(mbr.Mbr_partitions[i+2].Part_start - size.SizeMBR())) / float32(mbr.Mbr_tamano)) * 100
									fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" bgcolor="#CCCCCC"><b>Libre</b><br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
									i = 2
									continue
								}
							} else {
								espacio := ((float32(mbr.Mbr_partitions[i+1].Part_start - size.SizeMBR())) / float32(mbr.Mbr_tamano)) * 100
								fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" bgcolor="#CCCCCC"><b>Libre</b><br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
								i = 1
								continue
							}
						}
					} else if i == 1 {
						if mbr.Mbr_partitions[1].Part_s == -1 && mbr.Mbr_partitions[1].Part_start == -1 {
							if mbr.Mbr_partitions[0].Part_s == -1 && mbr.Mbr_partitions[0].Part_start == -1 {
								if mbr.Mbr_partitions[2].Part_s == -1 && mbr.Mbr_partitions[2].Part_start == -1 {
									if mbr.Mbr_partitions[3].Part_s == -1 && mbr.Mbr_partitions[3].Part_start == -1 {
										espacio := ((float32(mbr.Mbr_tamano - size.SizeMBR())) / float32(mbr.Mbr_tamano)) * 100
										fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" bgcolor="#CCCCCC"><b>Libre</b><br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
										break
									} else {
										espacio := ((float32(mbr.Mbr_partitions[3].Part_start - size.SizeMBR())) / float32(mbr.Mbr_tamano)) * 100
										fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" bgcolor="#CCCCCC"><b>Libre</b><br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
										//saltar ya que i es 2 más
										i = 3
										continue
									}
								} else {
									espacio := ((float32(mbr.Mbr_partitions[2].Part_start - size.SizeMBR())) / float32(mbr.Mbr_tamano)) * 100
									fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" bgcolor="#CCCCCC"><b>Libre</b><br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
									i = 2
									continue
								}
							} else {
								if mbr.Mbr_partitions[2].Part_s == -1 && mbr.Mbr_partitions[2].Part_start == -1 {
									if mbr.Mbr_partitions[3].Part_s == -1 && mbr.Mbr_partitions[3].Part_start == -1 {
										espacio := ((float32(mbr.Mbr_tamano - size.SizeMBR())) / float32(mbr.Mbr_tamano)) * 100
										fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" bgcolor="#CCCCCC"><b>Libre</b><br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
										break
									} else {
										espacio := ((float32(mbr.Mbr_partitions[3].Part_start - (mbr.Mbr_partitions[0].Part_start + mbr.Mbr_partitions[0].Part_s))) / float32(mbr.Mbr_tamano)) * 100
										fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" bgcolor="#CCCCCC"><b>Libre</b><br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
										i = 3
										continue
									}
								} else {
									espacio := ((float32(mbr.Mbr_partitions[2].Part_start - (mbr.Mbr_partitions[0].Part_start + mbr.Mbr_partitions[0].Part_s))) / float32(mbr.Mbr_tamano)) * 100
									fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" bgcolor="#CCCCCC"><b>Libre</b><br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
									i = 2
									continue
								}
							}
						}
					} else if i == 2 {
						if mbr.Mbr_partitions[2].Part_s == -1 && mbr.Mbr_partitions[2].Part_start == -1 {
							if mbr.Mbr_partitions[3].Part_s == -1 && mbr.Mbr_partitions[3].Part_start == -1 {
								espacio := ((float32(mbr.Mbr_tamano - (mbr.Mbr_partitions[1].Part_start - mbr.Mbr_partitions[1].Part_s))) / float32(mbr.Mbr_tamano)) * 100
								fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" bgcolor="#CCCCCC"><b>Libre</b><br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
								break
							} else {
								espacio := ((float32(mbr.Mbr_partitions[3].Part_start - (mbr.Mbr_partitions[1].Part_start + mbr.Mbr_partitions[1].Part_s))) / float32(mbr.Mbr_tamano)) * 100
								fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="3" bgcolor="#CCCCCC"><b>Libre</b><br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
								i = 3
								continue
							}
						}
					}
				}
			}
		} else { //333333 334FFF
			fmt.Fprintln(dot, "\t\t\t"+`<td rowspan="1" bgcolor="#333333"><font color="white"><b>Extendida  </b> `+utils.Returnstring(utils.ToString(disco.Part_name[:]))+` </font></td>`)
			i++
			continue
		}
	}

	fmt.Fprintln(dot, "\t\t"+`</tr>`)

	//Segunda parte verificar ebr
	for _, disco := range mbr.Mbr_partitions {
		if disco.Part_type == 'E' {
			ebr := structures.EBR{}
			if _, err := file.Seek(int64(disco.Part_start), 0); err != nil {
				color.Red("[REP]: Error en mover puntero")
				return
			}
			if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
				color.Red("[REP]: Error en la lectura del EBR")
				return
			}

			fmt.Fprintln(dot, "\t\t"+`<tr>`)
			fmt.Fprintln(dot, "\t\t\t"+`<td colspan="1">`)
			fmt.Fprintln(dot, "\t\t\t\t"+`<table border="1" cellspacing="1" cellpadding="5" color="gray">`)
			fmt.Fprintln(dot, "\t\t\t\t\t"+`<tr>`)
			if !(ebr.Part_start == -1 && ebr.Part_next == -1) {
				for {
					if ebr.Part_next == -1 {
						if (ebr.Part_start + ebr.Part_s) < (disco.Part_s + disco.Part_start) {
							espacio := ((float32(ebr.Part_s)) / float32(mbr.Mbr_tamano)) * 100
							espacio_libre := ((float32(disco.Part_start+disco.Part_s) - float32(ebr.Part_start+ebr.Part_s)) / float32(mbr.Mbr_tamano)) * 100
							fmt.Fprintln(dot, "\t\t\t\t\t\t"+`<td bgcolor="#3366FF"><b>Logica </b> `+utils.Returnstring(utils.ToString(ebr.Name[:]))+`<br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
							fmt.Fprintln(dot, "\t\t\t\t\t\t"+`<td bgcolor="#CCCCCC"><b>Libre</b><br/>`+fmt.Sprint(utils.Redondeo(espacio_libre))+"%"+` del disco</td>`)

						} else if (ebr.Part_start + ebr.Part_s) == disco.Part_s {
							espacio := (float32(ebr.Part_s) / float32(mbr.Mbr_tamano)) * 100
							fmt.Fprintln(dot, "\t\t\t\t\t\t"+`<td bgcolor="#FF5733"><b>EBR</b></td>`)
							fmt.Fprintln(dot, "\t\t\t\t\t\t"+`<td bgcolor="#3366FF"><b>Logica </b> `+utils.Returnstring(utils.ToString(ebr.Name[:]))+`<br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
						}
						break
					}
					espacio := (float32(ebr.Part_s) / float32(mbr.Mbr_tamano)) * 100
					fmt.Fprintln(dot, "\t\t\t\t\t\t"+`<td bgcolor="#FF5733"><b>EBR</b></td>`)
					fmt.Fprintln(dot, "\t\t\t\t\t\t"+`<td bgcolor="#3366FF"><b>Logica </b> `+utils.Returnstring(utils.ToString(ebr.Name[:]))+`<br/>`+fmt.Sprint(utils.Redondeo(espacio))+"%"+` del disco</td>`)
					if _, err := file.Seek(int64(ebr.Part_next), 0); err != nil {
						color.Red("[REP]: Error en mover puntero")
						return
					}
					if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
						color.Red("[REP]: Error en la lectura del EBR")
						return
					}
				}
			}
			fmt.Fprintln(dot, "\t\t\t\t\t"+`</tr>`)
			fmt.Fprintln(dot, "\t\t\t\t"+`</table>`)
			fmt.Fprintln(dot, "\t\t\t"+`</td>`)
			fmt.Fprintln(dot, "\t\t"+`</tr>`)
		} else {
			continue
		}
	}

	//final

	fmt.Fprintln(dot, "\t</table>>];")
	fmt.Fprintln(dot, "}")
	dot.Close()

	// Generacion del reporte
	imagePath := path + "/" + name

	cmd := exec.Command("/usr/bin/dot", "-T"+nombre_sin_extension[1], rutaB, "-o", imagePath)
	err = cmd.Run()
	if err != nil {
		color.Red("[REP]: Error al generar imagen")
		return
	}

	color.Green("[REP]: Disk «" + name + "» generated Sucessfull")
}
