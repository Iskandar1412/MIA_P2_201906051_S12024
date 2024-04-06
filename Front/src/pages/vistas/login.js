import React, { useState } from 'react';
import frame from '../../img/frame1.jpg'

function Login( props ) {
    const [password, setPassword] = useState("")
    const [user, setUser] = useState("");
    

    const enviarDirectorios = () => {
        //se usará un fetch para obtener los directorios y archivos
        var dir = [
            { nombre: "archivo.txt", tipo: "archivo", contenido: "Contenido del archivo 1" },
            { nombre: "archivo2", tipo: "archivo", contenido: "Contenido del archivo 1" },
                { nombre: "carpeta1", tipo: "carpeta", contenido: [
                    { nombre: "archivo2", tipo: "archivo", contenido: "Contenido del archivo 2" },
                    { nombre: "carpeta2", tipo: "carpeta", contenido: [
                        { nombre: "archivo3", tipo: "archivo", contenido: "Contenido del archivo 3" }
                    ]}
                ]
            }
        ]
        /*
            */
        props.cambiarDirectorios(dir)
    }

    const handleLogin = (event) => {
        event.preventDefault();
        if (user === "" && password === "") {
            alert("Campos vacios");
            return
        }
        if (user === "" ) { 
            alert("Casilla de Usuario vacio")
            return
        } else if (password === "") {
            alert("Casilla de Password vacio")
            return
        }
        
        //Verificar inicio seción
        //console.log(user, password)
        enviarDirectorios()
        setUser("")
        setPassword("")
        localStorage.setItem('user', user)
        props.onSeleccionar('dashboard')
    }

    return (
        <>  
            <div 
                className="login"
                style={{ 
                    background: `url(${frame})`,
                    backgroundSize: 'cover',
                    backgroundPosition: 'center',
                    backgroundRepeat: 'no-repeat',
                }}
            >
                <div className='login-form' style={{ backgroundColor: 'rgba(0, 0, 0, 0.6)' }} >
                    <button className='button-ant2' onClick={() => props.onSeleccionar('particiones')} />
                    <div className='usuario-logo' />
                    <div className='form-users'>
                        <div className='div-form'>
                            <input type='text' className='usuario-form' placeholder="Usuario" value={user} onChange={(e) => setUser(e.target.value)} />
                            <input type='password' className='password-form' placeholder="Contraseña" value={password} onChange={(e) => setPassword(e.target.value)} />

                        </div>
                        <button className='login-button' onClick={handleLogin}>Login</button>
                    </div>
                </div>

            </div>
            
        </>
    );

}

export default Login;