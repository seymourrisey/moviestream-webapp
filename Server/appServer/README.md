# MovieStream Backend 

README ini merangkum fitur, arsitektur, environment variables, instruksi menjalankan, endpoint, dan catatan penting terkait backend yang ada di `Server/appServer`.

---

## Ringkasan singkat
- Bahasa: Go, framework: Gin.
- Database: MongoDB.
- Otentikasi: JWT (access + refresh), disimpan di cookie (Secure, HttpOnly).
- Role-based authorization: `ADMIN` dan `USER`.
- Integrasi LLM (via Groq/OpenAI-compatible endpoint) untuk menentukan ranking dari admin review.
- Jika MongoDB tidak tersedia, protected routes akan dinonaktifkan; server tetap berjalan untuk endpoint publik.

---

## Struktur utama
- `main.go` — entrypoint, konfigurasi CORS, inisialisasi DB, pendaftaran routes, dan menjalankan server.
- `database/` — helper koneksi MongoDB dan helper buka collection.
- `routes/` — pendaftaran routes (protected & unprotected).
- `controllers/` — handler HTTP: user & movie logic.
- `middleware/` — auth middleware (memvalidasi JWT dari cookie).
- `models/` — definisi struktur data (User, Movie, Genre, Ranking).
- `utils/` — utilitas token/JWT dan helper untuk mengambil user/role dari context.

---

## Fitur logic (ringkasan)
- User management:
  - Registrasi user dengan validasi input dan hashing password (`bcrypt`).
  - Login: verifikasi password, generate access & refresh token, simpan token ke DB, set cookie.
  - Logout: clear token di DB dan hapus cookie pada response.
  - Refresh token: validasi refresh token dari cookie, generate token baru, update DB, set cookie.
- Movie management:
  - Ambil daftar film (`GET /movies`) — publik.
  - Ambil film berdasarkan `imdb_id` — protected.
  - Tambah film (`POST /addmovie`) — protected.
  - Update admin review (`PATCH /updatereview/:imdb_id`) — hanya `ADMIN`. Menggunakan LLM untuk menentukan `ranking_name` dan mengupdate `ranking_value`.
  - Rekomendasi film (`GET /recommendedmovies`) — protected; berdasarkan favourite genres user, diurutkan menurut `ranking.ranking_value`.
- Genre & rankings:
  - Ambil daftar genres (`GET /genres`) — publik.
  - Ambil collection `rankings` untuk mapping nama-nilai ranking (dipakai oleh LLM flow).
- Integrasi LLM:
  - `GetReviewRankings` membangun prompt menggunakan `BASE_PROMPT_TEMPLATE` + daftar ranking yang tersedia, lalu mengirim ke LLM (client `tmc/langchaingo/llms/openai`) memakai `GROQ_API_KEY` dan base URL `https://api.groq.com/openai/v1`.
  - Response LLM di-mapping ke `ranking_value` berdasarkan dokumen di collection `rankings`.
- Safety behavior: Protected routes tidak terdaftar jika MongoDB tidak dapat di-ping saat startup.

---

## Environment variables (.env)
- `MONGODB_URI` — connection string MongoDB (wajib).
- `DATABASE_NAME` — nama database (digunakan `OpenCollection`).
- `SECRET_KEY` — secret untuk sign access token JWT.
- `SECRET_REFRESH_KEY` — secret untuk sign refresh token JWT.
- `GROQ_API_KEY` — API key untuk Groq/OpenAI endpoint LLM.
- `BASE_PROMPT_TEMPLATE` — template prompt untuk LLM (harus berisi placeholder `{rankings}`).
- `RECOMMENDED_MOVIES_COUNT` — jumlah maksimum rekomendasi film (default 5).
- `ALLOWED_ORIGINS` — daftar origin untuk CORS (comma-separated). Jika kosong => `*`.
- `PORT` — port untuk menjalankan server (default `8080`).


---

