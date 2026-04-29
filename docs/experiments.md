# Katsed

See dokument kirjeldab reprodutseeritavaid teststsenaariume lõputöö praktilise osa jaoks. Katsete eesmärk on kontrollida, kas OpenTelemetry abil on võimalik jälgida ühe kasutajapäringu liikumist üle mitme mikroteenuse ning tuvastada, millises teenuses või kihis tekib viivitus või viga.

## Eeldused

Käivita kogu süsteem repositooriumi juurkaustast:

```bash
docker compose up --build
```

Kontrolli, et teenused vastavad health endpointidel:

```bash
curl -i http://localhost:8082/health
curl -i http://localhost:8080/health
curl -i http://localhost:8081/health
```

Oodatud tulemus kõigil teenustel: `200 OK`.

## Põhivoog

Rakenduse põhipäring tehakse läbi `mobile-api-service` teenuse:

```bash
curl -i "http://localhost:8082/api/v1/profile/chris"
```

Oodatud tulemus: `200 OK` ja profiiliandmed.

Põhivoog teenuste vahel:

```text
Client
  → mobile-api-service
  → profile-service
  → users-api-service
  → users MongoDB lookup by username
  → profile-service MongoDB lookup by UUID
  → response returned to client
```

New Relicus peaks sama päring olema nähtav ühe distributed trace’ina, kus on näha `mobile-api-service`, `profile-service`, `users-api-service` ning MongoDB span’id.

## Puuduva kasutaja stsenaarium

```bash
curl -i "http://localhost:8082/api/v1/profile/unknown"
```

Oodatud tulemus: `404 Not Found`.

Selle stsenaariumi eesmärk on kontrollida, kas trace näitab, et kasutajat ei leitud `users-api-service` kaudu ning `profile-service` tagastas kontrollitud 404 vastuse.

## Viivitus users-service’is

```bash
curl -i "http://localhost:8082/api/v1/profile/chris?usersDelayMs=800"
```

Oodatud tulemus: `200 OK`, kuid vastus tuleb tavapärasest aeglasemalt.

Selle stsenaariumi eesmärk on kontrollida, kas New Relicu trace’is on näha, et viivitus tekkis `users-api-service` harus.

## Viivitus profile-service’is

```bash
curl -i "http://localhost:8082/api/v1/profile/chris?profileDelayMs=800"
```

Oodatud tulemus: `200 OK`, kuid vastus tuleb tavapärasest aeglasemalt.

Selle stsenaariumi eesmärk on kontrollida, kas New Relicu trace’is on näha, et viivitus tekkis `profile-service` kihis.

## Viga users-service’is

```bash
curl -i "http://localhost:8082/api/v1/profile/chris?usersFail=true"
```

Oodatud tulemus: veavastus, sest `profile-service` ei saa `users-api-service` käest kasutaja UUID-d.

Selle stsenaariumi eesmärk on kontrollida, kas distributed trace’is on näha vea algpõhjus `users-api-service` poolel.

## Viga profile-service’is

```bash
curl -i "http://localhost:8082/api/v1/profile/chris?profileFail=true"
```

Oodatud tulemus: veavastus, sest `profile-service` simuleerib tõrget.

Selle stsenaariumi eesmärk on kontrollida, kas distributed trace’is on näha vea algpõhjus `profile-service` poolel.

## Otsesed teenusepõhised kontrollpäringud

Neid päringuid saab kasutada teenuste eraldi kontrollimiseks, kuid lõputöö põhikatsetes tuleks eelistada päringuid läbi `mobile-api-service`, sest need näitavad otspunktist-lõpp-punktini päringu liikumist.

### users-api-service

```bash
curl -i "http://localhost:8081/api/v1/users/chris"
curl -i "http://localhost:8081/api/v1/users/chris?delayMs=500"
curl -i "http://localhost:8081/api/v1/users/chris?fail=true"
```

### profile-service

```bash
curl -i "http://localhost:8080/api/v1/profile/chris"
curl -i "http://localhost:8080/api/v1/profile/chris?usersDelayMs=500"
curl -i "http://localhost:8080/api/v1/profile/chris?usersFail=true"
curl -i "http://localhost:8080/api/v1/profiles/11111111-1111-1111-1111-111111111111"
```

## Kinnitus New Relicus

Katsete järel kontrolli New Relicus järgmisi vaateid:

- **Service map / dependencies**: põhivoog peab näitama seost `mobile-api-service → profile-service → users-api-service`.
- **Distributed traces**: ühe päringu trace peab näitama teenustevahelist päringuahelat.
- **Span duration**: viivituse stsenaariumites peab pikim span asuma vastavas teenuses.
- **Errors**: vea stsenaariumites peab trace näitama, millises teenuses viga tekkis.
- **MongoDB span’id**: users- ja profile-teenustes peavad olema nähtavad andmebaasipäringutega seotud span’id.