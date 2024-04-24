let dir = "localhost";

const cambiarDir = () => {
  const nuevoDir = prompt("Ingrese IP Server Backend (Default: Localhost):", dir);
  if (nuevoDir !== null) {
    dir = nuevoDir;
    alert(`New IP Server Backend: --> http://${dir}:8080`);
  }
};


const inicializar = () => {
  cambiarDir(); 
};

inicializar();

export const pathbackend = "http://" + dir + ":8080";