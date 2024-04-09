import React, { useEffect, useState } from 'react';
import axios from 'axios'
import { pathbackend } from '../path';

function getCommand(comm, ...commands) {
    comm = comm.toLowerCase();
    for (let c of commands) {
        if (comm.startsWith(c)) {
            return c;
        }
    }
    return "";
}


function Pantalla1({ info, carpetas, cambiarDot }) {
    const [command, setCommand] = useState('');
    const [commandsSaved, setCommandSaved] = useState([]);
    

    const HandleCommandChange = (event) => {
        setCommand(event.target.value);
    };

    const handleEnter = (e) => {
        if (e.key === "Enter") {
            //console.log("Entere presionado")
            handlePostCommand()
        }
    }

    const ObtenerInformacionMBR2 = async () => {
        try {
            const res = await axios.get(pathbackend+"/obtainmbr")

            if (res.status === 200) {
                console.log(res.data.datos)
                const jsonData = JSON.parse(res.data.datos)
                info(jsonData)
            }
        } catch (e) { }
    }

    const EnviarInformacionCarpetas = async () => {
        try {
            const res = await axios.get(pathbackend+'/obtain-carpetas-archivos')
            if (res.status === 200) {
                const jsonData = JSON.parse(res.data.datos);
                // console.log(jsonData)
                // console.log(jsonData)
                carpetas(jsonData)
            }
        } catch (e) { }
    }

    useEffect(() => {
        const ObtenerInformacionMBR = async () => {
            try {
                const res = await axios.get(pathbackend+"/obtainmbr")

                if (res.status === 200) {
                    // console.log(res.data.datos)
                    const jsonData = JSON.parse(res.data.datos)
                    info(jsonData)
                }
            } catch (e) { }
        }
        ObtenerInformacionMBR()
        return () => {
            
        };
    },  [])

    const handleReports = async () => {
        try {
            const res = await axios.get(pathbackend+'/reportesobtener')
            if (res.status === 200) {
                const jsonData = JSON.parse(res.data.datos);
                // console.log(jsonData)
                cambiarDot(jsonData)
                // console.log(jsonData)
            }
        } catch (e) { }
    }

    const postInformacion = async (objeto) => {
        try {
            const res = await axios.post(pathbackend+'/command', objeto);
            if (res.status === 200) {
                ObtenerInformacionMBR2()
            }
        } catch (e) { }
    }

    const postContenido = async (objeto) => {
        try {
            const res = await axios.post(pathbackend+'/command', objeto);
            if (res.status === 200) {
                EnviarInformacionCarpetas()
            }
        } catch (e) { }
    }

    const postReportes = async (objeto) => {
        try {
            const res = await axios.post(pathbackend+'/command', objeto);
            if (res.status === 200) {
                handleReports()
            }
        } catch (e) { }
    }

    const handlePostCommand =  () => {
        let objeto = {
            comando: ""
        }
        if (command.trim() !== '') {
            //Hacer operacion para enviar
            let com = getCommand(command.toLowerCase(), "mkdisk", "fdisk", "rmdisk", "mount", "unmount", "mkfs");
            let com2 = getCommand(command.toLowerCase(), "mkfile", "cat", "remove", "edit", "rename", "mkdir", "copy", "move", "find", "chown", "chgrp", "chmod", "mkgrp", "rmgrp", "mkusr", "rmusr");
            let com3 = getCommand(command.toLowerCase(), "rep");
            let com4 = getCommand(command.toLowerCase(), "login", "logout");
            if (["mkdisk", "fdisk", "rmdisk", "mount", "unmount", "mkfs"].includes(com)) {
                //actualizar los discos
                objeto.comando = command
                postInformacion(objeto)
                setCommandSaved([...commandsSaved, command]);
                setCommand('');
                ObtenerInformacionMBR2()
            } else if (["mkfile", "remove", "edit", "rename", "mkdir", "copy", "move", "chown", "chgrp", "chmod", "mkgrp", "rmgrp", "mkusr", "rmusr"].includes(com2)) {
                //actualizar informacion
                objeto.comando = command
                //postear siguientes comandos
                postContenido(objeto)
                setCommandSaved([...commandsSaved, command]);
                setCommand('');
                EnviarInformacionCarpetas()
            } else if (["rep"].includes(com3)){
                //actualizar reportes
                //postear os grafos ---faltante
                objeto.comando = command
                postReportes(objeto)
                handleReports()
                setCommandSaved([...commandsSaved, command]);
                setCommand('');
            } else if (["login", "logout"].includes(com4)){
                alert("Login/Logout se hace desde pestaña inicio seción")
                return
            } else {
                alert("Comando no reconocido")
                return
            }
        } else { 
            alert("No hay ningun comando puesto")
            return
        }
        
    };

    return (
        <>
            <div>
                <textarea
                    className='textObtain'
                    value={commandsSaved.join('\n')}
                    readOnly
                />
                <input
                    type="text"
                    value={command}
                    onChange={HandleCommandChange}
                    className='subComm'
                    onKeyDown={handleEnter}
                />
                <button onClick={handlePostCommand} className='subEnv'>Enviar</button>

            </div>
        </>
    );
}

export default Pantalla1;