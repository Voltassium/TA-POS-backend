# Dokumentasi API SeliPOS (Point of Sale)

Dokumen ini berisi spesifikasi teknis lengkap dari API SeliPOS. Dokumentasi ini disusun secara terstruktur untuk menjelaskan endpoint, format request/response, autentikasi, serta aturan hak akses pengguna. Informasi ini juga dapat digunakan sebagai referensi untuk penulisan laporan Tugas Akhir (TA).

---

## 1. Informasi Umum

- **Base URL**:
  - Pengembangan Lokal: `http://localhost:3000/v1`
  - Produksi: `https://[domain-produksi]/v1`
- **Format Data**: Seluruh request dan response menggunakan format **JSON** (`Content-Type: application/json`).
- **Autentikasi**: Menggunakan **JWT (JSON Web Token)** dengan mekanisme Bearer Token pada HTTP Header:
  ```http
  Authorization: Bearer <access_token>
  ```

---

## 2. Format Response Global

API SeliPOS menerapkan struktur response standar baik untuk response sukses maupun response error.

### A. Response Sukses (Standard Success)
```json
{
  "status": 200,
  "message": "Success message description",
  "data": { ... }
}
```
*Catatan: Nilai `status` mencerminkan HTTP Status Code, `message` berisi deskripsi singkat, dan `data` berisi payload response utama.*

### B. Response Error (Standard Error)
```json
{
  "status": 400,
  "message": "Error message description",
  "error": "Detailed validation or system error info"
}
```
*Catatan: Field `error` dapat berisi detail error validasi atau string error sistem.*

---

## 3. Matriks Hak Akses (Role & Authorization)

Sistem POS ini menggunakan Role-Based Access Control (RBAC). Berikut adalah 4 tingkatan hak akses yang terdaftar pada sistem:
1. **Superadmin**: Memiliki hak akses penuh ke seluruh data toko, manajemen pengguna global, serta pengaturan toko.
2. **Owner (Pemilik)**: Memiliki hak penuh untuk mengelola data tokonya sendiri (produk, transaksi, pengeluaran, staf, dan statistik dashboard).
3. **Chef (Koki)**: Bertanggung jawab atas pengelolaan antrean dapur dan penyajian makanan (modul dapur/kitchen).
4. **Staff (Kasir)**: Bertanggung jawab melakukan transaksi penjualan (kasir), mencatat pengeluaran harian, dan mencatat pembayaran.

| Modul / Fitur | Endpoint Prefix | Superadmin | Owner | Chef | Staff |
| :--- | :--- | :---: | :---: | :---: | :---: |
| Autentikasi & Registrasi | `/v1/authentications` | V | V | V | V |
| Manajemen Toko (Store) | `/v1/stores` | V | - | - | - |
| Manajemen User / Staf | `/v1/users` | V | V | - | - |
| Kategori Produk | `/v1/categories` | V | V | V (Read) | V (Read) |
| Produk & Stok | `/v1/products` | V | V | V (Read) | V (Read) |
| Transaksi / Pesanan | `/v1/orders` | V | V | V | V |
| Pembayaran (Kasir) | `/v1/payments` | V | V | - | V |
| Pengeluaran Operasional | `/v1/pengeluaran` | V | V | V | V |
| Dashboard & Statistik | `/v1/statistics` | V | V | - | - |

---

## 4. Parameter Pagination Global (Query Parameter)

Untuk endpoint yang mengembalikan daftar data dalam jumlah besar (list), API SeliPOS menggunakan struktur pagination berikut sebagai query parameter:

- `page` (int, default: 1): Halaman data yang ingin diambil.
- `page_size` (int, default: 10): Jumlah data per halaman.
- `order_by` (string, default: `updated_at`): Kolom pengurutan data.
- `order_dir` (string, default: `desc`, pilihan: `asc` atau `desc`): Arah pengurutan.
- `search` (string, opsional): Kata kunci untuk pencarian data global.

### Response Data Terpaginasi:
```json
{
  "status": 200,
  "message": "Success",
  "data": {
    "current_page": 1,
    "page_size": 10,
    "total_items": 25,
    "total_pages": 3,
    "has_previous": false,
    "has_next": true,
    "data": [ ... ]
  }
}
```

---

## 5. Daftar Endpoint API

### 5.1. Modul Autentikasi (`/v1/authentications`)

