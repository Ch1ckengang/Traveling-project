# Database Design - Traveling App

## I. ERD (Entity Relationship Diagram)

`tblThanhVien` -> `tblKhachHang` -> `tblPDTour` -> `tblLichTour` -> `tblTour`

Additional links:
- `tblThanhVien` -> `tblNhanVien`
- `tblThanhVien` -> `tblHoaDon`
- `tblPDTour` -> `tblHoaDon`
- `tblTour` -> `tblTourDiaDiem` -> `tblDiaDiem`
- `tblTourDiaDiem` -> `tblDichvuDiaDiem` -> `tblDichvu`

---

## II. Table Details

### 1. `tblThanhVien` (Members)
Stores account information for all users (customers and employees).

| Column | Type | Constraint | Description |
|---|---|---|---|
| ID | INT(10) | PRIMARY KEY, AI | Primary key |
| username | VARCHAR(25) | NOT NULL | Login name |
| password | VARCHAR(25) | NOT NULL | Encrypted password |
| ngaysinh | DATE | NOT NULL | Date of birth |
| email | VARCHAR(25) | NOT NULL | Email address |

Relationships:
- 1 `tblThanhVien` -> 0..1 `tblKhachHang`
- 1 `tblThanhVien` -> 0..1 `tblNhanVien`
- 1 `tblThanhVien` -> N `tblHoaDon`

### 2. `tblKhachHang` (Customers)
Stores extended customer data linked to member account.

| Column | Type | Constraint | Description |
|---|---|---|---|
| ID | INT(10) | PRIMARY KEY, AI | Primary key |
| maKH | VARCHAR(25) | NOT NULL | Customer code |
| tblThanhVienID | INT(10) | NOT NULL, FK | Foreign key -> `tblThanhVien(ID)` |

Relationships:
- N `tblKhachHang` -> 1 `tblThanhVien`
- 1 `tblKhachHang` -> N `tblPDTour`

### 3. `tblNhanVien` (Employees)
Stores extended employee data linked to member account.

| Column | Type | Constraint | Description |
|---|---|---|---|
| ID | INT(10) | PRIMARY KEY, AI | Primary key |
| maNV | VARCHAR(25) | NOT NULL | Employee code |
| chucvu | VARCHAR(25) | NOT NULL | Role/position |
| tblThanhVienID | INT(10) | NOT NULL, FK | Foreign key -> `tblThanhVien(ID)` |

Relationships:
- N `tblNhanVien` -> 1 `tblThanhVien`

### 4. `tblTour` (Tours)
Stores detailed tour information.

| Column | Type | Constraint | Description |
|---|---|---|---|
| ID | INT(10) | PRIMARY KEY, AI | Primary key |
| maTour | VARCHAR(25) | NOT NULL | Tour code |
| tenTour | VARCHAR(25) | NOT NULL | Tour name |
| ThoiGian | VARCHAR(25) | NOT NULL | Duration |
| PhuongTien | VARCHAR(25) | NOT NULL | Transportation |
| SLKhachMax | INT(20) | NOT NULL | Max guests |
| Mota | VARCHAR(255) | NOT NULL | Tour description |
| chiPhi | INT(25) | NOT NULL | Tour price/cost |

Relationships:
- 1 `tblTour` -> N `tblLichTour`
- 1 `tblTour` -> N `tblTourDiaDiem`

### 5. `tblLichTour` (Tour Schedule)
Stores departure/schedule information for each tour.

| Column | Type | Constraint | Description |
|---|---|---|---|
| ID | INT(10) | PRIMARY KEY, AI | Primary key |
| ngayVe | DATE | NOT NULL | Return/end date |
| tblTourID | INT(10) | NOT NULL, FK | Foreign key -> `tblTour(ID)` |

Relationships:
- N `tblLichTour` -> 1 `tblTour`
- 1 `tblLichTour` -> N `tblPDTour`

### 6. `tblPDTour` (Tour Booking Form)
Stores booking details submitted by customers.

