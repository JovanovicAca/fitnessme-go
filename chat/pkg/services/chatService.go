package services

import (
	"encoding/json"
	"fitnessme/chat/pkg/dto"
	"fitnessme/chat/pkg/models"
	"fitnessme/chat/pkg/repository"
	"fitnessme/utils"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type ChatService interface {
	HandleMessage(w http.ResponseWriter, r *http.Request, message models.Message) error
	GetChat(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	GetMessages(w http.ResponseWriter, r *http.Request)
	GetUnreadMessages(w http.ResponseWriter, r *http.Request)
	ReadAll(w http.ResponseWriter, r *http.Request)
}

type chatService struct {
	repo       repository.ChatRepository
	jwtWrapper *utils.JwtWrapper
}

func NewChatService(repo repository.ChatRepository, jwtWrapper *utils.JwtWrapper) ChatService {
	return &chatService{repo: repo, jwtWrapper: jwtWrapper}
}

func (c *chatService) ReadAll(w http.ResponseWriter, r *http.Request) {
	var chatDTO dto.ChatDTO
	if err := json.NewDecoder(r.Body).Decode(&chatDTO); err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusInternalServerError)
		return
	}

	c.repo.ReadAll(chatDTO.User1.String(), chatDTO.User2.String())
	w.WriteHeader(http.StatusOK)
}

func (c *chatService) GetUnreadMessages(w http.ResponseWriter, r *http.Request) {
	user1ID := r.URL.Query().Get("user1")
	user2ID := r.URL.Query().Get("user2")

	num, err := c.repo.UnreadMessages(user1ID, user2ID)
	if err != nil {
		http.Error(w, "Failed to check chat existence", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]int64{"unread": num}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *chatService) GetMessages(w http.ResponseWriter, r *http.Request) {
	chatID := r.URL.Query().Get("chatid")

	messages, err := c.repo.GetChatMessages(chatID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (c *chatService) Create(w http.ResponseWriter, r *http.Request) {
	var chatDTO dto.ChatDTO
	if err := json.NewDecoder(r.Body).Decode(&chatDTO); err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusInternalServerError)
		return
	}

	var chat models.Chat
	chat.Id = uuid.New()
	chat.User1 = chatDTO.User1
	chat.User2 = chatDTO.User2

	if err := c.repo.CreateChat(chat); err != nil {
		http.Error(w, "Failed to save chat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *chatService) GetChat(w http.ResponseWriter, r *http.Request) {
	user1ID := r.URL.Query().Get("user1")
	user2ID := r.URL.Query().Get("user2")

	chatID, err := c.repo.ChatExists(uuid.MustParse(user1ID), uuid.MustParse(user2ID))
	if err != nil {
		http.Error(w, "Failed to check chat existence", http.StatusInternalServerError)
		return
	}
	fmt.Println(chatID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"chat_id": chatID}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *chatService) HandleMessage(w http.ResponseWriter, r *http.Request, message models.Message) error {
	log.Println("Received message: ", message.Text)
	message.Id = uuid.New()
	message.Status = "delivered"

	err := c.repo.InsertMessage(message)
	if err != nil {
		http.Error(w, "Failed to save", http.StatusInternalServerError)
		return nil
	}

	return nil
}

func (c *chatService) getUserId(r *http.Request) (bool, string) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false, "Authorization header is required"
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return false, "Invalid Authorization header format"
	}

	token := headerParts[1]
	claims, err := c.jwtWrapper.ValidateToken(token)
	if err != nil {
		return false, "Invalid token: " + err.Error()
	}

	return true, claims.Id.String()
}
