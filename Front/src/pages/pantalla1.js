import React, { useState } from 'react';


function Pantalla1() {
    const [command, setCommand] = useState('');
    const [commandsSaved, setCommandSaved] = useState([]);

    const HandleCommandChange = (event) => {
        setCommand(event.target.value);
    };

    const handlePostCommand = () => {
        if (command.trim() !== '') {
            //Hacer operacion para enviar
            setCommandSaved([...commandsSaved, command]);
            setCommand('');
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
                <input type="text" value={command} onChange={HandleCommandChange} className='subComm' />
                <button onClick={handlePostCommand} className='subEnv'>Enviar</button>

            </div>
        </>
    );
}

export default Pantalla1;