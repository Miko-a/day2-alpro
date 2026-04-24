# Laporan Penjelasan Challenge A dan Challenge B - Day 2 ALPRO

## Pendahuluan

Repository `day2-alpro` merupakan project backend berbasis Go yang menggunakan pola arsitektur berlapis. Struktur utama pada fitur user dipisahkan menjadi beberapa layer, yaitu `routes`, `controller`, `service`, `repository`, `dto`, dan `entities`.

Pada bagian bawah README utama, terdapat dua challenge yang harus diimplementasikan:

1. **Challenge A**: `GET /users/:id`
2. **Challenge B**: `GET /users`

Kedua challenge ini berfokus pada proses pengambilan data user dari database melalui REST API.

---

## Struktur Program yang Berkaitan

Fitur user berada pada folder:

```text
modules/user
├── controller
│   └── user_controller.go
├── dto
│   └── user_dto.go
├── repository
│   └── user_repository.go
├── service
│   └── user_service.go
├── validation
└── routes.go
```

Selain itu, struktur data user didefinisikan pada:

```text
database/entities/user_entity.go
```

Entity `User` menyimpan data utama user, yaitu:

```go
type User struct {
    Common
    Name     string `gorm:"not null" json:"name"`
    Email    string `gorm:"unique;not null" json:"email"`
    Password string `gorm:"not null" json:"-"`
    Role     string `gorm:"default:'user'" json:"role"`
}
```

Field `Password` menggunakan tag `json:"-"`, sehingga password tidak ikut dikirim pada response JSON. Ini penting karena password merupakan data sensitif.

---

## Alur Umum Request

Kedua challenge mengikuti alur backend yang sama:

```text
Client
  ↓
Route
  ↓
Controller
  ↓
Service
  ↓
Repository
  ↓
Database
```

Penjelasan tiap layer:

| Layer | Tanggung Jawab |
|---|---|
| Route | Menentukan endpoint API dan menghubungkannya ke controller |
| Controller | Membaca request HTTP dan mengirim response JSON |
| Service | Menjalankan business logic |
| Repository | Melakukan query ke database menggunakan GORM |
| Entity | Merepresentasikan struktur tabel database |

---

# Challenge A - GET /users/:id

## Tujuan Challenge

Challenge A bertujuan untuk membuat endpoint yang mengambil satu data user berdasarkan ID.

Endpoint yang diminta:

```http
GET /users/:id
```

Contoh request:

```http
GET /users/1
```

Jika user dengan ID tersebut ditemukan, API mengembalikan data user. Jika tidak ditemukan, API harus mengembalikan status `404 Not Found`.

---

## Konsep URL Parameter

Pada endpoint `GET /users/:id`, bagian `:id` adalah parameter dinamis.

Contoh:

```text
/users/1  → id = 1
/users/5  → id = 5
/users/10 → id = 10
```

Nilai ID tersebut dibaca oleh controller menggunakan Gin:

```go
idParam := c.Param("id")
```

Karena nilai dari URL selalu berbentuk string, ID perlu dikonversi ke integer sebelum digunakan untuk query database.

---

## Implementasi pada Repository

Repository bertugas mengambil data user dari database.

Method yang dibutuhkan:

```go
func (r *UserRepository) FindByID(id uint) (*entities.User, error) {
    var user entities.User
    err := r.db.First(&user, id).Error
    return &user, err
}
```

Penjelasan:

- `var user entities.User` membuat variabel kosong untuk menampung hasil query.
- `r.db.First(&user, id)` mencari user berdasarkan primary key.
- Jika data ditemukan, GORM mengisi variabel `user`.
- Jika data tidak ditemukan, GORM mengembalikan error.
- Function mengembalikan pointer ke user dan error.

Query yang secara konsep dijalankan:

```sql
SELECT * FROM users WHERE id = ? LIMIT 1;
```

---

## Implementasi pada Service

Pada layer service, kode yang digunakan:

```go
func (s *UserService) GetUserByID(id uint) (*entities.User, error) {
    return s.repo.FindByID(id)
}
```

Service tidak langsung melakukan query database. Service hanya memanggil repository. Pemisahan ini membuat kode lebih rapi karena business logic dan akses database tidak dicampur.

---

## Implementasi pada Controller

Controller bertugas membaca ID dari URL, melakukan validasi, memanggil service, lalu mengirim response.

Contoh implementasi:

```go
func (ctrl *UserController) GetUserByID(c *gin.Context) {
    idParam := c.Param("id")

    id, err := strconv.Atoi(idParam)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "ID user tidak valid")
        return
    }

    user, err := ctrl.service.GetUserByID(uint(id))
    if err != nil {
        utils.ErrorResponse(c, http.StatusNotFound, "User tidak ditemukan")
        return
    }

    utils.SuccessResponse(c, http.StatusOK, "User berhasil ditemukan", user)
}
```

Penjelasan:

1. `c.Param("id")` mengambil ID dari URL.
2. `strconv.Atoi(idParam)` mengubah string menjadi integer.
3. Jika ID tidak valid, server mengembalikan `400 Bad Request`.
4. Jika user tidak ditemukan, server mengembalikan `404 Not Found`.
5. Jika berhasil, server mengembalikan `200 OK`.

---

## Response Berhasil

Contoh response ketika user ditemukan:

```json
{
  "status": "success",
  "message": "User berhasil ditemukan",
  "data": {
    "id": 1,
    "name": "Budi",
    "email": "budi@example.com",
    "role": "user"
  }
}
```

---

## Response Gagal

Jika ID tidak valid:

```json
{
  "status": "error",
  "message": "ID user tidak valid"
}
```

