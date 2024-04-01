package admindisk

import (
	"proyecto/comandos/partition"
	"proyecto/comandos/utils"
	"strings"

	"github.com/fatih/color"
)

func DiskCommandProps(command string, instructions []string) {
	var _size int32       //mkdisk  fdisk
	var _fit byte         //mkdisk  fdisk
	var _unit byte        //mkdisk  fdisk
	var _driveletter byte //rmdisk  fdisk mount
	var _name [16]byte    //fdisk   mount
	var _type byte        //fdisk
	var _delete string    //fdisk
	var _add int32        //fdisk
	var _id string        //unmount mkfs
	var _type_mkfs string //mkfs
	var _fs string        //mkfs
	/*
	 */

	if strings.ToUpper(command) == "MKDISK" {
		_size, _fit, _unit = Values_MKDISK(instructions)
		if _size <= 0 || _fit == '0' || _unit == '0' {
			color.Yellow("[MKDISK]: Error to asign values for '" + string(_name[:]) + "'")
		} else {
			MKDISK_Create(_size, _fit, _unit)
		}
	} else if strings.ToUpper(command) == "FDISK" {
		_size, _driveletter, _name, _unit, _type, _fit, _delete, _add = partition.Values_FDISK(instructions)
		if _size <= 0 || utils.ToString(_name[:]) == "" || _driveletter == '0' {
			if utils.ToString(_name[:]) == "" {
				color.Yellow("[FDISK]: Error to asign values for (unamed) disk")
			} else {
				if _delete == "FULL" {
					partition.FDISK_Create(_size, _driveletter, _name[:], _unit, _type, _fit, _delete, _add)
				} else if _add != 0 {
					partition.FDISK_Create(_size, _driveletter, _name[:], _unit, _type, _fit, _delete, _add)
				} else {
					color.Yellow("[FDISK]: Error to asign values for disk '" + string(_name[:]) + "'")
				}
			}
		} else {
			partition.FDISK_Create(_size, _driveletter, _name[:], _unit, _type, _fit, _delete, _add)
		}
		//fmt.Println(_size, _driveletter, _name, _unit, _type, _fit, _delete, _add)
	} else if strings.ToUpper(command) == "RMDISK" {
		_driveletter, _error := Values_RMDISK(instructions)
		if _driveletter == '0' && !_error {
			color.Yellow("[RMDISK]: Error to asign values")
		} else {
			RMDISK_EXECUTE(_driveletter)
		}
	} else if strings.ToUpper(command) == "MOUNT" {
		_driveletter, _name, err := Values_Mount(instructions)
		if err && (_driveletter == '0') {
			color.Yellow("[MOUNT]: Error to asign values")
		} else {
			MOUNT_EXECUTE(_driveletter, _name[:])
			//fmt.Println("Mount", _driveletter, _name)
		}
		//fmt.Println(Partitions_Mounted)
	} else if strings.ToUpper(command) == "UNMOUNT" {
		var err bool
		_id, err = Values_Unmount(instructions)
		if !err {
			color.Yellow("[UNMOUNT]: Error to asign values")
		} else {
			UNMOUNT_EXECUTE("UNMOUNT", _id)
			// fmt.Println(_id)
		}
	} else if strings.ToUpper(command) == "MKFS" {
		var err bool
		_id, _type_mkfs, _fs, err = Values_MKFS(instructions)
		if !err {
			color.Yellow("[MKFS]: Error to asign values")
		} else {
			MKFS_EXECUTE(_id, _type_mkfs, _fs)
			//fmt.Println(_id, _type_mkfs, _fs)
		}
	} else {
		color.Red("[DiskComandProps]: Internal Error")
	}
}