#### 1. Registrasi User Baru (`POST /v1/authentications/register`)
- **Tingkat Akses**: Publik
- **Deskripsi**: Pendaftaran pengguna baru. Jika parameter `store_name` disertakan, sistem akan otomatis membuat satu toko (store) baru dan menautkannya ke user tersebut sebagai Owner.
- **Request Body**:
  ```json
  {
    "email": "owner@example.com",
    "password": "securepassword",
    "role": "owner",
    "store_name": "Kopi Makmur",
    "store_address": "Jl. Merdeka No. 45"
  }
  ```
- **Response Sukses (201 Created)**:
  ```json
  {
    "status": 201,
    "message": "User created successfully",
    "data": {
      "id": "e3057b32-8413-4c91-b3b3-8c467a9ad21b",
      "email": "owner@example.com",
      "role": "owner",
      "created_at": "2026-07-10T12:00:00Z",
      "updated_at": "2026-07-10T12:00:00Z"
    }
  }
  ```

#### 2. Login Pengguna (`POST /v1/authentications/login`)
- **Tingkat Akses**: Publik
- **Deskripsi**: Verifikasi kredensial pengguna untuk mendapatkan Access Token (JWT) dan Refresh Token.
- **Request Body**:
  ```json
  {
    "email": "owner@example.com",
    "password": "securepassword"
  }
  ```
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Login successful",
    "data": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
  }
  ```

#### 3. Refresh Access Token (`POST /v1/authentications/refresh-token`)
- **Tingkat Akses**: Publik
- **Deskripsi**: Memperbarui Access Token yang sudah kedaluwarsa dengan menggunakan Refresh Token yang valid.
- **Request Body**:
  ```json
  {
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Token refreshed successfully",
    "data": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
  }
  ```

#### 4. Logout Pengguna (`POST /v1/authentications/logout`)
- **Tingkat Akses**: Privat (Membutuhkan Bearer Token)
- **Deskripsi**: Menghapus sesi login dan menonaktifkan token yang digunakan saat ini.
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Logout successful",
    "data": null
  }
  ```

---

### 5.2. Modul Toko / Stores (`/v1/stores`)

*Catatan: Seluruh endpoint manajemen toko di bawah ini hanya dapat diakses oleh **Superadmin**.*

#### 1. Membuat Toko Baru (`POST /v1/stores`)
- **Tingkat Akses**: Superadmin
- **Request Body**:
  ```json
  {
    "name": "Kopi Makmur Cabang 2",
    "address": "Jl. Sudirman No. 12"
  }
  ```
- **Response Sukses (201 Created)**:
  ```json
  {
    "status": 201,
    "message": "Store created successfully",
    "data": {
      "id": 2,
      "name": "Kopi Makmur Cabang 2",
      "address": "Jl. Sudirman No. 12",
      "created_at": "2026-07-10T12:05:00Z",
      "updated_at": "2026-07-10T12:05:00Z"
    }
  }
  ```

#### 2. Mendapatkan Daftar Toko (`GET /v1/stores`)
- **Tingkat Akses**: Superadmin
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Success",
    "data": [
      {
        "id": 1,
        "name": "Kopi Makmur",
        "address": "Jl. Merdeka No. 45",
        "created_at": "2026-07-10T12:00:00Z",
        "updated_at": "2026-07-10T12:00:00Z"
      }
    ]
  }
  ```

#### 3. Mendapatkan Detail Toko (`GET /v1/stores/:id`)
- **Tingkat Akses**: Superadmin
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Success",
    "data": {
      "id": 1,
      "name": "Kopi Makmur",
      "address": "Jl. Merdeka No. 45",
      "created_at": "2026-07-10T12:00:00Z",
      "updated_at": "2026-07-10T12:00:00Z"
    }
  }
  ```

#### 4. Memperbarui Data Toko (`PUT /v1/stores/:id`)
- **Tingkat Akses**: Superadmin
- **Request Body**:
  ```json
  {
    "name": "Kopi Makmur Premium",
    "address": "Jl. Merdeka Baru No. 46"
  }
  ```
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Store updated successfully",
    "data": {
      "id": 1,
      "name": "Kopi Makmur Premium",
      "address": "Jl. Merdeka Baru No. 46",
      "created_at": "2026-07-10T12:00:00Z",
      "updated_at": "2026-07-10T12:10:00Z"
    }
  }
  ```

#### 5. Menghapus Toko (`DELETE /v1/stores/:id`)
- **Tingkat Akses**: Superadmin
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Store deleted successfully",
    "data": null
  }
  ```

---

