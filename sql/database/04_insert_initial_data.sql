-- =====================================================
-- EXPENSES API - INITIAL DATA SCRIPT
-- =====================================================
-- Descripción: Datos iniciales para el nuevo modelo simplificado
-- =====================================================

-- =====================================================
-- DATOS INICIALES: pockets
-- Bolsillos básicos para organización
-- =====================================================
INSERT INTO pockets (name, description) VALUES
('Hogar', 'Gastos relacionados con el hogar y servicios básicos'),
('Alimentación', 'Comida, supermercado y restaurantes'),
('Transporte', 'Transporte público, gasolina, mantenimiento vehículo'),
('Salud', 'Medicina, consultas médicas, seguros de salud'),
('Entretenimiento', 'Cine, streaming, salidas, hobbies'),
('Educación', 'Cursos, libros, capacitaciones'),
('Ropa', 'Vestimenta y accesorios'),
('Otros', 'Gastos varios no categorizados')
ON DUPLICATE KEY UPDATE 
    description = VALUES(description);

-- =====================================================
-- CONFIGURACIÓN INICIAL: salario para el mes actual
-- =====================================================
INSERT INTO salaries (monthly_amount, month) VALUES
(0.00, DATE_FORMAT(CURRENT_DATE, '%Y-%m'))
ON DUPLICATE KEY UPDATE 
    monthly_amount = VALUES(monthly_amount);

-- =====================================================
-- CONFIGURACIÓN INICIAL: presupuesto diario para el mes actual
-- =====================================================
INSERT INTO daily_expenses_configs (monthly_budget, month) VALUES
(0.00, DATE_FORMAT(CURRENT_DATE, '%Y-%m'))
ON DUPLICATE KEY UPDATE 
    monthly_budget = VALUES(monthly_budget);

-- =====================================================
-- DATOS DE EJEMPLO: gastos fijos comunes (comentado por defecto)
-- Descomenta para agregar ejemplos
-- =====================================================
/*
-- Obtener el mes actual
SET @current_month = DATE_FORMAT(CURRENT_DATE, '%Y-%m');

-- Gastos fijos de ejemplo
INSERT INTO fixed_expenses (pocket_id, concept_name, amount, payment_day, month) VALUES
-- Hogar (pocket_id = 1)
(1, 'Arriendo', 1200000.00, 5, @current_month),
(1, 'Servicios - Luz', 120000.00, 15, @current_month),
(1, 'Servicios - Agua', 60000.00, 20, @current_month),
(1, 'Servicios - Gas', 40000.00, 25, @current_month),
(1, 'Internet', 80000.00, 10, @current_month),

-- Salud (pocket_id = 4)
(4, 'Seguro de Salud', 150000.00, 1, @current_month),
(4, 'Medicina', 50000.00, 15, @current_month),

-- Entretenimiento (pocket_id = 5)
(5, 'Netflix', 35000.00, 15, @current_month),
(5, 'Spotify', 15000.00, 20, @current_month),
(5, 'Gimnasio', 80000.00, 1, @current_month)

ON DUPLICATE KEY UPDATE 
    amount = VALUES(amount),
    payment_day = VALUES(payment_day);
*/

-- =====================================================
-- DATOS DE EJEMPLO: gastos diarios (comentado por defecto)
-- =====================================================
/*
INSERT INTO daily_expenses (description, amount, date) VALUES
('Almuerzo restaurante', 25000.00, DATE_FORMAT(CURRENT_DATE, '%Y-%m-%d')),
('Transporte público', 8000.00, DATE_FORMAT(CURRENT_DATE, '%Y-%m-%d')),
('Café', 5000.00, DATE_FORMAT(DATE_SUB(CURRENT_DATE, INTERVAL 1 DAY), '%Y-%m-%d')),
('Supermercado', 85000.00, DATE_FORMAT(DATE_SUB(CURRENT_DATE, INTERVAL 2 DAY), '%Y-%m-%d'))
ON DUPLICATE KEY UPDATE 
    amount = VALUES(amount);
*/