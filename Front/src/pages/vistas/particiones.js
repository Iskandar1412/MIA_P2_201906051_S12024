import React from 'react';
import particion from '../../img/Partition.png';

function Partitions( props ) {

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
                if (props.disco.Mbr_partitions[i].Status !== 1) {
                    alert("Disco no formateado como EXT2 o EXT3")
                    return
                } 
                props.seleccionParticion(props.disco.Mbr_partitions[i].Particion)
                props.onSeleccionar('login')
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