### 5.3. Modul Pengguna / Users (`/v1/users`)

#### 1. Mendapatkan Profil User Sedang Login (`GET /v1/users/profile`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Success",
    "data": {
      "id": "e3057b32-8413-4c91-b3b3-8c467a9ad21b",
      "email": "owner@example.com",
      "role": "owner",
      "created_at": "2026-07-10T12:00:00Z",
      "updated_at": "2026-07-10T12:00:00Z"
    }
  }
  ```

#### 2. Membuat User Baru oleh Admin/Owner (`POST /v1/users`)
- **Tingkat Akses**: Superadmin, Owner
- **Deskripsi**: Menambahkan staf baru (misal kasir atau koki) ke dalam toko yang bersangkutan.
- **Request Body**:
  ```json
  {
    "email": "kasir1@example.com",
    "password": "kasirpassword",
    "role": "staff"
  }
  ```
- **Response Sukses (201 Created)**:
  ```json
  {
    "status": 201,
    "message": "User created successfully",
    "data": {
      "id": "18f95c43-4e89-4bc7-b84a-9ef85cf901ac",
      "email": "kasir1@example.com",
      "role": "staff",
      "created_at": "2026-07-10T12:15:00Z",
      "updated_at": "2026-07-10T12:15:00Z"
    }
  }
  ```

#### 3. Mendapatkan Daftar User (`GET /v1/users`)
- **Tingkat Akses**: Superadmin, Owner
- **Query Parameter**: Mendukung standard pagination (`page`, `page_size`, `order_by`, `order_dir`, `search`).
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Success",
    "data": {
      "current_page": 1,
      "page_size": 10,
      "total_items": 2,
      "total_pages": 1,
      "has_previous": false,
      "has_next": false,
      "data": [
        {
          "id": "e3057b32-8413-4c91-b3b3-8c467a9ad21b",
          "email": "owner@example.com",
          "role": "owner",
          "created_at": "2026-07-10T12:00:00Z",
          "updated_at": "2026-07-10T12:00:00Z"
        },
        {
          "id": "18f95c43-4e89-4bc7-b84a-9ef85cf901ac",
          "email": "kasir1@example.com",
          "role": "staff",
          "created_at": "2026-07-10T12:15:00Z",
          "updated_at": "2026-07-10T12:15:00Z"
        }
      ]
    }
  }
  ```

#### 4. Mendapatkan Detail User (`GET /v1/users/:id`)
- **Tingkat Akses**: Superadmin, Owner
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Success",
    "data": {
      "id": "18f95c43-4e89-4bc7-b84a-9ef85cf901ac",
      "email": "kasir1@example.com",
      "role": "staff",
      "created_at": "2026-07-10T12:15:00Z",
      "updated_at": "2026-07-10T12:15:00Z"
    }
  }
  ```

#### 5. Memperbarui Data User (`PUT /v1/users/:id`)
- **Tingkat Akses**: Superadmin, Owner
- **Request Body**:
  ```json
  {
    "email": "kasir_baru@example.com",
    "role": "staff"
  }
  ```
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "User updated successfully",
    "data": {
      "id": "18f95c43-4e89-4bc7-b84a-9ef85cf901ac",
      "email": "kasir_baru@example.com",
      "role": "staff",
      "created_at": "2026-07-10T12:15:00Z",
      "updated_at": "2026-07-10T12:20:00Z"
    }
  }
  ```

#### 6. Menghapus User (`DELETE /v1/users/:id`)
- **Tingkat Akses**: Superadmin, Owner
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "User deleted successfully",
    "data": null
  }
  ```

---

### 5.4. Modul Kategori Produk (`/v1/categories`)

#### 1. Mendapatkan Daftar Kategori (`GET /v1/categories`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Success",
    "data": [
      {
        "id": "8f8b0e8b-d7d8-4f8e-bd3a-59b4c09d5eb6",
        "name": "Makanan Utama",
        "created_at": "2026-07-10T12:00:00Z",
        "updated_at": "2026-07-10T12:00:00Z"
      }
    ]
  }
  ```

