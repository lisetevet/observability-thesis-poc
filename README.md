# Lõputöö PoC: mikroteenused + OpenTelemetry + New Relic

See repositoorium sisaldab väikest mikroteenustel põhinevat tõenduskontseptsiooni (proof-of-concept), mida kasutatakse bakalaureusetöös, et demonstreerida otspunktist-lõpp-punktini jälgitavust (end-to-end observability) OpenTelemetry ja New Relicu abil.

## Teenused

- **users-api-service**  
  Tagastab kasutaja UUID kasutajanime (username) põhjal.  
  - `GET /health`  
  - `GET /api/v1/users/:username`

- **mobile-api-service**  
  Avalik sissepääsupunkt (API gateway / BFF), mis kutsub users-api-service’i ja tagastab vastuse.  
  - `GET /health`  
  - `GET /api/v1/profile/:username`

## Eeldused

- Go (1.25.x)
- Docker & Docker Compose
- New Relicu konto + **license key** (ingest key)

## Seadistus

Loo repositooriumi juurkausta `.env` fail (ära commiti seda). Malli saad kopeerida:

```bash
cp .env.example .env