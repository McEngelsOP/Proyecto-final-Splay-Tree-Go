# Explicación del Código - Simulador de Splay Tree (Árbol Biselado)

Este documento explica de forma detallada la arquitectura, estructura y funcionamiento del simulador de **Splay Tree** implementado en Go y Vue.js.

---

## 1. El Problema que Resuelve este Proyecto

Este proyecto aborda un reto clásico de optimización en ciencias de la computación: **cómo organizar y acceder de la forma más rápida posible a un conjunto masivo de información en memoria, aprovechando los patrones de comportamiento y consulta de los usuarios**.

### A. El Desafío del Balanceo vs. Costo
En los árboles binarios de búsqueda tradicionales (BST), la inserción de claves en orden sucesivo (ej. 1, 2, 3...) genera una estructura totalmente desbalanceada (similar a una lista enlazada), donde el tiempo de búsqueda se degrada a un ineficiente $O(n)$. Para solucionar esto, estructuras como los árboles AVL o Red-Black mantienen el árbol estrictamente balanceado mediante rotaciones en cada inserción/eliminación. Sin embargo, estas operaciones imponen un costo de balanceo estricto e ignoran por completo **cuáles** elementos son los más accedidos.

### B. La Localidad de Referencia (Locality)
En sistemas reales (como tiendas en línea o almacenamiento en caché), **los accesos a los datos son altamente sesgados**. Rara vez se consulta cada elemento con la misma frecuencia; en su lugar, se cumple el **Principio de Pareto (regla del 80/20)**: el 80% de las consultas se dirigen a solo el 20% de los elementos (los productos populares, ofertas del día o consultas recientes).

### C. La Propuesta del Splay Tree
El **Splay Tree** aprovecha este comportamiento moviendo dinámicamente cualquier nodo accedido directamente a la raíz. Si una clave (por ejemplo, el código de producto `"22423"`) es muy popular y se consulta constantemente:
1. En la primera consulta, se recorre el árbol (tiempo de búsqueda promedio $O(\log n)$).
2. Se realiza la operación de **splay** (rotaciones) y la clave `"22423"` se convierte en la nueva raíz del árbol.
3. En las consultas consecutivas a esa misma clave, el tiempo de búsqueda cae a **$O(1)$** (tiempo constante), ya que está ubicada justo en la cima del árbol.

---

## 2. ¿Qué es un Splay Tree y cómo funciona?

Un **Splay Tree** (Árbol Biselado) es un árbol de búsqueda binario auto-balanceable. Su característica principal es que **cada vez que se accede a un nodo (al buscar, insertar o eliminar), este nodo se mueve automáticamente a la raíz** mediante una secuencia de rotaciones específicas llamadas **Splay**.

### Ventaja Principal:
Al subir los elementos más accedidos a la raíz:
* Las consultas repetidas sobre las mismas claves se realizan de forma inmediata.
* La estructura posee un rendimiento amortizado excelente de $O(\log n)$ por operación.
* Es ideal para sistemas de caché, índices de bases de datos o enrutamiento de red.

---

## 3. Estructura del Código

El proyecto está organizado en las siguientes partes:

```text
Proyecto/
├── splaytree/
│   ├── splay_tree.go       # Estructura del árbol, rotaciones e inserción/búsqueda/eliminación.
│   └── splay_tree_test.go  # Pruebas unitarias y pruebas de localidad (locality).
├── web/
│   ├── index.html          # Interfaz de usuario (SPA en Vue.js 3 y Tailwind CSS).
│   └── style.css           # Estilos personalizados para animaciones y transiciones de nodos.
├── go.mod                  # Declaración del módulo Go.
├── online_retail_clean.csv # Dataset con transacciones reales de productos.
├── server.go               # Servidor API HTTP que gestiona el árbol y genera snapshots.
└── package main.go         # Punto de entrada de la aplicación (Modo Consola / Modo Web).
```

---

## 4. Funcionamiento del Backend (Go)

### A. Estructura de Datos e Instrumentación ([splay_tree.go](file:///c:/Users/julor/Downloads/Proyecto/Proyecto/splaytree/splay_tree.go))

El nodo almacena la clave (`Key`), el valor (`Value`) y punteros al padre y a sus hijos izquierdo y derecho:

```go
type Node[K cmp.Ordered, V any] struct {
    Key    K
    Value  V
    left   *Node[K, V]
    right  *Node[K, V]
    parent *Node[K, V]
}
```

La estructura `SplayTree` cuenta con un campo especial llamado `OnSplayStep`. Este callback se ejecuta en cada rotación intermedia dentro de la función de balanceo `splay`:

```go
type SplayTree[K cmp.Ordered, V any] struct {
    root        *Node[K, V]
    size        int
    OnSplayStep func(xKey K, stepDescription string) // Callback para simulación paso a paso
}
```