#### 2. Mendapatkan Detail Kategori (`GET /v1/categories/:id`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Success",
    "data": {
      "id": "8f8b0e8b-d7d8-4f8e-bd3a-59b4c09d5eb6",
      "name": "Makanan Utama",
      "created_at": "2026-07-10T12:00:00Z",
      "updated_at": "2026-07-10T12:00:00Z"
    }
  }
  ```

#### 3. Membuat Kategori Baru (`POST /v1/categories`)
- **Tingkat Akses**: Superadmin, Owner
- **Request Body**:
  ```json
  {
    "name": "Minuman"
  }
  ```
- **Response Sukses (201 Created)**:
  ```json
  {
    "status": 201,
    "message": "Category created successfully",
    "data": {
      "id": "2b314d11-53ab-41c1-92b0-8e6b91fa01ab",
      "name": "Minuman",
      "created_at": "2026-07-10T12:30:00Z",
      "updated_at": "2026-07-10T12:30:00Z"
    }
  }
  ```

#### 4. Memperbarui Kategori (`PUT /v1/categories/:id`)
- **Tingkat Akses**: Superadmin, Owner
- **Request Body**:
  ```json
  {
    "name": "Minuman Dingin"
  }
  ```
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Category updated successfully",
    "data": {
      "id": "2b314d11-53ab-41c1-92b0-8e6b91fa01ab",
      "name": "Minuman Dingin",
      "created_at": "2026-07-10T12:30:00Z",
      "updated_at": "2026-07-10T12:35:00Z"
    }
  }
  ```

#### 5. Menghapus Kategori (`DELETE /v1/categories/:id`)
- **Tingkat Akses**: Superadmin, Owner
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Category deleted successfully",
    "data": null
  }
  ```

---

### 5.5. Modul Produk & Stok (`/v1/products` & `/v1/stock-histories`)

API SeliPOS membedakan dua tipe produk:
- **Kulakan**: Produk yang dibeli dalam keadaan siap jual (memiliki modal harga beli `harga_beli` dan stok langsung bertambah saat dibeli).
- **Olahan**: Produk yang dimasak/diracik langsung di tempat (tidak memiliki `harga_beli` modal pembelian secara langsung dan stok dikelola terpisah).

#### 1. Mendapatkan Daftar Produk (`GET /v1/products`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Query Parameter**:
  - Standard pagination (`page`, `page_size`, `order_by`, `order_dir`, `search`)
  - `category_id` (string/UUID, opsional): Memfilter produk berdasarkan kategori tertentu.
  - `product_type` (string, opsional, pilihan: `Kulakan`, `Olahan`): Memfilter jenis produk.
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Success",
    "data": {
      "current_page": 1,
      "page_size": 10,
      "total_items": 1,
      "total_pages": 1,
      "has_previous": false,
      "has_next": false,
      "data": [
        {
          "id": "c1f7a14e-0a56-42ab-ba41-a189b87df8b1",
          "category_id": "8f8b0e8b-d7d8-4f8e-bd3a-59b4c09d5eb6",
          "category_name": "Makanan Utama",
          "product_type": "Kulakan",
          "sku": "K001",
          "harga_beli": 12000,
          "name": "Nasi Goreng Kotak",
          "price": 20000,
          "is_available": true,
          "stock": 50,
          "created_at": "2026-07-10T12:00:00Z",
          "updated_at": "2026-07-10T12:00:00Z"
        }
      ]
    }
  }
  ```

#### 2. Mendapatkan Detail Produk (`GET /v1/products/:id`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Success",
    "data": {
      "id": "c1f7a14e-0a56-42ab-ba41-a189b87df8b1",
      "category_id": "8f8b0e8b-d7d8-4f8e-bd3a-59b4c09d5eb6",
      "category_name": "Makanan Utama",
      "product_type": "Kulakan",
      "sku": "K001",
      "harga_beli": 12000,
      "name": "Nasi Goreng Kotak",
      "price": 20000,
      "is_available": true,
      "stock": 50,
      "created_at": "2026-07-10T12:00:00Z",
      "updated_at": "2026-07-10T12:00:00Z"
    }
  }
  ```

#### 3. Membuat Produk Baru (`POST /v1/products`)
- **Tingkat Akses**: Superadmin, Owner
- **Request Body**:
  ```json
  {
    "category_id": "8f8b0e8b-d7d8-4f8e-bd3a-59b4c09d5eb6",
    "product_type": "Kulakan",
    "sku": "K001",
    "harga_beli": 12000,
    "name": "Nasi Goreng Kotak",
    "price": 20000,
    "is_available": true,
    "stock": 50
  }
  ```
- **Response Sukses (201 Created)**:
  ```json
  {
    "status": 201,
    "message": "Product created successfully",
    "data": {
      "id": "c1f7a14e-0a56-42ab-ba41-a189b87df8b1",
      "category_id": "8f8b0e8b-d7d8-4f8e-bd3a-59b4c09d5eb6",
      "category_name": "Makanan Utama",
      "product_type": "Kulakan",
      "sku": "K001",
      "harga_beli": 12000,
      "name": "Nasi Goreng Kotak",
      "price": 20000,
      "is_available": true,
      "stock": 50,
      "created_at": "2026-07-10T12:40:00Z",
      "updated_at": "2026-07-10T12:40:00Z"
    }
  }
  ```

#### 4. Memperbarui Produk (`PUT /v1/products/:id`)
- **Tingkat Akses**: Superadmin, Owner
- **Request Body**:
  ```json
  {
    "price": 22000,
    "is_available": true
  }
  ```
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Product updated successfully",
    "data": {
      "id": "c1f7a14e-0a56-42ab-ba41-a189b87df8b1",
      "category_id": "8f8b0e8b-d7d8-4f8e-bd3a-59b4c09d5eb6",
      "category_name": "Makanan Utama",
      "product_type": "Kulakan",
      "sku": "K001",
      "harga_beli": 12000,
      "name": "Nasi Goreng Kotak",
      "price": 22000,
      "is_available": true,
      "stock": 50,
      "created_at": "2026-07-10T12:40:00Z",
      "updated_at": "2026-07-10T12:45:00Z"
    }
  }
  ```

