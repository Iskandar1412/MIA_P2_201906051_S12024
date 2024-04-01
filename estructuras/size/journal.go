package size

import (
	"proyecto/estructuras/structures"
	"unsafe"
)

func SizeJournal() int32 { //68 bytes
	a01 := unsafe.Sizeof(structures.Journal{}.J_Tipo_Operacion)
	a01 += unsafe.Sizeof(structures.Journal{}.J_Tipo)
	a01 += unsafe.Sizeof(structures.Journal{}.J_Path)
	a01 += unsafe.Sizeof(structures.Journal{}.J_Contenido)
	a01 += unsafe.Sizeof(structures.Journal{}.J_Fecha)
	a01 += unsafe.Sizeof(structures.Journal{}.J_Size)
	a01 += unsafe.Sizeof(structures.Journal{}.J_Sig)
	a01 += unsafe.Sizeof(structures.Journal{}.J_Start)
	return int32(a01)
}
