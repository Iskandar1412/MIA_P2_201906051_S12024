import { useState } from "react";
import Discos from "./vistas/discos";
import Partitions from "./vistas/particiones";
import Login from "./vistas/login";
import Dashboard from "./vistas/dashboard";


function Pantalla2({ command, carpetasOb }) {
    const [componenteActual, setComponenteActual] = useState('discos');

    const cambiarComponente = (nuevoComponente) => {
        setComponenteActual(nuevoComponente);
    };

    const commandAnterior = command
    const [newCommandDisk, setNewCommandDisk] = useState([])
    const cambiarNuevoCommandDisk = (nuevocommando) => {
        setNewCommandDisk(nuevocommando)
    }

    const [particionSeleccionada, setParticionSeleccionada] = useState("")
    const cambiarParticionSeleccionada = (particion) => {
        setParticionSeleccionada(particion)
    }

    const [directorios, setDirectorios] = useState([]);
    const cambiarDirectorios = (dir) => {
        
        setDirectorios(dir)
    }

    let componenteMostrar;
    if (componenteActual === 'discos') {
        componenteMostrar = <Discos onSeleccionar={cambiarComponente} command={commandAnterior} cambiarDiscos={cambiarNuevoCommandDisk} />
    } else if (componenteActual === 'particiones') {
        componenteMostrar = <Partitions  onSeleccionar={cambiarComponente} disco={newCommandDisk} seleccionParticion={cambiarParticionSeleccionada} />
    } else if (componenteActual === 'login') {
        componenteMostrar = <Login  onSeleccionar={cambiarComponente} particion={particionSeleccionada} cambiarDirectorios={cambiarDirectorios} />
    } else if (componenteActual === 'dashboard') {
        componenteMostrar = <Dashboard onSeleccionar={cambiarComponente} dir={directorios} capetas={carpetasOb} />
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