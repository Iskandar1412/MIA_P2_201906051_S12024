package adminusers

import (
	"proyecto/comandos/global"
	"proyecto/comandos/utils"

	"github.com/fatih/color"
)

func LOGOUT_EXECUTE() {
	if !global.UsuarioLogeado.Logged_in {
		color.Red("[LOGOUT]: No hay una sesion activa")
		return
	}

	usuario_temp := global.UsuarioLogeado
	global.UsuarioLogeado = global.DefaultUser
	global.GrupoUsuarioLoggeado = global.DefaultGrupoUsuario

	color.Green("[LOGOUT]: Usuario «" + utils.ToString(usuario_temp.User[:]) + "» cerro secion exitosamente - (ID): -> " + utils.ToString(usuario_temp.ID_Particion[:]))
}
