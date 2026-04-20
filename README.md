# people-finder

Imports people from Wikidata into SQLite by iterating over death years.

## Build

```bash
cd /home/boris/dev/stillrunning/people-finder
go mod tidy
go build -o bin/people-finder ./people-finder
```

## Run

Default range is `-5000` to `2026` (inclusive):

```bash
cd /home/boris/dev/stillrunning/people-finder
./bin/people-finder
```

Custom range:

```bash
./bin/people-finder -start-year -5000 -end-year 2026
```

Custom range with multi-year step:

```bash
./bin/people-finder -start-year -5000 -end-year 2026 -step 10
```

Notes:
- SQLite DB is created at `./people.db` (repo root).
- Import may take a long time for the full range.
- You can tune pacing/retries with:
  - `-step` (default `1`, years per query interval)
  - `-request-delay-ms` (default `150`)
  - `-retries` (default `3`)
  - `-retry-delay-ms` (default `1500`)
