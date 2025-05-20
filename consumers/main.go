package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Deleted   string    `json:"__deleted"`
}

type Message struct {
	Payload User `json:"payload"`
}

var messages []Message

func main() {
	brokers := []string{os.Getenv("KAFKA_URL")}
	topic := os.Getenv("KAFKA_TOPIC")
	groupID := os.Getenv("KAFKA_GROUP_ID")

	connStr := os.Getenv("POSTGRES_CONN")
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Erro ao conectar ao PostgreSQL: %v", err)
	}
	defer pool.Close()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	defer reader.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("\nEncerrando o consumidor...")
		cancel()
	}()

	r := mux.NewRouter()
	r.HandleFunc("/messages", getMessages).Methods("GET")
	r.HandleFunc("/health", healthCheck).Methods("GET")

	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8000"
		}
		fmt.Printf("Servidor da API rodando em http://localhost:%s\n", port)
		log.Fatal(http.ListenAndServe(":"+port, r))
	}()

	fmt.Println("Iniciando o consumidor...")

	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("Consumidor encerrado")
				break
			}
			log.Fatalf("Erro ao consumir mensagem: %v", err)
		}

		var msg Message
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Printf("Erro ao desserializar mensagem: %v", err)
			continue
		}

		messages = append(messages, msg)
		fmt.Printf("Mensagem recebida: %+v\n", msg)

		insertUser(ctx, pool, msg.Payload)
	}

	fmt.Println("Consumidor finalizado.")
}

func insertUser(ctx context.Context, pool *pgxpool.Pool, user User) {
	_, err := pool.Exec(ctx, `INSERT INTO users (id, name, email, password, created_at, updated_at, __deleted) VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (id) DO UPDATE SET name = $2, email = $3, password = $4, created_at = $5, updated_at = $6, __deleted = $7`,
		user.ID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt, user.Deleted)
	if err != nil {
		log.Printf("Erro ao inserir/atualizar usuário no PostgreSQL: %v", err)
	}
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
