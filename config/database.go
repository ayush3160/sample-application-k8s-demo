package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	PostgresDB *sql.DB
	MySQLDB    *sql.DB
	MongoDB    *mongo.Client
	MongoDBCtx context.Context
)

func InitDatabases() {
	initPostgres()
	initMySQL()
	initMongoDB()
}

func initPostgres() {
	postgresHost := getEnv("POSTGRES_HOST", "localhost")
	postgresPort := getEnv("POSTGRES_PORT", "5432")
	postgresUser := getEnv("POSTGRES_USER", "postgres")
	postgresPass := getEnv("POSTGRES_PASSWORD", "postgres")
	postgresDB := getEnv("POSTGRES_DB", "ecommerce")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		postgresHost, postgresPort, postgresUser, postgresPass, postgresDB)

	var err error
	PostgresDB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	PostgresDB.SetMaxOpenConns(25)
	PostgresDB.SetMaxIdleConns(5)
	PostgresDB.SetConnMaxLifetime(5 * time.Minute)

	if err = PostgresDB.Ping(); err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}

	log.Println("Connected to PostgreSQL")
	createPostgresTables()
}

func initMySQL() {
	mysqlHost := getEnv("MYSQL_HOST", "localhost")
	mysqlPort := getEnv("MYSQL_PORT", "3306")
	mysqlUser := getEnv("MYSQL_USER", "root")
	mysqlPass := getEnv("MYSQL_PASSWORD", "root")
	mysqlDB := getEnv("MYSQL_DB", "ecommerce")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlDB)

	var err error
	MySQLDB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	MySQLDB.SetMaxOpenConns(25)
	MySQLDB.SetMaxIdleConns(5)
	MySQLDB.SetConnMaxLifetime(5 * time.Minute)

	if err = MySQLDB.Ping(); err != nil {
		log.Fatalf("Failed to ping MySQL: %v", err)
	}

	log.Println("Connected to MySQL")
	createMySQLTables()
}

func initMongoDB() {
	mongoHost := getEnv("MONGO_HOST", "localhost")
	mongoPort := getEnv("MONGO_PORT", "27017")
	mongoUser := getEnv("MONGO_USER", "")
	mongoPass := getEnv("MONGO_PASSWORD", "")

	var uri string
	if mongoUser != "" && mongoPass != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", mongoUser, mongoPass, mongoHost, mongoPort)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	MongoDB, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err = MongoDB.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	MongoDBCtx = context.Background()
	log.Println("Connected to MongoDB")
}

func createPostgresTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			address TEXT,
			phone VARCHAR(50),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			total_amount DECIMAL(10, 2) NOT NULL,
			status VARCHAR(50) DEFAULT 'pending',
			payment_method VARCHAR(50),
			shipping_address TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS order_items (
			id SERIAL PRIMARY KEY,
			order_id INTEGER REFERENCES orders(id),
			product_id VARCHAR(100) NOT NULL,
			quantity INTEGER NOT NULL,
			price DECIMAL(10, 2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS cart (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			product_id VARCHAR(100) NOT NULL,
			quantity INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		if _, err := PostgresDB.Exec(query); err != nil {
			log.Printf("Error creating PostgreSQL table: %v", err)
		}
	}
	log.Println("PostgreSQL tables created/verified")
}

func createMySQLTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS inventory (
			id INT AUTO_INCREMENT PRIMARY KEY,
			product_id VARCHAR(100) UNIQUE NOT NULL,
			quantity INT NOT NULL DEFAULT 0,
			warehouse_location VARCHAR(255),
			last_restocked TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			low_stock_threshold INT DEFAULT 10,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS sales_analytics (
			id INT AUTO_INCREMENT PRIMARY KEY,
			product_id VARCHAR(100) NOT NULL,
			quantity_sold INT NOT NULL,
			revenue DECIMAL(10, 2) NOT NULL,
			sale_date DATE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		if _, err := MySQLDB.Exec(query); err != nil {
			log.Printf("Error creating MySQL table: %v", err)
		}
	}
	log.Println("MySQL tables created/verified")
}

func CloseDatabases() {
	if PostgresDB != nil {
		PostgresDB.Close()
		log.Println("PostgreSQL connection closed")
	}
	if MySQLDB != nil {
		MySQLDB.Close()
		log.Println("MySQL connection closed")
	}
	if MongoDB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		MongoDB.Disconnect(ctx)
		log.Println("MongoDB connection closed")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetMongoDatabase() *mongo.Database {
	dbName := getEnv("MONGO_DB", "ecommerce")
	return MongoDB.Database(dbName)
}
