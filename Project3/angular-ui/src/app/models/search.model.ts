export interface SearchRequest {
  query: string;
  limit?: number;
  threshold?: number;
  document_ids?: string[];
}

export interface SearchResult {
  chunk: {
    id: string;
    document_id: string;
    chunk_index: number;
    content: string;
    page_number?: number;
    metadata?: {
      type?: string;
      section?: string;
    };
  };
  similarity: number;
  document: {
    id: string;
    filename: string;
    title?: string;
  };
}

export interface SearchResponse {
  results: SearchResult[];
  total: number;
  query: string;
  processing_time_ms: number;
}

export interface HybridSearchRequest {
  query: string;
  semantic_weight?: number;
  keyword_weight?: number;
  limit?: number;
}

export interface HybridSearchResult {
  chunk: {
    id: string;
    document_id: string;
    chunk_index: number;
    content: string;
    page_number?: number;
    metadata?: {
      type?: string;
      section?: string;
    };
  };
  similarity: number;
  keyword_score: number;
  combined_score: number;
  document: {
    id: string;
    filename: string;
    title?: string;
  };
}

export interface HybridSearchResponse {
  results: HybridSearchResult[];
  total: number;
  query: string;
  processing_time_ms: number;
}
