import { useState } from "react";
import Discos from "./vistas/discos";
import Partitions from "./vistas/particiones";
import Login from "./vistas/login";
import Dashboard from "./vistas/dashboard";


function Pantalla2() {
    const [componenteActual, setComponenteActual] = useState('discos');

    const cambiarComponente = (nuevoComponente) => {
        setComponenteActual(nuevoComponente);
    };

    let componenteMostrar;
    if (componenteActual === 'discos') {
        componenteMostrar = <Discos onSeleccionar={cambiarComponente} />
    } else if (componenteActual === 'particiones') {
        componenteMostrar = <Partitions  onSeleccionar={cambiarComponente} />
    } else if (componenteActual === 'login') {
        componenteMostrar = <Login  onSeleccionar={cambiarComponente} />
    } else if (componenteActual === 'dashboard') {
        componenteMostrar = <Dashboard onSeleccionar={cambiarComponente} />
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