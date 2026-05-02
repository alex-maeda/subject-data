### Marple sent this
@Roman Shteinberg and @Alex Maeda, I setup time tomorrow to go through Roman's work. Hopefully that time works across all time zones, but we can also share resources offline in this thread.

Tentative Agenda:
1. Roman, show enter and exit points for what you have built in https://github.com/sovraai/data. Do you have an example of the output it produces?
2. Discuss Evidence (https://github.com/sovraai/data/blob/6780bece4842e4aaa8ba8f33aadc716bd75f6753/src/ingestion/normalized/evidence.py) so we can incorporate it into https://github.com/sovraai/subject-data/tree/main/data-model and https://github.com/sovraai/architecture/blob/main/sequence.md .
3. Discuss the write pattern between data and subject-data.
 - Is data writing Subjects? Records? Attributes? All of the above?
4. Discuss the write and search pattern between data, ElasticSearch, and subject-search.
 - Is data executing ElasticSearch updates? If so, what index is it writing, and is it expecting to be the sole writer to that index?
 - Or should ElasticSearch updates be triggered by subject-data after write? If so, what do we need to include to support "internal subjects" as distinguishable from "composite/session subjects"?

---
A few notes as I'm skimming ahead of time...

It looks like the output of the pipeline you built is documents being batch-written into ES.
It looks like documents are written to the ES index sovra_persons : https://github.com/sovraai/data/blob/main/scripts/create_index.sh as well as a few other references.
... but I only see the sovra index on http://dev.tail5ab057.ts.net:5601/app/elasticsearch/index_management/indices

https://github.com/sovraai/data/blob/6780bece4842e4aaa8ba8f33aadc716bd75f6753/ingest/DESIGN.md gives a summary of the structure of the ingestion logic

### Roman sent this
Ingestion pipeline — current status

Purpose. Bring heterogeneous external sources to a single indexable representation: from raw files to records in Elasticsearch. Canonical stage sequence: `raw → versioned → parsed → normalized → indexed`. Every stage produces its own artifact with an explicit contract, lineage, and verifiable membership in a dataset and its version.

Stage status

Labels: `[done]` — implemented and closed; `[in progress]` — has an active task; `[not started]` — designed but no code yet.

• `[done]` *raw* — raw files as received (often `*.zst`, multiple formats). Storage and backup closed (TASK-008); profiling utilities for unpacked content exist (TASK-030).
• `[done on one dataset]` *versioned* — admission of a file as a valid dataset version. Contract and execution for `peopledatalabs_v1` are done (TASK-035). A general runner across all datasets does not exist yet.
• `[done on one dataset]` *parsed* — Parquet with row-level parsing and field categorization (`indexed / stored / extra / sensitive / uninterpretable`), lineage, and schema hashes. Contract — TASK-031, execution for `peopledatalabs_v1` — TASK-036, move of part-constant control fields into file metadata — TASK-043.
• `[in progress]` *normalized* — Parquet with normalized columns plus a nested array of atomic evidence. The PDL normalization baseline and the evidence contract are done (TASK-040, TASK-042); derivation and validation policy itself is the open TASK-047. In parallel: parts-level parallelization (TASK-044) and PyArrow-compute vectorization (TASK-045).
• `[not started]` *indexed* — final load into Elasticsearch. The design is fixed in the conceptual model; implementation for the v1 dataset subset (TASK-033) is intentionally deferred.

Supporting subsystems

• `[in progress]` *catalog.db* — SQLite source of stage state and scan materialization (TASK-028); for now the pipeline still relies on the JSON inventory.
• `[done]` *datasets.yaml + dataset_versions/* — dataset registry contract and the binding from parser configs to versions (TASK-025/026/027).
• `[in progress]` *inventory / dashboard* — unified source-state dashboard (TASK-010); baseline multi-source inventory refresh is closed (TASK-009).
• `[in progress]` *manual reassignment and invalidation* — operator correction flow (TASK-029).

Atomic evidence model

Why we need it. On the `normalized` stage a record carries two fundamentally different things: values that came over from `parsed` as-is or after light value-normalization, and *new* facts that did not exist in `parsed` — they are computed from normalized values of the same row. Without an explicit split between these two types, explainability is lost: it is unclear which column is whose and where `first_name` came from when the source only had `full_name`. Evidence captures *which normalized field*, *which attribute*, and *with which value* was derived. This is needed for debugging, for downstream validation, and for the final ES materialization — which can rely on evidence instead of re-parsing values.

How they are shaped. Evidence lives as a nested `evidence` array inside the same normalized record (not a separate companion artifact). One element of the array is exactly one derived fact from exactly one normalized column, and consists of three fields:
• `source_field` — name of the normalized column in the same record from which the fact was derived;
• `normalized_attribute` — name of the derived attribute;
• `normalized_value` — the value itself.

Multiple facts derived from the same column of the same record share that column's name as their `source_field` — that is how a downstream consumer tells «one source, multiple observations» from «independent sources». The dataset version and the evidence schema version are stored at artifact level, not per record. In-record dedup uses the triple `(source_field, normalized_attribute, normalized_value)`; cross-record deduplication and entity resolution are *outside* the evidence contract.

*How we obtain evidence.* The evidence emitter is a separate callback in the normalizing runner (delivered by TASK-042; the `evidence` column ships empty for now). Real rules arrive in TASK-047:
• *derivation runners* — a runner is launched on a normalized column and produces candidates. The first set: `full_name → first_name / middle_name / last_name`, `phone → country` and canonical E.164, `address → city / state / postal_code / county / country / street_address / po_box` via `usaddress + pgeocode` (US-only).
• *per-attribute candidate validation* — every candidate goes through a semantic check bound to the attribute (name-like, email, phone, address). Accepted candidates become evidence; rejected ones land in counters and samples in the Parquet metadata and in the JSON run report.
• *per-value validity check on the source value* — if the normalized value itself is not a valid instance of its declared type (for example, `full_name="12345"`), scalar fields are nulled, list-shaped fields (`addresses`, `emails`, `phones`) keep the element at its position with `status="rejected"` and a `reject_reason`.
• *non-redundancy rule* — at the emitter, not inside runners: do not emit a candidate if the same attribute already exists as a column or has already been emitted; decompositions are emitted when no same-attribute column exists; normalizations and corrections are emitted when they are not substrings of the source value.
• *record-status gate* — for records with `_record_status == "rejected"` the emitter is not invoked: deriving evidence from unreliable values would only inject noise.

What the model does not do. It does not build record-level summaries or `best_*` fields, does not define the ES mapping, does not resolve cross-record conflicts, and does not perform entity resolution — those are separate downstream steps.