package handler

import (
	"net/http"

	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

type Handler interface {
	CreateTopic(http.ResponseWriter, *http.Request)
	ListTopics(http.ResponseWriter, *http.Request)
	GetTopicByID(http.ResponseWriter, *http.Request)
	DeleteTopicByID(http.ResponseWriter, *http.Request)

	CreateMessage(http.ResponseWriter, *http.Request)
	ListMessagesByTopicID(http.ResponseWriter, *http.Request)
	DeleteMessageByID(http.ResponseWriter, *http.Request)
}

type handler struct {
	topics   repo.Topics
	messages repo.Messages
}

func New(repo repo.Repository) Handler {
	return &handler{
		topics:   repo.Topic,
		messages: repo.Message,
	}
}
