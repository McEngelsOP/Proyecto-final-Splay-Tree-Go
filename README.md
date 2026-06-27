Markdown
# 🌳 Splay Tree - Caché en Memoria (Go + SQLite)

Proyecto final del curso **Algoritmos y Estructura de Datos** (2026-1), implementando un Splay Tree desde cero en Go, con pruebas de rendimiento y un caso de uso real conectado a una base de datos.

## 📌 Descripción General
El Splay Tree (Árbol Biselado) es un árbol binario de búsqueda autoajustable. Basándose en el principio de localidad de referencia, las claves consultadas más recientemente se mueven a la raíz del árbol mediante operaciones de rotación (zig, zig-zig, zig-zag). 

En este proyecto utilizamos la estructura para simular un **Caché en Memoria para un Catálogo de Retail**, interceptando las consultas a una base de datos SQLite y acelerando el acceso repetitivo a los productos más vendidos (bestsellers).

## ⚙️ Instrucciones de Ejecución

### Requisitos Previos
* Go 1.20 o superior.
* Archivo de datos `online_retail_clean.csv` en la raíz del proyecto.
* Compilador de C en tu sistema (necesario para el driver `go-sqlite3`).

### 1. Inicializar el entorno
Abre tu terminal en la carpeta del proyecto y ejecuta:
*
```bash
# Inicializar el módulo de Go
go mod init tu_modulo

# Instalar el driver de SQLite
go get [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
*
Análisis de Complejidad (Big-O)El beneficio principal del Splay Tree no está en su peor caso individual, sino en su complejidad amortizada, ideal para sistemas donde unas pocas claves agrupan la mayoría de las consultas.Espacio en Memoria: $\mathcal{O}(n)$ — Se reserva un nodo estructural por cada clave almacenada en la RAM.Búsqueda / Inserción / Eliminación (Peor caso individual): $\mathcal{O}(n)$ — En el improbable caso de que el árbol se degenere en una lista enlazada pura.Complejidad Amortizada (Promedio): $\mathcal{O}(\log n)$ — Gracias al mecanismo de "splaying" (biselado), al realizar una secuencia de $m$ operaciones, el tiempo por operación se balancea logarítmicamente.Localidad de Referencia: Si consultamos un mismo elemento (ej. el producto estrella de la tienda) múltiples veces seguidas, el tiempo de búsqueda tiende a $\mathcal{O}(1)$ porque el nodo reposa en la raíz.
