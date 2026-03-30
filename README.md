# Lõputöö PoC: mikroteenused + OpenTelemetry + New Relic

See repositoorium sisaldab mikroteenustel põhinevat tõenduskontseptsiooni (proof-of-concept), mille eesmärk on demonstreerida otspunktist-lõpp-punktini jälgitavust (end-to-end observability) OpenTelemetry (distributed tracing) ja New Relicu abil.

## Arhitektuur (high-level)

- **mobile-api-service** (entrypoint / BFF)  
  Orkestreerib päringu: kasutajanimi → UUID → profiil.
- **users-api-service**  
  Tagastab kasutaja UUID kasutajanime (username) alusel (MongoDB).
- **profile-service**  
  Tagastab profiili UUID alusel (MongoDB).
- **MongoDB**  
  Andmete püsivus. Teenused seedivad algandmed käivitamisel (upsert).

## Endpointid

**mobile-api-service (port 8082)**
- `GET /health`
- `GET /api/v1/profile/:username`

**users-api-service (port 8081)**
- `GET /health`
- `GET /api/v1/users/:username`

**profile-service (port 8080)**
- `GET /health`
- `GET /api/v1/profiles/:uuid`

## Eeldused

- Go (1.25.x)
- Docker & Docker Compose
- New Relicu konto + **license key** (ingest key)

## Seadistus

Loo repositooriumi juurkausta `.env` fail (ära commiti seda). Mall:

```bash
cp .env.example .env