#### 5. Menghapus Produk (`DELETE /v1/products/:id`)
- **Tingkat Akses**: Superadmin, Owner
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Product deleted successfully",
    "data": null
  }
  ```

#### 6. Restock Produk Kulakan (`POST /v1/products/:id/restock`)
- **Tingkat Akses**: Superadmin, Owner
- **Deskripsi**: Endpoint khusus untuk melakukan restok produk berjenis 'Kulakan'. Request ini secara otomatis mencatat riwayat perubahan stok dan menambah entri pengeluaran modal pembelian secara otomatis.
- **Request Body**:
  ```json
  {
    "harga_beli": 12500,
    "jumlah_stok": 20
  }
  ```
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Product restocked successfully",
    "data": {
      "id": "c1f7a14e-0a56-42ab-ba41-a189b87df8b1",
      "category_id": "8f8b0e8b-d7d8-4f8e-bd3a-59b4c09d5eb6",
      "category_name": "Makanan Utama",
      "product_type": "Kulakan",
      "sku": "K001",
      "harga_beli": 12500,
      "name": "Nasi Goreng Kotak",
      "price": 22000,
      "is_available": true,
      "stock": 70,
      "created_at": "2026-07-10T12:40:00Z",
      "updated_at": "2026-07-10T12:50:00Z"
    }
  }
  ```

#### 7. Mendapatkan Riwayat Perubahan Stok Produk Tertentu (`GET /v1/products/:id/stock-histories`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Success",
    "data": {
      "current_page": 1,
      "page_size": 10,
      "total_items": 2,
      "total_pages": 1,
      "has_previous": false,
      "has_next": false,
      "data": [
        {
          "id": 12,
          "product_id": "c1f7a14e-0a56-42ab-ba41-a189b87df8b1",
          "product_name": "Nasi Goreng Kotak",
          "change": 20,
          "initial_stock": 50,
          "final_stock": 70,
          "reason": "Restock product",
          "created_at": "2026-07-10T12:50:00Z"
        }
      ]
    }
  }
  ```

#### 8. Mendapatkan Riwayat Perubahan Seluruh Stok Global (`GET /v1/stock-histories`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Response Sukses (200 OK)**:
  *(Struktur response sama seperti riwayat stok produk tertentu namun berisi seluruh log produk di toko)*

---

### 5.6. Modul Transaksi / Pesanan (`/v1/orders`)

Modul ini mengelola siklus pesanan pelanggan mulai dari pembuatan pesanan, modifikasi item di kasir/meja, hingga pemrosesan antrean dapur.
Daftar Status Pesanan (`OrderStatus`):
- `New`: Pesanan baru dibuat, belum dibayar.
- `Paid`: Pesanan telah dibayar lunas.
- `Cancelled`: Pesanan dibatalkan (stok dikembalikan).
- `Completed`: Pesanan selesai disajikan (makanan terkirim dan lunas).

#### 1. Mendapatkan Daftar Transaksi/Pesanan (`GET /v1/orders`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Query Parameter**:
  - Standard pagination (`page`, `page_size`, `order_by`, `order_dir`, `search`)
  - `status` (string, opsional): Memfilter pesanan dengan status tertentu (`New`, `Paid`, `Cancelled`, `Completed`).
  - `exclude_status` (string, opsional): Mengecualikan pesanan dengan status tertentu.
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Success",
    "data": {
      "current_page": 1,
      "page_size": 10,
      "total_items": 1,
      "total_pages": 1,
      "has_previous": false,
      "has_next": false,
      "data": [
        {
          "id": "e0bfa99b-6db3-43bb-81ef-7bf2b12361ac",
          "order_code": "ORD-20260710-001",
          "customer_name": "Anto",
          "table_id": 4,
          "staff_id": "18f95c43-4e89-4bc7-b84a-9ef85cf901ac",
          "staff_name": "kasir1@example.com",
          "total_amount": 40000,
          "status": "New",
          "created_at": "2026-07-10T13:00:00Z",
          "updated_at": "2026-07-10T13:00:00Z"
        }
      ]
    }
  }
  ```

