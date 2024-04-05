import { useState } from "react";
import Discos from "./vistas/discos";
import Partitions from "./vistas/particiones";


function Pantalla2() {
    const [componenteActual, setComponenteActual] = useState('discos');

    const cambiarComponente = (nuevoComponente) => {
        setComponenteActual(nuevoComponente);
    };

    let componenteMostrar;
    if (componenteActual === 'discos') {
        componenteMostrar = <Discos onSeleccionar={cambiarComponente}/>
    } else if (componenteActual === 'particiones') {
        componenteMostrar = <Partitions  onSeleccionar={cambiarComponente}/>
    } else if (componenteActual === 'login') {
        //login
    } else if (componenteMostrar === 'dashboard') {
        //dashboard
    }

    return (
        <>
            <div className="vistas">
                {componenteMostrar}
            </div>
        </>
    );

}

export default Pantalla2;