import { useEffect, useState } from "react";
import carpeta from '../../img/Carpeta1.png'
import archivo from '../../img/File.png'

function Dashboard (props) {
    const [userDash, setUserDash] = useState("")
    const [path, setPath] = useState("/")
    const [historial, setHistorial] = useState([]);
    const [contenidoArchivo, setContenidoArchivo] = useState("");
    const [data, setData] = useState([
        { nombre: "archivo1", tipo: "archivo", contenido: "Contenido del archivo 1" },
        { nombre: "archivo2", tipo: "archivo", contenido: "Contenido del archivo 1" },
        { nombre: "carpeta1", tipo: "carpeta", contenido: [
            { nombre: "archivo2", tipo: "archivo", contenido: "Contenido del archivo 2" },
            { nombre: "carpeta2", tipo: "carpeta", contenido: [
                { nombre: "archivo3", tipo: "archivo", contenido: "Contenido del archivo 3" }
            ]}
        ]}
    ]);
    
    useEffect(() => {
        const uss = localStorage.getItem('user')
        if (uss === "")  {
            props.onSeleccionar('particiones')
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
                {data.map((item, index) => (
                    item.tipo === "archivo" ? (
                        <button 
                            key={index}
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
                            key={index}
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