| Column | Type | Constraint | Description |
|---|---|---|---|
| ID | INT(10) | PRIMARY KEY, AI | Primary key |
| soKhachNL | INT(10) | NOT NULL | Number of adults |
| soKhachTreEm | INT(10) | NOT NULL | Number of children |
| tblKhachHangID | INT(10) | NOT NULL, FK | Foreign key -> `tblKhachHang(ID)` |
| tblLichTourID | INT(10) | NOT NULL, FK | Foreign key -> `tblLichTour(ID)` |

Relationships:
- N `tblPDTour` -> 1 `tblKhachHang`
- N `tblPDTour` -> 1 `tblLichTour`
- 1 `tblPDTour` -> 0..1 `tblHoaDon`

### 7. `tblHoaDon` (Invoice)
Stores payment invoice data for tour bookings.

| Column | Type | Constraint | Description |
|---|---|---|---|
| ID | INT(10) | PRIMARY KEY, AI | Primary key |
| maHD | VARCHAR(25) | NOT NULL | Invoice code |
| tblPDTourID | INT(10) | NOT NULL, FK | Foreign key -> `tblPDTour(ID)` |
| tblThanhVienID | INT(10) | NOT NULL, FK | Foreign key -> `tblThanhVien(ID)` |

Relationships:
- N `tblHoaDon` -> 1 `tblPDTour`
- N `tblHoaDon` -> 1 `tblThanhVien` (employee who created invoice)

### 8. `tblDiaDiem` (Locations)
Stores travel destination information.

| Column | Type | Constraint | Description |
|---|---|---|---|
| ID | INT(10) | PRIMARY KEY, AI | Primary key |
| tenDiaDiem | VARCHAR(25) | NOT NULL | Location name |
| tinhThanhpho | VARCHAR(25) | NOT NULL | Province/City |
| QuanHuyen | VARCHAR(25) | NOT NULL | District |

Relationships:
- 1 `tblDiaDiem` -> N `tblTourDiaDiem`

### 9. `tblTourDiaDiem` (Tour - Location)
Junction table for many-to-many relation between Tours and Locations.

| Column | Type | Constraint | Description |
|---|---|---|---|
| ID | INT(10) | PRIMARY KEY, AI | Primary key |
| tblTourID | INT(10) | NOT NULL, FK | Foreign key -> `tblTour(ID)` |
| tblDiaDiemID | INT(10) | NOT NULL, FK | Foreign key -> `tblDiaDiem(ID)` |

Relationships:
- N `tblTourDiaDiem` -> 1 `tblTour`
- N `tblTourDiaDiem` -> 1 `tblDiaDiem`
- 1 `tblTourDiaDiem` -> N `tblDichvuDiaDiem`

### 10. `tblDichvu` (Services)
Stores supplemental services (hotel, food, entertainment, etc.).

| Column | Type | Constraint | Description |
|---|---|---|---|
| ID | INT(10) | PRIMARY KEY, AI | Primary key |
| tenDV | VARCHAR(25) | NOT NULL | Service name |
| loaiHinh | VARCHAR(25) | NOT NULL | Service type |
| giaDV | INT(25) | NOT NULL | Service price |

Relationships:
- 1 `tblDichvu` -> N `tblDichvuDiaDiem`

### 11. `tblDichvuDiaDiem` (Service - Location)
Junction table linking services with tour stop locations.

| Column | Type | Constraint | Description |
|---|---|---|---|
| ID | INT(10) | PRIMARY KEY, AI | Primary key |
| tblDichvuID | INT(10) | NOT NULL, FK | Foreign key -> `tblDichvu(ID)` |
| tblTourDiaDiemID | INT(10) | NOT NULL, FK | Foreign key -> `tblTourDiaDiem(ID)` |

Relationships:
- N `tblDichvuDiaDiem` -> 1 `tblDichvu`
- N `tblDichvuDiaDiem` -> 1 `tblTourDiaDiem`

---

## III. Relationship Summary

