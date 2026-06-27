# Splay Tree - Caché en Memoria (Go + SQLite)

Proyecto final del curso **Algoritmos y Estructura de Datos** (2026-1), implementando un Splay Tree desde cero en Go, con pruebas de rendimiento y un caso de uso real conectado a una base de datos.

## Descripción General
El Splay Tree (Árbol Biselado) es un árbol binario de búsqueda autoajustable. Basándose en el principio de localidad de referencia, las claves consultadas más recientemente se mueven a la raíz del árbol mediante operaciones de rotación (zig, zig-zig, zig-zag). 

En este proyecto utilizamos la estructura para simular un **Caché en Memoria para un Catálogo de Retail**, interceptando las consultas a una base de datos SQLite y acelerando el acceso repetitivo a los productos más vendidos (bestsellers).

## Instrucciones de Ejecución

### Requisitos Previos
* Go 1.20 o superior.
* Archivo de datos `online_retail_clean.csv` en la raíz del proyecto.
* Compilador de C en tu sistema (necesario para el driver `go-sqlite3`).

### 1. Inicializar el entorno
Abre tu terminal en la carpeta del proyecto y ejecuta:
```bash
# Inicializar el módulo de Go
go mod init tu_modulo

# Instalar el driver de SQLite
go get [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
2. Ejecutar la Simulación (Entregable 3)
El script principal verificará si la base de datos retail.db existe. Si no, leerá los datos del archivo CSV, los cargará, y finalmente ejecutará una simulación mostrando la drástica diferencia de tiempos entre buscar en la Base de Datos vs. buscar en nuestro Splay Tree.
