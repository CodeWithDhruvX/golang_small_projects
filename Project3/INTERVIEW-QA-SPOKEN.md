# Interview Q&A (Spoken Format) – Project3 (Go + Angular + RAG)

Below are concise, interview-ready questions with answers phrased in a conversational, spoken style. Each answer aims for roughly a 60–90 second response.

---

## 1) Give me a quick tour of your system.
In short, it’s a full-stack RAG application. The Angular frontend handles auth, uploads, search, and chat. The Go backend exposes REST endpoints, authenticates via JWT middleware, ingests and chunks documents, creates embeddings, stores both text and vectors in Postgres with pgvector, retrieves top‑k matches, and constructs prompts for the model. Observability is wired through Prometheus metrics and Grafana dashboards. Locally it runs with Docker Compose; there are Kubernetes manifests for cluster deployment. I kept the service boundary simple—one Go service—for speed and clarity, but the seams are clean if we ever need to split ingestion, retrieval, or chat into separate services.

## 2) How does authentication work end‑to‑end?
The user logs in from Angular; we receive a token and store it securely in the client. An interceptor attaches the token to requests. On the server, a middleware validates and parses the JWT on each protected endpoint. I keep expiries short and can add refresh tokens for longer sessions. Claims carry the minimal user identity; any authorization checks can happen downstream. Errors are standardized so the UI can handle redirects or error toasts consistently.

## 3) Why did you choose Go for the backend and Angular for the frontend?
Go gives me performance, great concurrency primitives, and a single static binary that’s easy to containerize. It’s a strong fit for CPU‑bound ingestion and I/O‑heavy retrieval. Angular gives a batteries‑included structure: routing, guards, interceptors, strong TypeScript patterns, and a component model that scales cleanly. Together they let me move quickly while keeping type safety and testability.

## 4) Walk me through the RAG flow.
A user query hits the chat endpoint. I embed the query, run a vector similarity search in Postgres via pgvector, and retrieve the most relevant chunks. I format a prompt with system instructions, the user question, and the retrieved context. The model returns an answer that the UI can stream to the user. This design keeps the app stateless and the data discoverable, while letting me tune chunk sizes, top‑k, and the prompt to balance latency and quality.

## 5) How do you ingest documents efficiently?
I handle different file types with dedicated parsers, chunk text into overlapping windows, and deduplicate via content hashes. Ingestion is structured to be idempotent: if the same document is re‑uploaded, I detect it and skip or update as appropriate. I control concurrency with a worker‑pool pattern to avoid CPU or memory spikes. Failures are surfaced with metrics and logs, and I can add retries with backoff where it makes sense.

## 6) Why Postgres + pgvector instead of a dedicated vector database?
For this project, Postgres plus pgvector hits a sweet spot: fewer moving parts, strong consistency, and straightforward operations. It’s great for prototyping to mid‑scale. If I needed very large scale or sub‑100ms latencies across billions of vectors, I’d consider specialized vector stores or approximate indexes (like HNSW) and potentially sharding. But here, Postgres keeps my data model unified and easy to manage.

## 7) How do you reduce hallucinations?
I focus on prompt construction and retrieval quality. The prompt instructs the model to ground answers in the retrieved context and to say “I don’t know” if the context doesn’t support an answer. I keep chunks semantically coherent, tune top‑k, and can re‑rank results. I also include citations or source snippets when possible so users can verify. Finally, I monitor feedback and adjust chunking and prompts iteratively.

## 8) How does the chat UI handle responses?
On the UI, I treat long‑running model calls as streams, so the user gets incremental updates and can cancel if needed. I show clear loading states and fall back to non‑streamed responses when the backend doesn’t support streaming. Errors bubble through a centralized interceptor so the UX remains consistent. The net effect is fast perceived performance and better control for the user.

## 9) What’s your approach to error handling and logging?
I use structured logging with levels and consistent error envelopes at the API layer. That means the UI can confidently parse an error code and message. On the backend, logs include request context where available, making it easy to correlate failures. I avoid leaking sensitive data and keep logs actionable. Between logs and metrics, I can quickly move from “something’s wrong” to the specific failing path.

