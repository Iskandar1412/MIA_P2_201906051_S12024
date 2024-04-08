import React, { useState } from 'react';

function getCommand(comm, ...commands) {
    comm = comm.toLowerCase();
    for (let c of commands) {
        if (comm.startsWith(c)) {
            return c;
        }
    }
    return "";
}


function Pantalla1({ info, carpetas, dots }) {
    const [command, setCommand] = useState('');
    const [commandsSaved, setCommandSaved] = useState([]);
    const path = "http://localhost:8080"
    const HandleCommandChange = (event) => {
        setCommand(event.target.value);
    };

    const handleEnter = (e) => {
        if (e.key === "Enter") {
            //console.log("Entere presionado")
            handlePostCommand()
        }
    }

    const ObtenerInformacionMBR = async () => {
        try {
            const res = await fetch(path + "/obtain-mbr", {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' },
            })

            if (res.ok) {
                const response = await res.json();
                console.log("siuuuuuuuu", response)
                info(command)
            }
        } catch (e) {
            console.error(e)
        }
    }

    const postInformacion = async (objeto) => {
        try {
            const response = await fetch(path + "/command", {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(objeto)
            })

            if (response.ok) {
                ObtenerInformacionMBR()
                console.log("ok")
            }
            
        } catch (e) {
            console.error("Error")
            info(command)
        }
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
            } else if (["mkfile", "remove", "edit", "rename", "mkdir", "copy", "move", "chown", "chgrp", "chmod", "mkgrp", "rmgrp", "mkusr", "rmusr"].includes(com2)) {
                //actualizar informacion
                carpetas(command)
                objeto.comando = command
                setCommandSaved([...commandsSaved, command]);
                setCommand('');
            } else if (["cat", "find"].includes(com2)) {
                //no suceden acciones
                objeto.comando = command
                setCommandSaved([...commandsSaved, command]);
                setCommand('');
            } else if (["rep"].includes(com3)){
                //actualizar reportes
                dots(command)
                objeto.comando = command
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