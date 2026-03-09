package storage

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockPostgresDB for testing
type MockPostgresDB struct {
	documents     map[uuid.UUID]Document
	chunks        map[uuid.UUID]DocumentChunk
	sessions      map[uuid.UUID]ChatSession
	messages      map[uuid.UUID]ChatMessage
}

func NewMockPostgresDB() *MockPostgresDB {
	return &MockPostgresDB{
		documents: make(map[uuid.UUID]Document),
		chunks:    make(map[uuid.UUID]DocumentChunk),
		sessions:  make(map[uuid.UUID]ChatSession),
		messages:  make(map[uuid.UUID]ChatMessage),
	}
}

func (m *MockPostgresDB) CreateDocument(ctx context.Context, doc *Document) error {
	m.documents[doc.ID] = *doc
	return nil
}

func (m *MockPostgresDB) GetDocument(ctx context.Context, id uuid.UUID) (*Document, error) {
	doc, exists := m.documents[id]
	if !exists {
		return nil, ErrDocumentNotFound
	}
	return &doc, nil
}

func (m *MockPostgresDB) ListDocuments(ctx context.Context, limit, offset int) ([]Document, error) {
	docs := make([]Document, 0, len(m.documents))
	for _, doc := range m.documents {
		docs = append(docs, doc)
	}
	return docs, nil
}

func (m *MockPostgresDB) DeleteDocument(ctx context.Context, id uuid.UUID) error {
	delete(m.documents, id)
	// Also delete associated chunks
	for chunkID, chunk := range m.chunks {
		if chunk.DocumentID == id {
			delete(m.chunks, chunkID)
		}
	}
	return nil
}

func (m *MockPostgresDB) CreateDocumentChunk(ctx context.Context, chunk *DocumentChunk) error {
	m.chunks[chunk.ID] = *chunk
	return nil
}

func (m *MockPostgresDB) GetDocumentChunks(ctx context.Context, documentID uuid.UUID) ([]DocumentChunk, error) {
	chunks := make([]DocumentChunk, 0)
	for _, chunk := range m.chunks {
		if chunk.DocumentID == documentID {
			chunks = append(chunks, chunk)
		}
	}
	return chunks, nil
}

func (m *MockPostgresDB) CreateChatSession(ctx context.Context, session *ChatSession) error {
	m.sessions[session.ID] = *session
	return nil
}

func (m *MockPostgresDB) GetChatSession(ctx context.Context, id uuid.UUID) (*ChatSession, error) {
	session, exists := m.sessions[id]
	if !exists {
		return nil, ErrSessionNotFound
	}
	return &session, nil
}

func (m *MockPostgresDB) CreateChatMessage(ctx context.Context, message *ChatMessage) error {
	m.messages[message.ID] = *message
	return nil
}

func (m *MockPostgresDB) GetChatMessages(ctx context.Context, sessionID uuid.UUID) ([]ChatMessage, error) {
	messages := make([]ChatMessage, 0)
	for _, message := range m.messages {
		if message.SessionID == sessionID {
			messages = append(messages, message)
		}
	}
	return messages, nil
}

func TestCreateDocument(t *testing.T) {
	ctx := context.Background()
	mockDB := NewMockPostgresDB()

	doc := &Document{
		ID:          uuid.New(),
		Filename:    "test.pdf",
		ContentType: "application/pdf",
		FileSize:    1024,
		UploadTime:  time.Now(),
		Processed:   false,
		Metadata:    `{"pages": 10}`,
	}

	err := mockDB.CreateDocument(ctx, doc)
	require.NoError(t, err)

	// Verify document was created
	retrieved, err := mockDB.GetDocument(ctx, doc.ID)
	require.NoError(t, err)
	assert.Equal(t, doc.Filename, retrieved.Filename)
	assert.Equal(t, doc.ContentType, retrieved.ContentType)
	assert.Equal(t, doc.FileSize, retrieved.FileSize)
}

func TestGetDocumentNotFound(t *testing.T) {
	ctx := context.Background()
	mockDB := NewMockPostgresDB()

	nonExistentID := uuid.New()
	_, err := mockDB.GetDocument(ctx, nonExistentID)
	assert.Error(t, err)
	assert.Equal(t, ErrDocumentNotFound, err)
}

func TestListDocuments(t *testing.T) {
	ctx := context.Background()
	mockDB := NewMockPostgresDB()

	// Create test documents
	docs := []Document{
		{
			ID:          uuid.New(),
			Filename:    "test1.pdf",
			ContentType: "application/pdf",
			FileSize:    1024,
			UploadTime:  time.Now(),
			Processed:   false,
			Metadata:    `{"pages": 10}`,
		},
		{
			ID:          uuid.New(),
			Filename:    "test2.txt",
			ContentType: "text/plain",
			FileSize:    512,
			UploadTime:  time.Now(),
			Processed:   true,
			Metadata:    `{"lines": 100}`,
		},
	}

	for _, doc := range docs {
		err := mockDB.CreateDocument(ctx, &doc)
		require.NoError(t, err)
	}

	// List documents
	retrieved, err := mockDB.ListDocuments(ctx, 10, 0)
	require.NoError(t, err)
	assert.Len(t, retrieved, 2)
}