#### 2. Mendapatkan Detail Lengkap Pesanan (`GET /v1/orders/:id`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Deskripsi**: Mengembalikan detail transaksi beserta daftar item pesanan dan data pembayaran (jika sudah bayar).
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Success",
    "data": {
      "id": "e0bfa99b-6db3-43bb-81ef-7bf2b12361ac",
      "order_code": "ORD-20260710-001",
      "customer_name": "Anto",
      "table_id": 4,
      "staff_id": "18f95c43-4e89-4bc7-b84a-9ef85cf901ac",
      "staff_name": "kasir1@example.com",
      "total_amount": 40000,
      "status": "New",
      "created_at": "2026-07-10T13:00:00Z",
      "updated_at": "2026-07-10T13:00:00Z",
      "items": [
        {
          "id": "a92b314d-fa12-4c91-b3b3-8c467a9ad21b",
          "order_id": "e0bfa99b-6db3-43bb-81ef-7bf2b12361ac",
          "product_id": "c1f7a14e-0a56-42ab-ba41-a189b87df8b1",
          "product_name": "Nasi Goreng Kotak",
          "quantity": 2,
          "unit_price": 20000,
          "subtotal": 40000,
          "served_qty": 0,
          "created_at": "2026-07-10T13:00:00Z",
          "updated_at": "2026-07-10T13:00:00Z"
        }
      ],
      "payment": null
    }
  }
  ```

#### 3. Membuat Transaksi Baru (`POST /v1/orders`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Request Body**:
  ```json
  {
    "table_id": 4,
    "customer_name": "Anto",
    "items": [
      {
        "product_id": "c1f7a14e-0a56-42ab-ba41-a189b87df8b1",
        "quantity": 2
      }
    ]
  }
  ```
- **Response Sukses (201 Created)**:
  *(Struktur response sama seperti detail lengkap pesanan)*

#### 4. Memperbarui Status Pesanan (`PATCH /v1/orders/:id/status`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Request Body**:
  ```json
  {
    "status": "Completed"
  }
  ```
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Order status updated successfully",
    "data": {
      "id": "e0bfa99b-6db3-43bb-81ef-7bf2b12361ac",
      "order_code": "ORD-20260710-001",
      "customer_name": "Anto",
      "table_id": 4,
      "staff_id": "18f95c43-4e89-4bc7-b84a-9ef85cf901ac",
      "total_amount": 40000,
      "status": "Completed",
      "created_at": "2026-07-10T13:00:00Z",
      "updated_at": "2026-07-10T13:10:00Z"
    }
  }
  ```

#### 5. Membatalkan Pesanan (`DELETE /v1/orders/:id`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Deskripsi**: Membatalkan pesanan. Mengubah status pesanan menjadi `Cancelled` dan mengembalikan stok produk yang telah dikurangi sebelumnya.
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Order cancelled successfully",
    "data": null
  }
  ```

#### 6. Menambah Item Baru ke Pesanan Aktif (`POST /v1/orders/:id/items`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Request Body**:
  ```json
  {
    "product_id": "2b314d11-53ab-41c1-92b0-8e6b91fa01ab",
    "quantity": 1
  }
  ```
- **Response Sukses (201 Created)**:
  ```json
  {
    "status": 201,
    "message": "Item added to order successfully",
    "data": {
      "id": "d9ab421b-41c1-53ab-92b0-8e6b91fa02bc",
      "order_id": "e0bfa99b-6db3-43bb-81ef-7bf2b12361ac",
      "product_id": "2b314d11-53ab-41c1-92b0-8e6b91fa01ab",
      "product_name": "Es Teh Manis",
      "quantity": 1,
      "unit_price": 5000,
      "subtotal": 5000,
      "served_qty": 0,
      "created_at": "2026-07-10T13:15:00Z",
      "updated_at": "2026-07-10T13:15:00Z"
    }
  }
  ```

#### 7. Menghapus Item dari Pesanan Aktif (`DELETE /v1/orders/:id/items/:item_id`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Deskripsi**: Menghapus item dari pesanan yang belum dibayar, membatalkan pengurangan stok item tersebut, dan merekrut ulang total nominal transaksi.
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Item removed from order successfully",
    "data": null
  }
  ```

