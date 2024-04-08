package adminusers

import (
	"proyecto/comandos/utils"
	"strings"

	"github.com/fatih/color"
)

func UserCommandProps(command string, instructions []string) bool {
	var _user string
	var _pass string
	var _id string
	if strings.ToUpper(command) == "LOGIN" {
		valor_usuario, err := Values_LOGIN(instructions)
		if !err {
			color.Red("[Login]: Error to assign values")
			return false
		} else {
			_user = utils.ToString(valor_usuario.User[:])
			_pass = utils.ToString(valor_usuario.Password[:])
			_id = utils.ToString(valor_usuario.ID_Particion[:])
			return LOGIN_EXECUTE(_user, _pass, _id)
		}
	} else if strings.ToUpper(command) == "LOGOUT" {
		return LOGOUT_EXECUTE()
	} else {
		color.Red("[UserComamand]: Internal Error")
		return false
	}
}

func GroupCommandProps(group string, instructions []string) {
	var _name string //mkgrp rmgrp
	var er bool
	var _user string //mkusr rmusr
	var _pass string //mkusr
	var _grp string  //mkusr
	/*
	 */
	if strings.ToUpper(group) == "MKGRP" {
		_user, er = Values_MKGRP(instructions)
		if !er {
			color.Red("[MKGRP]: Error to assing values")
		} else {
			MKGRP_EXECUTE(_user)
		}
	} else if strings.ToUpper(group) == "RMGRP" {
		_name, er = Values_RMGRP(instructions)
		if !er {
			color.Red("[RMGRP]: Error to assing values")
		} else {
			RMGRP_EXECUTE(_name)
		}
	} else if strings.ToUpper(group) == "MKUSR" {
		_name, _pass, _grp, er = Values_MKUSR(instructions)
		if !er {
			color.Red("[MKUSR]: Error to asign values")
		} else {
			// fmt.Println(_name, _pass, _grp, er)
			MKUSR_EXECUTE(_name, _pass, _grp)
		}

	} else if strings.ToUpper(group) == "RMUSR" {
		_name, er := Values_RMUSR(instructions)
		if !er {
			color.Red("[RMUSR]: Error to asign values")
		} else {
			RMUSR_EXECUTE(_name)
		}
		// fmt.Println("Eliminando usuario en la parcicion")
	} else {
		color.Red("[GroupCommandProps]: Internal Error")
	}
}
