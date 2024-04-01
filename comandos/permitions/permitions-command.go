package permitions

import (
	"strings"

	"github.com/fatih/color"
)

func PermissionsCommandProps(permission string, instructions []string) {
	var _path string //chown chmod
	var _user string //chown
	//_r			 //chown chmod
	var _grp string //chgrp
	var _ugo string //chmod
	var _r bool
	var er bool
	/*
		var _user string //chgrp
	*/
	if strings.ToUpper(permission) == "CHOWN" {
		_path, _user, _r, er = Values_CHOWN(instructions)
		if !er {
			return
		} else {
			CHOWN_EXECUTE(_path, _user, _r)
		}
	} else if strings.ToUpper(permission) == "CHGRP" {
		_user, _grp, er = Values_CHGRP(instructions)
		if !er {
			return
		} else {
			CHGRP_EXECUTE(_user, _grp)
		}
	} else if strings.ToUpper(permission) == "CHMOD" {
		_path, _ugo, _r, er = Values_CHMOD(instructions)
		if !er {
			return
		} else {
			CHMOD_EXECUTE(_path, _ugo, _r)
		}
	} else {
		color.Red("[PermissionsCommandProps]: Internal Error")
	}
}