#### 8. Memperbarui Kuantitas Item yang Disajikan (`PATCH /v1/orders/:id/items/:item_id/served`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Deskripsi**: Digunakan oleh modul Dapur (Chef) untuk memperbarui kuantitas produk yang sudah disajikan ke meja pelanggan.
- **Request Body**:
  ```json
  {
    "served_qty": 2
  }
  ```
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Served quantity updated successfully",
    "data": {
      "id": "a92b314d-fa12-4c91-b3b3-8c467a9ad21b",
      "order_id": "e0bfa99b-6db3-43bb-81ef-7bf2b12361ac",
      "product_id": "c1f7a14e-0a56-42ab-ba41-a189b87df8b1",
      "product_name": "Nasi Goreng Kotak",
      "quantity": 2,
      "unit_price": 20000,
      "subtotal": 40000,
      "served_qty": 2,
      "created_at": "2026-07-10T13:00:00Z",
      "updated_at": "2026-07-10T13:20:00Z"
    }
  }
  ```

---

### 5.7. Modul Pembayaran (`/v1/payments`)

#### 1. Melakukan Pembayaran (`POST /v1/payments`)
- **Tingkat Akses**: Superadmin, Owner, Staff
- **Deskripsi**: Memproses pembayaran pesanan dengan metode `Cash`, `Card`, atau `Digital Wallet`. Sistem memvalidasi jumlah nominal pembayaran (`amount_paid`) harus bernilai sama atau lebih besar dari total pesanan. Pembayaran sukses akan mengubah status pesanan menjadi `Paid`.
- **Request Body**:
  ```json
  {
    "order_id": "e0bfa99b-6db3-43bb-81ef-7bf2b12361ac",
    "payment_method": "Cash",
    "amount_paid": 50000
  }
  ```
- **Response Sukses (201 Created)**:
  ```json
  {
    "status": 201,
    "message": "Payment created successfully",
    "data": {
      "id": "p72b314d-fa12-4c91-b3b3-8c467a9ad21b",
      "order_id": "e0bfa99b-6db3-43bb-81ef-7bf2b12361ac",
      "payment_method": "Cash",
      "amount_paid": 50000,
      "timestamp": "2026-07-10T13:25:00Z"
    }
  }
  ```

#### 2. Mendapatkan Info Pembayaran Pesanan (`GET /v1/payments/:order_id`)
- **Tingkat Akses**: Superadmin, Owner, Staff
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Success",
    "data": {
      "id": "p72b314d-fa12-4c91-b3b3-8c467a9ad21b",
      "order_id": "e0bfa99b-6db3-43bb-81ef-7bf2b12361ac",
      "payment_method": "Cash",
      "amount_paid": 50000,
      "timestamp": "2026-07-10T13:25:00Z"
    }
  }
  ```

---

### 5.8. Modul Pengeluaran / Expenses (`/v1/pengeluaran`)

Modul ini mengelola pencatatan pengeluaran operasional toko non-produk (misal: listrik, air, gaji staf harian, sewa, dll.).

#### 1. Mencatat Pengeluaran Baru (`POST /v1/pengeluaran`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Request Body**:
  ```json
  {
    "tanggal": "2026-07-10",
    "category": "Operasional Toko",
    "description": "Pembayaran Listrik Bulanan",
    "amount": 250000
  }
  ```
- **Response Sukses (201 Created)**:
  ```json
  {
    "status": 201,
    "message": "Pengeluaran created successfully",
    "data": {
      "id": "ex95c43-4e89-4bc7-b84a-9ef85cf901ac",
      "store_id": 1,
      "tanggal": "2026-07-10",
      "category": "Operasional Toko",
      "description": "Pembayaran Listrik Bulanan",
      "amount": 250000,
      "created_by": "owner@example.com",
      "created_at": "2026-07-10T13:30:00Z"
    }
  }
  ```

