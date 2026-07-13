# Panduan Menulis Commit Message untuk AI & Developer

Dokumen ini berisi panduan agar AI atau developer dapat membuat pesan commit baru yang konsisten dengan gaya (style) history commit pada repositori `TA-POS-backend`.

---

## 1. Aturan Format Commit Message
Setiap pesan commit wajib mengikuti format berikut:
```
[tag] lowercase commit message description
```

### Karakteristik Utama:
- **Menggunakan Tag**: Setiap commit dimulai dengan tag dalam kurung siku `[...]`.
- **Huruf Kecil (Lowercase)**: Seluruh isi deskripsi setelah tag harus ditulis menggunakan huruf kecil.
- **Ringkas & Jelas**: Deskripsi tidak boleh terlalu panjang, namun harus mencerminkan perubahan utama.

---

## 2. Daftar Tag yang Digunakan
Gunakan tag-tag berikut sesuai dengan tipe perubahan yang dilakukan:

| Tag | Deskripsi Penggunaan | Contoh |
| :--- | :--- | :--- |
| `[doc]` | Digunakan untuk perubahan, penambahan, atau perbaikan file di dalam folder `doc/` atau file dokumentasi lainnya (seperti `README.md`). | `[doc] add commit guide for ai` |
| `[fix]` | Digunakan untuk perbaikan bug, error, race condition, atau isu kalkulasi logika. | `[fix] race condition order, cogs statistics, and pagination offset` |
| `[update]` | Digunakan untuk pembaruan kode, konfigurasi, optimasi, atau refactoring kode yang sudah ada tanpa mengubah fungsionalitas inti secara drastis. | `[update] remove product description and transaction discount` |
| `[add]` | Digunakan untuk penambahan fitur baru, file baru, endpoint baru, atau model baru. | `[add] stock buy for kulakan product` |

---

## 3. Panduan Khusus untuk Perubahan di Folder `doc/`
Jika Anda melakukan perubahan pada dokumentasi di dalam folder `doc/` (seperti `dokumentasi_api.md`, `algoritma_statistik.md`, dll), ikuti langkah berikut:

1. **Gunakan Tag `[doc]`** sebagai penanda utama.
2. Jika ada penambahan dokumen baru, tulis deskripsi seperti `add <nama dokumen>`.
3. Jika memperbarui isi dokumen, tulis deskripsi seperti `update API documentation` atau `update statistics algorithm doc`.

### Contoh Commit Message:
* Penambahan panduan commit baru:
  ```bash
  git commit -m "[doc] add commit guide for ai"
  ```
* Pembaruan dokumentasi API:
  ```bash
  git commit -m "[doc] update api documentation details"
  ```
* Perbaikan typo pada dokumen algoritma:
  ```bash
  git commit -m "[doc] fix typo in statistics algorithm documentation"
  ```
