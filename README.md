# Programmierworkshop am 19.6.2026

## Namen

Tareq Daoud-Ghadieh, Yannik Bachmann

## Link zum Git-Repository

[https://github.com/Ycnik/suprise.git](https://github.com/Ycnik/suprise.git)

## KI-Werkzeuge

### Agenten

OpenAI Codex wurde zur Analyse der Aufgabenstellung, zur Auswahl geeigneter Go-Bibliotheken und zur Formulierung dieser Dokumentation verwendet.

### Chat-URLs, z.B. https://chatgpt.com

ChatGPT / Codex: https://chatgpt.com

## Frameworks und Bibliotheken

### REST-Schnittstelle (Lesen und Neuanlegen)

Verwendet wird Go mit `net/http` und `github.com/go-chi/chi/v5`.

`chi` eignet sich gut fuer eine kleine prototypische REST-Schnittstelle, weil Routing, Middleware und HTTP-Handler schlank bleiben und trotzdem sauber strukturiert werden koennen.

Umgesetzte Endpunkte:

* `GET /health`
* `GET /rest`
* `GET /rest/{id}`
* `POST /rest`

### Validierung (nur Neuanlegen)

Verwendet wird `github.com/go-playground/validator/v10`.

Beim Neuanlegen von Datensaetzen werden die eingehenden JSON-Daten ueber Struct-Tags validiert, z.B. Pflichtfelder, Mindestlaengen oder gueltige Wertebereiche.

### OR-Mapping (für PostgreSQL)

Verwendet wird `gorm.io/gorm` mit `gorm.io/driver/postgres`.

GORM wird als OR-Mapper fuer PostgreSQL verwendet. Dadurch koennen Entitaeten als Go-Structs modelliert und Datenbankoperationen wie Lesen und Neuanlegen uebersichtlich umgesetzt werden. Die Implementierung nutzt das vorhandene Datenbankschema `soldat` aus den vorherigen Abgaben.

### Optional: OIDC mit Keycloak

Keycloak ist als OIDC-Provider vorbereitet.

Verwendete Bibliotheken:

* `github.com/coreos/go-oidc/v3/oidc`
* `golang.org/x/oauth2`

Die REST-Schnittstelle kann damit ueber Bearer Tokens abgesichert werden. Die Middleware prueft dabei das vom Client mitgesendete Access Token gegen die Keycloak-Konfiguration, insbesondere Issuer, Client-ID/Audience und Signatur.

Bei `AUTH_ENABLED=true` wird das Neuanlegen ueber `POST /rest` mit der Keycloak-Middleware abgesichert. Lesen ueber `GET /rest`, `GET /rest/{id}` und `GET /health` bleibt ohne Token moeglich.

### Einfacher Integrationstest

Ein einfacher Integrationstest kann mit Go-Bordmitteln umgesetzt werden:

* `testing`
* `net/http/httptest`

Fuer Tests mit echter PostgreSQL-Datenbank kann optional `github.com/testcontainers/testcontainers-go` verwendet werden. Damit laesst sich fuer den Testlauf ein PostgreSQL-Container starten, gegen den die REST-Endpunkte getestet werden.

## Ausfuehren

Voraussetzung ist ein laufender PostgreSQL-Server mit dem Schema `soldat` aus den vorherigen Abgaben.

Beispiel fuer lokale Umgebungsvariablen unter Linux/macOS:

```bash
export HTTP_ADDR=":8080"
export DATABASE_URL="host=localhost user=soldat password=p dbname=soldat port=5432 sslmode=disable"
```

Beispiel fuer lokale Umgebungsvariablen unter PowerShell:

```powershell
$env:HTTP_ADDR=":8080"
$env:DATABASE_URL="host=localhost user=soldat password=p dbname=soldat port=5432 sslmode=disable"
```

Optionale Keycloak-Konfiguration:

```bash
export AUTH_ENABLED="true"
export OIDC_ISSUER_URL="http://localhost:8843/realms/soldat"
export OIDC_CLIENT_ID="soldat-client"
```

PowerShell:

```powershell
$env:AUTH_ENABLED="true"
$env:OIDC_ISSUER_URL="http://localhost:8843/realms/soldat"
$env:OIDC_CLIENT_ID="soldat-client"
```

Start des Servers:

```bash
go run ./cmd/server
```

Tests und statische Pruefungen:

```bash
gofmt -l cmd internal
go vet ./...
go test ./...
```

## Bruno-Collection

Zur manuellen Pruefung der REST-Schnittstelle liegt eine Bruno-Collection im Repository:

```text
bruno/suprise-http
```

Enthaltene Requests:

* `GET /health`
* `GET /rest`
* `GET /rest/{id}`
* `POST /rest` mit gueltigen Soldat-Daten
* `POST /rest` mit ungueltigen Soldat-Daten
* Keycloak-Token fuer `admin` abrufen
* `POST /rest` ohne Token pruefen
* `POST /rest` mit Bearer Token pruefen

Die lokale Bruno-Umgebung verwendet:

```text
baseUrl=http://localhost:8080
```

Bei aktivierter Keycloak-Absicherung (`AUTH_ENABLED=true`) ist die Reihenfolge:

1. `Keycloak Token admin` ausfuehren.
2. `Soldat anlegen ohne Token` ausfuehren und `401` erwarten.
3. `Soldat anlegen mit Token` ausfuehren und `201` erwarten.

## Prompts/Requests an KI-Agent/en

Die KI wurde schrittweise und unterstuetzend eingesetzt. Die Ergebnisse wurden geprueft, angepasst und in kleinen Git-Commits umgesetzt. Die wesentlichen Prompts bzw. Arbeitsauftraege waren:

### Aufgabenanalyse und Technologieauswahl

```text
Analysiere die Aufgabenstellung fuer den Programmierworkshop und fasse zusammen, welche Bestandteile fuer die Abgabe dokumentiert und implementiert werden muessen.
```

```text
Schlage eine realistische Go-Projektstruktur fuer eine prototypische REST-API mit PostgreSQL, Validierung, optionaler Keycloak-OIDC-Absicherung und einfachen Tests vor.
```

```text
Welche Go-Bibliotheken eignen sich fuer Routing, Validierung, OR-Mapping mit PostgreSQL, OIDC mit Keycloak und HTTP-Integrationstests?
```

### Implementierung der Basis

```text
Implementiere die Basis fuer ein Go-Projekt mit Konfiguration ueber Umgebungsvariablen und PostgreSQL-Anbindung ueber GORM. Nutze das vorhandene Datenbankschema soldat aus den vorherigen Abgaben.
```

```text
Modelliere die Tabellen soldat, ausruestung und verletzung als Go-Structs fuer GORM. Die Implementierung soll zum vorhandenen PostgreSQL-Schema passen.
```

```text
Erstelle ein Repository-Interface fuer Soldaten mit Funktionen zum Auflisten, Finden nach ID und Neuanlegen. Implementiere eine GORM-Variante fuer PostgreSQL und eine In-Memory-Variante fuer Tests.
```

### Arbeitsauftrag an Yannik fuer HTTP/Routing

```text
Du arbeitest im Repo https://github.com/Ycnik/suprise.git auf einem eigenen Branch.

Ausgangslage:
- Der Branch origin/workshop-prototype existiert bereits.
- Darin sind Go-Modul, Config, PostgreSQL/GORM-Anbindung, Soldat-Datenmodell und Repository-Schicht vorhanden.
- Die Datenbank kommt aus den vorherigen Abgaben und nutzt das Schema soldat mit Tabellen soldat, ausruestung und verletzung.
- Bitte nicht auf main oder workshop-prototype direkt arbeiten.

Deine Aufgabe:
1. Hole den Basisbranch:
   git fetch origin
   git switch -c yanik-http origin/workshop-prototype

2. Implementiere den HTTP-/REST-Teil in Go:
   - cmd/server/main.go
   - internal/httpapi/router.go
   - internal/handler/soldat_handler.go

3. REST-Endpunkte:
   - GET /health
   - GET /rest
   - GET /rest/{id}
   - POST /rest

4. Anforderungen:
   - chi als Router verwenden
   - Repository-Interface aus internal/repository nutzen
   - Keine direkte DB-Logik im Handler
   - POST mit go-playground/validator validieren
   - JSON Request/Response
   - Bei GET /rest/{id} einen ETag-Header mit der Version setzen
   - Fehler sauber als JSON zurueckgeben

5. Kleine Commits erstellen und den Branch nach origin pushen.
```

### Keycloak, CI und Dokumentation

```text
Arbeite auf einem eigenen Branch fuer Auth, CI und README. Implementiere Keycloak/OIDC getrennt vom HTTP-Branch.

Aufgaben:
- Go-Dependencies mit go mod tidy aufloesen
- Keycloak/OIDC-Middleware vorbereiten
- GitHub Actions fuer go test ./... einrichten
- README mit Start-, Test- und ENV-Hinweisen aktualisieren
- Kleine Commits erstellen und Branch nach origin pushen
```

```text
Schliesse Keycloak optional an. Wenn AUTH_ENABLED=true ist, soll POST /rest einen Bearer Token benoetigen. GET /health, GET /rest und GET /rest/{id} sollen ohne Token erreichbar bleiben.
```

### Tests und API-Polish

```text
Schreibe Tests fuer die HTTP-Schicht mit net/http/httptest und dem vorhandenen In-Memory-Repository.

Teste mindestens:
- GET /health liefert 200 und {"status":"ok"}
- POST /rest mit gueltigen Soldat-Daten liefert 201
- POST /rest mit ungueltigen Daten liefert 400
- GET /rest/{id} liefert einen angelegten Soldaten
- GET /rest/{id} setzt einen ETag-Header
```

```text
Pruefe, ob POST /rest mit geburtsdatum im Format "2000-01-01" funktioniert. Falls nicht, parse das Datum im Handler mit time.Parse("2006-01-02", req.Geburtsdatum) und gib bei ungueltigem Datum 400 zurueck.
```

### Bruno-Collection

```text
Erstelle eine Bruno HTTP Collection fuer die vorhandenen REST-Endpunkte.

Wichtig:
- Arbeite auf einem eigenen Branch.
- Fuege nur Dateien unter bruno/ hinzu.
- Keine Go-Dateien, README oder Keycloak-Konfiguration aendern.

Die Collection soll Requests enthalten fuer:
- GET /health
- GET /rest
- GET /rest/{id}
- POST /rest mit gueltigem Soldat-JSON
- POST /rest mit ungueltigem Soldat-JSON

Nutze eine lokale Environment-Datei mit baseUrl=http://localhost:8080.
```

### Abschliessende Pruefung

```text
Pruefe den finalen main-Stand mit gofmt, go vet und go test. Stelle sicher, dass Keycloak optional angeschlossen ist, die Bruno-Collection vorhanden ist und die README den implementierten Stand beschreibt.
```
