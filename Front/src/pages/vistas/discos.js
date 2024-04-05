import React from 'react';
import disco from '../../img/DiscoDuro1.png';

function Discos(props) {

    const handleClickButton = (event) => {
        const valorBoton = event.target.getAttribute('data-value');
        console.log("Valor del botón:", valorBoton);
        
        // A este se le hace verificación de informacion
        props.onSeleccionar('particiones')
    }

    return (
        <>
            <div className="discos">
                <button 
                    className="buttonDisk"
                    data-value='a.dsk'
                    onClick={handleClickButton}
                >
                    <img
                        src={disco}
                        alt="Imagen del botón" 
                        data-value='a.dsk'
                    />
                    <span className='valor-button' data-value='a.dsk'>a.dsk</span>
                </button>

                <button 
                    className="buttonDisk"
                    data-value='b.dsk'
                    onClick={handleClickButton}
                >
                    <img
                        src={disco}
                        alt="Imagen del botón"
                        data-value='b.dsk'
                    />
                    <span className='valor-button' data-value='b.dsk'>b.dsk</span>
                </button>

                <button 
                    className="buttonDisk"
                    data-value='c.dsk'
                    onClick={handleClickButton}
                >
                    <img
                        src={disco}
                        alt="Imagen del botón"
                        data-value='c.dsk'
                    />
                    <span className='valor-button' data-value='c.dsk'>c.dsk</span>
                </button>

                


            </div>
            
        </>
    );

}

export default Discos;