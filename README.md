# Lõputöö PoC: mikroteenused + OpenTelemetry + New Relic

See repositoorium sisaldab mikroteenustel põhinevat tõenduskontseptsiooni (*proof-of-concept*), mille eesmärk on demonstreerida otspunktist-lõpp-punktini jälgitavust OpenTelemetry ja New Relicu abil.

Rakendus on loodud bakalaureusetöö praktilise osa jaoks, et hinnata hajusjälgimise mõju mikroteenustel põhineva süsteemi tõrkeotsingule, teenustevaheliste sõltuvuste nähtavusele ja opereerimiskulude analüüsimisele.

## Arhitektuur

Rakendus koosneb kolmest Go-põhisest mikroteenusest ja MongoDB andmebaasist.

- **mobile-api-service**  
  Toimib frontend-rakenduse sisenemispunktina. Teenus võtab vastu kliendi päringu ja edastab selle `profile-service` teenusele.

- **profile-service**  
  Vastutab kasutajaprofiili andmete eest. Kasutajanime põhjal profiili küsimisel pöördub teenus esmalt `users-api-service` poole, et saada kasutajale vastav UUID. Seejärel otsib `profile-service` saadud UUID alusel profiiliandmed MongoDB-st.

- **users-api-service**  
  Vastutab kasutaja tehniliste andmete eest. Teenus tagastab kasutajanime alusel kasutaja UUID.

- **MongoDB**  
  Kasutatakse andmete püsivaks hoidmiseks. Teenused lisavad käivitumisel algandmed andmebaasi `upsert` loogika abil.

## Päringu põhivoog

Rakenduse keskne kasutusjuht on kasutajaprofiili küsimine kasutajanime alusel.

```text
Client
  → mobile-api-service: GET /api/v1/profile/:username
  → profile-service: GET /api/v1/profile/:username
  → users-api-service: GET /api/v1/users/:username
  → users MongoDB lookup by username
  → profile-service MongoDB lookup by UUID
  → response returned to client
```

Selline voog võimaldab New Relicus näha ühe päringu liikumist üle mitme teenusepiiri. OpenTelemetry abil lisatakse trace’id ja span’id HTTP päringutele, teenusekihtidele ning repository tasemel tehtavatele MongoDB operatsioonidele.

## Endpointid

### mobile-api-service, port 8082

- `GET /health`
- `GET /api/v1/profile/:username`

Näide:

```bash
curl -i http://localhost:8082/api/v1/profile/chris
```

### users-api-service, port 8081

- `GET /health`
- `GET /api/v1/users/:username`

Näide:

```bash
curl -i http://localhost:8081/api/v1/users/chris
```

### profile-service, port 8080

- `GET /health`
- `GET /api/v1/profile/:username`
- `GET /api/v1/profiles/:uuid`

Näited:

```bash
curl -i http://localhost:8080/api/v1/profile/chris
curl -i http://localhost:8080/api/v1/profiles/11111111-1111-1111-1111-111111111111
```

## Eeldused

- Go 1.25.x
- Docker ja Docker Compose
- New Relicu konto ja license key

## Seadistus

Loo repositooriumi juurkausta `.env` fail. Seda faili ei tohi committida.

```bash
cp .env.example .env
```

Seejärel lisa `.env` faili New Relicu license key ja OpenTelemetry ekspordi seadistus.

## Käivitamine

Rakenduse käivitamiseks kasuta Docker Compose’i:

```bash
docker compose up --build
```

Teenused käivituvad järgmistel portidel:

```text
mobile-api-service: http://localhost:8082
profile-service:    http://localhost:8080
users-api-service:  http://localhost:8081
MongoDB:            localhost:27017
```

## Smoke testid

Põhivoo kontrollimiseks:

```bash
curl -i http://localhost:8082/api/v1/profile/chris
```

Oodatud tulemus: `200 OK` ja profiiliandmed.

Puuduva kasutaja kontrollimiseks:

```bash
curl -i http://localhost:8082/api/v1/profile/unknown
```

Oodatud tulemus: `404 Not Found`.

## Tõrke ja viivituse süstimine

Eksperimentide jaoks saab päringutele lisada testparameetreid.

Users-service viivituse simuleerimine:

```bash
curl -i "http://localhost:8082/api/v1/profile/chris?usersDelayMs=800"
```

Profile-service viivituse simuleerimine:

```bash
curl -i "http://localhost:8082/api/v1/profile/chris?profileDelayMs=800"
```

Users-service vea simuleerimine:

```bash
curl -i "http://localhost:8082/api/v1/profile/chris?usersFail=true"
```

Profile-service vea simuleerimine:

```bash
curl -i "http://localhost:8082/api/v1/profile/chris?profileFail=true"
```

Need stsenaariumid on mõeldud hajusjälgimise tulemuste võrdlemiseks New Relicus.

## Katsed

Reprodutseeritavad katse- ja tõrkestsenaariumid on kirjeldatud failis:

```text
docs/experiments.md
```