#### 2. Mendapatkan Daftar Pengeluaran (`GET /v1/pengeluaran`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Query Parameter**:
  - Standard pagination (`page`, `page_size`, `order_by`, `order_dir`, `search`)
  - `start_date` (string, YYYY-MM-DD, opsional): Batas awal pencarian tanggal.
  - `end_date` (string, YYYY-MM-DD, opsional): Batas akhir pencarian tanggal.
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Success",
    "data": {
      "current_page": 1,
      "page_size": 10,
      "total_items": 1,
      "total_pages": 1,
      "has_previous": false,
      "has_next": false,
      "data": [
        {
          "id": "ex95c43-4e89-4bc7-b84a-9ef85cf901ac",
          "store_id": 1,
          "tanggal": "2026-07-10",
          "category": "Operasional Toko",
          "description": "Pembayaran Listrik Bulanan",
          "amount": 250000,
          "created_by": "owner@example.com",
          "created_at": "2026-07-10T13:30:00Z"
        }
      ]
    }
  }
  ```

#### 3. Mendapatkan Detail Pengeluaran (`GET /v1/pengeluaran/:id`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Response Sukses (200 OK)**:
  *(Struktur response sama seperti penambahan pengeluaran)*

#### 4. Memperbarui Data Pengeluaran (`PUT /v1/pengeluaran/:id`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Request Body**:
  ```json
  {
    "amount": 275000
  }
  ```
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Pengeluaran updated successfully",
    "data": {
      "id": "ex95c43-4e89-4bc7-b84a-9ef85cf901ac",
      "store_id": 1,
      "tanggal": "2026-07-10",
      "category": "Operasional Toko",
      "description": "Pembayaran Listrik Bulanan",
      "amount": 275000,
      "created_by": "owner@example.com",
      "created_at": "2026-07-10T13:30:00Z"
    }
  }
  ```

#### 5. Menghapus Data Pengeluaran (`DELETE /v1/pengeluaran/:id`)
- **Tingkat Akses**: Superadmin, Owner, Chef, Staff
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Pengeluaran deleted successfully",
    "data": null
  }
  ```

---

### 5.9. Modul Statistik Dashboard (`/v1/statistics`)

#### 1. Mendapatkan Data Statistik Ringkasan & Grafik (`GET /v1/statistics/dashboard`)
- **Tingkat Akses**: Superadmin, Owner
- **Query Parameter**:
  - `range` (string, opsional, pilihan: `weekly`, `monthly`, `all`, default: `daily`): Menentukan filter rentang waktu kalkulasi.
- **Response Sukses (200 OK)**:
  ```json
  {
    "status": 200,
    "message": "Success fetching dashboard data",
    "data": {
      "stats": {
        "total_orders": 12,
        "total_revenue": 240000,
        "total_profit": 115000,
        "total_expenses": 125000
      },
      "sales_chart": [
        {
          "date": "2026-07-10",
          "total": 240000
        }
      ],
      "finance_chart": [
        {
          "date": "2026-07-10",
          "revenue": 240000,
          "expenses": 125000,
          "profit": 115000
        }
      ],
      "top_products": [
        {
          "product_id": "c1f7a14e-0a56-42ab-ba41-a189b87df8b1",
          "product_name": "Nasi Goreng Kotak",
          "category_name": "Makanan Utama",
          "quantity": 12
        }
      ]
    }
  }
  ```
  *(Catatan: Detail formula matematika dan proses query untuk perhitungan statistik ini tercatat lengkap pada [algoritma_statistik.md](algoritma_statistik.md).)*

---

## 6. Penanganan Error & Kode Status HTTP

API SeliPOS mengembalikan HTTP status codes standar berikut:

- **200 OK**: Permintaan berhasil diproses.
- **201 Created**: Entitas baru berhasil dibuat (misalnya: pembuatan produk, pesanan, registrasi).
- **400 Bad Request**: Payload request tidak valid, field wajib tidak diisi, atau format salah.
- **401 Unauthorized**: Autentikasi gagal, Token JWT tidak ada, malformed, atau kedaluwarsa.
- **403 Forbidden**: Token valid tetapi hak akses (role) pengguna tidak diizinkan untuk mengakses modul tersebut.
- **404 Not Found**: Sumber daya (resource) yang diminta tidak ditemukan di database.
- **500 Internal Server Error**: Terjadi kegagalan pemrosesan internal pada server.
