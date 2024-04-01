package utils

import (
	"bytes"
	"os"

	"github.com/fatih/color"
)

func ToByte(str string) []byte {
	result := make([]byte, 1)
	copy(result[:], str)
	return result
}

func ToString(b []byte) string {
	nullIndex := bytes.IndexByte(b, 0)
	if nullIndex == -1 {
		return string(b)
	}
	return string(b[:nullIndex])
}

func ExisteArchivo(comando string, archivo string) bool {
	if _, err := os.Stat(archivo); os.IsNotExist(err) {
		color.Red("[" + comando + "]: Archivo No Encontrado")
		return false
	}
	return true
}
