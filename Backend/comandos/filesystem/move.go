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

func Values_MOVE(instructions []string) (string, string, bool) {
	var path, destino string

	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "path") {
			var value = utils.TienePathFilePermitions("MOVE", valor)
			path = value
		} else if strings.HasPrefix(strings.ToLower(valor), "destino") {
			var value = utils.TieneDestinoFile("MOVE", valor)
			destino = value
		} else {
			color.Yellow("[MOVE]: Atributo no reconocido")
			return "", "", false
		}
	}
	if path == "" || len(path) == 0 {
		color.Red("[MOVE]: No hay path")
		return "", "", false
	}
	if destino == "" || len(destino) == 0 {
		color.Red("[MOVE]: No hay path")
		return "", "", false
	}
	return path, destino, true
}

func MOVE_EXECUTE(path string, destino string) {
	if !global.UsuarioLogeado.Logged_in {
		color.Red("[MOVE]: Usuario no loggeado")
		return
	}

	nodo := global.UsuarioLogeado.Mounted
	file, err := os.OpenFile(nodo.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[MOVE]: Error al abrir archivo")
		return
	}
	defer file.Close()

	var start int32
	if nodo.Es_Particion_L {
		if nodo.Particion_L.Part_mount != 1 {
			color.Red("[MOVE]: El disco (logico) no se ha formateado -> " + utils.ToString(nodo.Particion_L.Name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
			return
		} else {
			start = nodo.Particion_L.Part_start + size.SizeEBR()
		}
	} else if nodo.Es_Particion_P {
		if nodo.Particion_P.Part_status != 1 {
			color.Red("[MOVE]: El disco (primario) no se ha formateado -> " + utils.ToString(nodo.Particion_P.Part_name[:]) + " - (ID): -> " + utils.ToString(nodo.ID_Particion[:]))
			return
		} else {
			start = nodo.Particion_P.Part_start
		}
	}

	if _, err := file.Seek(int64(start), 0); err != nil {
		color.Red("[MOVE]: Error en mover puntero")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &utils.Sb_System); err != nil {
		color.Red("[MOVE]: Error en la lectura del archivo")
		return
	}

	rutaS := utils.SplitRuta(path)
	rutaD := utils.SplitRuta(destino)

	if len(rutaS) == 0 {
		color.Red("[MOVE]: Ruta inicial no valida")
		return
	}

	if len(rutaD) == 0 {
		color.Red("[MOVE]: Ruta destino no valida")
		return
	}

	posInodoI := utils.Sb_System.S_inode_start
	posInodoO := posInodoI
	posInodoD := posInodoI
	var existP = true
	var inodoO, inodoD structures.TablaInodo
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
			color.Red("[MOVE]: No se encontro el directorio -> " + path)
			return
		}
	}

	existP = true
	// posInodoI = utils.Sb_System.S_inode_start
	for i := 0; i < len(rutaD); i++ {
		if existP {
			aux := posInodoD
			posInodoD = utils.ExistPathSystem(rutaD, int32(i), posInodoD, nodo.Path)
			if posInodoD == aux {
				existP = false
			}
		}
		if !existP {
			color.Red("[MOVE]: No se encontro el directorio -> " + destino)
			return
		}
	}

	if _, err := file.Seek(int64(posInodoD), 0); err != nil {
		color.Red("[MOVE]: Error en mover puntero")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &inodoD); err != nil {
		color.Red("[MOVE]: Error en la lectura del archivo")
		return
	}

	if inodoD.I_type == 1 {
		color.Red("[MOVE]: No es un inodo de carpeta")
		return
	}

	if posInodoD != utils.ExistPathSystem(rutaS, int32(len(rutaS)-1), posInodoD, nodo.Path) {
		color.Red("[MOVE]: El directorio >" + rutaS[len(rutaS)-1] + " ya esta en uso >>" + rutaD[len(rutaD)-1])
		return
	}

	if _, err := file.Seek(int64(posInodoO), 0); err != nil {
		color.Red("[MOVE]: Error en mover puntero")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &inodoO); err != nil {
		color.Red("[MOVE]: Error en la lectura del archivo")
		return
	}

	if !utils.ValidarPermisoWSystem(posInodoO, nodo.Path) {
		color.Red("[MOVE]: No se puede sobreescribir el archivo >" + path + " por falta de permisos")
		return
	}

	if !utils.ValidarPermisoWSystem(posInodoD, nodo.Path) {
		color.Red("[MOVE]: No se puede sobreescribi el archivo >" + destino + " por falta de permisos")
		return
	}

	posCambio := utils.GetPosInodo(inodoO, rutaS[len(rutaS)-1])
	inodoO.I_mtime = utils.ObFechaInt()
	if _, err := file.Seek(int64(posInodoO), 0); err != nil {
		color.Red("[MOVE]: Error en mover puntero")
		return
	}
	if err := binary.Write(file, binary.LittleEndian, &inodoO); err != nil {
		color.Red("[MOVE]: Error en la escritura del archivo")
		return
	}

	atc := utils.AgregarCarpetaSystem(posCambio, posInodoD, rutaS[len(rutaS)-1])

	if atc == -1 {
		return
	}

	inodoD.I_mtime = utils.ObFechaInt()
	if _, err := file.Seek(int64(posInodoD), 0); err != nil {
		color.Red("[MOVE]: Error en mover puntero")
		return
	}
	if err := binary.Write(file, binary.LittleEndian, &inodoD); err != nil {
		color.Red("[MOVE]: Error en la escritura del archivo")
		return
	}

	var inodoCambio structures.TablaInodo
	if _, err := file.Seek(int64(posCambio), 0); err != nil {
		color.Red("[MOVE]: Error en mover puntero")
		return
	}
	if err := binary.Read(file, binary.LittleEndian, &inodoCambio); err != nil {
		color.Red("[MOVE]: Error en la lectura del archivo")
		return
	}

	if inodoCambio.I_type == 0 {
		var carpeta structures.BloqueCarpeta
		if _, err := file.Seek(int64(inodoCambio.I_block[0]), 0); err != nil {
			color.Red("[MOVE]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &carpeta); err != nil {
			color.Red("[MOVE]: Error en la lectura del archivo")
			return
		}
		carpeta.B_content[1].B_inodo = posInodoD
		if _, err := file.Seek(int64(inodoCambio.I_block[0]), 0); err != nil {
			color.Red("[MOVE]: Error en mover puntero")
			return
		}
		if err := binary.Write(file, binary.LittleEndian, &carpeta); err != nil {
			color.Red("[MOVE]: Error en la escritura del archivo")
			return
		}
	}

	if utils.Sb_System.S_filesistem_type == 3 {
		utils.EscribirJournalSystem("move", '1', path, destino, nodo)
	}
	color.Green("[MOVE]: Se movio el archivo")
}
