# Search API

## Request Examples

### Plain Search

```json
{
  "query": "john smith",
  "providers": ["elasticsearch"]
}
```

### Search with Filters

```json
{
  "query": "john smith",
  "providers": ["elasticsearch"],
  "filters": {
    "filter_list": [
      {
        "key": "first_name",
        "match": "alex alexander"
        // TODO: do we want to indicate that this is a "match"? Or infer it based on the keys? Is "match" the word we want to use?
      },
      {
        "key": "last_name",
        "term": "marple"
        // TODO: do we want to indicate that this is a "term"? Or infer it based on the keys? Is "term" the word we want to use?
      },
      {
        "key": "birth_year",
        "gte": "1990-01-01",
        "lte": "2000-01-01"
        // TODO: do we want to indicate that this is a "range"? Or infer it based on the keys? Is "range" the word we want to use?
      }
      // TODO: leave room for boolean queries in the API, but don't define them until we have a need to
    ]
  }
}
```

## Response Example

```json
{
  "query_id": "abcde",
  "query": {
    "resolved_query": {
      // ... TBD, IMO this should reflect either the query above, or the "interpreted" query, which matches the request format
    },
    "raw_query": { /* ... same format as resolved_query */ },
    "parse_log": [
      "How did I parse the query?",
      "What did I decide to refine?",
      "????",
      "profit"
    ],
    "confidence": "medium"
  },
  "timing": {
    "total_ms": 2642,
    "parse_ms": 0,
    "search_ms": 2642,
    "resolve_ms": 1,
    "cache_ms": 0
  },
  "resolved_persons": [
    {
      "id": "12345",
      "confidence": 1,
      "recors": [],
      "first_name": "Alex",
      "last_name": "Marple",
      "city": "Philadelphia",
      "state": "PA",
      "dob": "1989-12-19"
    }
  ],
  "raw_results": [
    {
      "source": "ssn_breach",
      "source_id": "ssn:/data/raw_unpacked/ssn/ssn.txt:530948310",
      /* ... other fields follow same format as resolved_persons ... */
      "raw_data": { /* ... I don't understand how this differs from the raw_results elements */ }
    }
  ],
  "facets": {
    "facet_list": [
      {
        "key": "first_name",
        "name": "first-name", // TODO: how do we want to handle I18N and L10N?
        "options": [
          {
            "term": "alex",
            "name": "Alex"
          },
          {
            "term": "alexander",
            "name": "Alexander"
          }
        ]
      },
      {
        "key": "birth_year",
        "name": "birth-year",
        "options": [
          {
            "gte": "1990-01-01",
            "lte": "1994-12-31"
          },
          {
            "gte": "1995-01-01",
            "lte": "1990-12-31"
          }
        ]
      }
    ]
  }
}
```
