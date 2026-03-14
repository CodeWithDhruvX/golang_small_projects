Below is a **structured Technical Requirements Document (TRD) + Feature Requirements** for building your **Local AI Recruiter Reply Assistant** using the exact stack you provided (Go + Gin + Angular + Ollama + PGVector + RAG).

This is similar to what teams use in **real production projects**.

---

# 1. Project Overview

## Project Name

**Local AI Recruiter Assistant**

## Objective

Build a **fully local AI-powered application** that can:

* ingest recruiter emails
* understand recruiter requirements
* retrieve candidate information using **RAG**
* generate **professional email replies**
* track applications

All AI inference must run **locally using Ollama** with no external APIs.

---

# 2. System Goals

The system should:

1. Automatically process recruiter emails
2. Extract requested candidate information
3. Use **semantic search (PGVector)** to retrieve relevant candidate data
4. Generate **professional email replies using local LLM**
5. Maintain an **application tracking system**
6. Operate **fully offline / private**
7. Provide **monitoring and observability**

---

# 3. Functional Requirements

## 3.1 Authentication

### Features

* User registration
* Login/logout
* JWT authentication
* Secure session management

### API

```
POST /auth/register
POST /auth/login
POST /auth/logout
```

### Requirements

* Password hashing (bcrypt)
* JWT token expiry
* Role support (optional)

---

# 3.2 User Profile Management

Stores candidate information used by AI.

### Features

User can store:

* Name
* Experience
* Skills
* Current salary
* Expected salary
* Notice period
* Location
* Resume upload
* LinkedIn
* GitHub

### API

```
GET /profile
PUT /profile
POST /profile/resume
```

### Resume Processing

PDF resumes will be processed using:

**unipdf**

Extract:

```
experience
skills
projects
education
```

Then convert to embeddings.

---

# 3.3 Email Ingestion System

System should ingest recruiter emails.

### Sources

1. Gmail IMAP
2. Outlook IMAP
3. Manual email upload

### Features

System must:

* Fetch new emails
* Parse email body
* Detect recruiter messages
* Store email content

### API

```
POST /emails/import
GET /emails
GET /emails/{id}
```

### Email Processing Pipeline

```
Email received
     ↓
Text extraction
     ↓
Embedding generation
     ↓
Vector storage
     ↓
Recruiter detection
```

---

# 3.4 Recruiter Email Detection

Use **LLM classification**.

Input:

```
email body
```

Output:

```
Recruiter Email: true/false
```

Example detection prompt:

```
Determine if this email is a recruiter requesting job candidate details.
```

---

# 3.5 Requirement Extraction

System must extract requested candidate information.

Example recruiter request:

```
Please share:
- resume
- notice period
- expected CTC
```

Output JSON:

```
{
resume: true,
experience: true,
expected_ctc: true,
notice_period: true
}
```

---

# 3.6 RAG Knowledge Retrieval

The system will use **Retrieval Augmented Generation**.

### Knowledge Sources

* Resume
* Candidate profile
* Past applications
* Skills database

### Process

```
Email query
   ↓
Generate embedding
   ↓
Search PGVector
   ↓
Retrieve context
   ↓
Send to LLM
```

---

# 3.7 AI Email Reply Generator

Generate a **professional recruiter reply**.

Input:

```
Recruiter email
Retrieved context
User profile
```

Output:

```
Professional email response
```

Example output:

```
Hello,

Thank you for reaching out.

Please find my details below:

Experience: 4 years
Expected CTC: 12 LPA
Notice Period: 30 days

Resume attached.

Best regards
John
```

---

# 3.8 Application Tracker

Track all recruiter interactions.

### Features

Track:

* Company
* Role
* Recruiter
* Application status

Statuses:

```
Applied
Interview Scheduled
Offer
Rejected
```

### API

```
GET /applications
POST /applications
PUT /applications/{id}
```

---

# 3.9 Duplicate Detection

