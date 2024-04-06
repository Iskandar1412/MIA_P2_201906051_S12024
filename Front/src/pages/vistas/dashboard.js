import { useEffect, useState } from "react";
import carpeta from '../../img/Carpeta1.png'
import archivo from '../../img/File.png'

function Dashboard ( props ) {
    const [userDash, setUserDash] = useState("")
    const [path, setPath] = useState("/")
    const [historial, setHistorial] = useState([]);
    const [contenidoArchivo, setContenidoArchivo] = useState("");
    const [data, setData] = useState([]);

    
    
    useEffect(() => {
        const uss = localStorage.getItem('user')
        if (uss === "")  {
            props.onSeleccionar('particiones')
        }
        //console.log("dir", props.dir,  "carpeta", props.capetas)
        setData(props.dir)
        if ((props.capetas === undefined) || (props.capetas.length === 0) || (props.capetas.length < props.dir.length)) {
            setData(props.dir)
        } else {
            setData(props.capetas)
        }
        setUserDash(uss)
    }, [props])

    const handleEndSession = () => {
        localStorage.removeItem('user')
        props.onSeleccionar('login')
    }

    const handleMostrarContenido = (contenido) => {
        setContenidoArchivo(contenido);
    };

    const handleClick = (contenido, nombre) => {
        setHistorial([...historial, {data, path}]);
        setData(contenido);
        if (path === "/") {
            setPath(`${path}${nombre}`)
        } else {
            setPath(`${path}/${nombre}`);
        }
    };

    const handleBack = () => {
        if (historial.length > 0) {
            const {data: prevData, path: prevPath} = historial.pop();
            setData(prevData);
            setHistorial([...historial]);
            setPath(prevPath);
        }
    };

    const handleCerrarContenido = () => {
        setContenidoArchivo(null);
    };

    return (
        <>
        <div className="dash">
            <div className="navegadores-dash">
                <button className='button-ant-dash' onClick={handleEndSession} />
                <p className="name_user">{userDash}</p>
                <button className='button-sal-carpet' onClick={handleBack} />
                <div className="div-path">{path}</div>
            </div>
            <div className="content-dashboard">
                {data.map((item) => (
                    item.tipo === "archivo" ? (
                        <button 
                            key={item.nombre}
                            className="buttonDisk"
                            onClick={() => handleMostrarContenido(item.contenido)}
                        >
                            <img
                                src={archivo}
                                alt="Imagen del botón"
                            />
                            <span>{item.nombre}</span>
                        </button>
                    ) : (
                        <button 
                            key={item.nombre}
                            className="buttonPartition"
                            onClick={() => handleClick(item.contenido, item.nombre)}
                        >
                            <img
                                src={carpeta}
                                alt="Imagen del botón"
                            />
                            <span>{item.nombre}</span>
                        </button>
                    )
                ))}
            
                {contenidoArchivo && (
                    <div className="modal" style={{ backgroundColor: 'rgba(255, 255, 255, 0.9)' }}>
                        <div className="boton-exit">
                            <button onClick={handleCerrarContenido} className="button-x" />
                        </div>
                        <div>{contenidoArchivo}</div>
                    </div>
                )}
            </div>
        </div>
        </>
    );

}

export default Dashboard;