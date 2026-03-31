package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogEntry struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Level     string             `json:"level" bson:"level"` // info, warn, error, debug
	Message   string             `json:"message" bson:"message"`
	Service   string             `json:"service" bson:"service"`
	UserID    *uint              `json:"user_id,omitempty" bson:"user_id,omitempty"`
	RequestID string             `json:"request_id" bson:"request_id"`
	Metadata  map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
}

type UserProfile struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      uint               `json:"user_id" bson:"user_id"`
	Avatar      string             `json:"avatar" bson:"avatar"`
	Bio         string             `json:"bio" bson:"bio"`
	Location    string             `json:"location" bson:"location"`
	Website     string             `json:"website" bson:"website"`
	SocialLinks map[string]string  `json:"social_links" bson:"social_links"`
	Preferences map[string]interface{} `json:"preferences" bson:"preferences"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type AnalyticsEvent struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EventType   string             `json:"event_type" bson:"event_type"`
	UserID      *uint              `json:"user_id,omitempty" bson:"user_id,omitempty"`
	SessionID   string             `json:"session_id" bson:"session_id"`
	Properties  map[string]interface{} `json:"properties" bson:"properties"`
	Timestamp   time.Time          `json:"timestamp" bson:"timestamp"`
	IPAddress   string             `json:"ip_address" bson:"ip_address"`
	UserAgent   string             `json:"user_agent" bson:"user_agent"`
}

type Document struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Content     string             `json:"content" bson:"content"`
	AuthorID    uint               `json:"author_id" bson:"author_id"`
	Tags        []string           `json:"tags" bson:"tags"`
	Category    string             `json:"category" bson:"category"`
	Status      string             `json:"status" bson:"status"` // draft, published, archived
	Metadata    map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Version     int                `json:"version" bson:"version"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type Session struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SessionID    string             `json:"session_id" bson:"session_id"`
	UserID       *uint              `json:"user_id,omitempty" bson:"user_id,omitempty"`
	IPAddress    string             `json:"ip_address" bson:"ip_address"`
	UserAgent    string             `json:"user_agent" bson:"user_agent"`
	LastActivity time.Time          `json:"last_activity" bson:"last_activity"`
	Data         map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	ExpiresAt    time.Time          `json:"expires_at" bson:"expires_at"`
}
