package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"proyecto/splaytree"
	"sync"
)

// SerializedNode representa un nodo para enviar al frontend
type SerializedNode struct {
	Key   string          `json:"key"`
	Value string          `json:"value"`
	Left  *SerializedNode `json:"left"`
	Right *SerializedNode `json:"right"`
}

// SplayStep representa un snapshot del estado del árbol durante una rotación
type SplayStep struct {
	ActiveKey   string          `json:"activeKey"`
	Description string          `json:"description"`
	Tree        *SerializedNode `json:"tree"`
}

// Global state
var (
	tree      = splaytree.New[string, string]()
	treeMutex sync.Mutex
)

// Helper para serializar recursivamente
func serializeNode(node *splaytree.Node[string, string]) *SerializedNode {
	if node == nil {
		return nil
	}
	return &SerializedNode{
		Key:   node.Key,
		Value: node.Value,
		Left:  serializeNode(node.Left()),
		Right: serializeNode(node.Right()),
	}
}

// Response genérica
type TreeResponse struct {
	Size   int             `json:"size"`
	Height int             `json:"height"`
	Tree   *SerializedNode `json:"tree"`
}

func getTreeState() TreeResponse {
	return TreeResponse{
		Size:   tree.Size(),
		Height: tree.Height(),
		Tree:   serializeNode(tree.RootNode()),
	}
}

// Middleware para habilitar CORS y permitir despliegue híbrido en Vercel
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func handleGetTree(w http.ResponseWriter, r *http.Request) {
	treeMutex.Lock()
	defer treeMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getTreeState())
}

func handleInsert(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Key == "" {
		http.Error(w, "La clave no puede estar vacía", http.StatusBadRequest)
		return
	}

	treeMutex.Lock()
	defer treeMutex.Unlock()

	var steps []SplayStep

	// Capturar el estado inicial
	steps = append(steps, SplayStep{
		ActiveKey:   req.Key,
		Description: fmt.Sprintf("Inicio de inserción para la clave %s", req.Key),
		Tree:        serializeNode(tree.RootNode()),
	})

	// Configurar callback de splay
	tree.OnSplayStep = func(xKey string, desc string) {
		steps = append(steps, SplayStep{
			ActiveKey:   xKey,
			Description: desc,
			Tree:        serializeNode(tree.RootNode()),
		})
	}

	// Ejecutar inserción en memoria
	tree.Insert(req.Key, req.Value)

	// Limpiar callback
	tree.OnSplayStep = nil

	// Sincronizar con la base de datos MySQL de XAMPP
	if err := DBInsertProducto(req.Key, req.Value); err != nil {
		log.Printf("Advertencia: No se pudo insertar la clave %s en MySQL: %v\n", req.Key, err)
	} else {
		log.Printf("✓ Clave %s sincronizada/guardada en la base de datos MySQL.\n", req.Key)
	}

	// Agregar paso final
	steps = append(steps, SplayStep{
		ActiveKey:   req.Key,
		Description: fmt.Sprintf("Inserción completada. Clave %s ahora es la raíz y se guardó en MySQL.", req.Key),
		Tree:        serializeNode(tree.RootNode()),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"steps":   steps,
		"final":   getTreeState(),
	})
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key string `json:"key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	treeMutex.Lock()
	defer treeMutex.Unlock()

	var steps []SplayStep

	// Capturar estado inicial
	steps = append(steps, SplayStep{
		ActiveKey:   req.Key,
		Description: fmt.Sprintf("Buscando la clave %s en el árbol", req.Key),
		Tree:        serializeNode(tree.RootNode()),
	})

	// Configurar callback
	tree.OnSplayStep = func(xKey string, desc string) {
		steps = append(steps, SplayStep{
			ActiveKey:   xKey,
			Description: desc,
			Tree:        serializeNode(tree.RootNode()),
		})
	}

	val, found := tree.Search(req.Key)

	tree.OnSplayStep = nil

	var desc string
	if found {
		desc = fmt.Sprintf("Búsqueda completada. Clave %s encontrada con descripción: '%s'. Nodo splayado a la raíz.", req.Key, val)
	} else {
		desc = fmt.Sprintf("Búsqueda completada. Clave %s NO encontrada en el árbol. El último nodo visitado fue splayado a la raíz.", req.Key)
	}

	steps = append(steps, SplayStep{
		ActiveKey:   req.Key,
		Description: desc,
		Tree:        serializeNode(tree.RootNode()),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"found":   found,
		"value":   val,
		"steps":   steps,
		"final":   getTreeState(),
	})
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key string `json:"key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	treeMutex.Lock()
	defer treeMutex.Unlock()

	var steps []SplayStep

	steps = append(steps, SplayStep{
		ActiveKey:   req.Key,
		Description: fmt.Sprintf("Iniciando eliminación de la clave %s (primero se busca y se splaya a la raíz)", req.Key),
		Tree:        serializeNode(tree.RootNode()),
	})

	// Para capturar las rotaciones durante la búsqueda que precede a la eliminación
	tree.OnSplayStep = func(xKey string, desc string) {
		steps = append(steps, SplayStep{
			ActiveKey:   xKey,
			Description: desc,
			Tree:        serializeNode(tree.RootNode()),
		})
	}

	deleted := tree.Delete(req.Key)

	tree.OnSplayStep = nil

	var desc string
	if deleted {
		// Sincronizar eliminación con MySQL de XAMPP
		if err := DBDeleteProducto(req.Key); err != nil {
			log.Printf("Advertencia: No se pudo eliminar la clave %s de MySQL: %v\n", req.Key, err)
		} else {
			log.Printf("✓ Clave %s eliminada de la base de datos MySQL.\n", req.Key)
		}
		desc = fmt.Sprintf("Clave %s eliminada exitosamente del árbol y de MySQL. Se unieron los subárboles.", req.Key)
	} else {
		desc = fmt.Sprintf("Clave %s no encontrada en el árbol, no se realizó ninguna eliminación.", req.Key)
	}

	steps = append(steps, SplayStep{
		ActiveKey:   req.Key,
		Description: desc,
		Tree:        serializeNode(tree.RootNode()),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"deleted": deleted,
		"steps":   steps,
		"final":   getTreeState(),
	})
}

