package obtencionpf

import (
	"encoding/json"
	"fmt"
	"proyecto/comandos/global"
	"proyecto/comandos/utils"
	"proyecto/estructuras/structures"
)

var ObtenerEstructuras []structures.MBR_Obtener

func Retorno_MBR() ([]byte, error) {
	//fmt.Println(ObtenerEstructuras)

	jsonData, err := json.Marshal(ObtenerEstructuras)
	if err != nil {
		return []byte{}, fmt.Errorf("error")
	}

	// fmt.Println(string(jsonData))
	return jsonData, nil
	//return "", fmt.Errorf("error")
}

func Retorno_Paths() ([]byte, error) {
	log := utils.ToString(global.UsuarioLogeado.Mounted.ID_Particion[:])
	// fmt.Println(log)
	ino, efno := ObtenerPaths(log)

	// ino, efno := ObtenerPaths(log)

	if efno != nil {
		return []byte{}, fmt.Errorf("Error")
	}
	// fmt.Println(ino)
	return []byte(ino), nil
}
