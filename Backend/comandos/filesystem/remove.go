package filesystem

import (
	"encoding/binary"
	"os"
	"proyecto/comandos/global"
	"proyecto/comandos/utils"
	"proyecto/estructuras/size"
	"proyecto/estructuras/structures"
	"strings"

	"github.com/fatih/color"
)

func Values_REMOVE(instructions []string) (string, bool) {
	var path string

	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "path") {
			var value = utils.TienePathFilePermitions("REMOVE", valor)
			path = value
		} else {
			color.Yellow("[REMOVE]: Atributo no reconocido")
			return "", false
		}
	}
	if path == "" || len(path) == 0 {
		color.Red("[REMOVE]: No hay path")
		return "", false
	}
	return path, true
}

func REMOVE_EXECUTE(path string) {
	if !global.UsuarioLogeado.Logged_in {
		color.Red("[REMOVE]: No hay usuario loggeado")
		return
	}

	nodo := global.UsuarioLogeado.Mounted
	file, err := os.OpenFile(nodo.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[REMOVE]: Error al abrir archivo")
		return
	}
	defer file.Close()

	var start int32
	if nodo.Es_Particion_L {
		if nodo.Particion_L.Part_mount != 1 {
			color.Red("[REMOVE]: El disco (logico) no se ha formateado -> " + utils.ToString(nodo.Particion_L.Name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
			return
		} else {
			start = nodo.Particion_L.Part_start + size.SizeEBR()
		}
	} else if nodo.Es_Particion_P {
		if nodo.Particion_P.Part_status != 1 {
			color.Red("[REMOVE]: El disco (primario) no se ha formateado -> " + utils.ToString(nodo.Particion_P.Part_name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
			return
		} else {
			start = nodo.Particion_P.Part_start
		}
	}

	if _, err := file.Seek(int64(start), 0); err != nil {
		color.Red("[REMOVE]: Error en mover puntero")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &utils.Sb_System); err != nil {
		color.Red("[REMOVE]: Error en la lectura del SuperBloque")
		return
	}

	rutaS := utils.SplitRuta(path)
	if len(rutaS) == 0 {
		color.Red("[REMOVE]: Ruta invalida")
		return
	}

	posInodoI := utils.Sb_System.S_inode_start
	posInodoO := utils.Sb_System.S_inode_start
	existP := true
	// var inodoO, inodoD structures.TablaInodo
	for i := 0; i < len(rutaS); i++ {
		if existP {
			aux := posInodoI
			posInodoI = utils.ExistPathSystem(rutaS, int32(i), posInodoI, nodo.Path)
			if i != (len(rutaS) - 1) {
				posInodoO = posInodoI
			}
			if posInodoI == aux {
				existP = false
			}
		}
		if !existP {
			color.Red("[REMOVE]: No se encontro el directorio -> " + path)
			return
		}
	}

	bitM1 := utils.BuscarPosicionInodoBM(posInodoI)
	removeI := utils.RemoveInodo(posInodoI, bitM1)

	start = 0
	if nodo.Es_Particion_L {
		if nodo.Particion_L.Part_mount != 1 {
			color.Red("[REMOVE]: El disco (logico) no se ha formateado -> " + utils.ToString(nodo.Particion_L.Name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
			return
		} else {
			start = nodo.Particion_L.Part_start + size.SizeEBR()
		}
	} else if nodo.Es_Particion_P {
		if nodo.Particion_P.Part_status != 1 {
			color.Red("[REMOVE]: El disco (primario) no se ha formateado -> " + utils.ToString(nodo.Particion_P.Part_name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
			return
		} else {
			start = nodo.Particion_P.Part_start
		}
	}

	if _, err := file.Seek(int64(start), 0); err != nil {
		color.Red("[REMOVE]: Error en mover puntero")
		return
	}

	if err := binary.Write(file, binary.LittleEndian, &utils.Sb_System); err != nil {
		color.Red("[/]: Error en la escritura de Bitmap de Bloques")
		return
	}

	if removeI {
		var carpeta structures.BloqueCarpeta
		posC := utils.ReturnCarpetaSystem(rutaS[len(rutaS)-1], posInodoO, nodo.Path)
		if _, err := file.Seek(int64(posC), 0); err != nil {
			color.Red("[REMOVE]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
			color.Red("[REMOVE]: Error en la lectura del archivo")
			return
		}

		for i := 0; i < 4; i++ {
			if utils.ToString(carpeta.B_content[i].B_name[:]) == rutaS[len(rutaS)-1] {
				carpeta.B_content[i].B_name = utils.NameCarpeta12("")
				carpeta.B_content[i].B_inodo = -1
			}
		}
		if _, err := file.Seek(int64(posC), 0); err != nil {
			color.Red("[REMOVE]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &carpeta); err != nil {
			color.Red("[/]: Error en la escritura del archivo")
			return
		}
	}

	utils.BuscarPrimerInodoVacio()
	utils.BuscarPrimerBloqueVacio()

	start = 0
	if nodo.Es_Particion_L {
		if nodo.Particion_L.Part_mount != 1 {
			color.Red("[REMOVE]: El disco (logico) no se ha formateado -> " + utils.ToString(nodo.Particion_L.Name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
			return
		} else {
			start = nodo.Particion_L.Part_start + size.SizeEBR()
		}
	} else if nodo.Es_Particion_P {
		if nodo.Particion_P.Part_status != 1 {
			color.Red("[REMOVE]: El disco (primario) no se ha formateado -> " + utils.ToString(nodo.Particion_P.Part_name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
			return
		} else {
			start = nodo.Particion_P.Part_start
		}
	}
	if _, err := file.Seek(int64(start), 0); err != nil {
		color.Red("[REMOVE]: Error en mover puntero")
		return
	}
	if err := binary.Write(file, binary.LittleEndian, &utils.Sb_System); err != nil {
		color.Red("[/]: Error en la escritura de Bitmap de Bloques")
		return
	}

	if utils.Sb_System.S_filesistem_type == 3 {
		utils.EscribirJournalSystem("remove", '1', path, "", nodo)
	}

	color.Green("[REMOVE]: Archivo o carpeta eliminado -> " + path)
}
