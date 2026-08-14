CREATE TABLE IF NOT EXISTS schools (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    udise       TEXT NOT NULL UNIQUE,
    school_name TEXT NOT NULL,
    ddo         TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS employees (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    school_id   INTEGER NOT NULL REFERENCES schools(id),
    shalarth_id TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    gender      TEXT,
    designation TEXT,
    pan         TEXT NOT NULL UNIQUE,
    gpf         TEXT,
    aadhaar     TEXT NOT NULL UNIQUE,
    mobile      TEXT
);

CREATE TABLE IF NOT EXISTS payslip_records (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id       INTEGER NOT NULL REFERENCES employees(id),
    month             INTEGER NOT NULL,
    year              INTEGER NOT NULL,

    basic_pay         REAL NOT NULL DEFAULT 0,
    da                REAL NOT NULL DEFAULT 0,
    hra               REAL NOT NULL DEFAULT 0,
    ta                REAL NOT NULL DEFAULT 0,
    ta_arrears        REAL NOT NULL DEFAULT 0,
    da_arrears        REAL NOT NULL DEFAULT 0,
    basic_arrears     REAL NOT NULL DEFAULT 0,
    nps_empr_allow    REAL NOT NULL DEFAULT 0,

    fa                REAL NOT NULL DEFAULT 0,
    gpf_deduction     REAL NOT NULL DEFAULT 0,
    gpf_advance       REAL NOT NULL DEFAULT 0,
    pt                REAL NOT NULL DEFAULT 0,
    gis               REAL NOT NULL DEFAULT 0,
    revenue_stamp     REAL NOT NULL DEFAULT 0,
    nps_empr_contri   REAL NOT NULL DEFAULT 0,
    nps_emp_contri    REAL NOT NULL DEFAULT 0,
    income_tax        REAL NOT NULL DEFAULT 0,
    ngr_society_loan  REAL NOT NULL DEFAULT 0,
    home_loan         REAL NOT NULL DEFAULT 0,

    gross_earnings    REAL NOT NULL DEFAULT 0,
    total_deduction   REAL NOT NULL DEFAULT 0,
    net_payable       REAL NOT NULL DEFAULT 0,

    UNIQUE (employee_id, month, year)
);
