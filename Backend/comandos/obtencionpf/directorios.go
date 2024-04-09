package obtencionpf

import (
	"io"
	"os"
)

func isDirEmpty(dirname string) (bool, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // trata de leer algun nombre de directorio
	if err == nil {
		// existe algun archivo en el directorio
		return false, nil
	}

	// directorio vacio
	if err == io.EOF {
		return true, nil
	}

	return false, err
}
