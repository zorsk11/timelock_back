package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func ConnectDB() {
	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден")
	}

	// Получаем строку подключения из переменной окружения
	log.Println("Получаем MONGODB_URI из переменных окружения")
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("Переменная окружения MONGODB_URI не установлена")
	}
	log.Println("MONGODB_URI:", mongoURI)

	log.Println("Создаём новый клиент MongoDB")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Ошибка при создании клиента MongoDB: %v", err)
	}
	log.Println("Клиент MongoDB создан")

	// Контекст с таймаутом для подключения
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Подключаемся к MongoDB")
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Ошибка при подключении к MongoDB: %v", err)
	}
	log.Println("Соединение с MongoDB установлено")

	// Проверка соединения
	log.Println("Проверяем соединение с MongoDB")
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Ошибка при проверке соединения с MongoDB: %v", err)
	}
	log.Println("Соединение с MongoDB успешно проверено")

	DB = client
}

func GetCollection(collectionName string) *mongo.Collection {
	// Получаем название базы данных из переменной окружения
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("Переменная окружения DB_NAME не установлена")
	}
	log.Printf("Получаем базу данных: %s\n", dbName)
	db := DB.Database(dbName)
	log.Printf("База данных получена: %s\n", dbName)

	log.Printf("Получаем коллекцию: %s\n", collectionName)
	collection := db.Collection(collectionName)
	log.Printf("Коллекция получена: %s\n", collectionName)
	return collection
}
