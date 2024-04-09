import { useState, useEffect } from "react";
import filedot from '../img/dot.png'
import filetxt from '../img/txt.png'
import * as d3Graphviz from "d3-graphviz";
import axios from "axios";
import { pathbackend } from "../path";

function Pantalla3({ dots, cambiarDot }) {
    const [grafica, setGrafica] = useState("");

    const obtainGraphs = (graph) => {
        //usar el fetch para obtener información del archivo
        // var gr = "digraph G {\nsubgraph cluster_0 {\nstyle=filled;\ncolor=lightgrey;\nnode [style=filled,color=white];\na0 -> a1 -> a2 -> a3;\nlabel = \"process #1\";\n}\nsubgraph cluster_1 {\nnode [style=filled];\nb0 -> b1 -> b2 -> b3;\nlabel = \"process #2\";\ncolor=blue\n}\nstart -> a0;\nstart -> b0;\na1 -> b3;\nb2 -> a3;\na3 -> a0;\na3 -> end;\nb3 -> end;\nstart [shape=Mdiamond];\nend [shape=Msquare];\n}"
        setGrafica(graph)
    }


    useEffect(() => {
        if (grafica) {
            d3Graphviz.graphviz('#graph-graphviz').width('100%').height('100%').renderDot(grafica);
        }
    }, [grafica]);

    

    useEffect(() => {
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
        handleReports()
        return () => {
            // Aquí puedes limpiar cualquier efecto secundario si es necesario
        };
    },  [])

    const handleRepButtonDot = async (event) => {
        const valorBoton = event.target.getAttribute('data-value');
        //console.log("Valor del botón:", valorBoton);
        try {
            const res = await axios.get(pathbackend + '/graphs', {
                params: {
                    id: valorBoton,
                }
            });
            // console.log(res.data.datos)
            obtainGraphs(res.data.datos)
        } catch(e) { }
        // A este se le hace verificación de informacion
    }

    const handleRepButtonTxt = async (event) => {
        const valorBoton = event.target.getAttribute('data-value');
        //console.log("Valor del botón:", valorBoton);
        try {
            const res = await axios.get(pathbackend + '/graphs', {
                params: {
                    id: valorBoton
                }
            });
            // console.log(res.data.datos)
            // TxtGraph(res.data.datos)
            var temp = 'digraph G {\n\tnode[shape=none, lblstyle="align=left"];'
            temp += '\n\tstart[label="' + res.data.datos + '"];\n'
            temp += '}'
            // console.log(temp)
            setGrafica(temp)
        } catch (e) { } 
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
                        item.extension === "dot" ? (

                            <button 
                                className="buttonDisk"
                                data-value={item.dot}
                                key={item.dot}
                                onClick={handleRepButtonDot}
                            >
                                <img
                                    src={filedot}
                                    alt="Imagen del botón" 
                                    data-value={item.dot}
                                />
                                <span className='valor-button' data-value={item.dot}>{item.dot}</span>
                            </button>
                        ) : (
                            <button 
                                className="buttonDisk"
                                data-value={item.dot}
                                key={item.dot}
                                onClick={handleRepButtonTxt}
                            >
                                <img
                                    src={filetxt}
                                    alt="Imagen del botón" 
                                    data-value={item.dot}
                                />
                                <span className='valor-button' data-value={item.dot}>{item.dot}</span>
                            </button>
                        )
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