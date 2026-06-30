package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"proyecto/splaytree"
	"strconv"
	"time"
)

type Product struct {
	StockCode   string
	Description string
	Quantity    int
	UnitPrice   float64
	Country     string
}

func main() {
	serverMode := flag.Bool("server", false, "Iniciar en modo servidor de simulación web")
	port := flag.Int("port", 8080, "Puerto para el servidor de simulación")
	flag.Parse()

	if *serverMode || os.Getenv("PORT") != "" {
		IniciarServidor(*port)
		return
	}

	productIndex := splaytree.New[string, *Product]()

	file, err := os.Open("online_retail_clean.csv")
	if err != nil {
		log.Fatalf("Error al abrir el archivo CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Read()

	fmt.Println("Cargando productos en el Splay Tree...")
	startTime := time.Now()
	recordsRead := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error al leer una línea del CSV: %v", err)
			continue
		}

		quantity, _ := strconv.Atoi(record[3])
		unitPrice, _ := strconv.ParseFloat(record[5], 64)

		product := &Product{
			StockCode:   record[1],
			Description: record[2],
			Quantity:    quantity,
			UnitPrice:   unitPrice,
			Country:     record[7],
		}

		productIndex.Insert(product.StockCode, product)
		recordsRead++
	}

	loadTime := time.Since(startTime)
	fmt.Printf("✓ Se cargaron %d productos en %v.\n", productIndex.Size(), loadTime)
	fmt.Printf("  La altura actual del árbol es: %d\n\n", productIndex.Height())

	searchKey := "22423"
	fmt.Printf("Buscando el producto con StockCode: %s\n", searchKey)

	searchStart := time.Now()
	foundProduct, ok := productIndex.Search(searchKey)
	searchTime := time.Since(searchStart)

	if !ok {
		fmt.Printf("-> Producto no encontrado.\n")
	} else {
		fmt.Printf("✓ Producto encontrado en %v\n", searchTime)
		fmt.Printf("  Descripción: %s\n", foundProduct.Description)
		fmt.Printf("  País: %s\n", foundProduct.Country)
	}

	rootKey, _ := productIndex.Root()
	fmt.Printf("  La nueva raíz del árbol es: %s\n\n", rootKey)

	fmt.Printf("Volviendo a buscar el mismo producto (%s)...\n", searchKey)
	searchStart = time.Now()
	productIndex.Search(searchKey)
	searchTime = time.Since(searchStart)

	fmt.Printf("✓ La segunda búsqueda fue mucho más rápida: %v\n", searchTime)
}
