# SOVRA Tech Standup — Meeting Notes

## Overview

- Weekly tech sync led by Jimmy
- Focus: team updates, post-offsite alignment, infrastructure direction
- New team members introduced

---

## Team Updates

### Jimmy (CTO)

- Offsite completed → consolidating into **PRD (Product Requirements Doc)**
- PRD will define:
  - MVP scope (July)
  - Launch scope (Sept/Oct)
- Plan:
  - Share PRD → team review → lock scope
  - Break into concrete tasks + estimates
- Exploring **Linear** for task tracking
  - Goal: reduce overhead using Claude integration
- Mobile direction update:
  - Likely need **native mobile app** (not just responsive web)
  - Timeline: Sept/Oct (not immediate)

---

### Roman (Data / Ingestion)

- Progress:
  - Nearly completed:
    - Parsing stage
    - Normalization stage
  - Running PDL dataset through pipeline
- Performance:
  - ~1 hour per 50GB (parse)
  - ~10–20 hours (normalize, no parallelism yet)
- Next steps:
  - Finish normalization
  - Ingest into Elasticsearch
  - Add parallelism
- Introduced concept:
  - **“Evidence”**
    - Derived data (e.g. country from phone, name parsing)
- Open questions:
  - Index design in Elasticsearch
  - Data modeling alignment

---

### Alex Marple (Backend / AI / Infra)

- Key focus:
  - **Chatbot evaluation system**
- Building:
  - Metrics for response quality:
    - Grounding
    - Correctness
    - Tool usage
    - Clarity
- Using:
  - **Braintrust** for evaluation
- Challenge:
  - No clear definition of “good” output
- Approach:
  - Define baseline manually
  - Iterate later with feedback
- Note:
  - Evaluation > prompt tuning (long-term importance)

---

### Grant (Scraping / Data)

- Dataset stats:
  - 38M+ lines
  - Texas mentions:
    - 961K (Texas)
    - 337K (Houston)
- Ongoing:
  - Court record scraping
- Need:
  - Expansion to more counties / regions

---

### Michael (Data Vendors)

- Building spreadsheet of:
  - Data vendors (PDL, PIPL, Trestle, Data Axle, etc.)
- Goal:
  - Evaluate buying data vs scraping
- Next steps:
  - Share doc
  - Prioritize vendors
  - Begin outreach after review

---

### Alex Maeda (You)

- Week 1 onboarding:
  - Joined offsite sessions
  - Reviewed:
    - Architecture repo
    - SDS repo
    - Data schemas (subjects, records)
- Observations:
  - SDS = central data layer
- Current work:
  - Planning **attribute layer**
- Next steps:
  - Sync with:
    - Roman
    - Alex Marple
  - Define:
    - Data ingestion → SDS structure

**Clarification from Alex Marple:**

- No separate “data engineering team”
- Core group:
  - Roman
  - Alex Marple
  - Alex Maeda

---

### Minh (Compliance)

- Role:
  - Chief Compliance Officer
- Focus:
  - Privacy policy requirements
- Next steps:
  - Sync with Jimmy on technical implications

---

## Key Technical Themes

### 1. SDS (Subject Data Service)

- Acts as **central system of record**
- Must handle:
  - Subjects
  - Records
  - Attributes (planned)
- Open problem:
  - How ingestion pipeline outputs map into SDS

---

### 2. Data Pipeline → SDS Integration

- Roman building ingestion → Elasticsearch
- Missing clarity:
  - How data flows into SDS
  - Whether SDS and Elasticsearch share responsibilities
- Action:
  - Alex Marple + Roman to align
  - Alex Maeda involved in defining structure

---

### 3. Terminology Conflict

- “Evidence” meaning mismatch:
  - Roman: derived data
  - Existing schema: already has “evidence”
- Risk:
  - Schema confusion
- Needs alignment

---

### 4. Identity Resolution (Upcoming)

- Future step after ingestion:
  - Combine records → subjects
- Critical for MVP

---

### 5. Chat / Analysis Layer

- Goal:
  - Not just data retrieval
  - Produce **interpreted outputs**
- Needs:
  - Evaluation framework (in progress)

---

### 6. Data Strategy

- Hybrid approach:
  - Scraping (Grant)
  - Vendor data (Michael)
- Goal:
  - Texas + major metro coverage by Sept

---

## Process / Org

- Very **senior team**
- Low process, high autonomy
- Communication:
  - Slack (tech channel)
- Meetings:
  - Weekly tech sync
  - Weekly full team sync
- No daily standups

---

## Risks / Open Questions

- SDS data model not finalized
- Pipeline → SDS integration unclear
- Terminology inconsistencies (evidence)
- Identity resolution not implemented yet
- Evaluation framework still subjective
- Data coverage gaps (especially younger users)

---

## Next Steps

- Jimmy:
  - Share PRD
- Roman:
  - Complete ingestion pipeline
- Alex Marple:
  - Continue evaluation framework
- Michael:
  - Finalize vendor list
- Alex Maeda:
  - Sync with Roman + Alex Marple
  - Define attribute layer + SDS structure

---

## Overall Signal

- Transitioning from exploration → execution
- Core architecture still forming
- SDS + data pipeline = critical path
