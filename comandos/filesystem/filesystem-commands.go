package filesystem

import (
	"strings"

	"github.com/fatih/color"
)

func FilesCommandProps(file string, instructions []string) {
	var _er bool
	var _path string //mkfile remove edit rename mkdir copy move find
	//existe -r 		//mkfile mkdir
	var _r bool
	var _size int32     //mkfile
	var _cont string    //mkfile
	var _file []string  //cat
	var _name string    //rename find
	var _destino string //move   copy
	/*
	 */
	if strings.ToUpper(file) == "MKFILE" {
		_path, _r, _size, _cont, _er = Values_MKFILE(instructions)
		if !_er || _path == "" || len(_path) == 0 {
			return
		} else {
			MKFILE_EXECUTE(_path, _r, _size, _cont)
		}

	} else if strings.ToUpper(file) == "CAT" {
		_file, _er = Values_CAT(instructions)
		if !_er {
			return
		} else {
			CAT_EXECUTE(_file)
		}
	} else if strings.ToUpper(file) == "REMOVE" {
		_path, _er = Values_REMOVE(instructions)
		if !_er {
			return
		} else {
			REMOVE_EXECUTE(_path)
		}
	} else if strings.ToUpper(file) == "EDIT" {
		_path, _cont, _er = Values_EDIT(instructions)
		if !_er {
			return
		} else {
			EDIT_EXECUTE(_path, _cont)
		}
	} else if strings.ToUpper(file) == "RENAME" {
		_path, _name, _er = Values_RENAME(instructions)
		if !_er {
			return
		} else {
			RENAME_EXECUTE(_path, _name)
		}
	} else if strings.ToUpper(file) == "MKDIR" {
		_path, _r, _er = Values_MKDIR(instructions)
		if !_er {
			return
		} else {
			MKDIR_EXECUTE(_path, _r)
		}
	} else if strings.ToUpper(file) == "COPY" {
		_path, _destino, _er = Values_COPY(instructions)
		if !_er {
			return
		} else {
			COPY_EXECUTE(_path, _destino)
		}
	} else if strings.ToUpper(file) == "MOVE" {
		_path, _destino, _er = Values_MOVE(instructions)
		if !_er {
			return
		} else {
			MOVE_EXECUTE(_path, _destino)
		}
	} else if strings.ToUpper(file) == "FIND" {
		_path, _name, _er = Values_FIND(instructions)
		if !_er {
			return
		} else {
			FIND_EXECUTE(_path, _name)
		}
	} else {
		color.Red("[FilesCommandProps]: Internal Error")
	}
}
