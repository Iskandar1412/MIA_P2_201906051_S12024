import React from 'react';
import disco from '../../img/DiscoDuro1.png';

function Discos( props ) {

    const handleClickButton = (event) => {
        const valorBoton = event.target.getAttribute('data-value');
        // console.log("Valor del botón:", valorBoton);
        // console.log(props.command)
        
        for (let i = 0; i < props.command.length; i++) {
            if (props.command[i].Disco === valorBoton) {
                console.log(props.command[i])
                props.cambiarDiscos(props.command[i])
                props.onSeleccionar('particiones')
                
            }
        }
        // console.log(props.comandoDiscos)
        // A este se le hace verificación de informacion
    }

    return (
        <>
            
            <div className="discos">
                {props.command.map((item, index) => (
                    <button 
                        key={index}
                        className="buttonDisk"
                        data-value={item.Disco}
                        onClick={handleClickButton}
                    >
                        <img
                            src={disco}
                            alt="Imagen del botón" 
                            data-value={item.Disco}
                        />
                        <span className='valor-button' data-value={item.Disco}>{item.Disco}</span>
                    </button>
                ))}
                
            </div>
            
        </>
    );

}

export default Discos;

/*
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
*/