#### El Algoritmo de Splay (Rotaciones)
La función `splay(x)` sube recursivamente el nodo `x` hasta la raíz evaluando tres casos:
1. **Zig (Rotación Simple)**: El padre del nodo es la raíz. Se hace una rotación derecha (si es hijo izquierdo) o izquierda (si es hijo derecho).
2. **Zig-Zig**: El nodo y su padre son ambos hijos izquierdos (o ambos derechos). Se rota primero sobre el abuelo y luego sobre el padre.
3. **Zig-Zag**: El nodo es hijo izquierdo y su padre es hijo derecho (o viceversa). Se rota sobre el padre y luego sobre el abuelo.

Para la simulación, capturamos qué rotación se ejecuta llamando a `OnSplayStep`:

```go
if g == nil {
    if x == p.left {
        t.rotateRight(p)
        stepDesc = fmt.Sprintf("Zig Derecho sobre padre %v", p.Key)
    } else { ... }
} else if x == p.left && p == g.left {
    t.rotateRight(g)
    t.rotateRight(p)
    stepDesc = fmt.Sprintf("Zig-Zig Derecho sobre abuelo %v y padre %v", g.Key, p.Key)
} else { ... }

if t.OnSplayStep != nil {
    t.OnSplayStep(x.Key, stepDesc)
}
```

### B. Servidor API y Generación de Snapshots ([server.go](file:///c:/Users/julor/Downloads/Proyecto/Proyecto/server.go))

El archivo `server.go` expone endpoints REST en formato JSON. Su tarea crucial es **capturar el estado intermedio del árbol durante el splay**:

Al realizar una acción como insertar:
1. Definimos una lista de pasos `[]SplayStep`.
2. Asignamos `OnSplayStep` para que serialice y guarde una copia completa del árbol (un snapshot) en cada rotación.
3. Ejecutamos la operación en el árbol.
4. Desactivamos el callback y devolvemos al frontend la secuencia completa de pasos.

```go
// Definimos el callback para registrar pasos
tree.OnSplayStep = func(xKey string, desc string) {
    steps = append(steps, SplayStep{
        ActiveKey:   xKey,
        Description: desc,
        Tree:        serializeNode(tree.RootNode()), // Copia serializada de la estructura actual
    })
}
```

---

## 5. Funcionamiento del Frontend (Vue.js + SVG)

La visualización interactiva se encuentra en [web/index.html](file:///c:/Users/julor/Downloads/Proyecto/Proyecto/web/index.html).

### A. Renderizado del Árbol Binario
El árbol se dibuja en un elemento `<svg>` de forma dinámica:
1. **Posiciones (x, y)**: Se calculan de forma recursiva en la propiedad computada `calculatedNodes`.
   - La coordenada `y` es proporcional al nivel de profundidad del nodo (`depth * 90 + 70`).
   - La coordenada `x` se calcula dividiendo el espacio disponible a la mitad (`(minX + maxX) / 2`) para el nodo actual, y delegando los subrangos a sus respectivos hijos izquierdo y derecho.
2. **Enlaces (Aristas)**: Se calculan en `calculatedEdges` y se dibujan como curvas Bezier cúbicas suaves:
   ```javascript
   const drawEdge = (x1, y1, x2, y2) => {
       const midY = (y1 + y2) / 2;
       return `M ${x1} ${y1} C ${x1} ${midY}, ${x2} ${midY}, ${x2} ${y2}`;
   };
   ```

### B. Transiciones Fluidas (Efecto Rotación)
En [web/style.css](file:///c:/Users/julor/Downloads/Proyecto/Proyecto/web/style.css) aplicamos transiciones CSS sobre las figuras de los nodos (`circle`) y los enlaces (`path`):

```css
.node-group {
    transition: transform 0.6s cubic-bezier(0.25, 1, 0.5, 1);
}
.edge-line {
    transition: d 0.6s cubic-bezier(0.25, 1, 0.5, 1), stroke 0.4s ease;
}
```

Cuando avanzamos de un paso de simulación a otro, Vue actualiza el JSON del árbol. Al cambiar las coordenadas de los elementos SVG, el motor del navegador realiza una **interpolación animada** del movimiento. Esto crea una ilusión fluida y premium de balanceo tridimensional de los nodos durante el splay.

### C. Reproductor Paso a Paso
El panel de simulación en Vue.js recibe la lista de pasos del backend. Puedes:
* Usar los botones multimedia para controlar un temporizador que avanza de paso en paso (`setInterval`).
* Arrastrar la barra deslizadora (slider) para saltar a un snapshot específico.
* Visualizar en el panel de texto la descripción de la rotación actual.
