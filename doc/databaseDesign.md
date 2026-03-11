# Database Design - Traveling App

## I. Sơ Đồ ERD (Entity Relationship Diagram)

tblThanhVien ──────── tblKhachHang ──────── tblPDTour ──────── tblLichTour ──────── tblTour
      │                                          │                                      │
      │                                          │                                 tblTourDiaDiem
      ├──── tblNhanVien                      tblHoaDon                                 │
      │                                                                           tblDiaDiem
      └──── tblHoaDon                                                                  │
                                                                              tblDichvuDiaDiem
                                                                                       │
                                                                                  tblDichvu

---

## II. Mô Tả Chi Tiết Các Bảng

---

### 1. tblThanhVien (Thành Viên)

> Lưu thông tin tài khoản của tất cả người dùng hệ thống (khách hàng và nhân viên)

| Tên Cột    | Kiểu Dữ Liệu | Ràng Buộc       | Mô Tả                          |
|------------|--------------|-----------------|-------------------------------|
| ID         | INT(10)      | PRIMARY KEY, AI | Khóa chính, tự tăng            |
| username   | VARCHAR(25)  | NOT NULL        | Tên đăng nhập                  |
| password   | VARCHAR(25)  | NOT NULL        | Mật khẩu (đã mã hóa)           |
| ngaysinh   | DATE         | NOT NULL        | Ngày sinh                      |
| email      | VARCHAR(25)  | NOT NULL        | Địa chỉ email                  |

**Quan hệ:**
- 1 tblThanhVien → 0..1 tblKhachHang
- 1 tblThanhVien → 0..1 tblNhanVien
- 1 tblThanhVien → N tblHoaDon

---

### 2. tblKhachHang (Khách Hàng)

> Lưu thông tin mở rộng của khách hàng, liên kết với tài khoản thành viên

| Tên Cột         | Kiểu Dữ Liệu | Ràng Buộc       | Mô Tả                            |
|-----------------|--------------|-----------------|----------------------------------|
| ID              | INT(10)      | PRIMARY KEY, AI | Khóa chính, tự tăng              |
| maKH            | VARCHAR(25)  | NOT NULL        | Mã khách hàng (định danh)        |
| tblThanhVienID  | INT(10)      | NOT NULL, FK    | Khóa ngoại → tblThanhVien(ID)    |

**Quan hệ:**
- N tblKhachHang → 1 tblThanhVien
- 1 tblKhachHang → N tblPDTour

---

### 3. tblNhanVien (Nhân Viên)

> Lưu thông tin mở rộng của nhân viên, liên kết với tài khoản thành viên

| Tên Cột         | Kiểu Dữ Liệu | Ràng Buộc       | Mô Tả                            |
|-----------------|--------------|-----------------|----------------------------------|
| ID              | INT(10)      | PRIMARY KEY, AI | Khóa chính, tự tăng              |
| maNV            | VARCHAR(25)  | NOT NULL        | Mã nhân viên (định danh)         |
| chucvu          | VARCHAR(25)  | NOT NULL        | Chức vụ nhân viên                |
| tblThanhVienID  | INT(10)      | NOT NULL, FK    | Khóa ngoại → tblThanhVien(ID)    |

**Quan hệ:**
- N tblNhanVien → 1 tblThanhVien

---

### 4. tblTour (Tour Du Lịch)

> Lưu thông tin chi tiết về các tour du lịch

| Tên Cột     | Kiểu Dữ Liệu | Ràng Buộc       | Mô Tả                               |
|-------------|--------------|-----------------|-------------------------------------|
| ID          | INT(10)      | PRIMARY KEY, AI | Khóa chính, tự tăng                 |
| maTour      | VARCHAR(25)  | NOT NULL        | Mã tour (định danh)                 |
| tenTour     | VARCHAR(25)  | NOT NULL        | Tên tour                            |
| ThoiGian    | VARCHAR(25)  | NOT NULL        | Thời gian (VD: "3 ngày 2 đêm")      |
| PhuongTien  | VARCHAR(25)  | NOT NULL        | Phương tiện di chuyển               |
| SLKhachMax  | INT(20)      | NOT NULL        | Số lượng khách tối đa               |
| Mota        | VARCHAR(255) | NOT NULL        | Mô tả chi tiết tour                 |
| chiPhi      | INT(25)      | NOT NULL        | Chi phí / giá tour                  |

