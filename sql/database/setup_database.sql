-- =====================================================
-- EXPENSES API - COMPLETE DATABASE SETUP
-- =====================================================
-- Descripción: Script maestro para configurar toda la base de datos
-- Ejecuta todos los scripts en el orden correcto
-- =====================================================

-- Crear base de datos
CREATE DATABASE IF NOT EXISTS expenses_db 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

USE expenses_db;

-- Configurar zona horaria y opciones
SET time_zone = '-05:00'; -- Colombia timezone
SET SQL_MODE = 'NO_ZERO_DATE,NO_ZERO_IN_DATE,ERROR_FOR_DIVISION_BY_ZERO';
SET FOREIGN_KEY_CHECKS = 1;

-- =====================================================
-- 1. CREAR TABLAS
-- =====================================================

-- 1. SALARIOS MENSUALES
CREATE TABLE IF NOT EXISTS salaries (
    id INT PRIMARY KEY AUTO_INCREMENT,
    monthly_amount DECIMAL(15,2) NOT NULL,
    month VARCHAR(7) NOT NULL, -- "2024-01" format
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_month (month),
    UNIQUE KEY uk_month (month)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2. BOLSILLOS
CREATE TABLE IF NOT EXISTS pockets (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. GASTOS FIJOS MENSUALES
CREATE TABLE IF NOT EXISTS fixed_expenses (
    id INT PRIMARY KEY AUTO_INCREMENT,
    pocket_id INT NOT NULL,
    concept_name VARCHAR(255) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    payment_day INT NOT NULL, -- día del mes (1-31)
    is_paid BOOLEAN DEFAULT FALSE,
    month VARCHAR(7) NOT NULL, -- "2024-01" format
    paid_date VARCHAR(10) NULL, -- "2024-01-15" format
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_month (month),
    INDEX idx_is_paid (is_paid),
    INDEX idx_payment_day (payment_day),
    INDEX idx_pocket_month (pocket_id, month),
    
    FOREIGN KEY (pocket_id) REFERENCES pockets(id) ON DELETE RESTRICT,
    CONSTRAINT chk_payment_day CHECK (payment_day BETWEEN 1 AND 31)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 4. GASTOS DIARIOS
CREATE TABLE IF NOT EXISTS daily_expenses (
    id INT PRIMARY KEY AUTO_INCREMENT,
    description VARCHAR(500) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    date VARCHAR(10) NOT NULL, -- "2024-01-15" format
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_date (date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 5. CONFIGURACIÓN DE PRESUPUESTO DIARIO MENSUAL
CREATE TABLE IF NOT EXISTS daily_expenses_configs (
    id INT PRIMARY KEY AUTO_INCREMENT,
    monthly_budget DECIMAL(15,2) NOT NULL,
    month VARCHAR(7) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- 2. CREAR VISTAS
-- =====================================================

-- Vista: Gastos fijos con estado
CREATE OR REPLACE VIEW v_fixed_expenses_with_status AS
SELECT 
    fe.id,
    p.name as pocket_name,
    fe.concept_name as name,
    fe.amount,
    fe.payment_day as due_date,
    fe.is_paid,
    fe.month,
    fe.paid_date,
    fe.created_at,
    
    -- Estado calculado
    CASE 
        WHEN fe.is_paid = TRUE THEN 'paid'
        WHEN DAY(CURRENT_DATE) > fe.payment_day 
             AND DATE_FORMAT(CURRENT_DATE, '%Y-%m') = fe.month THEN 'overdue'
        ELSE 'pending'
    END as status,
    
    fe.pocket_id,
    p.description as pocket_description
FROM fixed_expenses fe
JOIN pockets p ON fe.pocket_id = p.id;

-- Vista: Resumen mensual
CREATE OR REPLACE VIEW v_monthly_summary AS
SELECT 
    months.month,
    COALESCE(s.monthly_amount, 0) as total_income,
    COALESCE(fe_summary.total_fixed_expenses, 0) as total_fixed_expenses,
    COALESCE(fe_summary.fixed_expenses_paid, 0) as fixed_expenses_paid,
    COALESCE(fe_summary.fixed_expenses_total, 0) as fixed_expenses_total,
    COALESCE(de_summary.total_daily_expenses, 0) as total_daily_expenses,
    COALESCE(dec.monthly_budget, 0) as daily_budget_total,
    
    (COALESCE(s.monthly_amount, 0) - 
     COALESCE(fe_summary.total_fixed_expenses, 0) - 
     COALESCE(de_summary.total_daily_expenses, 0)) as remaining_budget,
     
    CASE 
        WHEN COALESCE(dec.monthly_budget, 0) > 0 THEN
            ROUND((COALESCE(de_summary.total_daily_expenses, 0) / dec.monthly_budget) * 100, 2)
        ELSE 0
    END as daily_budget_used_percentage

FROM (
    SELECT DISTINCT month FROM salaries
    UNION SELECT DISTINCT month FROM fixed_expenses
    UNION SELECT DISTINCT DATE_FORMAT(STR_TO_DATE(date, '%Y-%m-%d'), '%Y-%m') FROM daily_expenses
    UNION SELECT DISTINCT month FROM daily_expenses_configs
) months

LEFT JOIN salaries s ON months.month = s.month
LEFT JOIN (
    SELECT month, SUM(amount) as total_fixed_expenses, COUNT(*) as fixed_expenses_total,
           SUM(CASE WHEN is_paid = TRUE THEN 1 ELSE 0 END) as fixed_expenses_paid
    FROM fixed_expenses GROUP BY month
) fe_summary ON months.month = fe_summary.month
LEFT JOIN (
    SELECT DATE_FORMAT(STR_TO_DATE(date, '%Y-%m-%d'), '%Y-%m') as month, SUM(amount) as total_daily_expenses
    FROM daily_expenses GROUP BY DATE_FORMAT(STR_TO_DATE(date, '%Y-%m-%d'), '%Y-%m')
) de_summary ON months.month = de_summary.month
LEFT JOIN daily_expenses_configs dec ON months.month = dec.month;

-- =====================================================
-- 3. INSERTAR DATOS INICIALES
-- =====================================================

-- Bolsillos básicos
INSERT INTO pockets (name, description) VALUES
('Hogar', 'Gastos relacionados con el hogar y servicios básicos'),
('Alimentación', 'Comida, supermercado y restaurantes'),
('Transporte', 'Transporte público, gasolina, mantenimiento vehículo'),
('Salud', 'Medicina, consultas médicas, seguros de salud'),
('Entretenimiento', 'Cine, streaming, salidas, hobbies'),
('Educación', 'Cursos, libros, capacitaciones'),
('Ropa', 'Vestimenta y accesorios'),
('Otros', 'Gastos varios no categorizados')
ON DUPLICATE KEY UPDATE description = VALUES(description);

-- Configuración inicial para el mes actual
INSERT INTO salaries (monthly_amount, month) VALUES
(0.00, DATE_FORMAT(CURRENT_DATE, '%Y-%m'))
ON DUPLICATE KEY UPDATE monthly_amount = VALUES(monthly_amount);

INSERT INTO daily_expenses_configs (monthly_budget, month) VALUES
(0.00, DATE_FORMAT(CURRENT_DATE, '%Y-%m'))
ON DUPLICATE KEY UPDATE monthly_budget = VALUES(monthly_budget);

-- =====================================================
-- VERIFICACIÓN FINAL
-- =====================================================
SELECT 'Database setup completed successfully!' as status;
SELECT COUNT(*) as pockets_created FROM pockets;
SELECT COUNT(*) as salary_configs FROM salaries;
SELECT COUNT(*) as daily_configs FROM daily_expenses_configs;
