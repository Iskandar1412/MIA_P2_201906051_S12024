package utils

import (
	"encoding/binary"
	"os"
	"proyecto/comandos/global"
	"proyecto/estructuras/structures"

	"github.com/fatih/color"
)

func RemoveInodo(posInodo int32, posBM int32) bool {
	nodo := global.UsuarioLogeado.Mounted
	file, err := os.OpenFile(nodo.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return false
	}
	defer file.Close()

	var banderaR = true
	var permiso = ValidarPermisoWSystem(posInodo, nodo.Path)
	var inodo structures.TablaInodo
	// var archivo structures.BloqueArchivo
	// var apuntador1, apuntador2, apuntador3 structures.BloqueApuntador
	if _, err := file.Seek(int64(posInodo), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return false
	}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[/]: Error en la lectura del archivo")
		return false
	}

	if !permiso && (inodo.I_type == '1') {
		color.Red("[/]: No se tienen o los permisos necesarios o el inodo no es carpeta")
		return false
	}

	for i := 0; i < 15; i++ {
		if inodo.I_block[i] != -1 {
			if i < 12 {
				var banderaAux = true
				var bM int32 = 0
				if inodo.I_type == 0 {
					bM = BuscarPosicionBloqueBM(inodo.I_block[i])
					banderaAux = RemoveBcarpeta(inodo.I_block[i], bM)
				} else if inodo.I_type == 1 {
					bM = BuscarPosicionBloqueBM(inodo.I_block[i])
					banderaAux = RemoveBarchivo(inodo.I_block[i], bM)
				}
				if banderaAux {
					inodo.I_block[i] = -1
				} else {
					banderaR = false
				}
			} else if i == 12 {
				var bM int32 = 0
				bM = BuscarPosicionBloqueBM(inodo.I_block[i])
				banderaAux := RemoveBapuntador1(inodo.I_block[i], bM, inodo.I_type)
				if banderaAux {
					inodo.I_block[i] = -1
				} else {
					banderaR = false
				}
			} else if i == 13 {
				var bM int32 = 0
				bM = BuscarPosicionBloqueBM(inodo.I_block[i])
				banderaAux := RemoveBapuntador2(inodo.I_block[i], bM, inodo.I_type)
				if banderaAux {
					inodo.I_block[i] = -1
				} else {
					banderaR = false
				}
			} else if i == 14 {
				var bM int32 = 0
				bM = BuscarPosicionBloqueBM(inodo.I_block[i])
				banderaAux := RemoveBapuntador3(inodo.I_block[i], bM, inodo.I_type)
				if banderaAux {
					inodo.I_block[i] = -1
				} else {
					banderaR = false
				}
			}
		}
	}

	if !permiso || !banderaR {
		return false
	}

	var cero byte = '0'
	if _, err := file.Seek(int64(Sb_System.S_bm_block_start+posBM), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return false
	}
	if err := binary.Write(file, binary.LittleEndian, &cero); err != nil {
		color.Red("[/]: Error en la escritura del archivo")
		return false
	}

	if _, err := file.Seek(int64(posInodo), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return false
	}
	if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[/]: Error en la escritura del archivo")
		return false
	}
	Sb_System.S_free_inodes_count += 1
	return true
}

func RemoveBarchivo(posBlock int32, posBM int32) bool {
	var cero byte = '0'
	nodo := global.UsuarioLogeado.Mounted
	file, err := os.OpenFile(nodo.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return false
	}
	defer file.Close()
	if _, err := file.Seek(int64(Sb_System.S_bm_block_start+posBM), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return false
	}
	if err := binary.Write(file, binary.LittleEndian, &cero); err != nil {
		color.Red("[/]: Error en la escritura del archivo")
		return false
	}
	Sb_System.S_free_blocks_count += 1
	return true
}

func RemoveBcarpeta(posBlock int32, posBM int32) bool {
	var banderaR = true
	var carpeta structures.BloqueCarpeta
	nodo := global.UsuarioLogeado.Mounted
	file, err := os.OpenFile(nodo.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return false
	}
	defer file.Close()
	if _, err := file.Seek(int64(posBlock), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return false
	}
	if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
		color.Red("[/]: Error en la lectura del archivo")
		return false
	}

	for i := 0; i < 4; i++ {
		var nameA string = ""
		nameA += ToString(carpeta.B_content[i].B_name[:])
		if (carpeta.B_content[i].B_inodo != -1) && (nameA != "." && nameA != "..") {
			var banderaAux = true
			var iM int32 = 0
			iM = BuscarPosicionInodoBM(carpeta.B_content[i].B_inodo)
			banderaAux = RemoveInodo(carpeta.B_content[i].B_inodo, iM)
			if banderaAux {
				carpeta.B_content[i].B_inodo = -1
				carpeta.B_content[i].B_name = NameCarpeta12("")
			} else {
				banderaR = false
			}
		}
	}

	if _, err := file.Seek(int64(posBlock), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return false
	}
	if err := binary.Write(file, binary.LittleEndian, &carpeta); err != nil {
		color.Red("[/]: Error en la escritura de Bitmap de Bloques")
		return false
	}
	if !banderaR {
		return false
	}

	var cero byte = '0'
	if _, err := file.Seek(int64(Sb_System.S_bm_block_start+posBM), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return false
	}
	if err := binary.Write(file, binary.LittleEndian, &cero); err != nil {
		color.Red("[/]: Error en la escritura de Bitmap de Bloques")
		return false
	}
	Sb_System.S_free_blocks_count += 1
	return true
}

