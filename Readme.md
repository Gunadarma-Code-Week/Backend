# Backend GCW

Panduan ini fokus ke cara menjalankan service backend secara lokal.

## Prasyarat

- Go `1.23.x` (sesuai `go.mod`)
- PostgreSQL `15+`
- Docker (opsional, jika mau menjalankan Postgres/container)
- Port `8000` (API) dan `5432` (Postgres) tidak dipakai proses lain

## 1) Setup Environment

Jalankan dari folder `Backend`:

```bash
cd Backend
cp .env.example .env
```

Nilai minimal yang wajib ada agar backend bisa boot:

```env
ENVIRONMENT=local
PORT=8000
CORS_ORIGIN=http://localhost:3000

DB_HOST=127.0.0.1
DB_PORT=5432
DB_NAME=gcw_db
DB_USER=gcw_user
DB_PASS=changeme

API_BASE_URL=/api/v1/gcw/resources
JWT_SECRET=change_me_secret
JWT_REFRESH_SECRET=change_me_refresh_secret
JWT_ISSUER=gcw
```

Catatan:
- `ENVIRONMENT=local` akan menjalankan auto-migrate tabel saat startup.
- File `.env.example` belum memuat `MIDTRANS_SERVER_KEY` dan `MIDTRANS_ENVIRONMENT`. Tambahkan manual jika butuh fitur pembayaran.

## 2) Jalankan Database

Gunakan salah satu opsi berikut.

### Opsi A - Postgres lokal (native install)

Pastikan database dan user berikut tersedia:
- Database: `gcw_db`
- User: `gcw_user`
- Password: `changeme`

### Opsi B - Postgres via Docker (direkomendasikan)

```bash
docker run --name gcw-postgres \
  -e POSTGRES_USER=gcw_user \
  -e POSTGRES_PASSWORD=changeme \
  -e POSTGRES_DB=gcw_db \
  -p 5432:5432 \
  -d postgres:15
```

## 3) Jalankan Backend

Dari folder `Backend`:

### Opsi A - Run normal

```bash
go mod download
go run main.go
```

### Opsi B - Hot reload dengan Air

File `.air.toml` sudah tersedia.

```bash
go install github.com/air-verse/air@latest
air
```

### Opsi C - Full Docker Compose (`db` + `app`)

1. Pastikan file `.env` ada.
2. Ubah `DB_HOST=db` (karena app dan DB berada di network compose yang sama).
3. Tambahkan juga:

```env
POSTGRES_USER=gcw_user
POSTGRES_DB=gcw_db
```

Lalu jalankan:

```bash
docker compose up --build
```

## 4) Cek Service Berjalan

- Health check: `http://localhost:8000/api/v1/gcw/resources/ping`
- Swagger docs: `http://localhost:8000/swagger/index.html`

## 5) Environment Opsional (Sesuai Fitur)

- `MIDTRANS_SERVER_KEY`, `MIDTRANS_ENVIRONMENT`: pembayaran
- `EMAIL_HOST`, `EMAIL_PORT`, `EMAIL_USERNAME`, `EMAIL_PASSWORD`, `EMAIL_FROM`: kirim email
- `DOMJUDGE_URL`, `DOMJUDGE_CONTEST_ID`, `DOMJUDGE_USERNAME`, `DOMJUDGE_PASSWORD`: integrasi CP
- `GOOGLE_CLIENT_ID`: login Google

## 6) Panduan DOMjudge

Dokumentasi detail DOMjudge (docker, curl testing, akun PIC admin, setup upload soal, test checklist submission/upload/judging) tersedia di:

- `Backend/DOMJUDGE_USAGE.md`