| Parent Table | Relation | Child Table | Foreign Key |
|---|---|---|---|
| tblThanhVien | 1 -> N | tblKhachHang | tblKhachHang.tblThanhVienID |
| tblThanhVien | 1 -> N | tblNhanVien | tblNhanVien.tblThanhVienID |
| tblThanhVien | 1 -> N | tblHoaDon | tblHoaDon.tblThanhVienID |
| tblKhachHang | 1 -> N | tblPDTour | tblPDTour.tblKhachHangID |
| tblLichTour | 1 -> N | tblPDTour | tblPDTour.tblLichTourID |
| tblTour | 1 -> N | tblLichTour | tblLichTour.tblTourID |
| tblPDTour | 1 -> 1 | tblHoaDon | tblHoaDon.tblPDTourID |
| tblTour | N <-> N | tblDiaDiem | via tblTourDiaDiem |
| tblDichvu | N <-> N | tblTourDiaDiem | via tblDichvuDiaDiem |

---

## IV. SQL Script to Create Database

```sql
CREATE DATABASE travel_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE travel_db;

-- 1. Members table
CREATE TABLE tblThanhVien (
    ID       INT(10)      NOT NULL AUTO_INCREMENT,
    username VARCHAR(25)  NOT NULL,
    password VARCHAR(25)  NOT NULL,
    ngaysinh DATE         NOT NULL,
    email    VARCHAR(25)  NOT NULL,
    PRIMARY KEY (ID)
);

-- 2. Customers table
CREATE TABLE tblKhachHang (
    ID              INT(10)     NOT NULL AUTO_INCREMENT,
    maKH            VARCHAR(25) NOT NULL,
    tblThanhVienID  INT(10)     NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (tblThanhVienID) REFERENCES tblThanhVien(ID)
);

-- 3. Employees table
CREATE TABLE tblNhanVien (
    ID              INT(10)     NOT NULL AUTO_INCREMENT,
    maNV            VARCHAR(25) NOT NULL,
    chucvu          VARCHAR(25) NOT NULL,
    tblThanhVienID  INT(10)     NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (tblThanhVienID) REFERENCES tblThanhVien(ID)
);

-- 4. Tours table
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

-- 5. Tour schedule table
CREATE TABLE tblLichTour (
    ID          INT(10) NOT NULL AUTO_INCREMENT,
    ngayVe      DATE    NOT NULL,
    tblTourID   INT(10) NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (tblTourID) REFERENCES tblTour(ID)
);

-- 6. Booking form table
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

-- 7. Invoice table
CREATE TABLE tblHoaDon (
    ID              INT(10)     NOT NULL AUTO_INCREMENT,
    maHD            VARCHAR(25) NOT NULL,
    tblPDTourID     INT(10)     NOT NULL,
    tblThanhVienID  INT(10)     NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (tblPDTourID)    REFERENCES tblPDTour(ID),
    FOREIGN KEY (tblThanhVienID) REFERENCES tblThanhVien(ID)
);

-- 8. Locations table
CREATE TABLE tblDiaDiem (
    ID            INT(10)     NOT NULL AUTO_INCREMENT,
    tenDiaDiem    VARCHAR(25) NOT NULL,
    tinhThanhpho  VARCHAR(25) NOT NULL,
    QuanHuyen     VARCHAR(25) NOT NULL,
    PRIMARY KEY (ID)
);

-- 9. Tour-location junction table
CREATE TABLE tblTourDiaDiem (
    ID            INT(10) NOT NULL AUTO_INCREMENT,
    tblTourID     INT(10) NOT NULL,
    tblDiaDiemID  INT(10) NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (tblTourID)    REFERENCES tblTour(ID),
    FOREIGN KEY (tblDiaDiemID) REFERENCES tblDiaDiem(ID)
);

-- 10. Services table
CREATE TABLE tblDichvu (
    ID        INT(10)     NOT NULL AUTO_INCREMENT,
    tenDV     VARCHAR(25) NOT NULL,
    loaiHinh  VARCHAR(25) NOT NULL,
    giaDV     INT(25)     NOT NULL,
    PRIMARY KEY (ID)
);

-- 11. Service-location junction table
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

## V. Notes
- All tables use `AUTO_INCREMENT` for primary key `ID`.
- Every column is currently marked as `NOT NULL` in this baseline script.
- Foreign keys enforce referential integrity.
- `tblTourDiaDiem` and `tblDichvuDiaDiem` are junction tables for many-to-many relationships.
- `tblThanhVien` is the central shared account table for both customers and employees.
