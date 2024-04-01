package comandos

import "strings"

func getCommand(comm string, commands ...string) string {
	comm = strings.ToLower(comm)
	for _, c := range commands {
		if strings.HasPrefix(comm, c) {
			return c
		}
	}
	return ""
}
