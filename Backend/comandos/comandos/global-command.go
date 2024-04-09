package comandos

import (
	"proyecto/comandos/admindisk"
	"proyecto/comandos/adminusers"
	"proyecto/comandos/filesystem"
	"proyecto/comandos/permitions"
	"proyecto/comandos/report"
	"strings"

	"github.com/fatih/color"
)

// var Particiones_Montadas []string

func GlobalCom(lista []string) bool {
	for _, comm := range lista {
		// Administracion de Discos
		if (strings.HasPrefix(strings.ToLower(comm), "mkdisk")) || (strings.HasPrefix(strings.ToLower(comm), "fdisk")) || (strings.HasPrefix(strings.ToLower(comm), "rmdisk")) || (strings.HasPrefix(strings.ToLower(comm), "mount")) || (strings.HasPrefix(strings.ToLower(comm), "unmount")) || (strings.HasPrefix(strings.ToLower(comm), "mkfs")) {
			comandos := ObtenerComandos(comm)
			command := getCommand(strings.ToLower(comm), "mkdisk", "fdisk", "rmdisk", "mount", "unmount", "mkfs")
			admindisk.DiskCommandProps(strings.ToUpper(command), comandos)
			// Reportes
		} else if strings.HasPrefix(strings.ToLower(comm), "rep") {
			comandos := ObtenerComandos(comm)
			report.ReportCommandProps("REP", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "pause") {
			Pause()
			// Files
		} else if (strings.HasPrefix(strings.ToLower(comm), "mkfile")) || (strings.HasPrefix(strings.ToLower(comm), "cat")) || (strings.HasPrefix(strings.ToLower(comm), "remove")) || (strings.HasPrefix(strings.ToLower(comm), "edit")) || (strings.HasPrefix(strings.ToLower(comm), "rename")) || (strings.HasPrefix(strings.ToLower(comm), "mkdir")) || (strings.HasPrefix(strings.ToLower(comm), "copy")) || (strings.HasPrefix(strings.ToLower(comm), "move")) || (strings.HasPrefix(strings.ToLower(comm), "find")) {
			comandos := ObtenerComandos(comm)
			command := getCommand(strings.ToLower(comm), "mkfile", "cat", "remove", "edit", "rename", "mkdir", "copy", "move", "find")
			filesystem.FilesCommandProps(strings.ToUpper(command), comandos)
			// Permisos
		} else if (strings.HasPrefix(strings.ToLower(comm), "chown")) || (strings.HasPrefix(strings.ToLower(comm), "chgrp")) || (strings.HasPrefix(strings.ToLower(comm), "chmod")) {
			comandos := ObtenerComandos(comm)
			command := getCommand(strings.ToLower(comm), "chown", "chgrp", "chmod")
			permitions.PermissionsCommandProps(strings.ToUpper(command), comandos)
			// Usuarios
		} else if (strings.HasPrefix(strings.ToLower(comm), "login")) || (strings.HasPrefix(strings.ToLower(comm), "logout")) {
			comandos := ObtenerComandos(comm)
			command := getCommand(strings.ToLower(comm), "login", "logout")
			return adminusers.UserCommandProps(strings.ToUpper(command), comandos)
			// Grupo
		} else if (strings.HasPrefix(strings.ToLower(comm), "mkgrp")) || (strings.HasPrefix(strings.ToLower(comm), "rmgrp")) || (strings.HasPrefix(strings.ToLower(comm), "mkusr")) || (strings.HasPrefix(strings.ToLower(comm), "rmusr")) || (strings.HasPrefix(strings.ToLower(comm), "mkgrp")) {
			comandos := ObtenerComandos(comm)
			command := getCommand(strings.ToLower(comm), "mkgrp", "rmgrp", "mkusr", "rmusr")
			adminusers.GroupCommandProps(strings.ToUpper(command), comandos)
		} else {
			color.Red("Comando no Reconocido")
		}
	}
	return true
}
