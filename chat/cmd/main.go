package cmd

import (
	"context"
	"encoding/json"
	"fitnessme/chat/pkg/db"
	"fitnessme/chat/pkg/models"
	"fitnessme/chat/pkg/repository"
	"fitnessme/chat/pkg/services"
	"fitnessme/cors"
	"fitnessme/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/gorilla/websocket"
)

var connections = make(map[string]*websocket.Conn)

type ChatApp struct {
	router http.Handler
	db.Handler
	repo repository.ChatRepository
	jw   utils.JwtWrapper
}

func New(jwtWrapper *utils.JwtWrapper) *ChatApp {
	handler := db.Init("localhost", "chat", "postgres", "123", 5432)
	repo := repository.NewChatRepository(handler)

	app := &ChatApp{
		router: loadRoutes(repo, jwtWrapper),
		repo:   repo,
		jw:     *jwtWrapper,
	}

	app.router = cors.Cors(app.router)
	return app
}

func loadRoutes(repo repository.ChatRepository, jwtWrapper *utils.JwtWrapper) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	chatService := services.NewChatService(repo, jwtWrapper)

	router.Get("/ws", handleWebSocket(chatService))

	router.Route("/chat", func(router chi.Router) {
		loadAuthRoutes(router, chatService)
	})

	return router
}

func handleWebSocket(chatService services.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Failed to upgrade connection to WebSocket:", err)
			return
		}
		defer conn.Close()

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading WebSocket message: %v", err)
				return
			}
			if messageType != websocket.TextMessage {
				log.Println("Received unsupported message type:", messageType)
				continue
			}

			var message models.Message
			if err := json.Unmarshal(p, &message); err != nil {
				log.Println("Error decoding message:", err)
				continue
			}
			chatID := message.ChatId
			connections[chatID+message.SentBy.String()] = conn

			chatService.HandleMessage(w, r, message)

			fmt.Println(message.SentBy.String())
			broadcastMessage(message)
		}
	}
}

func broadcastMessage(message models.Message) {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	for chatID, conn := range connections {
		if conn != nil {
			err = conn.WriteMessage(websocket.TextMessage, messageJSON)
			if err != nil {
				log.Printf("Failed to write message to WebSocket connection for chatID %s: %v", chatID, err)
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func loadAuthRoutes(router chi.Router, chatService services.ChatService) {
	router.Get("/", chatService.GetChat)
	router.Post("/", chatService.Create)
	router.Get("/messages", chatService.GetMessages)
	router.Get("/unread", chatService.GetUnreadMessages)
	router.Post("/readAll", chatService.ReadAll)
}

func (a *ChatApp) Start(ctw context.Context) error {
	server := &http.Server{
		Addr:    ":3003",
		Handler: a.router,
	}
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}