## 10) What metrics do you export and how do you use them?
I export request counts, latencies, error rates, ingestion timings, and embedding/DB timings. Prometheus scrapes these, and Grafana visualizes golden signals—throughput, latency, errors, and saturation. With that, I set simple SLOs for endpoints and can spot regressions or spikes immediately. For example, if 99th‑percentile latency creeps up, I can see whether it’s embedding, vector search, or prompt processing.

## 11) How do you test the system?
I split tests into unit and integration. Unit tests cover JWT parsing, chunking logic, and storage boundaries with fakes. Integration tests bring up dependencies, seed sample data, and run full request flows. For LLMs, I stub the embedding/model layer so tests are deterministic. That way, I can confidently refactor internals without breaking behavior.

## 12) How do you deploy it locally and to Kubernetes?
Locally, Docker Compose brings up the backend, Postgres with pgvector, Prometheus, and Grafana. For clusters, I have Kubernetes manifests for the app, database, and monitoring stack. The main changes are around configuration, secrets, and persistence. Scaling is straightforward in K8s: bump replicas for stateless services and size the database appropriately.

## 13) How would you scale this system?
Horizontally scale the Go service, keep it stateless, and handle sessions via JWTs. Move embedding jobs to a background worker with a queue to decouple ingestion from user traffic. Add caching for frequent queries and reuse embeddings where possible. On the DB side, introduce approximate indexes or sharding for vectors, and add read replicas for analytical reads. I’d also profile end‑to‑end to ensure I’m solving the bottlenecks that actually matter.

## 14) How do you handle backpressure and rate limiting?
On ingestion, I use a bounded worker pool so we never overload memory or CPU. For APIs, I’d introduce a token‑bucket rate limiter keyed by user or IP and surface “retry‑after” headers. Internally, I bound queue sizes and apply exponential backoff on transient failures. Combined with metrics, this keeps the system stable under bursty or abusive traffic.

## 15) What security measures did you take for uploads and APIs?
I validate file types and sizes early, avoid writing untrusted files to arbitrary paths, and sanitize content. JWTs secure protected endpoints; I keep tokens short‑lived and can add refresh rotations. On the API side, I validate inputs, use parameterized queries, and return minimal error details to avoid information leaks. For the UI, I avoid storing tokens in places vulnerable to XSS and lean on the browser’s security model.

## 16) How do you manage configuration and secrets?
I separate config from code via environment variables, and mount secrets through the runtime environment or Kubernetes secrets rather than committing them. Locally, I keep non‑sensitive defaults for fast starts. For production, I integrate with a secret manager and rotate keys regularly. The goal is reproducible deploys without leaking sensitive data.

## 17) Describe your API design patterns.
I keep endpoints predictable and resource‑oriented, return standardized error envelopes, and support pagination for list endpoints. Responses are typed and consistent, which makes the Angular client simple. If we needed public clients, I’d add OpenAPI and generate SDKs. Versioning is easy to introduce with a /v1 prefix when breaking changes arise.

## 18) If p99 latency spikes, how do you triage?
I’d start with the Grafana dashboard: is the spike global or on a specific endpoint? Then I’d break down by component timing—embedding calls, vector search, prompt assembly, or model response. I’d correlate with error rates and saturation metrics like CPU, memory, and DB connections. From there, I’d reproduce with a load test and add targeted profiling—often the fix is a hot path optimization or adding a small cache.

## 19) If a bad embedding model hurts search quality, what’s your rollback plan?
I pin embedding model versions and store that metadata with vectors. If quality drops, I can freeze the current model, re‑embed in the background, and atomically switch queries to the new index or table. If needed, I roll back to the previous embeddings quickly. I also keep a small offline evaluation suite so I can validate quality before switching traffic.

## 20) How would you answer this for a service‑based client vs a product company?
For service clients, I emphasize adaptability: clean seams, clear SDLC, documentation, and how I’d re‑platform this to a different stack. I talk about handover quality—dashboards, runbooks, and API docs. For product companies, I emphasize SLOs, performance, incremental experimentation, and cost/efficiency trade‑offs. I frame decisions in terms of user‑visible latency, reliability, and the roadmap to scale.

---

If you’d like, I can tailor this further to specific companies or tighten each answer to a strict 60‑second delivery.

