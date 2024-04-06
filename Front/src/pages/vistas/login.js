import React, { useState } from 'react';
import frame from '../../img/frame1.jpg'

function Login(props) {
    const [password, setPassword] = useState("")
    const [user, setUser] = useState("");
    
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

        console.log(user, password)
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