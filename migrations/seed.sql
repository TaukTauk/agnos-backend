-- ============================================
-- Create Tables
-- ============================================

CREATE TABLE IF NOT EXISTS hospitals (
    id          VARCHAR(36) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    code        VARCHAR(100) NOT NULL UNIQUE,
    api_base_url VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP,
    updated_at  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS staff (
    id            VARCHAR(36) PRIMARY KEY,
    username      VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    hospital_id   VARCHAR(36) NOT NULL REFERENCES hospitals(id),
    created_at    TIMESTAMP,
    updated_at    TIMESTAMP
);

CREATE TABLE IF NOT EXISTS patients (
    id             VARCHAR(36) PRIMARY KEY,
    hospital_id    VARCHAR(36) NOT NULL REFERENCES hospitals(id),
    first_name_th  VARCHAR(255),
    middle_name_th VARCHAR(255),
    last_name_th   VARCHAR(255),
    first_name_en  VARCHAR(255),
    middle_name_en VARCHAR(255),
    last_name_en   VARCHAR(255),
    date_of_birth  DATE,
    patient_hn     VARCHAR(100),
    national_id    VARCHAR(20),
    passport_id    VARCHAR(20),
    phone_number   VARCHAR(20),
    email          VARCHAR(255),
    gender         VARCHAR(10),
    created_at     TIMESTAMP,
    updated_at     TIMESTAMP,
    UNIQUE (hospital_id, national_id),
    UNIQUE (hospital_id, passport_id)
);

-- ============================================
-- Hospitals
-- ============================================
INSERT INTO hospitals (id, name, code, api_base_url, created_at, updated_at)
VALUES
    ('a0000000-0000-0000-0000-000000000001', 'Hospital A', 'HOSPITAL_A', 'https://hospital-a.api.co.th', NOW(), NOW()),
    ('a0000000-0000-0000-0000-000000000002', 'Hospital B', 'HOSPITAL_B', 'https://hospital-b.api.co.th', NOW(), NOW());

-- ============================================
-- Staff (password is "password123" for all)
-- ============================================
INSERT INTO staff (id, username, password_hash, hospital_id, created_at, updated_at)
VALUES
    (
        'b0000000-0000-0000-0000-000000000001',
        'staff.hospital.a',
        '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
        'a0000000-0000-0000-0000-000000000001',
        NOW(), NOW()
    ),
    (
        'b0000000-0000-0000-0000-000000000002',
        'staff.hospital.b',
        '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
        'a0000000-0000-0000-0000-000000000002',
        NOW(), NOW()
    );

-- ============================================
-- Patients for Hospital A
-- ============================================
INSERT INTO patients (
    id, hospital_id,
    first_name_th, middle_name_th, last_name_th,
    first_name_en, middle_name_en, last_name_en,
    date_of_birth, patient_hn, national_id, passport_id,
    phone_number, email, gender,
    created_at, updated_at
)
VALUES
    (
        'c0000000-0000-0000-0000-000000000001',
        'a0000000-0000-0000-0000-000000000001',
        'สมชาย', '', 'ใจดี',
        'Somchai', '', 'Jaidee',
        '1990-05-15', 'HN-A-001', '1100100011001', NULL,
        '0811111111', 'somchai@email.com', 'male',
        NOW(), NOW()
    ),
    (
        'c0000000-0000-0000-0000-000000000002',
        'a0000000-0000-0000-0000-000000000001',
        'มาลี', '', 'สุขใจ',
        'Malee', '', 'Sukjai',
        '1985-08-20', 'HN-A-002', '1100100011002', NULL,
        '0822222222', 'malee@email.com', 'female',
        NOW(), NOW()
    ),
    (
        'c0000000-0000-0000-0000-000000000003',
        'a0000000-0000-0000-0000-000000000001',
        'วิชัย', '', 'รักดี',
        'Wichai', '', 'Rakdee',
        '1995-12-01', 'HN-A-003', '1100100011003', NULL,
        '0833333333', 'wichai@email.com', 'male',
        NOW(), NOW()
    ),
    (
        'c0000000-0000-0000-0000-000000000004',
        'a0000000-0000-0000-0000-000000000001',
        'นภา', '', 'แสงดาว',
        'Napa', '', 'Saengdao',
        '2000-03-10', 'HN-A-004', NULL, 'PA123456',
        '0844444444', 'napa@email.com', 'female',
        NOW(), NOW()
    ),
    (
        'c0000000-0000-0000-0000-000000000005',
        'a0000000-0000-0000-0000-000000000001',
        'ประยุทธ', '', 'มั่นคง',
        'Prayuth', '', 'Mankong',
        '1978-07-25', 'HN-A-005', '1100100011005', NULL,
        '0855555555', 'prayuth@email.com', 'male',
        NOW(), NOW()
    );

-- ============================================
-- Patients for Hospital B
-- ============================================
INSERT INTO patients (
    id, hospital_id,
    first_name_th, middle_name_th, last_name_th,
    first_name_en, middle_name_en, last_name_en,
    date_of_birth, patient_hn, national_id, passport_id,
    phone_number, email, gender,
    created_at, updated_at
)
VALUES
    (
        'c0000000-0000-0000-0000-000000000006',
        'a0000000-0000-0000-0000-000000000002',
        'สมหญิง', '', 'ดีงาม',
        'Somying', '', 'Deengam',
        '1992-04-18', 'HN-B-001', '1100100011006', NULL,
        '0866666666', 'somying@email.com', 'female',
        NOW(), NOW()
    ),
    (
        'c0000000-0000-0000-0000-000000000007',
        'a0000000-0000-0000-0000-000000000002',
        'อนันต์', '', 'พรมดี',
        'Anan', '', 'Promdee',
        '1988-11-30', 'HN-B-002', '1100100011007', NULL,
        '0877777777', 'anan@email.com', 'male',
        NOW(), NOW()
    );