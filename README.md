# Nourishment — Backend

REST API do zarządzania żywieniem. Obsługuje produkty, przepisy, dziennik posiłków, listy zakupów oraz autoryzację użytkowników.

## Funkcjonalności

- **Produkty** — CRUD bazy produktów spożywczych z wartościami odżywczymi i kategoriami
- **Posiłki** — zarządzanie przepisami, automatyczne przeliczanie makroskładników
- **Dziennik posiłków** — planowanie posiłków na konkretne dni z podziałem na porcje
- **Lista zakupów** — tworzenie i zarządzanie listami zakupów
- **Optymalizacja posiłków** — sugestie z wykorzystaniem AI (OpenRouter)
- **Autoryzacja** — JWT z konfigurowalnymi parametrami (ważność, issuer, audience)
- **Dokumentacja API** — Swagger UI dostępny pod `/swagger/index.html`

## Technologie

- Go 1.24
- Gin (HTTP framework)
- Firebird 2.5 (baza danych)
- JWT (autoryzacja)
- Swagger (dokumentacja API)
- Zerolog (logowanie)

## Uruchomienie

```bash
# Z katalogu głównego backendu
cd cmd/nourishment
go run main.go
```

Serwer startuje domyślnie na porcie **8080**.

### Konfiguracja

Skopiuj `.env.example` do `.env` w katalogu `cmd/nourishment/` i uzupełnij wartości:

```
DB_USER=sysdba
DB_PASSWORD=masterkey
DB_ADDRESS=localhost
DB_NAME=ścieżka/do/NOURISHMENT.FDB
SERVER_PORT=8080
JWT_SECRET=twój_sekret
JWT_ISSUER=nourishment
JWT_AUDIENCE=nourishment-app
CORS_ALLOW_ORIGINS_LIST=http://localhost:3000
LOG_LEVEL=debug
SERVER_MODE=debug
```

## Struktura projektu

```
├── cmd/nourishment/     # Punkt wejścia aplikacji (main.go)
├── docs/                # Swagger (generowane automatycznie)
├── internal/
│   ├── api/             # Handlery HTTP
│   ├── auth/            # Logika autoryzacji (JWT)
│   ├── database/        # Repozytoria bazodanowe
│   ├── logic/           # Logika biznesowa
│   ├── mealDomain/      # Modele domenowe
│   ├── mealOptimizer/   # Optymalizacja posiłków (AI)
│   ├── modules/         # Rejestracja modułów (Kernel Pattern)
│   └── logging/         # Konfiguracja logowania
├── kernel/              # Interfejsy i kontrakty modułów
└── postman/             # Kolekcja Postman do testowania API
```

## Dokumentacja API

Po uruchomieniu serwera dokumentacja Swagger jest dostępna pod:

```
http://localhost:8080/swagger/index.html
```

## Powiązane repozytoria

- [Frontend (React)](https://github.com/Marekt94/Nourishment-2-0-Frontend) — interfejs użytkownika