**Quan hệ:**
- 1 tblTour → N tblLichTour
- 1 tblTour → N tblTourDiaDiem

---

### 5. tblLichTour (Lịch Tour)

> Lưu thông tin lịch trình khởi hành của từng tour

| Tên Cột     | Kiểu Dữ Liệu | Ràng Buộc       | Mô Tả                         |
|-------------|--------------|-----------------|-------------------------------|
| ID          | INT(10)      | PRIMARY KEY, AI | Khóa chính, tự tăng           |
| ngayVe      | DATE         | NOT NULL        | Ngày về (ngày kết thúc tour)  |
| tblTourID   | INT(10)      | NOT NULL, FK    | Khóa ngoại → tblTour(ID)      |

**Quan hệ:**
- N tblLichTour → 1 tblTour
- 1 tblLichTour → N tblPDTour

---

### 6. tblPDTour (Phiếu Đặt Tour)

> Lưu thông tin đặt tour của khách hàng

| Tên Cột           | Kiểu Dữ Liệu | Ràng Buộc       | Mô Tả                            |
|-------------------|--------------|-----------------|----------------------------------|
| ID                | INT(10)      | PRIMARY KEY, AI | Khóa chính, tự tăng              |
| soKhachNL         | INT(10)      | NOT NULL        | Số khách người lớn               |
| soKhachTreEm      | INT(10)      | NOT NULL        | Số khách trẻ em                  |
| tblKhachHangID    | INT(10)      | NOT NULL, FK    | Khóa ngoại → tblKhachHang(ID)    |
| tblLichTourID     | INT(10)      | NOT NULL, FK    | Khóa ngoại → tblLichTour(ID)     |

**Quan hệ:**
- N tblPDTour → 1 tblKhachHang
- N tblPDTour → 1 tblLichTour
- 1 tblPDTour → 0..1 tblHoaDon

---

### 7. tblHoaDon (Hóa Đơn)

> Lưu thông tin hóa đơn thanh toán của phiếu đặt tour

| Tên Cột           | Kiểu Dữ Liệu | Ràng Buộc       | Mô Tả                            |
|-------------------|--------------|-----------------|----------------------------------|
| ID                | INT(10)      | PRIMARY KEY, AI | Khóa chính, tự tăng              |
| maHD              | VARCHAR(25)  | NOT NULL        | Mã hóa đơn (định danh)           |
| tblPDTourID       | INT(10)      | NOT NULL, FK    | Khóa ngoại → tblPDTour(ID)       |
| tblThanhVienID    | INT(10)      | NOT NULL, FK    | Khóa ngoại → tblThanhVien(ID)    |

**Quan hệ:**

- N tblHoaDon → 1 tblPDTour
- N tblHoaDon → 1 tblThanhVien (nhân viên lập hóa đơn)

---

### 8. tblDiaDiem (Địa Điểm)

> Lưu thông tin các địa điểm du lịch

| Tên Cột       | Kiểu Dữ Liệu | Ràng Buộc       | Mô Tả                        |
|---------------|--------------|-----------------|-------------------------------|
| ID            | INT(10)      | PRIMARY KEY, AI | Khóa chính, tự tăng           |
| tenDiaDiem    | VARCHAR(25)  | NOT NULL        | Tên địa điểm                  |
| tinhThanhpho  | VARCHAR(25)  | NOT NULL        | Tỉnh / Thành phố              |
| QuanHuyen     | VARCHAR(25)  | NOT NULL        | Quận / Huyện                  |

**Quan hệ:**

- 1 tblDiaDiem → N tblTourDiaDiem

---

### 9. tblTourDiaDiem (Tour - Địa Điểm)

> Bảng trung gian thể hiện mối quan hệ nhiều-nhiều giữa Tour và Địa Điểm

| Tên Cột        | Kiểu Dữ Liệu | Ràng Buộc       | Mô Tả                          |
|----------------|--------------|-----------------|-------------------------------|
| ID             | INT(10)      | PRIMARY KEY, AI | Khóa chính, tự tăng            |
| tblTourID      | INT(10)      | NOT NULL, FK    | Khóa ngoại → tblTour(ID)       |
| tblDiaDiemID   | INT(10)      | NOT NULL, FK    | Khóa ngoại → tblDiaDiem(ID)    |

