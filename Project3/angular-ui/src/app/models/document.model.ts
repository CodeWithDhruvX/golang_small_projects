export interface Document {
  id: string;
  filename: string;
  content_type: string;
  file_size: number;
  upload_time: string;
  processed: boolean;
  processing_status: 'pending' | 'processing' | 'completed' | 'failed';
  chunk_count: number;
  metadata?: {
    pages?: number;
    title?: string;
    author?: string;
    created_date?: string;
  };
}

export interface DocumentUploadResponse {
  document_id: string;
  filename: string;
  size: number;
  content_type: string;
  status: string;
  processing_status: string;
}

export interface DocumentListResponse {
  documents: Document[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}

export interface DocumentChunk {
  id: string;
  document_id: string;
  chunk_index: number;
  content: string;
  page_number?: number;
  metadata?: {
    type?: string;
    section?: string;
  };
}

export interface DocumentChunksResponse {
  chunks: DocumentChunk[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}
