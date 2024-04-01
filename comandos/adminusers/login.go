package adminusers

import (
	"encoding/binary"
	"os"
	"proyecto/comandos/global"
	"proyecto/comandos/utils"
	"proyecto/estructuras/size"
	"proyecto/estructuras/structures"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func Values_LOGIN(instructions []string) (global.Usuario, bool) {
	var usuario string
	var password string
	var id string

	user_temp := global.Usuario{}
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "id") {
			var value = utils.TieneID("LOGIN", valor)
			id = value
			continue
		} else if strings.HasPrefix(strings.ToLower(valor), "pass") {
			var value = utils.TienePassword("LOGIN", valor)
			password = value
			continue
		} else if strings.HasPrefix(strings.ToLower(valor), "user") {
			var value = utils.TieneUser("LOGIN", valor)
			usuario = value
			continue
		} else {
			color.Yellow("[LOGIN]: Atributo no reconocido")
			return user_temp, false
		}
	}
	if id == "" || len(id) == 0 || len(id) > 4 || len(id) < 4 {
		return user_temp, false
	} else if password == "" || len(password) == 0 || len(password) > 10 {
		return user_temp, false
	} else if usuario == "" || len(usuario) == 0 || len(usuario) > 10 {
		return user_temp, false
	}
	user_temp.ID_Particion = global.Global_ID(id)
	user_temp.User = global.Global_Data(usuario)
	user_temp.Password = global.Global_Data(password)
	return user_temp, true
}

func LOGIN_EXECUTE(usuario string, password string, id_disco string) {
	if global.UsuarioLogeado.Logged_in {
		color.Red("[LOGIN]: Ya hay una secion activa")
		return
	}

	particion_montada, epm := utils.Buscar_ID_Montada(id_disco)
	if !epm {
		return
	}

	file, err := os.OpenFile(particion_montada.Path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[LOGIN]: Error al abrir archivo")
		return
	}
	defer file.Close()

	superbloque := structures.SuperBloque{}
	if particion_montada.Es_Particion_P {
		mbr_full, _ := utils.Obtener_FULL_MBR_FDISK(particion_montada.Path)
		for _, mbr := range mbr_full.Mbr_partitions {
			if utils.ToString(mbr.Part_name[:]) == utils.ToString(particion_montada.Particion_P.Part_name[:]) {
				if mbr.Part_status != 1 {
					color.Red("[LOGIN]: Particion no formateada")
					return
				}

				//Obtener SuperBloque
				if _, err := file.Seek(int64(particion_montada.Particion_P.Part_start), 0); err != nil {
					color.Red("[LOGIN]: Error en mover puntero")
					return
				}
				if err := binary.Read(file, binary.LittleEndian, &superbloque); err != nil {
					color.Red("[LOGIN]: Error en la lectura del superbloque")
					return
				}
				break
			}
		}
	} else if particion_montada.Es_Particion_L {
		ebr := structures.EBR{}
		if _, err := file.Seek(int64(particion_montada.Particion_L.Part_start), 0); err != nil {
			color.Red("[LOGIN]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
			color.Red("[LOGIN]: Error en la lectura del superbloque")
			return
		}
		if ebr.Part_mount != 1 {
			color.Red("[LOGIN]: Particion no formateada")
			return
		}

		//Obtener SuperBloque
		if _, err := file.Seek(int64(particion_montada.Particion_L.Part_start+size.SizeSuperBloque()), 0); err != nil {
			color.Red("[LOGIN]: Error en mover puntero")
			return
		}
		if err := binary.Read(file, binary.LittleEndian, &superbloque); err != nil {
			color.Red("[LOGIN]: Error en la lectura del superbloque")
			return
		}
	}

	contenido, eco := utils.GetContent(superbloque.S_inode_start+size.SizeTablaInodo(), particion_montada.Path)
	if !eco {
		return
	}
	contenido_split := strings.Split(contenido, "\n")
	for _, con := range contenido_split {
		if strings.Contains(con, ",U,") {
			if strings.Contains(con, usuario) {
				if strings.Contains(con, password) {
					//uid,tipo,grupo,usuario,contra
					usuario := strings.Split(con, ",")
					uid, _ := strconv.Atoi(usuario[0])
					global.UsuarioLogeado.UID = int32(uid)
					global.UsuarioLogeado.Tipo = 'U'
					global.UsuarioLogeado.Grupo = global.Global_Data(usuario[2])
					global.UsuarioLogeado.User = global.Global_Data(usuario[3])
					global.UsuarioLogeado.Password = global.Global_Data(usuario[4])
					global.UsuarioLogeado.ID_Particion = global.Global_ID(id_disco)
					// global.UsuarioLogeado.Logged_in = true
					goto c1
				}
			}
		}
	}
	global.UsuarioLogeado = global.DefaultUser
	color.Red("[LOGIN]: Usuario o Contraseña o Id de Disco incorrectos")
	return

c1:
	for _, con := range contenido_split {
		if strings.Contains(con, ",G,") {
			if strings.Contains(con, utils.ToString(global.UsuarioLogeado.Grupo[:])) {
				grupo := strings.Split(con, ",")
				gid, _ := strconv.Atoi(grupo[0])
				global.UsuarioLogeado.GID = int32(gid)
				global.GrupoUsuarioLoggeado.GID = int32(gid)
				global.GrupoUsuarioLoggeado.Tipo = 'G'
				global.GrupoUsuarioLoggeado.Nombre = grupo[2]
				goto c2
			}
		}
	}
	global.GrupoUsuarioLoggeado = global.DefaultGrupoUsuario
	global.UsuarioLogeado = global.DefaultUser
	color.Red("[LOGIN]: Grupo no encontrado")
	return

c2:
	// Ya que si paso por todo
	global.UsuarioLogeado.Mounted = particion_montada
	global.UsuarioLogeado.Logged_in = true
	color.Green("[LOGIN]: Usuario «" + usuario + "» Loggeado exitosamente (id): -> " + id_disco)
}