**Quan hệ:**

- N tblTourDiaDiem → 1 tblTour
- N tblTourDiaDiem → 1 tblDiaDiem
- 1 tblTourDiaDiem → N tblDichvuDiaDiem

---

### 10. tblDichvu (Dịch Vụ)

> Lưu thông tin các dịch vụ bổ sung (khách sạn, nhà hàng, vui chơi, ...)

| Tên Cột   | Kiểu Dữ Liệu | Ràng Buộc       | Mô Tả                      |
|-----------|--------------|-----------------|---------------------------|
| ID        | INT(10)      | PRIMARY KEY, AI | Khóa chính, tự tăng        |
| tenDV     | VARCHAR(25)  | NOT NULL        | Tên dịch vụ                |
| loaiHinh  | VARCHAR(25)  | NOT NULL        | Loại hình dịch vụ          |
| giaDV     | INT(25)      | NOT NULL        | Giá dịch vụ                |

**Quan hệ:**

- 1 tblDichvu → N tblDichvuDiaDiem

---

### 11. tblDichvuDiaDiem (Dịch Vụ - Địa Điểm)

> Bảng trung gian liên kết Dịch Vụ với điểm dừng của Tour

| Tên Cột             | Kiểu Dữ Liệu | Ràng Buộc       | Mô Tả                               |
|---------------------|--------------|-----------------|-------------------------------------|
| ID                  | INT(10)      | PRIMARY KEY, AI | Khóa chính, tự tăng                 |
| tblDichvuID         | INT(10)      | NOT NULL, FK    | Khóa ngoại → tblDichvu(ID)          |
| tblTourDiaDiemID    | INT(10)      | NOT NULL, FK    | Khóa ngoại → tblTourDiaDiem(ID)     |

**Quan hệ:**
- N tblDichvuDiaDiem → 1 tblDichvu
- N tblDichvuDiaDiem → 1 tblTourDiaDiem

---

## III. Quan Hệ Giữa Các Bảng (Relationships)

| Bảng Cha         | Quan Hệ  | Bảng Con           | Khóa Ngoại                          |
|------------------|----------|--------------------|--------------------------------------|
| tblThanhVien     | 1 → N    | tblKhachHang       | tblKhachHang.tblThanhVienID          |
| tblThanhVien     | 1 → N    | tblNhanVien        | tblNhanVien.tblThanhVienID           |
| tblThanhVien     | 1 → N    | tblHoaDon          | tblHoaDon.tblThanhVienID             |
| tblKhachHang     | 1 → N    | tblPDTour          | tblPDTour.tblKhachHangID             |
| tblLichTour      | 1 → N    | tblPDTour          | tblPDTour.tblLichTourID              |
| tblTour          | 1 → N    | tblLichTour        | tblLichTour.tblTourID                |
| tblPDTour        | 1 → 1    | tblHoaDon          | tblHoaDon.tblPDTourID                |
| tblTour          | N ↔ N    | tblDiaDiem         | qua tblTourDiaDiem                   |
| tblDichvu        | N ↔ N    | tblTourDiaDiem     | qua tblDichvuDiaDiem                 |

---

## IV. Script SQL Tạo Database

