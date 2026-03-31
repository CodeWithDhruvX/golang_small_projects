package database

import (
	"context"
	"fmt"
	"log"
	"time"
	"project7/internal/config"
	"project7/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB(cfg *config.AppConfig) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := cfg.GetMongoURI()
	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Successfully connected to MongoDB database")

	database := client.Database(cfg.Database.MongoDB)

	// Create indexes for better performance
	mongoDB := &MongoDB{
		Client:   client,
		Database: database,
	}

	err = mongoDB.createIndexes(ctx)
	if err != nil {
		log.Printf("Warning: Failed to create MongoDB indexes: %v", err)
	}

	return mongoDB, nil
}

func (m *MongoDB) createIndexes(ctx context.Context) error {
	// Index for LogEntry
	logCollection := m.Database.Collection("log_entries")
	_, err := logCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "timestamp", Value: -1},
			{Key: "level", Value: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create log entries index: %w", err)
	}

	// Index for UserProfile
	profileCollection := m.Database.Collection("user_profiles")
	_, err = profileCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create user profile index: %w", err)
	}

	// Index for AnalyticsEvent
	analyticsCollection := m.Database.Collection("analytics_events")
	_, err = analyticsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "timestamp", Value: -1},
			{Key: "event_type", Value: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create analytics events index: %w", err)
	}

	// Index for Session
	sessionCollection := m.Database.Collection("sessions")
	_, err = sessionCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "session_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create session index: %w", err)
	}

	log.Println("MongoDB indexes created successfully")
	return nil
}

func (m *MongoDB) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}

// Helper methods for common operations
func (m *MongoDB) CreateLogEntry(entry *models.LogEntry) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.Database.Collection("log_entries").InsertOne(ctx, entry)
	return err
}

func (m *MongoDB) GetLogEntries(limit int64) ([]models.LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(limit)
	cursor, err := m.Database.Collection("log_entries").Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	var entries []models.LogEntry
	err = cursor.All(ctx, &entries)
	return entries, err
}

func (m *MongoDB) CreateUserProfile(profile *models.UserProfile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.Database.Collection("user_profiles").InsertOne(ctx, profile)
	return err
}

func (m *MongoDB) GetUserProfile(userID uint) (*models.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var profile models.UserProfile
	err := m.Database.Collection("user_profiles").FindOne(ctx, bson.M{"user_id": userID}).Decode(&profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (m *MongoDB) UpdateUserProfile(userID uint, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.Database.Collection("user_profiles").UpdateOne(
		ctx,
		bson.M{"user_id": userID},
		bson.M{"$set": update},
	)
	return err
}

func (m *MongoDB) CreateAnalyticsEvent(event *models.AnalyticsEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.Database.Collection("analytics_events").InsertOne(ctx, event)
	return err
}

func (m *MongoDB) GetAnalyticsEvents(eventType string, limit int64) ([]models.AnalyticsEvent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if eventType != "" {
		filter["event_type"] = eventType
	}

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(limit)
	cursor, err := m.Database.Collection("analytics_events").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var events []models.AnalyticsEvent
	err = cursor.All(ctx, &events)
	return events, err
}
