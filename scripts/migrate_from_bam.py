#!/usr/bin/env python3
"""
Migrate subjects and their trace data (records) from BAM Prototype Beta
to Subject Data.

Usage:
    python3 scripts/migrate_from_bam.py                          # local Subject Data
    SD_URL=http://subject-data-beta-1... SD_AUTH_TOKEN=xxx python3 scripts/migrate_from_bam.py
"""

import json
import os
import sys
import urllib.request

BAM_URL = os.environ.get("BAM_URL", "http://bam-prototype-beta.tail5ab057.ts.net:8080")
SD_URL = os.environ.get("SD_URL", "http://localhost:8380")
SD_AUTH_TOKEN = os.environ.get("SD_AUTH_TOKEN", "")

SUBJECTS_TO_MIGRATE = [
    ("843400dd-bcdc-42ea-ae46-db838712c287", "Bunkeddeko, Adem"),
    ("b7a02d5a-04be-4419-9927-df6bb566151a", "Zherelyev, Anton"),
    ("91462ab1-a166-476b-b549-0d0541297107", "Marks, Ana"),
]


def bam_get(path: str):
    req = urllib.request.Request(f"{BAM_URL}{path}", headers={"Accept": "application/json"})
    with urllib.request.urlopen(req) as resp:
        return json.loads(resp.read())


def sd_request(method: str, path: str, body=None):
    data = json.dumps(body).encode() if body else None
    req = urllib.request.Request(
        f"{SD_URL}{path}",
        data=data,
        method=method,
        headers={"Content-Type": "application/json"},
    )
    if SD_AUTH_TOKEN:
        req.add_header("Authorization", f"Bearer {SD_AUTH_TOKEN}")
    with urllib.request.urlopen(req) as resp:
        return json.loads(resp.read())


def main():
    print("=== Migrating from BAM Prototype Beta to Subject Data ===")
    print(f"  BAM: {BAM_URL}")
    print(f"  SD:  {SD_URL}")
    print()

    # Step 1: Create subjects
    id_map = {}  # BAM ID -> new SD ID
    for bam_id, name in SUBJECTS_TO_MIGRATE:
        print(f"Creating subject: {name}")
        resp = sd_request("POST", "/v1/subjects", {"subject_name": name})
        new_id = resp["id"]
        print(f"  BAM ID: {bam_id} -> SD ID: {new_id}")
        id_map[bam_id] = new_id

    print()

    # Step 2: Migrate trace data as records
    for bam_id, name in SUBJECTS_TO_MIGRATE:
        new_id = id_map[bam_id]
        print(f"Fetching trace data for {name} ({bam_id})...")

        trace_data = bam_get(f"/api/subjects/{bam_id}/trace-data")
        print(f"  Found {len(trace_data)} trace data items. Migrating as records...")

        success = 0
        errors = 0
        for i, item in enumerate(trace_data):
            # Remove BAM id, set new subject_id
            item.pop("id", None)
            item["subject_id"] = new_id

            try:
                sd_request("POST", "/v1/records", item)
                success += 1
            except Exception as e:
                errors += 1
                if errors <= 3:
                    print(f"  WARNING: Failed to create record {i}: {e}")
                elif errors == 4:
                    print("  (suppressing further warnings)")

            if (i + 1) % 100 == 0:
                print(f"  Progress: {i + 1}/{len(trace_data)}")

        print(f"  Done: {success} created, {errors} errors.")

    print()
    print("=== Migration complete ===")
    print()
    print("Subject ID mapping:")
    for bam_id, name in SUBJECTS_TO_MIGRATE:
        print(f"  {name}: {id_map[bam_id]}")


if __name__ == "__main__":
    main()
