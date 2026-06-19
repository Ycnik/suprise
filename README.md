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

### Einfacher Integrationstest

Ein einfacher Integrationstest kann mit Go-Bordmitteln umgesetzt werden:

* `testing`
* `net/http/httptest`

Fuer Tests mit echter PostgreSQL-Datenbank kann optional `github.com/testcontainers/testcontainers-go` verwendet werden. Damit laesst sich fuer den Testlauf ein PostgreSQL-Container starten, gegen den die REST-Endpunkte getestet werden.

## Ausfuehren

Voraussetzung ist ein laufender PostgreSQL-Server mit dem Schema `soldat` aus den vorherigen Abgaben.

Beispiel fuer lokale Umgebungsvariablen:

```bash
export HTTP_ADDR=":8080"
export DATABASE_URL="host=localhost user=soldat password=p dbname=soldat port=5432 sslmode=disable"
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

## Prompts/Requests an KI-Agent/en

Die KI wurde schrittweise und unterstuetzend eingesetzt. Dabei wurden sinngemaess folgende Prompts bzw. Requests verwendet:

1. Analysiere die Aufgabenstellung und fasse zusammen, welche Bestandteile fuer die Abgabe dokumentiert werden muessen.

2. Schlage eine einfache technische Struktur fuer eine prototypische Go-Anwendung mit REST-Schnittstelle vor.

3. Welche Go-Bibliothek eignet sich fuer Routing und Middleware bei einer kleinen REST-API?

4. Welche Bibliothek kann fuer die Validierung eingehender JSON-Daten beim Neuanlegen von Datensaetzen verwendet werden?

5. Schlage eine passende OR-Mapping-Bibliothek fuer die Anbindung an PostgreSQL vor.

6. Wie kann eine REST-Schnittstelle optional ueber OIDC mit Keycloak abgesichert werden?

7. Welche Bibliotheken eignen sich fuer einfache Integrationstests einer Go-REST-API?

8. Formuliere die Ergebnisse strukturiert fuer eine kurze Workshop-Dokumentation in Markdown.

Die Antworten der KI wurden geprueft und fuer die Dokumentation zusammengefasst.