func handleClear(w http.ResponseWriter, r *http.Request) {
	treeMutex.Lock()
	defer treeMutex.Unlock()

	// Reiniciar en memoria
	tree = splaytree.New[string, string]()

	// Vaciar tabla MySQL
	if err := DBClearProductos(); err != nil {
		log.Printf("Advertencia: No se pudo vaciar la tabla productos en MySQL: %v\n", err)
	} else {
		log.Println("✓ Tabla 'productos' vaciada en MySQL.")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"final":   getTreeState(),
	})
}

func handleLoadSample(w http.ResponseWriter, r *http.Request) {
	treeMutex.Lock()
	defer treeMutex.Unlock()

	// Obtener los productos desde la base de datos MySQL (migrados previamente del CSV)
	log.Println("Cargando productos de muestra desde MySQL...")
	productos, err := DBGetProductos()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al leer productos desde MySQL: %v", err), http.StatusInternalServerError)
		return
	}

	// Si MySQL no tiene datos, retornar 0 cargados
	if len(productos) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"loaded":  0,
			"final":   getTreeState(),
		})
		return
	}

	// Limpiar el Splay Tree actual para cargar los datos frescos de la BD
	tree = splaytree.New[string, string]()

	count := 0
	for stockCode, description := range productos {
		tree.Insert(stockCode, description)
		count++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"loaded":  count,
		"final":   getTreeState(),
	})
}

// IniciarServidor arranca el servidor web en el puerto especificado
func IniciarServidor(puerto int) {
	// Inicializar la conexión a MySQL en XAMPP
	if err := InitDB(); err != nil {
		log.Fatalf("Error crítico: No se pudo inicializar la base de datos MySQL en XAMPP. Asegúrate de tener XAMPP/MySQL encendido. Detalles: %v", err)
	}

	// Rutas API envueltas con enableCORS para compatibilidad de despliegue en Vercel
	http.HandleFunc("/api/tree", enableCORS(handleGetTree))
	http.HandleFunc("/api/insert", enableCORS(handleInsert))
	http.HandleFunc("/api/search", enableCORS(handleSearch))
	http.HandleFunc("/api/delete", enableCORS(handleDelete))
	http.HandleFunc("/api/clear", enableCORS(handleClear))
	http.HandleFunc("/api/load-sample", enableCORS(handleLoadSample))

	// Servir archivos estáticos del frontend
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	log.Printf("Servidor de simulación iniciado en http://localhost:%d\n", puerto)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", puerto), nil))
}