func RemoveBapuntador1(posBlock int32, posBM int32, _type int32) bool {
	var banderaR = true
	var apuntador structures.BloqueApuntador
	nodo := global.UsuarioLogeado.Mounted
	file, err := os.OpenFile(nodo.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return false
	}
	defer file.Close()
	if _, err := file.Seek(int64(posBlock), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return false
	}
	if err := binary.Read(file, binary.LittleEndian, &apuntador); err != nil {
		color.Red("[/]: Error en la lectura del archivo")
		return false
	}

	for i := 0; i < 16; i++ {
		if apuntador.B_pointers[i] != -1 {
			var banderaAux = true
			var bM int32 = 0
			if _type == 0 {
				bM = BuscarPosicionBloqueBM(apuntador.B_pointers[i])
				banderaAux = RemoveBcarpeta(apuntador.B_pointers[i], bM)
			} else if _type == 1 {
				bM = BuscarPosicionBloqueBM(apuntador.B_pointers[i])
				banderaAux = RemoveBarchivo(apuntador.B_pointers[i], bM)
			}
			if banderaAux {
				apuntador.B_pointers[i] = -1
			} else {
				banderaR = false
			}
		}
	}

	if _, err := file.Seek(int64(posBlock), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return false
	}
	if err := binary.Write(file, binary.LittleEndian, &apuntador); err != nil {
		color.Red("[/]: Error en la escritura del archivo")
		return false
	}
	if !banderaR {
		return false
	}

	var cero byte = '0'
	if _, err := file.Seek(int64(Sb_System.S_bm_block_start+posBM), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return false
	}
	if err := binary.Write(file, binary.LittleEndian, &cero); err != nil {
		color.Red("[/]: Error en la escritura de Bitmap de Bloques")
		return false
	}
	Sb_System.S_free_blocks_count += 1

	return true
}

func RemoveBapuntador2(posBlock int32, posBM int32, _type int32) bool {
	var banderaR = true
	var apuntador structures.BloqueApuntador
	nodo := global.UsuarioLogeado.Mounted
	file, err := os.OpenFile(nodo.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return false
	}
	defer file.Close()
	if _, err := file.Seek(int64(posBlock), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return false
	}
	if err := binary.Read(file, binary.LittleEndian, &apuntador); err != nil {
		color.Red("[/]: Error en la lectura del archivo")
		return false
	}

	for i := 0; i < 16; i++ {
		if apuntador.B_pointers[i] != -1 {
			var banderaAux = true
			var bM int32 = 0
			bM = BuscarPosicionBloqueBM(apuntador.B_pointers[i])
			banderaAux = RemoveBapuntador1(apuntador.B_pointers[i], bM, _type)
			if banderaAux {
				apuntador.B_pointers[i] = -1
			} else {
				banderaR = false
			}
		}
	}

	if _, err := file.Seek(int64(posBlock), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return false
	}
	if err := binary.Write(file, binary.LittleEndian, &apuntador); err != nil {
		color.Red("[/]: Error en la escritura del archivo")
		return false
	}
	if !banderaR {
		return false
	}

	var cero byte = '0'
	if _, err := file.Seek(int64(Sb_System.S_bm_block_start+posBM), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return false
	}
	if err := binary.Write(file, binary.LittleEndian, &cero); err != nil {
		color.Red("[/]: Error en la escritura de Bitmap de Bloques")
		return false
	}
	Sb_System.S_free_blocks_count += 1

	return true
}

func RemoveBapuntador3(posBlock int32, posBM int32, _type int32) bool {
	var banderaR = true
	var apuntador structures.BloqueApuntador
	nodo := global.UsuarioLogeado.Mounted
	file, err := os.OpenFile(nodo.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[/]: Error al abrir archivo")
		return false
	}
	defer file.Close()
	if _, err := file.Seek(int64(posBlock), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return false
	}
	if err := binary.Read(file, binary.LittleEndian, &apuntador); err != nil {
		color.Red("[/]: Error en la lectura del archivo")
		return false
	}

	for i := 0; i < 16; i++ {
		if apuntador.B_pointers[i] != -1 {
			var banderaAux = true
			var bM int32 = 0
			bM = BuscarPosicionBloqueBM(apuntador.B_pointers[i])
			banderaAux = RemoveBapuntador2(apuntador.B_pointers[i], bM, _type)
			if banderaAux {
				apuntador.B_pointers[i] = -1
			} else {
				banderaR = false
			}
		}
	}

	if _, err := file.Seek(int64(posBlock), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return false
	}
	if err := binary.Write(file, binary.LittleEndian, &apuntador); err != nil {
		color.Red("[/]: Error en la escritura del archivo")
		return false
	}
	if !banderaR {
		return false
	}

	var cero byte = '0'
	if _, err := file.Seek(int64(Sb_System.S_bm_block_start+posBM), 0); err != nil {
		color.Red("[/]: Error en mover puntero")
		return false
	}
	if err := binary.Write(file, binary.LittleEndian, &cero); err != nil {
		color.Red("[/]: Error en la escritura de Bitmap de Bloques")
		return false
	}
	Sb_System.S_free_blocks_count += 1

	return true
}