Prevent duplicate applications.

Check:

```
company + recruiter_email
```

If exists:

```
warning: already applied
```

---

# 3.10 AI Follow-up Generator

Generate follow-up emails.

Example:

```
Hi,

Just following up on my previous email regarding the opportunity.

Looking forward to hearing from you.
```

---

# 4. Non-Functional Requirements

## Performance

| Metric           | Requirement |
| ---------------- | ----------- |
| API response     | < 200ms     |
| AI generation    | < 5 seconds |
| Email processing | < 3 seconds |

---

## Security

System must implement:

* JWT authentication
* Password hashing
* HTTPS support
* Role based access (optional)

---

## Scalability

Architecture should support:

```
1000+ users
100k+ emails
```

---

## Privacy

All AI processing must remain:

```
local
offline
private
```

No external API calls allowed.

---

# 5. AI / ML Requirements

## Models (Ollama)

Required models:

```
llama3.1:8b
phi3
qwen2.5-coder
nomic-embed-text
```

---

## AI Tasks

| Task                | Model       |
| ------------------- | ----------- |
| Email generation    | Llama 3.1   |
| Text classification | Phi-3       |
| Reasoning           | Qwen2.5     |
| Embeddings          | Nomic Embed |

---

# 6. Database Requirements

## PostgreSQL + PGVector

Enable extension:

```
CREATE EXTENSION vector;
```

---

### Emails Table

```
emails
-----
id
user_id
subject
body
recruiter_email
created_at
embedding
```

---

### Applications Table

```
applications
-----
id
company
role
status
recruiter_email
created_at
```

---

### Documents Table

```
documents
---------
id
content
embedding
source
```

---

# 7. API Documentation

Swagger must auto-generate docs.

Example endpoint:

```
GET /emails
```

Response:

```
[
 {
   id:1,
   subject:"React Developer Opportunity",
   recruiter:"hr@company.com"
 }
]
```

---

# 8. Monitoring & Observability

## Prometheus Metrics

Metrics to expose:

```
api_requests_total
ai_generation_time
email_processing_time
vector_search_latency
```

---

## Grafana Dashboards

Visualize:

* API latency
* AI model performance
* CPU usage
* memory usage

---

# 9. Infrastructure Requirements

## Docker Services

```
frontend
backend
postgres
redis
ollama
prometheus
grafana
pgadmin
```

---

## Example docker-compose services

```
backend
frontend
ollama
postgres
redis
grafana
prometheus
```

---

# 10. Folder Structure

Backend:

```
backend/
├── cmd/
├── internal/
│   ├── api/
│   ├── auth/
│   ├── email/
│   ├── ai/
│   ├── rag/
│   ├── embeddings/
│   └── services/
├── pkg/
│   ├── database/
│   ├── cache/
│   └── logger/
```

---

# 11. Frontend Requirements

Angular 18 app must provide:

Pages:

```
Login
Dashboard
Recruiter Emails
AI Generated Replies
Applications
Profile
Settings
```

UI must support:

* dark mode
* responsive layout
* email preview
* reply editing

---

# 12. Hardware Requirements

Minimum:

```
16GB RAM
```

Recommended:

```
32GB RAM
GPU optional
```

---

# 13. Future Features

Planned enhancements:

### AI Resume Improvement

Suggest improvements.

---

### Job Description Matching

Compare:

```
job description
vs
resume
```

---

### AI Interview Preparation

Generate interview questions.

---

### Chrome Extension

Reply directly from Gmail.

---

# 14. MVP Scope

Initial MVP should include:

✔ authentication
✔ profile management
✔ resume upload
✔ email ingestion
✔ AI reply generation
✔ application tracker

---

# Final Assessment

Your stack is **perfect for building this system**.

It provides:

* local AI
* vector search
* scalable backend
* enterprise monitoring
* containerized infrastructure

This is **very close to production-grade AI architecture**.


