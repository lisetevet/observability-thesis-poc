# Katsed (PoC)

See dokument kirjeldab reprodutseeritavaid teststsenaariume lõputöö jaoks.

## Eeldused

Käivita kogu süsteem:
```bash
docker compose up --build
```
## Baasstsenaarium

```bash
curl -i "http://localhost:8082/api/v1/profile/chris"
```
## Viivitus (delay) users-teenuses
```bash
curl -i "http://localhost:8081/api/v1/users/chris?delayMs=500"
```
## Viivitus (delay) profile-teenuses
```bash
curl -i "http://localhost:8080/api/v1/profiles/11111111-1111-1111-1111-111111111111?delayMs=500"
```
## Viga (fail) profile-teenuses
```bash
curl -i "http://localhost:8080/api/v1/profiles/11111111-1111-1111-1111-111111111111?fail=true"
```
## Viga (fail) sissepääsupunktis (mobile)
```bash
curl -i "http://localhost:8082/api/v1/profile/chris?fail=true"
```
## Kinnitus New Relicus (mida vaadata)
- Service map / dependencies: mobile-api-service -> users-api-service -> profile-service
- Traces / distributed tracing: näha, kus viivitus või viga tekkis (span duration / error)
- MongoDB span’id users- ja profile-teenustes DB päringute ajal