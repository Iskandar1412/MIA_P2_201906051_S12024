package admindisk

import (
	"encoding/binary"
	"math"
	"os"
	"proyecto/comandos/global"
	"proyecto/comandos/utils"
	"proyecto/estructuras/size"
	"proyecto/estructuras/structures"
	"strings"

	"github.com/fatih/color"
)

func Values_MKFS(instructions []string) (string, string, string, bool) {
	var _id string
	var _type = "FULL"
	var _fs = "2fs"
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "id") {
			var value = utils.TieneID("MKFS", valor)
			_id = value
		} else if strings.HasPrefix(strings.ToLower(valor), "type") {
			var value = utils.TieneTypeMKFS(strings.ToLower(valor))
			_type = value
		} else if strings.HasPrefix(strings.ToLower(valor), "fs") {
			var value = utils.TieneFS(strings.ToLower(valor))
			_fs = value
		} else {
			color.Yellow("[MKFS]: Atributo no reconocido")
			_id = ""
			break
		}
	}
	if (_id == "" || len(_id) == 0 || len(_id) > 4) || (_type == "" || _fs == "") {
		return "", "", "", false
	} else {
		return _id, _type, _fs, true
	}
}

func MKFS_EXECUTE(id_disco string, tipo_formateo string, fs string) {
	nodoM, enodoM := utils.Buscar_ID_Montada(id_disco)
	if !enodoM {
		color.Red("[MKFS]: Particion " + id_disco + " no existente o no montada")
		return
	}

	//Buscar disco
	file, err := os.OpenFile(nodoM.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[MKFS]: Error al abrir archivo")
		return
	}
	defer file.Close()

	sb := structures.SuperBloque{}
	inodo := structures.TablaInodo{}
	carpeta := structures.BloqueCarpeta{}
	journal := structures.Journal{}

	if nodoM.Es_Particion_P {
		tamanio := nodoM.Particion_P.Part_s
		n := float32(0)
		if strings.ToLower(fs) == "2fs" {
			n = float32(tamanio-size.SizeSuperBloque()) / (float32(4 + size.SizeTablaInodo() + (3 * size.SizeBloqueArchivo())))
			sb.S_filesistem_type = 2
		} else if strings.ToLower(fs) == "3fs" {
			n = float32(tamanio-size.SizeSuperBloque()) / (float32(4 + size.SizeJournal() + size.SizeTablaInodo() + (3 * size.SizeBloqueArchivo())))
			sb.S_filesistem_type = 3
		}

		numeroEstructuras := int32(math.Floor(float64(n)))
		nBloques := 3 * numeroEstructuras
		inoding := numeroEstructuras * size.SizeTablaInodo()

		sb.S_inodes_count = numeroEstructuras
		sb.S_blocks_count = nBloques
		sb.S_free_blocks_count = nBloques - 2
		sb.S_free_inodes_count = numeroEstructuras - 2
		sb.S_mtime = utils.ObFechaInt()
		sb.S_umtime = 0
		sb.S_mnt_count = 1
		sb.S_magic = 0xEF53
		sb.S_inode_s = size.SizeTablaInodo()
		sb.S_block_s = size.SizeBloqueArchivo()
		sb.S_first_ino = 1
		sb.S_first_blo = 2
		if strings.ToLower(fs) == "2fs" {
			sb.S_bm_inode_start = nodoM.Particion_P.Part_start + size.SizeSuperBloque()
		} else if strings.ToLower(fs) == "3fs" {
			journal := numeroEstructuras * size.SizeJournal()
			sb.S_bm_inode_start = nodoM.Particion_P.Part_start + size.SizeSuperBloque() + journal
		}
		sb.S_bm_block_start = sb.S_bm_inode_start + numeroEstructuras
		sb.S_inode_start = sb.S_bm_block_start + nBloques
		sb.S_block_start = sb.S_bm_inode_start + inoding

		//borrar datos del disco
		inicio := nodoM.Particion_P.Part_start
		estruc_tam := int(nodoM.Particion_P.Part_s - 2)
		estructura := make([]byte, estruc_tam)
		for i := range estructura {
			estructura[i] = '\x00'
		}
		if _, err := file.Seek(int64(inicio), 0); err != nil {
			color.Red("[MKFS]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &estructura); err != nil {
			color.Red("[MKFS]: Error en la escritura de la parcicion")
			return
		}

		//inodo
		inodo.I_uid = 1
		inodo.I_gid = 1
		inodo.I_atime = utils.ObFechaInt()
		inodo.I_ctime = utils.ObFechaInt()
		inodo.I_mtime = utils.ObFechaInt()
		inodo.I_perm = 664
		for i := range inodo.I_block {
			inodo.I_block[i] = -1
		}
		inodo.I_block[0] = sb.S_block_start
		inodo.I_type = 0 //carpeta
		inodo.I_s = 0

		if _, err := file.Seek(int64(sb.S_inode_start), 0); err != nil {
			color.Red("[MKFS]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
			color.Red("[MKFS]: Error en la escritura del MBR")
			return
		}

		// contenido carpeta
		carpeta.B_content[0].B_name = utils.NameCarpeta12(".")
		carpeta.B_content[0].B_inodo = sb.S_inode_start
		carpeta.B_content[1].B_name = utils.NameCarpeta12("..")
		carpeta.B_content[1].B_inodo = sb.S_inode_start
		carpeta.B_content[2].B_name = utils.NameCarpeta12("users.txt")
		carpeta.B_content[2].B_inodo = sb.S_inode_start + size.SizeTablaInodo()
		carpeta.B_content[3].B_name = utils.NameCarpeta12("")
		carpeta.B_content[3].B_inodo = -1

		if _, err := file.Seek(int64(sb.S_block_start), 0); err != nil {
			color.Red("[MKFS]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &carpeta); err != nil {
			color.Red("[MKFS]: Error en la escritura de Carpeta")
			return
		}

		// Crear users.txt
		inodoU := structures.TablaInodo{}
		archivoU := structures.BloqueArchivo{}

		inodoU.I_uid = 1
		inodoU.I_gid = 1
		inodoU.I_atime = utils.ObFechaInt()
		inodoU.I_ctime = utils.ObFechaInt()
		inodoU.I_mtime = utils.ObFechaInt()
		inodoU.I_perm = 700
		for i := range inodoU.I_block {
			inodoU.I_block[i] = -1
		}
		inodoU.I_block[0] = sb.S_block_start + size.SizeBloqueCarpeta()
		s := "1,G,root\n1,U,root,root,123\n"
		inodoU.I_s = int32(len(s))
		inodoU.I_type = 1 //archivo

		if _, err := file.Seek(int64(sb.S_inode_start+size.SizeTablaInodo()), 0); err != nil {
			color.Red("[MKFS]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &inodoU); err != nil {
			color.Red("[MKFS]: Error en la escritura de Carpeta")
			return
		}

		//Archivo usuario
		archivoU.B_content = utils.DevolverContenidoArchivo(s)
		if _, err := file.Seek(int64(sb.S_block_start+size.SizeBloqueCarpeta()), 0); err != nil {
			color.Red("[MKFS]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &archivoU); err != nil {
			color.Red("[MKFS]: Error en la escritura de Archivo")
			return
		}

		// Terminar el proceso
		if strings.ToLower(fs) == "3fs" {
			journal.J_Sig = -1
			journal.J_Tipo = '0'
			journal.J_Size = 0
			journal.J_Fecha = utils.ObFechaInt()
			journal.J_Tipo_Operacion = utils.NameArchivosByte("mkfs")
			journal.J_Start = nodoM.Particion_P.Part_start + size.SizeSuperBloque()
		}

		// Escritura MBR
		nodoM.Particion_P.Part_status = 1
		part := nodoM.Particion_P
		mbr_full, embr := utils.Obtener_FULL_MBR_FDISK(nodoM.Path)
		if !embr {
			return
		}
		pos := int32(0)
		for i, co := range mbr_full.Mbr_partitions {
			if utils.ToString(co.Part_name[:]) == utils.ToString(part.Part_name[:]) {
				pos = int32(i)
				break
			}
		}
		mbr_full.Mbr_partitions[pos] = part
		if _, err := file.Seek(0, 0); err != nil {
			color.Red("[MKFS]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &mbr_full); err != nil {
			color.Red("[MKFS]: Error en la escritura de SuperBloque")
			return
		}

		for z, us := range global.Mounted_Partitions {
			if utils.ToString(us.Particion_P.Part_id[:]) == id_disco {
				global.Mounted_Partitions[z].Particion_P.Part_status = 1
				break
			}
		}

		//escritura SuperBloque
		if _, err := file.Seek(int64(nodoM.Particion_P.Part_start), 0); err != nil {
			color.Red("[MKFS]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &sb); err != nil {
			color.Red("[MKFS]: Error en la escritura de SuperBloque")
			return
		}

		// Bitmap
		var ch0 byte = '0'
		var ch1 byte = '1'
		for i := 0; i < int(numeroEstructuras); i++ {
			if _, err := file.Seek((int64(sb.S_bm_inode_start) + int64(i)), 0); err != nil {
				color.Red("[MKFS]: Error en mover puntero")
				return
			}
			if err := binary.Write(file, binary.LittleEndian, &ch0); err != nil {
				color.Red("[MKFS]: Error en la escritura del Bitmap de Inodos")
				return
			}
		}
		if _, err := file.Seek(int64(sb.S_bm_inode_start), 0); err != nil {
			color.Red("[MKFS]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &ch1); err != nil {
			color.Red("[MKFS]: Error en la escritura de Bitmap de Inodos")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &ch1); err != nil {
			color.Red("[MKFS]: Error en la escritura de Bitmap de Inodo")
			return
		}

		for i := 0; i < int(nBloques); i++ {
			if _, err := file.Seek((int64(sb.S_bm_block_start) + int64(i)), 0); err != nil {
				color.Red("[MKFS]: Error en mover puntero")
				return
			}
			if err := binary.Write(file, binary.LittleEndian, &ch0); err != nil {
				color.Red("[MKFS]: Error en la escritura del Bitmap de Bloques")
				return
			}
		}
		if _, err := file.Seek(int64(sb.S_bm_block_start), 0); err != nil {
			color.Red("[MKFS]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &ch1); err != nil {
			color.Red("[MKFS]: Error en la escritura de Bitmap de Bloques")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &ch1); err != nil {
			color.Red("[MKFS]: Error en la escritura de Bitmap de Bloques")
			return
		}

		// Escribir journal
		if strings.ToLower(fs) == "3fs" {
			if _, err := file.Seek(int64(nodoM.Particion_P.Part_start+size.SizeSuperBloque()), 0); err != nil {
				color.Red("[MKFS]: Error en mover puntero")
				return
			}
			if err := binary.Write(file, binary.LittleEndian, &journal); err != nil {
				color.Red("[MKFS]: Error en la escritura de Journal")
				return
			}
		}
		file.Close()
		color.Blue("[MKFS]: Formateada la parcicion " + utils.ToString(nodoM.Particion_P.Part_name[:]) + " exitosamente - (id): -> " + id_disco)
		return
	} else if nodoM.Es_Particion_L { // Caso de ser logica
		tamanio := nodoM.Particion_L.Part_s
		n := float32(0)
		if strings.ToLower(fs) == "2fs" {
			n = float32(tamanio-size.SizeSuperBloque()) / float32(4+size.SizeTablaInodo()+(3*size.SizeBloqueArchivo()))
			sb.S_filesistem_type = 2
		} else if strings.ToLower(fs) == "3fs" {
			n = float32(tamanio-size.SizeSuperBloque()) / float32(4+size.SizeJournal()+size.SizeTablaInodo()+(3*size.SizeBloqueArchivo()))
			sb.S_filesistem_type = 3
		}

		numeroEstructuras := int32(math.Floor(float64(n)))
		nBloques := 3 * numeroEstructuras
		inoding := numeroEstructuras * size.SizeTablaInodo()

		sb.S_inodes_count = numeroEstructuras
		sb.S_blocks_count = nBloques
		sb.S_free_blocks_count = nBloques - 2
		sb.S_free_inodes_count = numeroEstructuras - 2
		sb.S_mtime = utils.ObFechaInt()
		sb.S_umtime = 0
		sb.S_mnt_count = 1
		sb.S_magic = 0xEF53
		sb.S_inode_s = size.SizeTablaInodo()
		sb.S_block_s = size.SizeBloqueArchivo()
		sb.S_first_ino = 1
		sb.S_first_blo = 1
		if strings.ToLower(fs) == "2fs" {
			sb.S_bm_inode_start = nodoM.Particion_L.Part_start + size.SizeEBR() + size.SizeSuperBloque()
		} else if strings.ToLower(fs) == "3fs" {
			journal := numeroEstructuras * size.SizeJournal()
			sb.S_bm_inode_start = nodoM.Particion_L.Part_start + size.SizeEBR() + size.SizeSuperBloque() + journal
		}
		sb.S_bm_block_start = sb.S_bm_inode_start + numeroEstructuras
		sb.S_inode_start = sb.S_bm_block_start + nBloques
		sb.S_block_start = sb.S_bm_inode_start + inoding

		//borrar datos del disco
		inicio := nodoM.Particion_L.Part_start + size.SizeEBR()
		estruc_tam := int(nodoM.Particion_L.Part_s - 2 - size.SizeEBR())
		estructura := make([]byte, estruc_tam)
		for i := range estructura {
			estructura[i] = '\x00'
		}
		if _, err := file.Seek(int64(inicio), 0); err != nil {
			color.Red("[MKFS]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &estructura); err != nil {
			color.Red("[MKFS]: Error en la escritura de la parcicion")
			return
		}

		//inodo
		inodo.I_uid = 1
		inodo.I_gid = 1
		inodo.I_atime = utils.ObFechaInt()
		inodo.I_ctime = utils.ObFechaInt()
		inodo.I_mtime = utils.ObFechaInt()
		inodo.I_perm = 664
		for i := range inodo.I_block {
			inodo.I_block[i] = -1
		}
		inodo.I_block[0] = sb.S_block_start
		inodo.I_type = 0 //carpeta
		inodo.I_s = 0
		if _, err := file.Seek(int64(sb.S_inode_start), 0); err != nil {
			color.Red("[MKFS]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
			color.Red("[MKFS]: Error en la escritura del MBR")
			return
		}

		// contenido carpeta
		carpeta.B_content[0].B_name = utils.NameCarpeta12(".")
		carpeta.B_content[0].B_inodo = sb.S_inode_start
		carpeta.B_content[1].B_name = utils.NameCarpeta12("..")
		carpeta.B_content[1].B_inodo = sb.S_inode_start
		carpeta.B_content[2].B_name = utils.NameCarpeta12("users.txt")
		carpeta.B_content[2].B_inodo = sb.S_inode_start + size.SizeTablaInodo()
		carpeta.B_content[3].B_name = utils.NameCarpeta12("")
		carpeta.B_content[3].B_inodo = -1

		if _, err := file.Seek(int64(sb.S_block_start), 0); err != nil {
			color.Red("[MKFS]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &carpeta); err != nil {
			color.Red("[MKFS]: Error en la escritura de Carpeta")
			return
		}

		// Crear users.txt
		inodoU := structures.TablaInodo{}
		archivoU := structures.BloqueArchivo{}

		inodoU.I_uid = 1
		inodoU.I_gid = 1
		inodoU.I_atime = utils.ObFechaInt()
		inodoU.I_ctime = utils.ObFechaInt()
		inodoU.I_mtime = utils.ObFechaInt()
		inodoU.I_perm = 700
		for i := range inodoU.I_block {
			inodoU.I_block[i] = -1
		}
		inodoU.I_block[0] = sb.S_block_start + size.SizeBloqueCarpeta()
		s := "1,G,root\n1,U,root,root,123\n"
		inodoU.I_s = int32(len(s))
		inodoU.I_type = 1 //archivo

		if _, err := file.Seek(int64(sb.S_inode_start+size.SizeTablaInodo()), 0); err != nil {
			color.Red("[MKFS]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &inodoU); err != nil {
			color.Red("[MKFS]: Error en la escritura de Carpeta")
			return
		}

		//Archivo usuario
		archivoU.B_content = utils.DevolverContenidoArchivo(s)
		if _, err := file.Seek(int64(sb.S_block_start+size.SizeBloqueCarpeta()), 0); err != nil {
			color.Red("[MKFS]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &archivoU); err != nil {
			color.Red("[MKFS]: Error en la escritura de Archivo")
			return
		}

		// Terminar el proceso
		if strings.ToLower(fs) == "3fs" {
			journal.J_Sig = 1
			journal.J_Tipo = '0'
			journal.J_Size = 0
			journal.J_Fecha = utils.ObFechaInt()
			journal.J_Tipo_Operacion = utils.NameArchivosByte("mkfs")
			journal.J_Start = nodoM.Particion_L.Part_start + size.SizeSuperBloque() + size.SizeEBR()
		}

		// Escritura EBR
		ebr := nodoM.Particion_L
		ebr.Part_mount = 1
		if _, err := file.Seek(int64(nodoM.Particion_L.Part_start), 0); err != nil {
			color.Red("[MKFS]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &ebr); err != nil {
			color.Red("[MKFS]: Error en la escritura de EBR")
			return
		}

		for z, us := range global.Mounted_Partitions {
			if utils.ToString(us.Particion_L.Name[:]) == utils.ToString(nodoM.Particion_L.Name[:]) {
				global.Mounted_Partitions[z].Particion_L.Part_mount = 1
				break
			}
		}

		// Escritura SuperBloque
		if _, err := file.Seek(int64(nodoM.Particion_L.Part_start+size.SizeEBR()), 0); err != nil {
			color.Red("[MKFS]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &sb); err != nil {
			color.Red("[MKFS]: Error en la escritura de SuperBloque")
			return
		}

		// Bitmap
		var ch0 byte = '0'
		var ch1 byte = '1'
		for i := 0; i < int(numeroEstructuras); i++ {
			if _, err := file.Seek((int64(sb.S_bm_inode_start) + int64(i)), 0); err != nil {
				color.Red("[MKFS]: Error en mover puntero")
				return
			}
			if err := binary.Write(file, binary.LittleEndian, &ch0); err != nil {
				color.Red("[MKFS]: Error en la escritura del Bitmap de Inodos")
				return
			}
		}
		if _, err := file.Seek(int64(sb.S_bm_inode_start), 0); err != nil {
			color.Red("[MKFS]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &ch1); err != nil {
			color.Red("[MKFS]: Error en la escritura de Bitmap de Inodos")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &ch1); err != nil {
			color.Red("[MKFS]: Error en la escritura de Bitmap de Inodo")
			return
		}

		for i := 0; i < int(nBloques); i++ {
			if _, err := file.Seek((int64(sb.S_bm_block_start) + int64(i)), 0); err != nil {
				color.Red("[MKFS]: Error en mover puntero")
				return
			}
			if err := binary.Write(file, binary.LittleEndian, &ch0); err != nil {
				color.Red("[MKFS]: Error en la escritura del Bitmap de Bloques")
				return
			}
		}
		if _, err := file.Seek(int64(sb.S_bm_block_start), 0); err != nil {
			color.Red("[MKFS]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &ch1); err != nil {
			color.Red("[MKFS]: Error en la escritura de Bitmap de Bloques")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &ch1); err != nil {
			color.Red("[MKFS]: Error en la escritura de Bitmap de Bloques")
			return
		}

		// Escribir journal
		if strings.ToLower(fs) == "3fs" {
			if _, err := file.Seek(int64(nodoM.Particion_L.Part_start+size.SizeSuperBloque()+size.SizeEBR()), 0); err != nil {
				color.Red("[MKFS]: Error en mover puntero")
				return
			}
			if err := binary.Write(file, binary.LittleEndian, &journal); err != nil {
				color.Red("[MKFS]: Error en la escritura de Journal")
				return
			}
		}
		file.Close()
		color.Blue("[MKFS]: Formateada la parcicion " + utils.ToString(nodoM.Particion_P.Part_name[:]) + " exitosamente - (id): -> " + id_disco)
		return

	} else { // Error
		color.Red("[MKFS]: Error al encontrar el disco")
		return
	}

}