## Persyaratan & dependency
- Go (versi yang kompatibel dengan `go.mod` di project).
- MongoDB instance (atau connection string ke cloud Mongo).
- Access ke API Groq/OpenAI jika Anda ingin fitur LLM berjalan.
- Pastikan `SECRET_KEY` dan `SECRET_REFRESH_KEY` ter-set di environment.

---

## Menjalankan secara lokal
1. Salin contoh `.env` (jika ada) dan set environment variables yang diperlukan.
2. Jalankan:
   - Di folder `Server/appServer`:
     - `go mod tidy`
     - `go run main.go`
3. Server akan menggunakan port di `PORT` atau `8080` jika tidak di-set.

---

## Endpoint (ringkasan)
Public:
- `GET /hello` — health check.
- `GET /movies` — ambil semua movies.
- `POST /register` — registrasi user.
  - Body: JSON sesuai `models.User` (harus include `favourite_genres`).
- `POST /login` — login user.
  - Body: `{ "email": "...", "password": "..." }`
  - Response: `UserResponse` JSON (token disimpan di cookie).
  - Cookie yang diset: `access_token` (24h), `refresh_token` (7d).
- `POST /logout` — logout user.
  - Body: `{ "user_id": "..." }`
  - Hapus token di DB dan clear cookies.
- `GET /genres` — daftar genres.
- `POST /refresh` — refresh token flow.
  - Mengambil `refresh_token` dari cookie, memvalidasi, dan mengeluarkan tokens baru.

Protected (membutuhkan cookie `access_token` yang valid):
- `GET /movie/:imdb_id` — ambil movie tertentu.
- `POST /addmovie` — tambah movie (validasi menggunakan `validator`).
- `GET /recommendedmovies` — rekomendasi film berdasarkan favourite genres user.
- `PATCH /updatereview/:imdb_id` — ADMIN only; body: `{ "admin_review": "..." }` — sistem akan memanggil LLM untuk mendapatkan ranking dan mengupdate dokumen movie.

Contoh singkat (curl-like, jalankan di shell yang mendukung cookie):
  - Login (mengembalikan cookie yang diset oleh server):
    $ curl -X POST -H "Content-Type: application/json" -d '{"email":"you@example.com","password":"secret"}' http://localhost:8080/login -c cookies.txt
  - Mengakses protected endpoint:
    $ curl -X GET http://localhost:8080/recommendedmovies -b cookies.txt

> Catatan: server menggunakan cookie `Secure: true` sehingga di environment non-HTTPS (local) Anda mungkin perlu menyesuaikan atau men-test dengan client yang bisa menerima cookie lewat HTTP.

---

## Models (singkat)
- `User`:
  - `user_id`, `first_name`, `last_name`, `email`, `password` (hashed), `role` (`ADMIN` atau `USER`), `favourite_genres` (array).
- `Movie`:
  - `imdb_id`, `title`, `poster_path` (url), `youtube_id`, `genre` (array of Genre), `admin_review`, `ranking` (Ranking struct).
- `Ranking`:
  - `ranking_name`, `ranking_value` (dipakai untuk sorting & pemetaan).

---

## Keamanan & catatan operasional
- Cookie di-set dengan `Secure: true` & `HttpOnly: true` — butuh HTTPS di production.
- Pastikan `SECRET_KEY`, `SECRET_REFRESH_KEY`, dan `GROQ_API_KEY` aman (environment variables / secret manager).
- LLM prompt & data: berhati-hati jika prompt mengandung data sensitif.
- Jika MongoDB tidak tersedia saat startup, protected endpoints akan dinonaktifkan — namun beberapa public endpoint (mis. `/movies`) dapat tetap error jika memerlukan DB saat dipanggil.

---

## Troubleshooting cepat
- Error `MONGODB_URI NOT SET` → set `MONGODB_URI`.
- Jika login selalu gagal setelah set cookie: periksa domain/path/security cookie dan environment HTTPS vs HTTP.
- Jika LLM calls gagal: periksa `GROQ_API_KEY`, `BASE_PROMPT_TEMPLATE`, dan akses ke `https://api.groq.com/openai/v1`.

---

 