Jika user tidak ditemukan:

```json
{
  "status": "error",
  "message": "User tidak ditemukan"
}
```

---

# Challenge B - GET /users

## Tujuan Challenge

Challenge B bertujuan untuk membuat endpoint yang mengambil semua data user dari database.

Endpoint yang diminta:

```http
GET /users
```

Endpoint ini mengembalikan array JSON berisi daftar user.

---

## Implementasi pada Repository

Repository membutuhkan method untuk mengambil seluruh data user.

```go
func (r *UserRepository) FindAll() ([]entities.User, error) {
    var users []entities.User
    err := r.db.Find(&users).Error
    return users, err
}
```

Penjelasan:

- `var users []entities.User` membuat slice untuk menampung banyak user.
- `r.db.Find(&users)` mengambil semua record dari tabel users.
- Function mengembalikan slice user dan error.

Query yang secara konsep dijalankan:

```sql
SELECT * FROM users;
```

---

## Implementasi pada Service

Pada layer service, kode yang digunakan:

```go
func (s *UserService) GetAllUsers() ([]entities.User, error) {
    return s.repo.FindAll()
}
```

Service meneruskan permintaan dari controller ke repository. Karena challenge ini hanya mengambil data tanpa business logic tambahan, service cukup memanggil `FindAll()`.

---

## Implementasi pada Controller

Controller untuk mengambil seluruh user:

```go
func (ctrl *UserController) GetAllUsers(c *gin.Context) {
    users, err := ctrl.service.GetAllUsers()
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data user")
        return
    }

    utils.SuccessResponse(c, http.StatusOK, "Data user berhasil diambil", users)
}
```

Penjelasan:

1. Controller memanggil `ctrl.service.GetAllUsers()`.
2. Jika terjadi error saat query database, API mengembalikan `500 Internal Server Error`.
3. Jika berhasil, API mengembalikan `200 OK` beserta array user.

---

## Response Berhasil

Contoh response ketika data user berhasil diambil:

```json
{
  "status": "success",
  "message": "Data user berhasil diambil",
  "data": [
    {
      "id": 1,
      "name": "Budi",
      "email": "budi@example.com",
      "role": "user"
    },
    {
      "id": 2,
      "name": "Siti",
      "email": "siti@example.com",
      "role": "user"
    }
  ]
}
```

Jika belum ada user, response tetap berhasil tetapi data berupa array kosong:

```json
{
  "status": "success",
  "message": "Data user berhasil diambil",
  "data": []
}
```

---

# Pendaftaran Route

Agar kedua endpoint dapat diakses, route perlu didaftarkan pada `modules/user/routes.go`.

Contoh implementasi route:

```go
func RegisterUserRoutes(r *gin.RouterGroup, ctrl *controller.UserController, jwtSvc *authService.JWTService) {
    users := r.Group("/users")
    {
        users.POST("", ctrl.CreateUser)
        users.GET("/:id", ctrl.GetUserByID)
        users.GET("", ctrl.GetAllUsers)
    }
}
```

Jika endpoint diwajibkan login, route dapat ditambahkan middleware authentication sesuai struktur project.

Contoh:

```go
users.GET("/:id", middlewares.Authentication(jwtSvc), ctrl.GetUserByID)
users.GET("", middlewares.Authentication(jwtSvc), ctrl.GetAllUsers)
```

---

# Perbandingan Challenge A dan Challenge B

| Aspek | Challenge A | Challenge B |
|---|---|---|
| Endpoint | `GET /users/:id` | `GET /users` |
| Fungsi | Mengambil satu user | Mengambil semua user |
| Parameter | Menggunakan `id` dari URL | Tidak menggunakan parameter |
| Repository | `FindByID(id uint)` | `FindAll()` |
| Service | `GetUserByID(id uint)` | `GetAllUsers()` |
| Response sukses | Object user | Array user |
| Error utama | `404 Not Found` jika user tidak ditemukan | `500 Internal Server Error` jika query gagal |

---

# Catatan terhadap Kode Saat Ini

Berdasarkan struktur kode program, layer service sudah memiliki method:

```go
func (s *UserService) GetUserByID(id uint) (*entities.User, error)
func (s *UserService) GetAllUsers() ([]entities.User, error)
```

Namun agar program berjalan penuh, repository, controller, dan route juga perlu memastikan method berikut tersedia dan terhubung:

```go
// Repository
FindByID(id uint)
FindAll()

// Controller
GetUserByID(c *gin.Context)
GetAllUsers(c *gin.Context)

// Routes
GET /users/:id
GET /users
```

Jika salah satu bagian belum dibuat, program akan gagal compile atau endpoint tidak bisa diakses.

---
# Screenshot Hasil

<img width="1117" height="792" alt="image" src="https://github.com/user-attachments/assets/c885b598-14ab-4821-854f-c74af0d471c3" />

<img width="1104" height="824" alt="image" src="https://github.com/user-attachments/assets/6e715b90-817f-49ac-a467-52162be7363a" />



# Kesimpulan

Challenge A dan Challenge B memperkenalkan konsep dasar pembuatan endpoint `GET` pada backend Go dengan Gin dan GORM.

Challenge A mengajarkan cara mengambil satu data berdasarkan ID menggunakan URL parameter dan menangani kondisi data tidak ditemukan dengan status `404`.

Challenge B mengajarkan cara mengambil banyak data sekaligus dari database dan mengembalikannya dalam bentuk array JSON.

Kedua challenge ini memperkuat pemahaman tentang alur request pada backend, yaitu dari route ke controller, dari controller ke service, dari service ke repository, lalu dari repository ke database.
