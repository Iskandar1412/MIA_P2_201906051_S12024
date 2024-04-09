import React, { useState } from 'react';
import frame from '../../img/frame1.jpg'
import axios from 'axios'

function Login( props ) {
    const [password, setPassword] = useState("")
    const [user, setUser] = useState("");
    const path = "http://localhost:8080"

    const enviarDirectorios = async () => {
        //se usará un fetch para obtener los directorios y archivos
        try {
            const res = await axios.get(path+'/obtain-carpetas-archivos')
            if (res.status === 200) {
                const jsonData = JSON.parse(res.data.datos);
                // console.log(jsonData)
                // console.log(jsonData)
                props.cambiarDirectorios(jsonData)
            }
        } catch (e) { }
        
    }

    const ingresarSecion = async (objeto) => {
        try { 
            console.log(objeto)
            const res = await axios.post(path + "/login", objeto)
            if (res.status === 200) {
                enviarDirectorios()
                setUser("")
                setPassword("")

                localStorage.setItem('user', user)
                alert("Bienvenido")
                props.onSeleccionar('dashboard')
            } else {
                alert("Usuario o contraseña incorrecta")
            }
        } catch (e) { 
            alert("Usuario o contraseña incorrectos")
        }
    }

    const handleLogin =  (event) => {
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
        let objeto = {
            comando: ""
        }
        var id = localStorage.getItem("id_particion")
        objeto.comando = "login -user=" + user + " -pass=" + password + " -id=" + id
        //Verificar inicio seción
        // login -user -password -id_particion
        ingresarSecion(objeto)
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