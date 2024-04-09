import React from 'react';
import particion from '../../img/Partition.png';
import axios from 'axios'
import { pathbackend } from '../../path';

function Partitions( props ) {
    const Carpetas = async () => {
        try {
            const res = await axios.get(pathbackend+'/obtain-carpetas-archivos')
            if (res.status === 200) {
                const jsonData = JSON.parse(res.data.datos);
                // console.log(jsonData)
                // console.log(jsonData)
                props.cambiarDirectorios(jsonData)
            }
        } catch (e) { }
    }

    const handleClickButton = (event) => {
        const valorBoton = event.target.getAttribute('data-value');
        // console.log(props.disco.particiones.length)
        for (let i = 0; i  < props.disco.Mbr_partitions.length; i++){
            //console.log(props.disco.particiones[i])
            if (props.disco.Mbr_partitions[i].Particion === valorBoton) {
                // console.log(props.disco.particiones[i])
                //verificar que sea ext2 o ext3
                if (props.disco.Mbr_partitions[i].Type === 'E') {
                    alert("Disco Extendido no se puede usar")
                    return
                }
                if (props.disco.Mbr_partitions[i].Status === -1) {
                    alert("Disco no montado")
                    return
                }
                if (props.disco.Mbr_partitions[i].Id_mounted === "") {
                    alert("No tiene ID el disco por lo que esta desmontado")
                    return
                }
                if (props.disco.Mbr_partitions[i].Status !== 1) {
                    alert("Disco no formateado como EXT2 o EXT3")
                    return
                }
                console.log(localStorage.getItem('user'))
                var usuario_logeado = localStorage.getItem('user')
                if (usuario_logeado !== null) {
                    props.seleccionParticion(props.disco.Mbr_partitions[i].Particion)
                    localStorage.setItem('id_particion', props.disco.Mbr_partitions[i].Id_mounted)
                    Carpetas()
                    props.onSeleccionar('dashboard')
                    break
                }
                props.seleccionParticion(props.disco.Mbr_partitions[i].Particion)
                localStorage.setItem('id_particion', props.disco.Mbr_partitions[i].Id_mounted)
                props.onSeleccionar('login')
                break
            }
        }
        // console.log("Valor del botón:", valorBoton);
        // console.log(props.disco)
        // A este se le hace verificación de informacion
        //
    }

    return (
        <>  
            <button className='button-ant' onClick={() => props.onSeleccionar('discos')} />
            
            <div className="particiones">

                {props.disco.Mbr_partitions.map((item, index) => (
                    item.Particion !== "" && (
                        <button 
                            key={index}
                            className="buttonPartition"
                            data-value={item.Particion}
                            onClick={handleClickButton}
                        >
                            <img
                                src={particion}
                                alt="Imagen del botón" 
                                data-value={item.Particion}
                            />
                            <span className='valor-button' data-value={item.Particion}>{item.Particion}</span>
                        </button>
                    )
                ))}
                
            </div>
            
        </>
    );

}

export default Partitions;

/*
                <button 
                    className="buttonPartition"
                    data-value='Particion1'
                    onClick={handleClickButton}
                >
                    <img
                        src={particion}
                        alt="Imagen del botón" 
                        data-value='Particion1'
                    />
                    <span className='valor-button' data-value='Particion1'>Particion1</span>
                </button>
*/