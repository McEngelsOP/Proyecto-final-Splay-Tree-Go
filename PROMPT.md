# Prompt para Gamma App - Generación de Presentación (PPT)

Copia y pega el siguiente prompt en **Gamma App** (en la sección "Generar con IA" -> "Texto a Presentación") para crear de forma automática tus diapositivas educativas e interactivas sobre el proyecto.

---

```text
Actúa como un profesor universitario y experto en algoritmos y estructuras de datos. Genera una presentación de diapositivas clara, visual y altamente educativa sobre el proyecto de implementación y simulación de un Splay Tree. 

Usa un tono profesional, tecnológico y académico. Organiza la información en viñetas cortas, evita párrafos extensos y utiliza analogías visuales o bloques destacados.

Estructura de la Presentación (8 Diapositivas):

Diapositiva 1: Portada de Proyecto
- Título: Simulación y Aplicación del Splay Tree (Árbol Biselado)
- Subtítulo: Optimización de Búsqueda mediante la Localidad de Referencia de Datos
- Elementos: Nombre del curso (Algoritmos y Estructuras de Datos), espacio para nombres de integrantes y profesor.

Diapositiva 2: El Problema de Búsqueda Convencional
- Título: El Desafío del Acceso a Datos y el Desbalanceo
- Contenido:
  * Explicar el problema de los Árboles de Búsqueda Binarios (BST) convencionales que degeneran en listas enlazadas (O(n)).
  * Limitaciones de los enfoques rígidos (AVL/Red-Black): alto costo de balanceo estricto sin considerar la frecuencia de consulta.

Diapositiva 3: Localidad de Referencia (Locality)
- Título: ¿Cómo Consultan los Usuarios Realmente?
- Contenido:
  * El Principio de Pareto (Regla del 80/20) en sistemas reales: el 80% de los accesos ocurren sobre el 20% de los elementos (los más populares).
  * La localidad temporal: si un dato se acaba de consultar, es muy probable que se consulte de nuevo muy pronto.

Diapositiva 4: La Solución: ¿Qué es un Splay Tree?
- Título: Splay Tree: Autobalanceo Adaptativo por Acceso
- Contenido:
  * Definición: Árbol binario de búsqueda auto-ajustable inventado por Sleator y Tarjan (1985).
  * Concepto clave: Cada vez que un nodo es accedido (búsqueda, inserción o eliminación), este se traslada físicamente a la raíz del árbol mediante la operación "Splay".
  * Rendimiento amortizado: O(log n) en el peor caso, pero O(1) inmediato para consultas repetidas.

Diapositiva 5: El Mecanismo Interno: Rotaciones
- Título: ¿Cómo Funciona la Operación Splay?
- Contenido:
  * Los tres tipos de rotaciones según la posición del nodo, su padre y su abuelo:
    1. Zig (Rotación simple): Cuando el padre es la raíz.
    2. Zig-Zig (Rotaciones consecutivas en la misma dirección): El nodo y su padre son ambos hijos izquierdos (o derechos).
    3. Zig-Zag (Rotación doble cruzada): El nodo es hijo izquierdo y su padre es hijo derecho (o viceversa).
  * Ilustrar que estas rotaciones reestructuran el árbol dinámicamente.

Diapositiva 6: Aplicación con Datos Reales (Caso de Estudio)
- Título: Caso Práctico: Indexación de Productos
- Contenido:
  * Carga masiva de transacciones comerciales desde un dataset real (online_retail_clean.csv) en Go.
  * Demostración empírica de velocidad:
    - 1ra búsqueda (recorrido profundo): Tarda unos microsegundos.
    - 2da búsqueda (mismo elemento, ya en la raíz): Tarda prácticamente cero microsegundos (O(1)).

Diapositiva 7: Arquitectura de la Simulación Visual (Go + Vue.js)
- Título: Arquitectura del Simulador Interactivo
- Contenido:
  * Backend en Go: Procesa las operaciones, instrumenta el árbol mediante callbacks para capturar "snapshots" (fotos del árbol en cada rotación) y expone una API REST JSON.
  * Frontend en Vue.js + SVG: Renderiza la topología del árbol en tiempo real.
  * Características: Reproductor multimedia paso a paso para observar el movimiento físico de las rotaciones durante el splay.

Diapositiva 8: Conclusiones y Demostración
- Título: Conclusiones del Estudio
- Contenido:
  * El Splay Tree optimiza las búsquedas basándose en el comportamiento del usuario en vez del orden estricto de los datos.
  * La simulación paso a paso facilita la comprensión visual y didáctica de las complejas rotaciones del árbol.
  * Espacio para preguntas y demostración en vivo.
```
