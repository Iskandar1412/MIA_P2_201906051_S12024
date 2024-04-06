import { useState, useEffect } from "react";
import filedot from '../img/dot.png'
import * as d3Graphviz from "d3-graphviz";

function Pantalla3({ dots }) {
    const [grafica, setGrafica] = useState("");
    const obtainGraphs = (graph) => {
        //usar el fetch para obtener información del archivo
        var gr = "digraph G {\nsubgraph cluster_0 {\nstyle=filled;\ncolor=lightgrey;\nnode [style=filled,color=white];\na0 -> a1 -> a2 -> a3;\nlabel = \"process #1\";\n}\nsubgraph cluster_1 {\nnode [style=filled];\nb0 -> b1 -> b2 -> b3;\nlabel = \"process #2\";\ncolor=blue\n}\nstart -> a0;\nstart -> b0;\na1 -> b3;\nb2 -> a3;\na3 -> a0;\na3 -> end;\nb3 -> end;\nstart [shape=Mdiamond];\nend [shape=Msquare];\n}"
        setGrafica(gr)
    }

    useEffect(() => {
        if (grafica) {
            d3Graphviz.graphviz('#graph-graphviz').width('100%').height('100%').renderDot(grafica);
        }
    }, [grafica]);

    const handleRepButton = (event) => {
        const valorBoton = event.target.getAttribute('data-value');
        //console.log("Valor del botón:", valorBoton);
        obtainGraphs(valorBoton)
        // A este se le hace verificación de informacion
        
    }

    const handleCerrarContenido = () => {
        setGrafica("");
    };

    return (
        <>
            <div className="vistas">
                <div className="discos">
                    
                    {dots.map((item, index) => (
                        <button 
                            className="buttonDisk"
                            data-value={item.dot}
                            key={item.dot}
                            onClick={handleRepButton}
                        >
                            <img
                                src={filedot}
                                alt="Imagen del botón" 
                                data-value={item.dot}
                            />
                            <span className='valor-button' data-value={item.dot}>{item.dot}</span>
                        </button>
                    ))}
                    
                    {grafica && (
                        <div className="modal2" style={{ backgroundColor: 'rgba(255, 255, 255, 0.9)' }}>
                            <div className="boton-exit">
                                <button onClick={handleCerrarContenido} className="button-x" />
                            </div>
                            <div className="graphica" id="graph-graphviz" />
                        </div>
                    )}


                </div>
            </div>
        </>
    );

}

export default Pantalla3;