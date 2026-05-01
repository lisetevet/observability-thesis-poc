# k6 koormustestide tulemused

See kaust sisaldab lõputöö praktilises osas kasutatud k6 koormustestide JSON-väljundeid.

Koormustestid viidi läbi kuue stsenaariumiga, mis vastavad lõputöö peatükis 3.5 kirjeldatud eksperimendi disainile ja peatükis 4.2 esitatud koondtulemustele. Põhiteksti tabelis „k6 koormustestide tulemused“ kasutati nendest JSON-failidest järgmisi mõõdikuid:

- päringute arv;
- check’ide õnnestumise määr;
- HTTP veamäär;
- keskmine päringu kestus;
- p95 päringu kestus;
- maksimaalne päringu kestus.

## Failide vastavus stsenaariumidele

| Fail | Stsenaarium | Kirjeldus |
|---|---|---|
| `e1-baseline.json` | E1 baasstsenaarium | Edukas tavapärane profiilipäring. |
| `e2-users-delay.json` | E2 users-service viivitus | Viivitus lisatakse users-api-service harus. |
| `e3-profile-delay.json` | E3 profile-service viivitus | Viivitus lisatakse profile-service tasemel. |
| `e4-users-fail.json` | E4 users-service viga | Simuleeritud viga users-api-service harus. |
| `e5-profile-fail.json` | E5 profile-service viga | Simuleeritud viga profile-service tasemel. |
| `e6-not-found.json` | E6 puuduv kasutaja | Kontrollitud 404 Not Found vastus puuduva kasutaja korral. |

## Märkus

JSON-failid on säilitatud algandmetena, et lõputöös esitatud koondtulemused oleksid kontrollitavad. Lõputöö põhitekstis ei esitata kogu JSON-väljundit, sest see on masinloetav ja mahukas. Selle asemel on tulemused koondatud tabelisse ning käesolev kaust sisaldab algseid k6 väljundeid.