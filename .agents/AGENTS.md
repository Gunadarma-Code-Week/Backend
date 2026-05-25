# 1. Response API
## 1.1 Response Status Code :
### 1.1.1 untuk method post:
- jika berhasil kembalikan status code 201 (wajib)
- jika validasi gagal kembalikan status code 400 (wajib)
- jika database conflict kembalikan status code 409 (wajib)
- jika data tidak ditemukan kembalikan status code 404 (wajib)

### 1.1.2 untuk method put:
- jika berhasil kembalikan status code 201 (wajib)
- jika validasi gagal kembalikan status code 400 (wajib)
- jika database conflict kembalikan status code 409 (wajib)
- jika data tidak ditemukan kembalikan status code 404 (wajib)

### 1.1.3 untuk method get:
- jika berhasil kembalikan status code 201 (wajib)
- jika validasi gagal kembalikan status code 400 (wajib)
- jika database conflict kembalikan status code 409 (wajib)
- jika data tidak ditemukan kembalikan status code 404 (wajib)

### 1.1.4 penjelasan atribut2 api:
- message: pesan informasi (wajib)
- data: data yang dikirimkan jika tidak ada data maka kembalikan nil (wajib)
- errors: data error yang dikirimkan wajib berupa dictionary atau nil (wajib)

## 1.2 Format Response API:
### 1.2.1 untuk status code 201:
```json
{
  "message": "string",
  "data": "interface{}",
  "errors": "nil"
}
```

contoh :
```json
{
	"message": "Permintaan berhasil diproses",
	"data": {
    "username": "adit"
  },
	"errors": null
}
```

### 1.2.2 untuk status code 400:
```json
{
  "message": "string",
  "data": "nil",
  "errors": "interface{}"
}
```

contoh :
```json
{
	"message": "gagal memproses permintaan",
	"data": null,
	"errors": {
    "key": "IS_REQUIRED"
  }
}
```

### 1.2.3 untuk status code 409:
```json
{
  "message": "string",
  "data": "nil",
  "errors": "nil"
}
```

### 1.2.4 untuk status code 404:
```json
{
  "message": "string",
  "data": "nil",
  "errors": "nil"
}
```

### 1.2.5 untuk status code 500:
```json
{
  "message": "string",
  "data": "nil",
  "errors": "nil"
}
```

## 1.3 Error Code List:
- di gunakan untuk field errors pada response api
- untuk error codenya: IS_REQUIRED, IS_INVALID, TOO_LONG, TOO_SHORT, MUST_LOWER, MUST_UPPER, MUST_SYMBOL, MUST_NUMBER, IS_ALREADY

## 1.4 Standarisasi Validation:
- untuk username, nama team, nama lengkap, dll (kecuali link) maksimal 50 karakter
- untuk seluruh data email maksimal 100 karakter
- untuk dateline jangan gunakan tipe data string
