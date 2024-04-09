import React, { useState } from 'react';

import './App.css';
import Pantalla1 from './pages/pantalla1';
import Pantalla2 from './pages/pantalla2';
import Pantalla3 from './pages/pantalla3';

function App() {

  const [selectedTab, setSelectedTab] = useState('contente1');
  const [command, setCommand] = useState([]);
  
  const handleTabChange = (event) => {
    setSelectedTab(event.target.id.replace('tab', 'content')); 
  };

  const obtainInfo = (ex) => {
    // console.log(ex)
    
    setCommand(ex)
    // console.log(command)
    //console.log("comando", command)
  }

  const [carpetas, setCarpetas] = useState([]);
  const obtainCarpetas = (carp) => {
    
    setCarpetas(carp)
  } 

  const [dots, setDots] = useState([]);
  const obtainDots = (dot) => {
    // var gd = [
    //   { dot: "grahicacas.dot", extension: "dot" },
    //   { dot: "grahicaca2.txt", extension: "txt" },
    // ]
    setDots(dot)
  }

  

  return (
    <div className='usuario-data3'>
      <div className='group-usuario-nombre'>
        <main className="container-x">
          <div className='barmenu'>
            <input id="tabe1" type="radio" name="tabs-1" defaultChecked onChange={handleTabChange} />
            <label htmlFor="tabe1" className="label-type">Pantalla 1</label>
            <input id="tabe2" type="radio" name="tabs-1" onChange={handleTabChange} />
            <label htmlFor="tabe2" className="label-type">Pantalla 2</label>
            <input id="tabe3" type="radio" name="tabs-1" onChange={handleTabChange} />
            <label htmlFor="tabe3" className="label-type">Pantalla 3</label>
          </div>

          <div className='content-type'>
            <section id="contente1" className={`tabs-contentype ${selectedTab === 'contente1' && 'active'}`}>
              <Pantalla1 info={obtainInfo} carpetas={obtainCarpetas} cambiarDot={obtainDots} />
            </section>
            <section id="contente2" className={`tabs-contentype ${selectedTab === 'contente2' && 'active'}`}>
              <Pantalla2 command={command} carpetasOb={carpetas} carpetas={obtainCarpetas} />
            </section>
            <section id="contente3" className={`tabs-contentype ${selectedTab === 'contente3' && 'active'}`}>
              <Pantalla3 dots={dots} cambiarDot={obtainDots} />
            </section>
          </div>
        </main>
      </div>
    </div>

  );
}

export default App;
