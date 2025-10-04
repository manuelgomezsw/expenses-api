-- =====================================================
-- 1. SALARIOS MENSUALES
-- Interface: Salary { id?, monthly_amount, month, created_at? }
-- =====================================================
CREATE TABLE salaries (
    id INT PRIMARY KEY AUTO_INCREMENT,
    monthly_amount DECIMAL(15, 2) NOT NULL,
    month VARCHAR(7) NOT NULL,
    -- "2024-01" format
    INDEX idx_month (month)
);
-- =====================================================
-- 2. BOLSILLOS (Solo para organización, sin presupuesto)
-- Interface: Pocket { id?, name, description?, created_at? }
-- =====================================================
CREATE TABLE pockets (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NULL
);
-- =====================================================
-- 3. GASTOS FIJOS MENSUALES
-- Interface: FixedExpense { id?, pocket_name, concept_name, amount, payment_day, is_paid, month, paid_date?, created_at? }
-- =====================================================
CREATE TABLE fixed_expenses (
    id INT PRIMARY KEY AUTO_INCREMENT,
    pocket_id INT NOT NULL,
    concept_name VARCHAR(255) NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    payment_day INT NOT NULL,
    -- día del mes (1-31)
    is_paid BOOLEAN DEFAULT FALSE,
    month VARCHAR(7) NOT NULL,
    -- "2024-01" format
    paid_date VARCHAR(10) NULL,
    -- "2024-01-15" format
    INDEX idx_month (month),
    INDEX idx_is_paid (is_paid),
    INDEX idx_payment_day (payment_day),
    INDEX idx_pocket_month (pocket_id, month),
    FOREIGN KEY (pocket_id) REFERENCES pockets(id)
);
-- =====================================================
-- 4. GASTOS DIARIOS
-- Interface: DailyExpense { id?, description, amount, date, created_at? }
-- =====================================================
CREATE TABLE daily_expenses (
    id INT PRIMARY KEY AUTO_INCREMENT,
    description VARCHAR(500) NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    date VARCHAR(10) NOT NULL,
    -- "2024-01-15" format
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_date (date)
);
-- =====================================================
-- 5. CONFIGURACIÓN DE PRESUPUESTO DIARIO MENSUAL
-- Interface: DailyExpensesConfig { id?, monthly_budget, month }
-- =====================================================
CREATE TABLE daily_expenses_configs (
    id INT PRIMARY KEY AUTO_INCREMENT,
    monthly_budget DECIMAL(15, 2) NOT NULL,
    month VARCHAR(7) NOT NULL UNIQUE
);