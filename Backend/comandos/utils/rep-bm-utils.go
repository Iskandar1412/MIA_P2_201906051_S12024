package utils

import "strings"

func EsImagen(comm string, commands ...string) (string, bool) {
	comm = strings.ToLower(comm)
	for _, c := range commands {
		if strings.HasPrefix(comm, c) {
			return c, true
		}
	}
	return "", false
}