```sql
CREATE DATABASE travel_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE travel_db;

-- 1. Bảng thành viên
CREATE TABLE tblThanhVien (
    ID       INT(10)      NOT NULL AUTO_INCREMENT,
    username VARCHAR(25)  NOT NULL,
    password VARCHAR(25)  NOT NULL,
    ngaysinh DATE         NOT NULL,
    email    VARCHAR(25)  NOT NULL,
    PRIMARY KEY (ID)
);

-- 2. Bảng khách hàng
CREATE TABLE tblKhachHang (
    ID              INT(10)     NOT NULL AUTO_INCREMENT,
    maKH            VARCHAR(25) NOT NULL,
    tblThanhVienID  INT(10)     NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (tblThanhVienID) REFERENCES tblThanhVien(ID)
);

-- 3. Bảng nhân viên
CREATE TABLE tblNhanVien (
    ID              INT(10)     NOT NULL AUTO_INCREMENT,
    maNV            VARCHAR(25) NOT NULL,
    chucvu          VARCHAR(25) NOT NULL,
    tblThanhVienID  INT(10)     NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (tblThanhVienID) REFERENCES tblThanhVien(ID)
);

-- 4. Bảng tour
CREATE TABLE tblTour (
    ID          INT(10)      NOT NULL AUTO_INCREMENT,
    maTour      VARCHAR(25)  NOT NULL,
    tenTour     VARCHAR(25)  NOT NULL,
    ThoiGian    VARCHAR(25)  NOT NULL,
    PhuongTien  VARCHAR(25)  NOT NULL,
    SLKhachMax  INT(20)      NOT NULL,
    Mota        VARCHAR(255) NOT NULL,
    chiPhi      INT(25)      NOT NULL,
    PRIMARY KEY (ID)
);

-- 5. Bảng lịch tour
CREATE TABLE tblLichTour (
    ID          INT(10) NOT NULL AUTO_INCREMENT,
    ngayVe      DATE    NOT NULL,
    tblTourID   INT(10) NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (tblTourID) REFERENCES tblTour(ID)
);

-- 6. Bảng phiếu đặt tour
CREATE TABLE tblPDTour (
    ID              INT(10) NOT NULL AUTO_INCREMENT,
    soKhachNL       INT(10) NOT NULL,
    soKhachTreEm    INT(10) NOT NULL,
    tblKhachHangID  INT(10) NOT NULL,
    tblLichTourID   INT(10) NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (tblKhachHangID) REFERENCES tblKhachHang(ID),
    FOREIGN KEY (tblLichTourID)  REFERENCES tblLichTour(ID)
);

-- 7. Bảng hóa đơn
CREATE TABLE tblHoaDon (
    ID              INT(10)     NOT NULL AUTO_INCREMENT,
    maHD            VARCHAR(25) NOT NULL,
    tblPDTourID     INT(10)     NOT NULL,
    tblThanhVienID  INT(10)     NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (tblPDTourID)    REFERENCES tblPDTour(ID),
    FOREIGN KEY (tblThanhVienID) REFERENCES tblThanhVien(ID)
);

-- 8. Bảng địa điểm
CREATE TABLE tblDiaDiem (
    ID            INT(10)     NOT NULL AUTO_INCREMENT,
    tenDiaDiem    VARCHAR(25) NOT NULL,
    tinhThanhpho  VARCHAR(25) NOT NULL,
    QuanHuyen     VARCHAR(25) NOT NULL,
    PRIMARY KEY (ID)
);

-- 9. Bảng tour - địa điểm (trung gian)
CREATE TABLE tblTourDiaDiem (
    ID            INT(10) NOT NULL AUTO_INCREMENT,
    tblTourID     INT(10) NOT NULL,
    tblDiaDiemID  INT(10) NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (tblTourID)    REFERENCES tblTour(ID),
    FOREIGN KEY (tblDiaDiemID) REFERENCES tblDiaDiem(ID)
);

-- 10. Bảng dịch vụ
CREATE TABLE tblDichvu (
    ID        INT(10)     NOT NULL AUTO_INCREMENT,
    tenDV     VARCHAR(25) NOT NULL,
    loaiHinh  VARCHAR(25) NOT NULL,
    giaDV     INT(25)     NOT NULL,
    PRIMARY KEY (ID)
);

-- 11. Bảng dịch vụ - địa điểm (trung gian)
CREATE TABLE tblDichvuDiaDiem (
    ID                INT(10) NOT NULL AUTO_INCREMENT,
    tblDichvuID       INT(10) NOT NULL,
    tblTourDiaDiemID  INT(10) NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (tblDichvuID)      REFERENCES tblDichvu(ID),
    FOREIGN KEY (tblTourDiaDiemID) REFERENCES tblTourDiaDiem(ID)
);
```

---

## V. Ghi Chú

- Tất cả các bảng sử dụng `AUTO_INCREMENT` cho khóa chính `ID`
- Mọi cột đều có ràng buộc `NOT NULL`
- Các khóa ngoại (`FK`) đảm bảo tính toàn vẹn tham chiếu giữa các bảng
- `tblTourDiaDiem` và `tblDichvuDiaDiem` là bảng trung gian xử lý quan hệ nhiều-nhiều
- `tblThanhVien` là bảng trung tâm dùng chung cho cả khách hàng và nhân viên
