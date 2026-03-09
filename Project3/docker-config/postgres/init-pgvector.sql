-- Initialize PGVector extension and create tables for the knowledge base

-- Enable the PGVector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Create documents table
CREATE TABLE IF NOT EXISTS documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    filename VARCHAR(255) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    file_size INTEGER NOT NULL,
    upload_time TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    processed BOOLEAN DEFAULT FALSE,
    metadata JSONB DEFAULT '{}'::jsonb
);

-- Create document chunks table for RAG
CREATE TABLE IF NOT EXISTS document_chunks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    chunk_index INTEGER NOT NULL,
    content TEXT NOT NULL,
    embedding vector(768), -- Nomic-Embed-Text dimension
    page_number INTEGER,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create chat sessions table
CREATE TABLE IF NOT EXISTS chat_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_name VARCHAR(255),
    model_name VARCHAR(100) NOT NULL DEFAULT 'llama3.1:8b',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create chat messages table
CREATE TABLE IF NOT EXISTS chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES chat_sessions(id) ON DELETE CASCADE,
    message_type VARCHAR(20) NOT NULL CHECK (message_type IN ('user', 'assistant')),
    content TEXT NOT NULL,
    citations JSONB DEFAULT '[]'::jsonb,
    model_used VARCHAR(100),
    tokens_used INTEGER,
    response_time_ms INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_document_chunks_document_id ON document_chunks(document_id);
CREATE INDEX IF NOT EXISTS idx_document_chunks_embedding ON document_chunks USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
CREATE INDEX IF NOT EXISTS idx_chat_messages_session_id ON chat_messages(session_id);
CREATE INDEX IF NOT EXISTS idx_chat_messages_created_at ON chat_messages(created_at);

-- Create a function to update last_activity timestamp
CREATE OR REPLACE FUNCTION update_chat_session_activity()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE chat_sessions 
    SET last_activity = NOW() 
    WHERE id = NEW.session_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically update last_activity
CREATE TRIGGER trigger_update_chat_session_activity
    AFTER INSERT ON chat_messages
    FOR EACH ROW
    EXECUTE FUNCTION update_chat_session_activity();

-- Insert sample data for testing (optional)
-- This will be removed in production
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM documents LIMIT 1) THEN
        INSERT INTO documents (filename, content_type, file_size, processed, metadata)
        VALUES 
            ('sample-document.txt', 'text/plain', 1024, true, '{"description": "Sample document for testing"}');
        
        -- Get the document ID for sample chunks
        DECLARE doc_id UUID;
        BEGIN
            SELECT id INTO doc_id FROM documents WHERE filename = 'sample-document.txt';
            
            -- Insert sample chunks with mock embeddings
            INSERT INTO document_chunks (document_id, chunk_index, content, embedding, metadata)
            VALUES 
                (doc_id, 0, 'This is a sample chunk of text for testing the vector search functionality.', 
                 array_fill(0.1::float, ARRAY[768]), 
                 '{"test": true}'),
                (doc_id, 1, 'Another sample chunk to demonstrate the RAG system capabilities.', 
                 array_fill(0.2::float, ARRAY[768]), 
                 '{"test": true}');
        END;
    END IF;
END $$;
