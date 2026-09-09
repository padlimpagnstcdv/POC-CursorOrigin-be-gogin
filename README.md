# be-go-gin

REST API CRUD sederhana dengan **Go + Gin**, disusun memakai arsitektur berlapis
**Presentation – Core – Data**. Penyimpanan datanya in-memory (data dummy), jadi
proyek ini bisa langsung dijalankan tanpa perlu database.

## Daftar Isi / testing conflict

- [Teknologi](#teknologi)
- [Arsitektur](#arsitektur)
- [Struktur Folder](#struktur-folder)
- [Alur Request](#alur-request)
- [Menjalankan Aplikasi](#menjalankan-aplikasi)
- [Dokumentasi API](#dokumentasi-api)
- [Data Dummy](#data-dummy)
- [Validasi & Penanganan Error](#validasi--penanganan-error)
- [Mengganti Data Layer](#mengganti-data-layer)
- [Menambah Resource Baru](#menambah-resource-baru)
- [Catatan / Batasan](#catatan--batasan)

## Teknologi

| Komponen | Versi |
|---|---|
| Go | 1.26.5 |
| Gin | v1.12.0 |
| Penyimpanan | in-memory (`map` + `sync.RWMutex`) |

Nama module: `be-gin-go` (lihat [go.mod](go.mod)) — semua import internal memakai
prefix ini, contohnya `be-gin-go/internal/core/domain`.

## Arsitektur

Tiga lapisan dengan arah dependency yang selalu menuju ke dalam (**core**):

```
┌─────────────────────────────────────────────────┐
│  PRESENTATION  (internal/presentation/http)     │
│  router → handler → dto                         │
│  Tugas: baca HTTP request, panggil service,     │
│         terjemahkan error jadi status HTTP      │
└───────────────────────┬─────────────────────────┘
                        │  bergantung pada interface
                        ▼
┌─────────────────────────────────────────────────┐
│  CORE  (internal/core)                          │
│  domain  → entity + error bisnis                │
│  port    → interface (kontrak)                  │
│  service → business logic + validasi            │
│                                                 │
│  Tidak mengimpor Gin, tidak mengimpor           │
│  repository. Murni logika bisnis.               │
└───────────────────────▲─────────────────────────┘
                        │  mengimplementasikan interface
                        │
┌───────────────────────┴─────────────────────────┐
│  DATA  (internal/data/repository)               │
│  Implementasi penyimpanan (saat ini in-memory)  │
└─────────────────────────────────────────────────┘
```

Kunci desainnya ada di [internal/core/port/product.go](internal/core/port/product.go).
Dua interface di situ jadi batas antar lapisan:

- `ProductRepository` — kontrak yang **wajib dipenuhi data layer**. Core memanggil
  interface ini, bukan struct konkretnya.
- `ProductService` — kontrak yang **dipakai presentation layer**. Handler tidak
  tahu isi implementasi service.

Efeknya: core tidak tahu-menahu soal HTTP maupun soal cara data disimpan. Mau
ganti Gin ke Echo, atau ganti in-memory ke PostgreSQL, core tetap utuh.

## Struktur Folder

```
be-go-gin/
├── main.go                                    # entry point + wiring dependency
├── go.mod
├── README.md
└── internal/
    ├── core/
    │   ├── domain/
    │   │   └── product.go                     # entity Product + error domain
    │   ├── port/
    │   │   └── product.go                     # interface Repository & Service
    │   └── service/
    │       └── product.go                     # implementasi business logic
    ├── data/
    │   └── repository/
    │       └── product_memory.go              # penyimpanan in-memory + seed dummy
    └── presentation/
        └── http/
            ├── dto/
            │   └── product.go                 # request & response payload
            ├── handler/
            │   └── product.go                 # Gin handler
            └── router/
                └── router.go                  # registrasi route
```

Penjelasan per file:

| File | Isi |
|---|---|
| [main.go](main.go) | Merangkai dependency: repository → service → handler → router, lalu menjalankan server di `:8080` |
| [internal/core/domain/product.go](internal/core/domain/product.go) | Struct `Product` dan error `ErrProductNotFound`, `ErrInvalidProduct` |
| [internal/core/port/product.go](internal/core/port/product.go) | Interface `ProductRepository` dan `ProductService` |
| [internal/core/service/product.go](internal/core/service/product.go) | Logika `List`, `GetByID`, `Create`, `Update`, `Delete` + fungsi `validate` |
| [internal/data/repository/product_memory.go](internal/data/repository/product_memory.go) | Implementasi repository pakai `map[int64]domain.Product`, aman untuk akses paralel |
| [internal/presentation/http/dto/product.go](internal/presentation/http/dto/product.go) | `ProductRequest` (payload masuk) dan `Response` (bentuk response standar) |
| [internal/presentation/http/handler/product.go](internal/presentation/http/handler/product.go) | Handler tiap endpoint + `writeError` untuk mapping error → status HTTP |
| [internal/presentation/http/router/router.go](internal/presentation/http/router/router.go) | Route group `/api/products` dan endpoint `/ping` |

## Alur Request

Contoh `POST /api/products`:

```
HTTP Request
   │
   ▼
router.New                      route cocok → handler.Create
   │
   ▼
handler.Create                  ShouldBindJSON ke dto.ProductRequest
   │                            gagal bind → 400, berhenti di sini
   ▼
service.Create(name, price, stock)
   │                            validate() → ErrInvalidProduct kalau tidak valid
   │                            set CreatedAt & UpdatedAt
   ▼
repository.Create(product)      lastID++ , simpan ke map
   │
   ▼
handler                         balikkan 201 + dto.Response
```

Perhatikan bahwa handler mengirim tipe primitif (`name`, `price`, `stock`) ke
service, bukan struct DTO-nya. Jadi `dto` tetap milik presentation layer dan core
tidak perlu mengenalnya.

## Menjalankan Aplikasi

Install dependency:

```bash
go mod tidy
```

Jalankan:

```bash
go run .
```

Server aktif di `http://localhost:8080`. Cek cepat:

```bash
curl http://localhost:8080/ping
# {"message":"pong"}
```

Perintah lain yang berguna:

```bash
go build ./...              # kompilasi seluruh package
go vet ./...                # analisis statis
GIN_MODE=release go run .   # mode release, log lebih ringkas
```

## Dokumentasi API

Base URL: `http://localhost:8080`

| Method | Endpoint | Keterangan | Sukses |
|---|---|---|---|
| GET | `/ping` | Health check | 200 |
| GET | `/api/products` | Ambil semua product (urut naik by `id`) | 200 |
| GET | `/api/products/:id` | Ambil satu product | 200 |
| POST | `/api/products` | Tambah product baru | 201 |
| PUT | `/api/products/:id` | Perbarui product | 200 |
| DELETE | `/api/products/:id` | Hapus product | 200 |

Semua response memakai bentuk yang sama ([dto.Response](internal/presentation/http/dto/product.go)):

```json
{
  "message": "success",
  "data": {}
}
```

Field `data` dihilangkan (`omitempty`) kalau tidak ada isinya, misalnya pada
response delete dan pada response error.

### Entity Product

| Field | Tipe | Keterangan |
|---|---|---|
| `id` | int64 | Auto-increment, dibuat oleh repository |
| `name` | string | Wajib, tidak boleh kosong |
| `price` | float64 | Tidak boleh negatif |
| `stock` | int | Tidak boleh negatif |
| `created_at` | time | Diisi saat create, tidak berubah |
| `updated_at` | time | Diperbarui setiap update |

### GET /api/products

```bash
curl http://localhost:8080/api/products
```

```json
{
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "Kopi Arabica 250gr",
      "price": 85000,
      "stock": 25,
      "created_at": "2026-09-09T09:36:47.8026468+08:00",
      "updated_at": "2026-09-09T09:36:47.8026468+08:00"
    }
  ]
}
```

### GET /api/products/:id

```bash
curl http://localhost:8080/api/products/1
```

```json
{
  "message": "success",
  "data": {
    "id": 1,
    "name": "Kopi Arabica 250gr",
    "price": 85000,
    "stock": 25,
    "created_at": "2026-09-09T09:36:47.8026468+08:00",
    "updated_at": "2026-09-09T09:36:47.8026468+08:00"
  }
}
```

`404` kalau id tidak ada, `400` kalau id bukan angka.

### POST /api/products

```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Susu Kental","price":22000,"stock":10}'
```

`201 Created`:

```json
{
  "message": "product created",
  "data": {
    "id": 4,
    "name": "Susu Kental",
    "price": 22000,
    "stock": 10,
    "created_at": "2026-09-09T09:36:54.7588366+08:00",
    "updated_at": "2026-09-09T09:36:54.7588366+08:00"
  }
}
```

### PUT /api/products/:id

Update bersifat menyeluruh — ketiga field dikirim semua.

```bash
curl -X PUT http://localhost:8080/api/products/4 \
  -H "Content-Type: application/json" \
  -d '{"name":"Susu Kental Manis","price":25000,"stock":8}'
```

```json
{
  "message": "product updated",
  "data": {
    "id": 4,
    "name": "Susu Kental Manis",
    "price": 25000,
    "stock": 8,
    "created_at": "2026-09-09T09:36:54.7588366+08:00",
    "updated_at": "2026-09-09T09:36:54.8309432+08:00"
  }
}
```

`created_at` dipertahankan, hanya `updated_at` yang berubah.

### DELETE /api/products/:id

```bash
curl -X DELETE http://localhost:8080/api/products/1
```

```json
{ "message": "product deleted" }
```

## Data Dummy

Repository in-memory menanam tiga product saat aplikasi start
(lihat `NewProductMemoryRepository` di [product_memory.go](internal/data/repository/product_memory.go)):

| id | name | price | stock |
|---|---|---|---|
| 1 | Kopi Arabica 250gr | 85000 | 25 |
| 2 | Teh Hijau Premium | 45000 | 40 |
| 3 | Gula Aren Cair | 30000 | 15 |

Untuk mengubah data awal, sunting slice `seeds` di fungsi tersebut.

## Validasi & Penanganan Error

Validasi berjalan di **dua tempat**, dan ini disengaja:

1. **Presentation** — binding tag Gin di `dto.ProductRequest`
   (`binding:"required"`, `binding:"gte=0"`). Menyaring payload HTTP yang jelas
   salah sedini mungkin, termasuk JSON yang formatnya rusak.
2. **Core** — fungsi `validate()` di service. Menjaga aturan bisnis tetap berlaku
   walau service dipanggil dari luar HTTP (unit test, CLI, message consumer).

Konsekuensinya: untuk kasus yang sama, format pesan `400` bisa berbeda tergantung
lapisan mana yang menangkap lebih dulu. Binding Gin membalas pesan detail per
field, sementara service membalas `invalid product data`.

Pemetaan error domain ke status HTTP terpusat di `writeError`
([handler/product.go](internal/presentation/http/handler/product.go)):

| Error / kondisi | Status | `message` |
|---|---|---|
| `domain.ErrProductNotFound` | 404 | `product not found` |
| `domain.ErrInvalidProduct` | 400 | `invalid product data` |
| Binding JSON gagal | 400 | pesan detail dari validator |
| Param `:id` bukan angka | 400 | `invalid id` |
| Error lain | 500 | `internal server error` |

Error tak terduga sengaja tidak dibocorkan isinya ke client — hanya dibalas
`internal server error`.

## Mengganti Data Layer

Karena core hanya mengenal interface, mengganti penyimpanan tidak menyentuh
service maupun handler. Misal mau pindah ke MongoDB (driver-nya sudah ada di
`go.mod`):

1. Buat `internal/data/repository/product_mongo.go`.
2. Implementasikan kelima method `port.ProductRepository`: `FindAll`, `FindByID`,
   `Create`, `Update`, `Delete`.
3. Ganti satu baris di [main.go](main.go):

   ```go
   // sebelum
   productRepo := repository.NewProductMemoryRepository()

   // sesudah
   productRepo := repository.NewProductMongoRepository(db)
   ```

Selesai. Tidak ada file lain yang perlu diubah.

Interface yang sama juga memudahkan testing: buat repository tiruan di dalam test
untuk menguji service tanpa I/O sama sekali.

## Menambah Resource Baru

Untuk resource lain (misalnya `Category`), ikuti urutan lapisan yang sama:

1. `internal/core/domain/category.go` — entity + error domain.
2. `internal/core/port/category.go` — interface `CategoryRepository` & `CategoryService`.
3. `internal/core/service/category.go` — business logic.
4. `internal/data/repository/category_memory.go` — implementasi penyimpanan.
5. `internal/presentation/http/dto/category.go` — payload.
6. `internal/presentation/http/handler/category.go` — handler.
7. Daftarkan route di `router.New`, lalu rangkai dependency-nya di `main.go`.

## Catatan / Batasan

Hal-hal yang memang belum ada karena proyek ini sengaja dijaga sederhana:

- **Data tidak persisten.** Semua tersimpan di memory, jadi ikut hilang setiap
  kali server restart, termasuk ID counter yang kembali dari 3.
- **Belum ada test.** Struktur berbasis interface sudah siap untuk itu, tapi
  belum ditulis.
- **Belum ada autentikasi, rate limiting, maupun pagination** pada endpoint list.
- **Port di-hardcode** `:8080` di `main.go`, belum lewat environment variable.
- **Belum ada graceful shutdown.**
