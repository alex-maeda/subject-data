# Attributes, Records, and Evidence — Data Model Analysis

**Date:** 2026-05-03
**Source:** `sovraai/data` repo (Roman's ingestion pipeline) + team discussions
**Purpose:** Define how evidence, records, and attributes relate across the Data Ingest Pipeline and Subject Data Service

---

## 0. Glossary — Terms, Definitions, and Relationships

### Subject

A **Subject** is the central identity entity — a person that Sovra builds a profile about. Everything else (records, attributes, evidence) attaches to a Subject.

- Has a stable UUID
- Has a `name` (human-readable identifier)
- Is **mutable**: new records and attributes accumulate over time
- Initially created either by Roman's clustering logic or by the subject-search orchestrator
- May have zero records at creation (empty shell) or be created from clustering many records together

**Relationships:**
- A Subject has many Records (linked via clustering)
- A Subject has many Attributes (written by multiple producers)

---

### Record

A **Record** is one row of normalized external data from a single source. It represents a raw observation about a person.

- Has a stable UUID + `subject_id` foreign key (nullable until clustering links it)
- Is **immutable** once written — never updated, only superseded by newer records
- Contains normalized field values (name, email, phone, address, etc.)
- Contains metadata (dataset, version, batch, source file, row_id, timestamp)
- Includes `_record_status`: eligible / ineligible / rejected
- Ingested via batch import (parquet from S3)

**Relationships:**
- A Record belongs to zero or one Subject (unlinked until clustering)
- A Record contains multiple Evidence items (derived facts)
- Multiple Records form the basis for computing Attributes

---

### Attribute

An **Attribute** is one piece of computed/aggregated information about a Subject. It's a first-class object with its own lifecycle.

- Has a stable UUID + `subject_id` foreign key
- Is **mutable**: updated as new evidence arrives
- Has a `type` (e.g., "demographics", "relationship_status", "behavioral")
- Has a `schema_version` for payload evolution
- Has a `confidence` score (numeric or interval)
- Has a `writer` identifier (which system wrote it)
- Links back to source evidence via `evidence_refs`
- Written by **multiple independent producers**: Roman's pipeline, S-VIP, S-BAM, subject-search

**Relationships:**
- An Attribute belongs to one Subject
- An Attribute references zero or more Records/Evidence as its source
- An Attribute is typed — its `payload` schema depends on its `type`

---

### Evidence

**Evidence** is an atomic derived fact — one normalized value computed from one source field within one record. It is the provenance link between raw data and derived attributes.

- Is **NOT a standalone table** in SDS — it's a logical concept
- Lives either as a nested array inside a Record, or as `evidence_refs` on an Attribute
- Represents: "from field X in record Y, we derived attribute Z with value V"
- Three fields: `source_field`, `normalized_attribute`, `normalized_value`
- Currently ships **empty** from the pipeline (TASK-047 pending)

**Relationships:**
- Evidence belongs to one Record
- Evidence supports one or more Attributes (via `evidence_refs`)
- Evidence is the audit trail: why does this Attribute have this value?

---

### Envelope

The **Envelope** is the wrapper structure for an Attribute. It separates the metadata (who wrote it, when, confidence) from the payload (the actual data). This pattern allows:
- Generic CRUD operations on any attribute type
- Schema evolution without changing the envelope
- Multiple writers with different payloads for the same type
- Querying attributes by type/writer/confidence without parsing payload

**Envelope fields:**
`id`, `subject_id`, `type`, `schema_version`, `payload`, `confidence`, `writer`, `evidence_refs`, `created_at`, `updated_at`

**Relationships:**
- An Envelope IS an Attribute (same thing, different perspective)
- The `payload` field contains the type-specific data
- The `type` + `schema_version` determine how to interpret `payload`

---

### Ingest Job

An **Ingest Job** tracks the lifecycle of a batch import from S3 into SDS.

- Created when pipeline calls `POST /v1/ingest-jobs`
- Idempotent by `manifest_sha256`
- Tracks status: `queued` → `running` → `succeeded` / `partial` / `failed`
- Tracks counts: rows_total, rows_processed, records_upserted, errors
- Writes `ingestion.json` back to S3 as completion signal

**Relationships:**
- An Ingest Job processes one manifest (one batch + dataset + version)
- An Ingest Job creates many Records
- The pipeline polls the job for status

---

### Relationship Diagram

```
┌─────────────────────────────────────────────────────────┐
│                     SUBJECT                              │
│  id, name, created_at, updated_at                       │
└──────────┬──────────────────────────────┬───────────────┘
           │ 1:N                          │ 1:N
           ▼                              ▼
┌─────────────────────┐      ┌───────────────────────────┐
│       RECORD        │      │     ATTRIBUTE (Envelope)   │
│  id                 │      │  id                        │
│  subject_id (FK)    │      │  subject_id (FK)           │
│  dataset            │      │  type                      │
│  dataset_version    │      │  schema_version            │
│  batch              │      │  payload (JSON)            │
│  platform_data      │◄─────│  confidence                │
│  metadata           │ ref  │  writer                    │
│  _record_status     │      │  evidence_refs []          │
│  created_at         │      │  created_at, updated_at    │
└────────┬────────────┘      └───────────────────────────┘
         │ contains
         ▼
┌─────────────────────┐
│     EVIDENCE        │
│  (nested in record) │
│  source_field       │
│  normalized_attr    │
│  normalized_value   │
└─────────────────────┘

┌─────────────────────────────────────────────────────────┐
│                   INGEST JOB                             │
│  id, import_id, manifest_uri, manifest_sha256           │
│  dataset, dataset_version, batch                        │
│  status, rows_total, rows_processed, errors[]           │
│  started_at, completed_at                               │
└─────────────────────────────────────────────────────────┘
```

---

## 1. Conceptual Hierarchy

```
Subject (identity entity in SDS)
  │
  ├── Records (immutable, from normalized parquet)
  │     │
  │     ├── Evidence (atomic derived facts per field)
  │     │     - source_field: which normalized column
  │     │     - normalized_attribute: derived attribute name
  │     │     - normalized_value: the value itself
  │     │
  │     └── Raw field values (direct from source)
  │
  └── Attributes (mutable, first-class objects)
        - Aggregated from multiple records
        - Written by multiple producers (pipeline, S-VIP, S-BAM)
        - Versioned, confidence-bearing
```

---

## 2. Records — What They Are

A **Record** is one row from a normalized Parquet file. It represents a single observation from a single data source about a person.

### Parquet Schema (normalized stage output)

| Column | Type | Example | Category |
|--------|------|---------|----------|
| `first_name` | string | `"JOHN"` | indexed |
| `last_name` | string | `"SMITH"` | indexed |
| `full_name` | string | `"JOHN SMITH"` | indexed |
| `emails` | list\<string\> | `["john@gmail.com"]` | indexed |
| `phones` | list\<string\> | `["5551234567"]` | indexed |
| `addresses` | list\<string\> | `["123 Main St, Houston TX"]` | indexed |
| `city` | string | `"HOUSTON"` | indexed |
| `state` | string | `"TX"` | indexed |
| `country` | string | `"US"` | indexed |
| `gender` | string | `"male"` | stored |
| `age` | string | `"35"` | stored |
| `_row_id` | int64 | `12345` | control |
| `_record_status` | string | `"eligible"` | control |
| `_batch` | string | `"raw_2026_05"` | control |
| `_dataset` | string | `"peopledatalabs"` | control |
| `_dataset_version` | string | `"v1"` | control |
| `_source_raw_file` | string | `"/data/raw/pdl.csv.zst"` | control |
| `_source_versioned_file` | string | `"/data/versioned/..."` | control |
| `_parsed_at` | string (ISO) | `"2026-05-01T..."` | control |

### Record Properties

- **Immutable** once written to S3
- **Source-specific**: one record = one row from one dataset version
- **Eligibility-filtered**: only `eligible` records (≥2 distinct PII values) should be ingested
- **Lineage-tracked**: every record traces back to source file + row number

### Record Status Values

| Status | Meaning |
|--------|---------|
| `eligible` | ≥2 distinct useful PII values — qualifies for ingestion |
| `ineligible` | Has PII but below threshold — skip for now |
| `rejected` | Missing critical data (e.g., no `_row_id`) — discard |

---

## 3. Evidence — What It Is

**Evidence** is not a standalone object in the pipeline. It's a logical concept: one atomic fact derived from one source field in one record.

### Evidence Model (from Roman's design)

```
evidence = [
  {
    "source_field": "full_name",
    "normalized_attribute": "first_name",
    "normalized_value": "JOHN"
  },
  {
    "source_field": "full_name",
    "normalized_attribute": "last_name",
    "normalized_value": "SMITH"
  },
  {
    "source_field": "phones",
    "normalized_attribute": "country",
    "normalized_value": "US"
  }
]
```

### How Evidence Is Produced

Evidence comes from **transforms** during normalization:

| Transform | Source Field | Derived Attributes |
|-----------|-------------|-------------------|
| `split_name` | `full_name` | `first_name`, `last_name`, `middle_name` |
| `parse_address` | `addresses` | `city`, `state`, `country` |
| `phone_country` | `phones` | `country` (from phone code) |

A fact is **evidence** when it didn't exist in the parsed record — it was computed from another field. Direct field values (email from email column) are just normalized values, not evidence.

### Evidence Rules

1. **Non-redundancy**: don't emit evidence if the attribute already exists as a column
2. **Record-status gate**: no evidence from `rejected` records
3. **Deduplication**: triple `(source_field, normalized_attribute, normalized_value)` is unique per record
4. **Multiple sources**: same `source_field` can produce multiple evidence items (one source → many observations)

### Current State

The evidence column ships **empty** in current normalized output. The emitter framework exists but actual derivation rules (TASK-047) are in progress. This means:
- Records have normalized field values ✓
- Records do NOT yet have nested evidence arrays
- SDS should design for evidence but expect it empty initially

---

## 4. Attributes — What They Are

An **Attribute** is a first-class object attached to a Subject in SDS. It represents one piece of computed/aggregated information.

### Attribute Envelope (from Maeda's design, approved by Marple)

```json
{
  "id": "attr_uuid",
  "subject_id": "subj_uuid",
  "type": "relationship_status",
  "schema_version": "v1",
  "payload": {
    "is_married": false,
    "is_in_relationship": true
  },
  "confidence": 0.8,
  "writer": "subject-search",
  "evidence_refs": ["record:pdl_12345", "record:breach_67890"],
  "created_at": "2026-05-03T...",
  "updated_at": "2026-05-03T..."
}
```

### Attribute Properties

- **Mutable**: updated as new evidence arrives
- **Multi-writer**: Roman's pipeline, S-VIP, S-BAM, subject-search can all write
- **Typed**: each attribute has a `type` that determines the `payload` schema
- **Confidence-bearing**: numeric confidence (or interval)
- **Evidence-linked**: `evidence_refs` point back to source records

### 9 Attribute Types Needed for July Beta (from Patrick's BFF doc)

| Type | Writer | Payload Shape |
|------|--------|--------------|
| `demographics` | IDS (Roman) | `{first_name, last_name, age, gender, location}` |
| `relationship_status` | S-VIP / Analysis | `{is_married, is_in_relationship}` |
| `public_records` | IDS / scraping | `{criminal: [...], civil: [...], marriage: [...]}` |
| `social_media_footprint` | Data Enricher | `{profiles: [{platform, url, handle}]}` |
| `in_the_news` | Data Enricher | `{articles: [{title, url, date, source}]}` |
| `data_leaks` | IDS (breach data) | `{breaches: [{source, date, categories: []}]}` |
| `behavioral` | S-BAM | `{ratings: [{trait, score, confidence, evidence}]}` |
| `timeline` | Analysis | `{events: [{date, type, description}]}` |
| `geographic_footprint` | IDS / enrichment | `{locations: [{address, lat, lng, period}]}` |

### Attribute vs Evidence vs Record

| Concept | Mutability | Scope | Source |
|---------|-----------|-------|--------|
| Record | Immutable | One row, one source | Data Ingest pipeline |
| Evidence | Immutable | One derived fact from one record | Normalization transforms |
| Attribute | Mutable | Aggregated from N records/evidence | Multiple writers |

---

## 5. How They Connect (End-to-End Flow)

```
[Data Ingest Pipeline]
    Raw CSV → Parsed Parquet → Normalized Parquet
    Each row = 1 Record with N field values
    Transforms produce Evidence (derived facts)

        ↓ S3: published/<batch>/<dataset>_<version>/

[Subject Data Service - Batch Ingest]
    POST /v1/ingest-jobs {manifest_uri, dataset, batch}
    → Read parquet from S3
    → Store each eligible row as a Record
    → Record is immutable, linked to dataset/batch/row_id

        ↓ Records stored in SDS DB

[Clustering / Entity Resolution]
    Read Records from SDS
    → Match records to existing Subjects (by shared PII)
    → Or create new Subjects
    → Link Records ↔ Subjects

        ↓ Subjects linked to Records

[Attribute Writers]
    Read Subject's Records
    → Aggregate field values across records
    → Apply logic (simple copy, confidence scoring, AI analysis)
    → Write Attribute to SDS

        ↓ Attributes stored, frontend can read

[Frontend / BFF]
    GET /v1/subjects/{id}/attributes
    → Render subject detail page
```

---

## 6. Field Registry (Canonical Vocabulary)

From `config/fields-registry.yaml`:

### Indexed Fields (searchable)
```
first_name, last_name, full_name, middle_name,
emails, phones, addresses,
city, state, country, postal_code,
linkedin_url, facebook_url, twitter_url,
company_name, job_title
```

### Stored Fields (non-searchable)
```
gender, age, birth_date, income_code,
carrier, education_level
```

### PII Fields (used for eligibility)
```
first_name, last_name, emails, phones, addresses,
city, state, postal_code
```

### Sensitive Fields
```
password, ssn
```

---

## 7. Normalizer Functions

| Normalizer | Logic | Example |
|-----------|-------|---------|
| `name` | uppercase + collapse whitespace | `"john smith"` → `"JOHN SMITH"` |
| `uppercase` | simple uppercase | `"john"` → `"JOHN"` |
| `lowercase` | simple lowercase | `"John@Gmail.COM"` → `"john@gmail.com"` |
| `phone` | extract digits, strip leading 1, validate ≥7 | `"+1 (555) 123-4567"` → `"5551234567"` |
| `state` | uppercase → 2-char abbreviation | `"texas"` → `"TX"` |
| `trim` | collapse whitespace | `"  hello   world "` → `"hello world"` |

---

## 8. Manifest & Ingestion Contract (Roman's Proposal)

### S3 Path Layout
```
<bucket>/published/<batch>/<dataset>_<dataset_version>/
    part-00001.parquet
    part-00002.parquet
    ...
    manifest.json        ← written by pipeline (commit marker)
    ingestion.json       ← written by SDS (completion marker)
```

### manifest.json
```json
{
  "schema_version": 1,
  "dataset": "peopledatalabs",
  "dataset_version": "v1",
  "batch": "raw_2026_05",
  "created_at": "2026-05-03T...",
  "parts": [
    {
      "key": "part-00001.parquet",
      "sha256": "abc123...",
      "size_bytes": 52428800,
      "row_count": 500000,
      "source_versioned_file": "...",
      "source_raw_file": "..."
    }
  ]
}
```

### ingestion.json (SDS writes back)
```json
{
  "schema_version": 1,
  "job_id": "uuid",
  "dataset": "peopledatalabs",
  "dataset_version": "v1",
  "batch": "raw_2026_05",
  "manifest_sha256": "hex",
  "status": "succeeded",
  "started_at": "2026-05-03T...",
  "completed_at": "2026-05-03T...",
  "rows_total": 500000,
  "rows_processed": 499000,
  "subjects_upserted": 0,
  "records_upserted": 499000,
  "errors": []
}
```

### POST /v1/ingest-jobs
```json
{
  "dataset": "peopledatalabs",
  "dataset_version": "v1",
  "batch": "raw_2026_05",
  "manifest_uri": "s3://subject-data-ingestion-beta/published/raw_2026_05/peopledatalabs_v1/manifest.json",
  "manifest_sha256": "hex"
}
```

**Idempotency key**: `manifest_sha256`

---

## 9. Design Decisions Summary

| Decision | Rationale |
|----------|-----------|
| Records are immutable | Enables replay, audit, multi-version analysis |
| Subjects are mutable | Accumulate evidence over time |
| Attributes are first-class | Multiple writers, independent lifecycle |
| Evidence is nested in records | Provenance without separate table |
| S3 is the handoff point | Decouples pipeline from SDS internals |
| Idempotency by manifest hash | Same content = same job, safe retries |
| No DeleteObject | Immutability enforced at IAM level |
| ingestion.json writeback | Pipeline's only completion signal |
| Eligibility threshold (≥2 PII) | Avoid ingesting useless/garbage records |

---

## 10. Open Questions

1. **Evidence column is empty today** — do we store placeholder or skip until TASK-047 lands?
2. **Record-to-Subject linking** — happens during clustering, not during ingest. SDS stores unlinked records initially.
3. **Attribute derivation timing** — simple attributes (name, email) could be written at ingest time; complex ones (behavioral) come later from S-BAM.
4. **Schema versioning** — how to handle when `peopledatalabs_v2` has different fields than v1?
5. **Cross-record deduplication** — not evidence's job; happens at clustering stage.

---

## 11. Suggested Schema & API Design

### 11.1 Subject

#### Database Schema

```sql
CREATE TABLE subjects (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    namespace   TEXT NOT NULL DEFAULT 'internal',  -- internal | composite | session
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_subjects_namespace ON subjects(namespace);
```

#### API Endpoints

```http
GET    /v1/subjects                          → []SubjectSummary
POST   /v1/subjects                          → {id, name}
GET    /v1/subjects/{subjectID}              → Subject (full, with linked record count)
PUT    /v1/subjects/{subjectID}              → {id, status: "updated"}
DELETE /v1/subjects/{subjectID}              → {id, status: "deleted"}
```

#### Response Shape

```json
// GET /v1/subjects/{subjectID}
{
  "id": "subj_e0b0fa6e-85b8-4f30-83e1-12d23cfdd1ad",
  "name": "John Smith",
  "namespace": "internal",
  "record_count": 47,
  "attribute_count": 6,
  "created_at": "2026-04-28T12:00:00.000Z",
  "updated_at": "2026-05-01T09:30:00.000Z"
}
```

---

### 11.2 Attributes (Envelope Pattern)

#### Database Schema

```sql
CREATE TABLE attributes (
    id              TEXT PRIMARY KEY,
    subject_id      TEXT NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    type            TEXT NOT NULL,           -- e.g., "demographics", "relationship_status"
    schema_version  TEXT NOT NULL DEFAULT 'v1',
    payload         JSONB NOT NULL,          -- type-specific data
    confidence      REAL,                    -- 0.0 to 1.0
    writer          TEXT NOT NULL,           -- e.g., "subject-search", "s-bam", "ids-pipeline"
    evidence_refs   JSONB DEFAULT '[]',      -- ["record:<id>", "record:<id>"]
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_attributes_subject ON attributes(subject_id);
CREATE INDEX idx_attributes_type ON attributes(subject_id, type);
CREATE UNIQUE INDEX idx_attributes_subject_type_writer ON attributes(subject_id, type, writer);
```

Note: `UNIQUE(subject_id, type, writer)` ensures one attribute per type per writer per subject. If S-VIP and S-BAM both write "behavioral", they get separate rows.

#### API Endpoints

```http
GET    /v1/subjects/{subjectID}/attributes                          → []Attribute
GET    /v1/subjects/{subjectID}/attributes?type=demographics,relationship_status  → []Attribute (filtered)
POST   /v1/subjects/{subjectID}/attributes                          → {id, type, status: "created"}
PUT    /v1/subjects/{subjectID}/attributes/{attributeID}            → {id, status: "updated"}
DELETE /v1/subjects/{subjectID}/attributes/{attributeID}            → {id, status: "deleted"}
```

#### Request Shape (POST / PUT)

```json
// POST /v1/subjects/{subjectID}/attributes
{
  "type": "demographics",
  "schema_version": "v1",
  "payload": {
    "first_name": "John",
    "last_name": "Smith",
    "age": 35,
    "gender": "male",
    "location": {
      "city": "Houston",
      "state": "TX",
      "country": "US"
    }
  },
  "confidence": 0.85,
  "writer": "ids-pipeline",
  "evidence_refs": [
    "record:rec_abc123",
    "record:rec_def456"
  ]
}
```

#### Response Shape

```json
// GET /v1/subjects/{subjectID}/attributes?type=relationship_status
[
  {
    "id": "attr_789ghi",
    "subject_id": "subj_e0b0fa6e",
    "type": "relationship_status",
    "schema_version": "v1",
    "payload": {
      "is_married": false,
      "is_in_relationship": true
    },
    "confidence": 0.72,
    "writer": "s-vip",
    "evidence_refs": ["record:rec_xyz789"],
    "created_at": "2026-05-02T14:00:00.000Z",
    "updated_at": "2026-05-02T14:00:00.000Z"
  }
]
```

#### Payload Schemas by Type

```json
// type: "demographics"
{
  "first_name": "string",
  "last_name": "string",
  "middle_name": "string | null",
  "age": "number | null",
  "gender": "string | null",
  "location": {
    "city": "string | null",
    "state": "string | null",
    "country": "string | null"
  }
}

// type: "relationship_status"
{
  "is_married": "boolean | null",
  "is_in_relationship": "boolean | null"
}

// type: "social_media_footprint"
{
  "profiles": [
    {
      "platform": "linkedin | twitter | facebook | instagram | tiktok",
      "url": "string",
      "handle": "string | null",
      "verified": "boolean"
    }
  ]
}

// type: "data_leaks"
{
  "breach_count": "number",
  "breaches": [
    {
      "source": "string",
      "date": "string (ISO date) | null",
      "categories": ["email", "password", "address", "ssn-partial"]
    }
  ]
}

// type: "in_the_news"
{
  "articles": [
    {
      "title": "string",
      "url": "string",
      "date": "string (ISO date) | null",
      "source": "string",
      "sentiment": "positive | negative | neutral | null"
    }
  ]
}

// type: "public_records"
{
  "criminal": [{"type": "string", "date": "string", "jurisdiction": "string"}],
  "civil": [{"type": "string", "date": "string", "jurisdiction": "string"}],
  "marriage": [{"date": "string", "spouse": "string | null", "jurisdiction": "string"}]
}

// type: "behavioral"
{
  "ratings": [
    {
      "category": "Normal Personality | Character | Values | ...",
      "trait": "string",
      "score": "1-5",
      "confidence": "Low | Medium | High",
      "evidence_summary": "string"
    }
  ]
}

// type: "geographic_footprint"
{
  "locations": [
    {
      "address": "string",
      "city": "string",
      "state": "string",
      "lat": "number | null",
      "lng": "number | null",
      "period": "string | null"
    }
  ]
}

// type: "timeline"
{
  "events": [
    {
      "date": "string (ISO date)",
      "type": "residence | employment | education | legal | relationship",
      "description": "string"
    }
  ]
}
```

---

### 11.3 Records

#### Database Schema

```sql
CREATE TABLE records (
    id                  TEXT PRIMARY KEY,
    subject_id          TEXT REFERENCES subjects(id),  -- nullable until clustering
    dataset             TEXT NOT NULL,
    dataset_version     TEXT NOT NULL,
    batch               TEXT NOT NULL,
    source_row_id       BIGINT NOT NULL,
    record_status       TEXT NOT NULL DEFAULT 'eligible',  -- eligible | ineligible | rejected
    platform_data       JSONB NOT NULL,       -- normalized field values
    metadata            JSONB NOT NULL,       -- source file, timestamps, lineage
    evidence            JSONB DEFAULT '[]',   -- nested evidence array (empty until TASK-047)
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
    -- no updated_at: records are immutable
);

CREATE INDEX idx_records_subject ON records(subject_id);
CREATE INDEX idx_records_dataset_batch ON records(dataset, batch);
CREATE UNIQUE INDEX idx_records_dedup ON records(dataset, dataset_version, batch, source_row_id);
```

Note: `UNIQUE(dataset, dataset_version, batch, source_row_id)` prevents duplicate ingest of the same row.

#### API Endpoints

```http
GET    /v1/records                                    → []Record (filter: ?subject_id=X&dataset=Y)
POST   /v1/records                                    → {id}
GET    /v1/records/{recordID}                         → Record
GET    /v1/subjects/{subjectID}/records               → []Record
PUT    /v1/subjects/{subjectID}/records/{recordID}    → {status: "linked"}  (link record to subject)
```

#### Response Shape

```json
// GET /v1/records/{recordID}
{
  "id": "rec_abc123",
  "subject_id": "subj_e0b0fa6e",
  "dataset": "peopledatalabs",
  "dataset_version": "v1",
  "batch": "raw_2026_05",
  "source_row_id": 12345,
  "record_status": "eligible",
  "platform_data": {
    "first_name": "JOHN",
    "last_name": "SMITH",
    "emails": ["john.smith@gmail.com"],
    "phones": ["5551234567"],
    "addresses": ["123 Main St, Houston TX 77001"],
    "city": "HOUSTON",
    "state": "TX",
    "country": "US",
    "gender": "male"
  },
  "metadata": {
    "source_raw_file": "/data/raw/pdl_global.csv.zst",
    "source_versioned_file": "/data/versioned/peopledatalabs_v1/pdl_global.csv.zst",
    "parsed_at": "2026-05-01T10:00:00.000Z"
  },
  "evidence": [],
  "created_at": "2026-05-02T14:30:00.000Z"
}
```

#### Link Record to Subject (used by Clustering)

```http
PUT /v1/subjects/{subjectID}/records/{recordID}
```

```json
// Request body (optional metadata)
{
  "linker": "clustering-v1",
  "confidence": 0.95,
  "reason": "matched on email + phone"
}
```

```json
// Response
{
  "record_id": "rec_abc123",
  "subject_id": "subj_e0b0fa6e",
  "status": "linked"
}
```

---

### 11.4 Evidence

Evidence is **not a standalone API resource**. It's accessed in two ways:

#### A) Nested inside a Record

```json
// Part of GET /v1/records/{id} response
{
  "evidence": [
    {
      "source_field": "full_name",
      "normalized_attribute": "first_name",
      "normalized_value": "JOHN"
    },
    {
      "source_field": "full_name",
      "normalized_attribute": "last_name",
      "normalized_value": "SMITH"
    },
    {
      "source_field": "phones",
      "normalized_attribute": "country",
      "normalized_value": "US"
    }
  ]
}
```

#### B) Referenced from Attributes

```json
// Part of Attribute response
{
  "evidence_refs": [
    "record:rec_abc123",
    "record:rec_def456",
    "record:rec_abc123#evidence[0]"   // optional: specific evidence item
  ]
}
```

#### Evidence Schema (nested in record JSON)

```json
{
  "source_field": "string",           // which normalized column produced this
  "normalized_attribute": "string",   // the derived attribute name
  "normalized_value": "string"        // the value
}
```

#### When Evidence Becomes Useful

- **Now**: evidence array is empty (pipeline TASK-047 pending)
- **Soon**: pipeline emits evidence for derived facts (split_name, parse_address, phone→country)
- **Later**: clustering uses evidence to explain why a record was linked to a subject
- **Eventually**: S-BAM/S-VIP cite specific evidence items when writing attributes

---

### 11.5 Envelope — The Generic Attribute Container

The Envelope is not a separate entity — it IS the Attribute. The term "envelope" describes the **design pattern**: metadata wraps an opaque typed payload.

#### Why This Pattern

```
┌─────────────────────────────── Envelope ───────────────────────────────┐
│                                                                         │
│  ┌─── Generic (same for all types) ───┐  ┌─── Type-specific ────────┐ │
│  │ id                                  │  │ payload: {               │ │
│  │ subject_id                          │  │   first_name: "John",    │ │
│  │ type: "demographics"                │  │   last_name: "Smith",    │ │
│  │ schema_version: "v1"                │  │   age: 35                │ │
│  │ confidence: 0.85                    │  │ }                        │ │
│  │ writer: "ids-pipeline"              │  │                          │ │
│  │ evidence_refs: [...]                │  │                          │ │
│  │ created_at, updated_at              │  │                          │ │
│  └─────────────────────────────────────┘  └──────────────────────────┘ │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

Benefits:
- **Generic CRUD**: one API handles all attribute types
- **Schema evolution**: bump `schema_version`, old payloads still readable
- **Multi-writer**: same `type` can have multiple writers (each gets own row)
- **Queryable metadata**: filter by type, writer, confidence without parsing payload
- **Validation deferred**: payload schema checked by type-specific validators, not the DB layer

#### Envelope API Contract

```http
# Write any attribute type with the same endpoint
POST /v1/subjects/{subjectID}/attributes
Content-Type: application/json

{
  "type": "<attribute_type>",
  "schema_version": "v1",
  "payload": { ... },            ← opaque to the envelope layer
  "confidence": 0.85,
  "writer": "<system_name>",
  "evidence_refs": ["record:<id>", ...]
}
```

```http
# Read with optional type filter
GET /v1/subjects/{subjectID}/attributes?type=demographics,data_leaks

# Response: array of envelopes
[
  { "id": "...", "type": "demographics", "payload": {...}, ... },
  { "id": "...", "type": "data_leaks", "payload": {...}, ... }
]
```

#### Upsert Semantics

The unique constraint `(subject_id, type, writer)` means:
- First POST → creates attribute (201)
- Subsequent POST with same subject + type + writer → updates payload (200, upsert)
- Different writer for same type → creates new row (allows S-VIP and S-BAM to both write "behavioral")

---

### 11.6 Ingest Jobs

#### Database Schema

```sql
CREATE TABLE import_jobs (
    id              TEXT PRIMARY KEY,
    dataset         TEXT NOT NULL,
    dataset_version TEXT NOT NULL,
    batch           TEXT NOT NULL,
    manifest_uri    TEXT NOT NULL,
    manifest_sha256 TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'queued',  -- queued|running|succeeded|partial|failed
    rows_total      BIGINT DEFAULT 0,
    rows_processed  BIGINT DEFAULT 0,
    records_upserted BIGINT DEFAULT 0,
    errors          JSONB DEFAULT '[]',
    started_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_import_jobs_idempotent ON import_jobs(manifest_sha256);
```

#### API Endpoints

```http
POST   /v1/ingest-jobs                    → {job_id, status: "queued"}
GET    /v1/ingest-jobs/{jobID}            → IngestJob (full status)
```

#### Request Shape

```json
// POST /v1/ingest-jobs
{
  "dataset": "peopledatalabs",
  "dataset_version": "v1",
  "batch": "raw_2026_05",
  "manifest_uri": "s3://subject-data-ingestion-beta/published/raw_2026_05/peopledatalabs_v1/manifest.json",
  "manifest_sha256": "a1b2c3d4e5f6..."
}
```

#### Idempotency Rules

| Existing state for `manifest_sha256` | Behavior |
|--------------------------------------|----------|
| No job exists | Create new job, return `queued` |
| Job exists, status `queued` or `running` | Return existing job (no-op) |
| Job exists, status `succeeded` | Return existing job (no-op) |
| Job exists, status `failed` or `partial` | Create new attempt (new job_id) |

#### Response Shape

```json
// GET /v1/ingest-jobs/{jobID}
{
  "job_id": "job_uuid",
  "dataset": "peopledatalabs",
  "dataset_version": "v1",
  "batch": "raw_2026_05",
  "manifest_uri": "s3://...",
  "manifest_sha256": "a1b2c3d4e5f6...",
  "status": "running",
  "rows_total": 500000,
  "rows_processed": 250000,
  "records_upserted": 249500,
  "errors": [
    {
      "code": "invalid_record",
      "message": "missing required field: emails",
      "object_uri": "s3://.../part-00003.parquet"
    }
  ],
  "started_at": "2026-05-03T10:00:00.000Z",
  "completed_at": null,
  "created_at": "2026-05-03T09:59:55.000Z"
}
```

---

## 12. Complete API Surface Summary

```
SUBJECTS
  GET    /v1/subjects
  POST   /v1/subjects
  GET    /v1/subjects/{id}
  PUT    /v1/subjects/{id}
  DELETE /v1/subjects/{id}

RECORDS
  GET    /v1/records?subject_id=X&dataset=Y
  POST   /v1/records
  GET    /v1/records/{id}
  GET    /v1/subjects/{id}/records
  PUT    /v1/subjects/{id}/records/{recordID}         ← link record to subject

ATTRIBUTES (Envelope)
  GET    /v1/subjects/{id}/attributes?type=X,Y
  POST   /v1/subjects/{id}/attributes                 ← upsert by (subject, type, writer)
  PUT    /v1/subjects/{id}/attributes/{attrID}
  DELETE /v1/subjects/{id}/attributes/{attrID}

INGEST JOBS
  POST   /v1/ingest-jobs                              ← register batch import
  GET    /v1/ingest-jobs/{jobID}                      ← poll status

HEALTH
  GET    /healthz
  GET    /readyz
```
