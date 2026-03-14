To implement the **features you described as Functional Requirements**, you need to convert each capability of the system into **clear, testable requirements** that developers can implement and QA can verify. Below is a **structured Functional Requirements Specification (FRS)** tailored to your stack (**Go + Gin, Angular, Ollama, PGVector, Redis, RAG pipeline**).

---

# 1. User Authentication

### FR-1: User Registration

**Description:**
The system shall allow users to create an account.

**Inputs**

* Email
* Password

**Process**

* Validate email format
* Hash password using bcrypt
* Store user in PostgreSQL

**Output**

* User account created
* JWT token issued

**API**

```
POST /auth/register
```

---

### FR-2: User Login

**Description:**
The system shall authenticate users using email and password.

**Process**

* Verify password
* Generate JWT token
* Return session token

**API**

```
POST /auth/login
```

---

# 2. User Profile Management

### FR-3: Manage Candidate Profile

**Description:**
The system shall allow users to create and update candidate information.

**Fields**

* Name
* Experience
* Skills
* Current Salary
* Expected Salary
* Notice Period
* Location
* LinkedIn
* GitHub

**API**

```
GET /profile
PUT /profile
```

---

### FR-4: Resume Upload

**Description:**
The system shall allow users to upload resume files.

**Process**

1. Upload PDF
2. Extract text using **unipdf**
3. Store raw text
4. Trigger embedding generation

**API**

```
POST /profile/resume
```

---

# 3. Data Ingestion (AI Knowledge Feeding)

### FR-5: Resume Processing Pipeline

**Description:**
The system shall process uploaded resumes and convert them into embeddings.

**Process Flow**

```
Upload Resume
     ↓
Extract Text
     ↓
Chunk Text
     ↓
Generate Embeddings
     ↓
Store in PGVector
```

**Output**

* Resume stored in vector database.

---

### FR-6: Email Ingestion

**Description:**
The system shall import recruiter emails into the system.

**Sources**

* Gmail IMAP
* Outlook IMAP
* Manual upload

**Process**

```
Email received
   ↓
Parse email
   ↓
Store email content
   ↓
Generate embedding
```

**API**

```
POST /emails/import
GET /emails
```

---

# 4. Recruiter Email Detection

### FR-7: Recruiter Email Classification

**Description:**
The system shall detect whether an email is a recruiter email.

**Input**

```
Email content
```

**Process**

1. Send email text to local LLM
2. Classify as recruiter / non recruiter

**Output**

```
{
  "isRecruiterEmail": true
}
```

---

# 5. Requirement Extraction

### FR-8: Extract Requested Information

**Description:**
The system shall extract candidate information requested in recruiter emails.

**Input**

```
Email body
```

**Process**

LLM analyzes and extracts fields.

**Output**

```
{
 "resume": true,
 "experience": true,
 "expected_ctc": true,
 "notice_period": true
}
```

---

# 6. RAG Retrieval

### FR-9: Retrieve Candidate Knowledge

**Description:**
The system shall retrieve relevant candidate data from the vector database.

**Process**

```
Email request
   ↓
Generate embedding
   ↓
Vector search
   ↓
Retrieve relevant context
```

**Database**

PostgreSQL + PGVector.

---

# 7. AI Reply Generation

### FR-10: Generate Email Reply

**Description:**
The system shall generate professional replies to recruiter emails.

**Input**

* recruiter email
* retrieved context
* candidate profile

**Process**

Send prompt to local LLM via Ollama.

**Output**

```
AI generated reply
```

**API**

```
POST /emails/generate-reply
```

---

# 8. Application Tracking

### FR-11: Track Job Applications

**Description:**
The system shall track job applications sent to recruiters.

**Fields**

* Company
* Role
* Recruiter Email
* Status
* Date

**Statuses**

```
Applied
Interview
Offer
Rejected
```

**API**

```
GET /applications
POST /applications
```

---

# 9. Duplicate Application Detection

### FR-12: Prevent Duplicate Applications

**Description:**
The system shall detect duplicate recruiter responses.

**Process**

Check:

```
company + recruiter_email
```

If exists:

```
duplicate warning
```

---

# 10. Follow-up Email Generation

### FR-13: Generate Follow-Up Emails

**Description:**
The system shall generate follow-up emails for pending recruiter responses.

**Trigger**

```
No recruiter response after X days
```

**Output**

AI generated follow-up email.

---

# 11. AI Knowledge Embedding

### FR-14: Embedding Generation

**Description:**
The system shall convert text data into vector embeddings.

**Model**

Nomic Embed Text

**Process**

```
Text chunk
   ↓
Embedding model
   ↓
Vector output
```

---

# 12. Vector Search

### FR-15: Semantic Search

**Description:**
The system shall retrieve relevant knowledge using vector similarity search.

**Database**

PostgreSQL + PGVector.

**Example query**

```
SELECT content
FROM knowledge_base
ORDER BY embedding <-> query_embedding
LIMIT 5
```

---

# 13. Performance Optimization

### FR-16: Cache AI Responses

**Description:**
The system shall cache generated AI responses.

**Technology**

Redis.

**Process**

```
Prompt hash → cached response
```

---

# 14. Monitoring

### FR-17: System Metrics

**Description:**
The system shall expose metrics for monitoring.

Metrics include:

```
API latency
AI inference time
Vector query time
Email ingestion rate
```

**Tools**

Prometheus + Grafana.

---

# 15. API Documentation

### FR-18: API Documentation

**Description:**
The system shall provide interactive API documentation.

**Tool**

Swagger.

---

# 16. Frontend Requirements

Angular UI must provide:

### Dashboard

View recruiter activity.

### Email Inbox

List recruiter emails.

### AI Draft Reply

View and edit AI generated replies.

### Application Tracker

Track job applications.

### Profile Settings

Manage candidate profile.

---

# 17. Non-Functional Requirements

### Performance

* API response < 200 ms
* AI response < 5 seconds

### Security

* JWT authentication
* encrypted passwords
* secure API access

### Privacy

* All AI processing must run locally via Ollama.

---

# Final Summary

Your **functional requirements define how each feature behaves**:

| Feature                 | Functional Requirement |
| ----------------------- | ---------------------- |
| Authentication          | FR-1, FR-2             |
| Profile management      | FR-3                   |
| Resume ingestion        | FR-4                   |
| Data ingestion pipeline | FR-5                   |
| Email processing        | FR-6                   |
| Recruiter detection     | FR-7                   |
| Information extraction  | FR-8                   |
| RAG retrieval           | FR-9                   |
| AI reply generation     | FR-10                  |
| Application tracking    | FR-11                  |
| Duplicate detection     | FR-12                  |
| Follow-ups              | FR-13                  |
| Embedding generation    | FR-14                  |
| Vector search           | FR-15                  |
| Caching                 | FR-16                  |
| Monitoring              | FR-17                  |
| API docs                | FR-18                  |

---

If you'd like, I can also show you **how to convert this into a proper Software Requirements Specification (SRS)** with:

* **Use Case diagrams**
* **Sequence diagrams**
* **System architecture diagrams**
* **User stories for Agile development**.