func TestDeleteDocument(t *testing.T) {
	ctx := context.Background()
	mockDB := NewMockPostgresDB()

	// Create document
	doc := &Document{
		ID:          uuid.New(),
		Filename:    "test.pdf",
		ContentType: "application/pdf",
		FileSize:    1024,
		UploadTime:  time.Now(),
		Processed:   false,
		Metadata:    `{"pages": 10}`,
	}

	err := mockDB.CreateDocument(ctx, doc)
	require.NoError(t, err)

	// Delete document
	err = mockDB.DeleteDocument(ctx, doc.ID)
	require.NoError(t, err)

	// Verify document is deleted
	_, err = mockDB.GetDocument(ctx, doc.ID)
	assert.Error(t, err)
	assert.Equal(t, ErrDocumentNotFound, err)
}

func TestCreateDocumentChunk(t *testing.T) {
	ctx := context.Background()
	mockDB := NewMockPostgresDB()

	chunk := &DocumentChunk{
		ID:         uuid.New(),
		DocumentID: uuid.New(),
		ChunkIndex: 0,
		Content:    "This is a test chunk",
		Embedding:  []float32{0.1, 0.2, 0.3},
		PageNumber: intPtr(1),
		Metadata:   `{"type": "paragraph"}`,
		CreatedAt:  time.Now(),
	}

	err := mockDB.CreateDocumentChunk(ctx, chunk)
	require.NoError(t, err)

	// Verify chunk was created
	chunks, err := mockDB.GetDocumentChunks(ctx, chunk.DocumentID)
	require.NoError(t, err)
	assert.Len(t, chunks, 1)
	assert.Equal(t, chunk.Content, chunks[0].Content)
	assert.Equal(t, chunk.ChunkIndex, chunks[0].ChunkIndex)
}

func TestCreateChatSession(t *testing.T) {
	ctx := context.Background()
	mockDB := NewMockPostgresDB()

	session := &ChatSession{
		ID:           uuid.New(),
		SessionName:  stringPtr("Test Session"),
		ModelName:    "llama3.1:8b",
		CreatedAt:    time.Now(),
		LastActivity: time.Now(),
	}

	err := mockDB.CreateChatSession(ctx, session)
	require.NoError(t, err)

	// Verify session was created
	retrieved, err := mockDB.GetChatSession(ctx, session.ID)
	require.NoError(t, err)
	assert.Equal(t, session.ModelName, retrieved.ModelName)
	assert.Equal(t, *session.SessionName, *retrieved.SessionName)
}

func TestCreateChatMessage(t *testing.T) {
	ctx := context.Background()
	mockDB := NewMockPostgresDB()

	sessionID := uuid.New()
	message := &ChatMessage{
		ID:           uuid.New(),
		SessionID:    sessionID,
		MessageType:  "user",
		Content:      "Hello, how are you?",
		Citations:    `[]`,
		ModelUsed:    nil,
		TokensUsed:   intPtr(50),
		ResponseTime: intPtr(1000),
		CreatedAt:    time.Now(),
	}

	err := mockDB.CreateChatMessage(ctx, message)
	require.NoError(t, err)

	// Verify message was created
	messages, err := mockDB.GetChatMessages(ctx, sessionID)
	require.NoError(t, err)
	assert.Len(t, messages, 1)
	assert.Equal(t, message.Content, messages[0].Content)
	assert.Equal(t, message.MessageType, messages[0].MessageType)
}

func TestGetChatMessages(t *testing.T) {
	ctx := context.Background()
	mockDB := NewMockPostgresDB()

	sessionID := uuid.New()
	messages := []ChatMessage{
		{
			ID:           uuid.New(),
			SessionID:    sessionID,
			MessageType:  "user",
			Content:      "Hello",
			Citations:    `[]`,
			CreatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			SessionID:    sessionID,
			MessageType:  "assistant",
			Content:      "Hi there!",
			Citations:    `[]`,
			CreatedAt:    time.Now(),
		},
	}

	for _, msg := range messages {
		err := mockDB.CreateChatMessage(ctx, &msg)
		require.NoError(t, err)
	}

	// Get messages
	retrieved, err := mockDB.GetChatMessages(ctx, sessionID)
	require.NoError(t, err)
	assert.Len(t, retrieved, 2)

	// Verify order (should be chronological)
	assert.Equal(t, "user", retrieved[0].MessageType)
	assert.Equal(t, "assistant", retrieved[1].MessageType)
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
