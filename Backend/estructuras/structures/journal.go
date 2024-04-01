package structures

type Journal struct {
	J_Tipo_Operacion [10]byte
	J_Tipo           byte
	J_Path           [100]byte
	J_Contenido      [100]byte
	J_Fecha          int32
	J_Size           int32
	J_Sig            int32
	J_Start          int32
}
