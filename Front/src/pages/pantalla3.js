import filepdf from '../img/pdf.png'
import filepng from '../img/png.png'
import filejpg from '../img/jpg.png'
import filetxt from '../img/txt.png'

function Pantalla3() {

    const handleRepButton = (event) => {
        const valorBoton = event.target.getAttribute('data-value');
        console.log("Valor del botón:", valorBoton);
        
        // A este se le hace verificación de informacion
        
    }

    return (
        <>
            <div className="vistas">
                <div className="discos">
                    <button 
                        className="buttonDisk"
                        data-value='report'
                        onClick={handleRepButton}
                    >
                        <img
                            src={filepdf}
                            alt="Imagen del botón" 
                            data-value='report'
                        />
                        <span className='valor-button' data-value='report'>report</span>
                    </button>

                    <button 
                        className="buttonDisk"
                        data-value='report png'
                        onClick={handleRepButton}
                        >
                        <img
                            src={filepng}
                            alt="Imagen del botón" 
                            data-value='report png'
                            />
                        <span className='valor-button' data-value='report png'>report png</span>
                    </button>

                    <button 
                        className="buttonDisk"
                        data-value='report jpg'
                        onClick={handleRepButton}
                    >
                        <img
                            src={filejpg}
                            alt="Imagen del botón" 
                            data-value='report jpg'
                        />
                        <span className='valor-button' data-value='report jpg'>report jpg</span>
                    </button>

                    <button 
                        className="buttonDisk"
                        data-value='report txt'
                        onClick={handleRepButton}
                    >
                        <img
                            src={filetxt}
                            alt="Imagen del botón" 
                            data-value='report txt'
                        />
                        <span className='valor-button' data-value='report txt'>report txt</span>
                    </button>

                </div>
            </div>
        </>
    );

}

export default Pantalla3;