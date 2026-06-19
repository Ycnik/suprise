# Programmierworkshop am 19.6.2026

## Namen

Tareq Daoud-Ghadieh, Yannik Bachmann

## Link zum Git-Repository

[https://github.com/Ycnik/suprise.git](https://github.com/Ycnik/suprise.git)

## KI-Werkzeuge

### Agenten

Für die Bearbeitung wurde OpenAI Codex als lokaler Entwicklungsagent eingesetzt.
Der Agent wurde nicht als Ersatz für die fachliche Entscheidung verwendet, sondern zur strukturierten Unterstützung bei Analyse, Implementierung, Fehlersuche und Dokumentation.

Konkret wurde Codex genutzt für:

* Analyse der Aufgabenstellung und Ableitung einer sinnvollen Projektstruktur
* Auswahl passender Go-Bibliotheken für REST, Validierung, PostgreSQL, OIDC und Tests
* Implementierungsvorschläge für einzelne, klar abgegrenzte Teilaufgaben
* Erstellung und Überarbeitung von Tests
* Unterstützung bei Git-Workflow, Branches, Pull Requests und kleinen Commits
* Unterstützung bei der Keycloak- und Bruno-Konfiguration

Die Ergebnisse wurden lokal getestet und schrittweise in Git übernommen.

Die Umsetzung wurde über mehrere kleine Teilaufgaben strukturiert. Die Teilergebnisse wurden über Branches und Pull Requests zusammengeführt. Vor dem Mergen wurden die GitHub-Actions-Checks betrachtet.

### Chat-URLs, z.B. https://chatgpt.com

Verwendet wurde OpenAI Codex über ChatGPT:

[https://chatgpt.com](https://chatgpt.com)

Ein öffentlicher Chat-Link wurde nicht verwendet, da die Arbeit in einer lokalen Codex-/Repository-Umgebung stattgefunden hat.

## Frameworks und Bibliotheken

### REST-Schnittstelle (Lesen und Neuanlegen)

Die REST-Schnittstelle ist in Go umgesetzt. Für das Routing werden `net/http` und `github.com/go-chi/chi/v5` verwendet.

`chi` wurde gewählt, weil es schlank ist, gut zu Go-Standardbibliotheken passt und Middleware sowie URL-Parameter übersichtlich abbildet.

Umgesetzte Endpunkte:

* `GET /health`
* `GET /rest`
* `GET /rest/{id}`
* `POST /rest`

Die Handler greifen nicht direkt auf die Datenbank zu, sondern verwenden ein Repository-Interface. Dadurch bleiben HTTP-Schicht und Persistenz getrennt.

### Validierung (nur Neuanlegen)

Für die Validierung beim Neuanlegen wird `github.com/go-playground/validator/v10` verwendet.

Validiert werden nur eingehende `POST /rest`-Requests. Beispiele:

* Pflichtfelder wie `vorname`, `nachname` und `username`
* Mindestlängen für Textfelder
* korrektes Datumsformat bei `geburtsdatum`
* Pflichtfelder innerhalb der optionalen Ausrüstung

Ungültige Requests werden mit `400 Bad Request` und einer JSON-Fehlermeldung beantwortet.

### OR-Mapping (für PostgreSQL)

Für PostgreSQL wird `gorm.io/gorm` mit `gorm.io/driver/postgres` verwendet.

Das vorhandene Datenbankschema aus den vorherigen Abgaben wird weiterverwendet. Die Tabellen werden als Go-Structs modelliert:

* `soldat.soldat`
* `soldat.ausruestung`
* `soldat.verletzung`

Die Repository-Schicht enthält eine GORM-Implementierung für die echte Datenbank und eine In-Memory-Implementierung für Tests.

### Optional: OIDC mit Keycloak

Keycloak ist optional als OIDC-Provider eingebunden.

Verwendete Bibliotheken:

* `github.com/coreos/go-oidc/v3/oidc`
* `golang.org/x/oauth2`

Die Absicherung ist über Umgebungsvariablen aktivierbar. Wenn `AUTH_ENABLED=true` gesetzt ist, benötigt `POST /rest` ein gültiges Bearer Token von Keycloak. Lesende Endpunkte bleiben bewusst frei erreichbar:

* `GET /health`
* `GET /rest`
* `GET /rest/{id}`

Damit ist die API auch ohne Keycloak testbar, kann aber für das Neuanlegen zusätzlich abgesichert werden.

### Einfacher Integrationstest

Die HTTP-Schicht wird mit Go-Bordmitteln getestet:

* `testing`
* `net/http/httptest`

Für die Tests wird die In-Memory-Repository-Implementierung verwendet. Dadurch sind die Tests schnell, unabhängig von Docker/PostgreSQL und trotzdem nah an der HTTP-Schnittstelle.

Getestet wird unter anderem:

* `GET /health` liefert `200 OK`
* `POST /rest` mit gültigen Daten liefert `201 Created`
* `POST /rest` mit ungültigen Daten liefert `400 Bad Request`
* `GET /rest/{id}` liefert einen zuvor angelegten Soldaten
* `GET /rest/{id}` setzt einen `ETag`-Header
* `POST /rest` kann bei aktivierter Authentifizierung geschützt werden

Zusätzlich gibt es optionale PostgreSQL-Integrationstests. Diese verwenden die echte Datenbankverbindung und laufen nur, wenn der Build-Tag `integration` gesetzt wird.

## Ausführen

### Voraussetzungen

Für den vollständigen lokalen Betrieb werden benötigt:

* Go
* PostgreSQL mit dem vorhandenen Schema `soldat`
* optional Keycloak für OIDC
* optional Bruno für manuelle API-Tests

### Lokale Umgebungsvariablen

PowerShell:

```powershell
$env:HTTP_ADDR=":8080"
$env:DATABASE_URL="host=localhost user=soldat password=p dbname=soldat port=5432 sslmode=disable"
```

Optional mit Keycloak:

```powershell
$env:AUTH_ENABLED="true"
$env:OIDC_ISSUER_URL="http://localhost:8880/realms/suprise"
$env:OIDC_CLIENT_ID="suprise-client"
```

Optionaler lokaler Reset-Endpunkt für Bruno:

```powershell
$env:DB_RESET_ENABLED="true"
```

Dann kann die Datenbank über `POST /dev/reset-db` zurückgesetzt und neu befüllt werden. Dieser Endpunkt ist standardmäßig deaktiviert.

Ohne Keycloak:

```powershell
$env:AUTH_ENABLED="false"
```

Server starten:

```powershell
go run ./cmd/server
```

### Qualitätssicherung

Normale Tests:

```powershell
gofmt -l cmd internal
go vet ./...
go test ./...
```

Diese Befehle werden auch im GitHub-Actions-Workflow ausgeführt.

Echte PostgreSQL-Integrationstests:

```powershell
$env:DATABASE_URL="host=localhost user=soldat password=p dbname=soldat port=5432 sslmode=disable"
go test -tags=integration ./internal/httpapi
```

### Datenbank zurücksetzen und befüllen

Für lokale Tests kann die PostgreSQL-Datenbank auf feste Beispieldaten zurückgesetzt werden. Dadurch können Bruno-Requests wieder mit bekannten IDs und eindeutigen Testdaten ausgeführt werden.

PowerShell:

```powershell
.\scripts\db\reset_populate.ps1
```

Alternativ kann bei gestartetem Server und `DB_RESET_ENABLED=true` der Bruno-Request `DB Reset Populate` ausgeführt werden.

Der Reset löscht die Daten in `soldat.soldat`, `soldat.ausruestung` und `soldat.verletzung`, setzt die Identity-Zähler zurück und fügt neue Beispieldaten ein.

Enthaltene Testdaten:

* `1000` - `eren-test`
* `1001` - `mikasa-test`
* `1002` - `armin-test`
* `1003` - `levi-test`
* `1004` - `historia-test`

Wenn ein anderer Containername verwendet wird:

```powershell
.\scripts\db\reset_populate.ps1 -Container postgres
```

## Bruno-Collection

Zur manuellen Prüfung liegt eine Bruno-Collection im Repository:

```text
bruno/suprise-http
```

Enthaltene Requests:

* `Health`
* `Soldaten auflisten`
* `Soldat mit ID`
* `Soldat anlegen`
* `Soldat ungueltig anlegen`
* `Keycloak Token admin`
* `Soldat anlegen ohne Token`
* `Soldat anlegen mit Token`

Die lokale Umgebung verwendet:

```text
baseUrl=http://localhost:8080
keycloakBaseUrl=http://localhost:8880
realm=suprise
clientId=suprise-client
username=workshop
password=p
```

Das `clientSecret` wird nicht im Repository gespeichert. Es muss lokal in Bruno in der Umgebung `local` eingetragen werden.

Reihenfolge bei aktivierter Keycloak-Absicherung:

1. Keycloak starten und Realm `suprise` verwenden.
2. Client `suprise-client` mit Client Secret verwenden.
3. Benutzer `workshop` mit Passwort `p` verwenden.
4. In Bruno das `clientSecret` in der Umgebung `local` eintragen.
5. Request `Keycloak Token admin` ausführen.
6. Bruno speichert den erhaltenen `access_token` automatisch als Variable `bearerToken`.
7. Request `Soldat anlegen mit Token` ausführen.

## Keycloak einrichten

Keycloak lokal öffnen:

```text
http://localhost:8880
```

Realm:

```text
suprise
```

Client:

```text
Client ID: suprise-client
Client authentication: On
Standard flow: On
Direct access grants: On
Service account roles: On
Valid redirect URIs: *
Web origins: +
```

Benutzer:

```text
Username: workshop
Email verified: On
Required user actions: leer
Password: p
Temporary: Off
```

Wichtig: Das Client Secret wird lokal aus Keycloak kopiert und nicht committed.

## Git-Workflow

Die Arbeit wurde bewusst in kleinen Schritten durchgeführt:

* Basisstruktur und Datenbankzugriff
* HTTP-Routing und Handler
* Tests
* optionale Keycloak-Absicherung
* Bruno-Collection
* GitHub Actions
* Dokumentation

Teilaufgaben wurden über Branches und Pull Requests zusammengeführt. Dadurch waren die Änderungen besser prüfbar und Konflikte leichter zu behandeln.

Wichtige Regel im Workflow: Änderungen wurden nicht direkt groß auf `main` gesammelt, sondern über thematische Arbeitsbranches und möglichst kleine Commits vorbereitet.

## Prompts/Requests an KI-Agent/en

Die folgenden Prompts sind sinngemäß wiedergegeben. Ziel war nicht, unstrukturierte Chat-Verläufe zu dokumentieren, sondern die relevanten Arbeitsaufträge nachvollziehbar zu machen.

Die Prompts wurden bewusst als Arbeitsaufträge formuliert. Nach der Ausgabe der KI wurden Code, Tests und Verhalten überprüft und bei Bedarf angepasst.

### Aufgabenanalyse

```text
Analysiere die Aufgabenstellung für den Programmierworkshop und leite daraus ab, welche Bestandteile implementiert und dokumentiert werden müssen.
```

```text
Schlage eine sinnvolle Go-Projektstruktur für eine kleine REST-API mit PostgreSQL, Validierung, optionaler Keycloak-Absicherung und Tests vor.
```

```text
Welche Go-Bibliotheken eignen sich für Routing, Validierung, OR-Mapping mit PostgreSQL, OIDC mit Keycloak und HTTP-Integrationstests?
```

### Basisimplementierung

```text
Implementiere die Grundstruktur des Go-Projekts mit Konfiguration über Umgebungsvariablen, PostgreSQL-Anbindung über GORM und einem Repository-Interface.
```

```text
Modelliere die vorhandenen Tabellen soldat, ausruestung und verletzung als Go-Structs für GORM. Die Implementierung soll zum bestehenden PostgreSQL-Schema passen.
```

```text
Erstelle eine GORM-Repository-Implementierung für PostgreSQL und zusätzlich eine In-Memory-Implementierung für automatisierte Tests.
```

### REST-Schnittstelle

```text
Implementiere die HTTP-Schicht mit chi. Es sollen GET /health, GET /rest, GET /rest/{id} und POST /rest bereitgestellt werden.
```

```text
Die Handler sollen das Repository-Interface verwenden, JSON Request/Response unterstützen und Fehler konsistent als JSON zurückgeben.
```

```text
Validiere POST /rest mit go-playground/validator. Ungültige Daten sollen mit 400 Bad Request beantwortet werden.
```

### Zusammenarbeit über Branches

```text
Arbeite auf einem eigenen Branch für den HTTP-Teil. Implementiere Routing, Handler und Validierung getrennt von der Basisstruktur. Erstelle kleine Commits und pushe den Branch für einen Pull Request.
```

```text
Arbeite parallel auf einem eigenen Branch für Authentifizierung, CI und Dokumentation. Implementiere Keycloak/OIDC optional, ohne die lesenden Endpunkte zu schützen.
```

```text
Prüfe vor dem Merge, ob der Branch aktuell zu main ist, ob die GitHub-Actions-Checks grün sind und ob keine Secrets oder Tokens im Commit enthalten sind.
```

### Keycloak/OIDC

```text
Schließe Keycloak optional an. Wenn AUTH_ENABLED=true ist, soll POST /rest einen Bearer Token benötigen. GET /health, GET /rest und GET /rest/{id} sollen ohne Token erreichbar bleiben.
```

```text
Prüfe, ob die Middleware Access Tokens aus Keycloak akzeptiert und Client-ID bzw. Audience korrekt validiert.
```

```text
Erstelle eine kurze lokale Anleitung für Realm, Client, Benutzer und Bruno-Token-Request, ohne Secrets ins Repository zu schreiben.
```

### Tests und CI

```text
Schreibe Tests für die HTTP-Schicht mit net/http/httptest und dem In-Memory-Repository.
```

```text
Teste Healthcheck, erfolgreiches Neuanlegen, Validierungsfehler, Finden per ID, ETag und Authentifizierungspflicht bei aktiviertem Keycloak-Modus.
```

```text
Richte GitHub Actions so ein, dass gofmt, go vet und go test bei Push und Pull Request ausgeführt werden.
```

```text
Ergänze optionale PostgreSQL-Integrationstests so, dass sie gegen die echte soldat-Datenbank laufen können, aber den normalen CI-Lauf mit go test ./... nicht blockieren.
```

```text
Verschiebe Tests mit echter PostgreSQL-Abhängigkeit hinter den Build-Tag integration und dokumentiere den Aufruf mit go test -tags=integration ./internal/httpapi.
```

### Docker und Lasttests

```text
Erstelle einen Dockerfile für die Go-Anwendung, sodass die API als Container gebaut und gestartet werden kann.
```

```text
Ergänze einfache Benchmarks für die HTTP-Endpunkte, um Healthcheck, Listen, Finden per ID und Neuanlegen unter Last grob prüfen zu können.
```

### Bruno-Collection

```text
Erstelle eine Bruno-Collection für die REST-Endpunkte. Sie soll Requests für Health, Auflisten, Finden per ID, gültiges Neuanlegen, ungültiges Neuanlegen und Keycloak-Token enthalten.
```

```text
Speichere Tokens und Client Secrets nicht im Repository. Nutze Bruno-Variablen, damit der Access Token nach dem Token-Request automatisch für den geschützten POST-Request verwendet werden kann.
```

```text
Passe die Bruno-Requests so an, dass gültige Beispielkörper verwendet werden und keine echten Bearer Tokens oder Client Secrets im Repository gespeichert werden.
```

```text
Prüfe, warum ein POST /rest mit 500 fehlschlägt. Vergleiche Request-Body, Enum-Werte und Unique Constraints der PostgreSQL-Datenbank und leite daraus eine saubere Korrektur ab.
```

### Fehlerbehandlung und Nachbesserungen

```text
Verbessere die Validierungsfehler so, dass die API nicht nur validierung fehlgeschlagen zurückgibt, sondern konkrete Details zu betroffenen Feldern und Gründen.
```

```text
Wenn ein Datenbankfehler durch falsche Enum-Werte entstehen kann, validiere diese Werte bereits vor dem Insert und antworte mit 400 Bad Request.
```

```text
Prüfe die Bruno-Collection auf versehentlich gespeicherte Tokens, Secrets oder veraltete Testdaten und korrigiere nur die betroffenen Request-Dateien.
```

### Gemeinsame Arbeitsaufträge

```text
Analysiere die vorhandenen vorherigen Abgaben und nutze die bestehende PostgreSQL-Datenbank mit dem Schema soldat als Grundlage für die neue Go-REST-API.
```

```text
Implementiere die Basisstruktur mit Konfiguration, GORM-Anbindung, Repository-Schicht und Datenmodellen. Achte darauf, dass die Struktur für weitere HTTP- und Auth-Branches geeignet ist.
```

```text
Ergänze optionale Keycloak-Authentifizierung. Bei AUTH_ENABLED=true soll POST /rest geschützt sein, während GET /health, GET /rest und GET /rest/{id} öffentlich bleiben.
```

```text
Richte GitHub Actions für gofmt, go vet und go test ein. Die Checks sollen bei Pull Requests nach main laufen.
```

```text
Erstelle eine README, die Namen, Repository-Link, verwendete KI-Werkzeuge, Bibliotheken, Tests, Keycloak, Bruno und die wichtigsten KI-Arbeitsaufträge nachvollziehbar dokumentiert.
```

```text
Prüfe bei jeder Änderung, ob ein eigener Branch sinnvoll ist, ob kleine Commits entstehen und ob keine Secrets oder lokalen Tokens ins Repository gelangen.
```

```text
Arbeite auf einem eigenen Branch für die HTTP-Schicht. Implementiere GET /health, GET /rest, GET /rest/{id} und POST /rest mit chi und dem vorhandenen Repository-Interface.
```

```text
Erstelle eine Bruno-Collection mit Requests für Healthcheck, Auflisten, Finden per ID, gültiges Neuanlegen, ungültiges Neuanlegen, Token abrufen und Neuanlegen mit Bearer Token.
```

```text
Ergänze echte PostgreSQL-Integrationstests für die HTTP-Schicht. Die Tests sollen die reale Datenbank verwenden können, aber nur mit dem Build-Tag integration laufen.
```

```text
Erstelle Benchmarks für die wichtigsten HTTP-Endpunkte und achte darauf, dass normale Tests weiterhin ohne Docker und ohne Datenbank laufen.
```

```text
Erstelle einen Dockerfile für den Build der Go-Anwendung und integriere ihn über einen eigenen Pull Request.
```

```text
Rebase den eigenen Branch regelmäßig auf origin/main, pushe mit --force-with-lease nur auf den eigenen Branch und merge erst, wenn die Checks grün sind.
```

### Abschlussprüfung

```text
Prüfe den finalen Stand mit gofmt, go vet und go test. Kontrolliere außerdem, ob README, Bruno-Collection, Keycloak-Hinweise und GitHub Actions zum implementierten Stand passen.
```
