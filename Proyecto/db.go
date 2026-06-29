package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDB inicializa la conexión, crea las tablas y realiza la migración inicial si es necesario.
// Soporta dos modos: local (XAMPP) y nube (MYSQL_DSN env var).
func InitDB() error {
	dsn := os.Getenv("MYSQL_DSN")

	if dsn == "" {
		// Modo local: XAMPP en localhost
		dsnRoot := "root:@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
		dbRoot, err := sql.Open("mysql", dsnRoot)
		if err != nil {
			return fmt.Errorf("error al conectar a MySQL: %v", err)
		}
		defer dbRoot.Close()

		_, err = dbRoot.Exec("CREATE DATABASE IF NOT EXISTS proyecto_retail")
		if err != nil {
			return fmt.Errorf("error al crear la base de datos proyecto_retail: %v", err)
		}

		dsn = "root:@tcp(127.0.0.1:3306)/proyecto_retail?charset=utf8mb4&parseTime=True&loc=Local"
		log.Println("Modo local: conectando a MySQL de XAMPP...")
	} else {
		log.Println("Modo nube: conectando a MySQL remoto...")
	}

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error al abrir base de datos: %v", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("error al hacer ping a la base de datos: %v", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS productos (
		stock_code VARCHAR(50) PRIMARY KEY,
		description VARCHAR(255),
		quantity INT,
		unit_price DECIMAL(10, 2),
		country VARCHAR(100)
	);`
	_, err = DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error al crear la tabla productos: %v", err)
	}

	log.Println("✓ Conexión establecida. Tabla 'productos' lista.")

	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM productos").Scan(&count)
	if err != nil {
		return fmt.Errorf("error al verificar tamaño de tabla productos: %v", err)
	}

	if count == 0 {
		log.Println("La tabla 'productos' está vacía. Iniciando migración de registros desde el CSV...")
		if err := migrateCSVToSQL(); err != nil {
			log.Printf("Advertencia en migración: %v (esto es normal en modo nube sin CSV)\n", err)
		}
	} else {
		log.Printf("✓ Base de datos inicializada con %d registros.\n", count)
	}

	return nil
}

// migrateCSVToSQL lee una porción del CSV e inserta los datos en la base de datos MySQL
func migrateCSVToSQL() error {
	file, err := os.Open("online_retail_clean.csv")
	if err != nil {
		return fmt.Errorf("no se pudo abrir el archivo CSV para migración: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Leer cabecera
	if _, err := reader.Read(); err != nil {
		return err
	}

	// Usar una transacción para insertar rápido
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT IGNORE INTO productos (stock_code, description, quantity, unit_price, country) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("error al preparar declaración de inserción: %v", err)
	}
	defer stmt.Close()

	count := 0
	const limit = 500 // Migrar 500 registros iniciales para tener una muestra robusta y rápida

	for {
		record, err := reader.Read()
		if err == io.EOF || count >= limit {
			break
		}
		if err != nil {
			continue
		}

		stockCode := record[1]
		description := record[2]
		quantity := record[3]
		unitPrice := record[5]
		country := record[7]

		// Insertar en lote
		_, err = stmt.Exec(stockCode, description, quantity, unitPrice, country)
		if err == nil {
			count++
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error al comprometer transacción de migración: %v", err)
	}

	log.Printf("✓ Migración finalizada: %d productos insertados en MySQL.\n", count)
	return nil
}

// DBGetProductos obtiene todos los productos de la base de datos MySQL
func DBGetProductos() (map[string]string, error) {
	rows, err := DB.Query("SELECT stock_code, description FROM productos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var code, desc string
		if err := rows.Scan(&code, &desc); err != nil {
			return nil, err
		}
		result[code] = desc
	}
	return result, nil
}

// DBInsertProducto guarda un producto en la base de datos MySQL
func DBInsertProducto(code, desc string) error {
	query := "INSERT INTO productos (stock_code, description, quantity, unit_price, country) VALUES (?, ?, 1, 1.0, 'Simulado') ON DUPLICATE KEY UPDATE description=?"
	_, err := DB.Exec(query, code, desc, desc)
	return err
}

// DBDeleteProducto elimina un producto de la base de datos MySQL
func DBDeleteProducto(code string) error {
	_, err := DB.Exec("DELETE FROM productos WHERE stock_code=?", code)
	return err
}

// DBClearProductos vacía la tabla productos
func DBClearProductos() error {
	_, err := DB.Exec("TRUNCATE TABLE productos")
	return err
}
