import React from 'react';
import particion from '../../img/Partition.png';

function Partitions(props) {

    const handleClickButton = (event) => {
        const valorBoton = event.target.getAttribute('data-value');
        console.log("Valor del botón:", valorBoton);

        // A este se le hace verificación de informacion
        props.onSeleccionar('login')
    }

    return (
        <>  
            <button className='button-ant' onClick={() => props.onSeleccionar('discos')} />
            <div className="particiones">
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

                <button 
                    className="buttonPartition"
                    data-value='Particion2'
                    onClick={handleClickButton}
                >
                    <img
                        src={particion}
                        alt="Imagen del botón"
                        data-value='Particion2'
                    />
                    <span className='valor-button' data-value='Particion2'>Particion2</span>
                </button>

                <button 
                    className="buttonPartition"
                    data-value='Particion3'
                    onClick={handleClickButton}
                >
                    <img
                        src={particion}
                        alt="Imagen del botón"
                        data-value='Particion3'
                    />
                    <span className='valor-button' data-value='Particion3'>Particion3</span>
                </button>

                


            </div>
            
        </>
    );

}

export default Partitions;