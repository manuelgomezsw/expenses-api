-- =====================================================
-- EXPENSES API - VIEWS CREATION SCRIPT
-- =====================================================
-- Descripción: Vistas optimizadas para el nuevo modelo simplificado
-- =====================================================

-- =====================================================
-- VISTA: v_fixed_expenses_with_status
-- Descripción: Gastos fijos con estado calculado para el frontend
-- Mapea directamente a FixedExpense interface
-- =====================================================
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
    
    -- Estado calculado basado en fecha actual y día de pago
    CASE 
        WHEN fe.is_paid = TRUE THEN 'paid'
        WHEN DAY(CURRENT_DATE) > fe.payment_day 
             AND DATE_FORMAT(CURRENT_DATE, '%Y-%m') = fe.month THEN 'overdue'
        ELSE 'pending'
    END as status,
    
    -- Información del bolsillo
    fe.pocket_id,
    p.description as pocket_description
FROM fixed_expenses fe
JOIN pockets p ON fe.pocket_id = p.id;

-- =====================================================
-- VISTA: v_daily_expenses_by_month
-- Descripción: Gastos diarios agrupados por mes
-- =====================================================
CREATE OR REPLACE VIEW v_daily_expenses_by_month AS
SELECT 
    DATE_FORMAT(STR_TO_DATE(de.date, '%Y-%m-%d'), '%Y-%m') as month,
    de.id,
    de.description,
    de.amount,
    de.date,
    de.created_at,
    
    -- Información adicional
    DAYNAME(STR_TO_DATE(de.date, '%Y-%m-%d')) as day_name,
    DAY(STR_TO_DATE(de.date, '%Y-%m-%d')) as day_number
FROM daily_expenses de;

-- =====================================================
-- VISTA: v_monthly_summary
-- Descripción: Resumen financiero mensual para el dashboard
-- =====================================================
CREATE OR REPLACE VIEW v_monthly_summary AS
SELECT 
    months.month,
    
    -- Ingresos del mes
    COALESCE(s.monthly_amount, 0) as total_income,
    
    -- Gastos fijos del mes
    COALESCE(fe_summary.total_fixed_expenses, 0) as total_fixed_expenses,
    COALESCE(fe_summary.fixed_expenses_paid, 0) as fixed_expenses_paid,
    COALESCE(fe_summary.fixed_expenses_total, 0) as fixed_expenses_total,
    
    -- Gastos diarios del mes
    COALESCE(de_summary.total_daily_expenses, 0) as total_daily_expenses,
    COALESCE(de_summary.daily_expenses_count, 0) as daily_expenses_count,
    
    -- Presupuesto diario configurado
    COALESCE(dec.monthly_budget, 0) as daily_budget_total,
    
    -- Cálculos derivados
    (COALESCE(s.monthly_amount, 0) - 
     COALESCE(fe_summary.total_fixed_expenses, 0) - 
     COALESCE(de_summary.total_daily_expenses, 0)) as remaining_budget,
     
    -- Porcentaje de presupuesto diario usado
    CASE 
        WHEN COALESCE(dec.monthly_budget, 0) > 0 THEN
            ROUND((COALESCE(de_summary.total_daily_expenses, 0) / dec.monthly_budget) * 100, 2)
        ELSE 0
    END as daily_budget_used_percentage

FROM (
    -- Generar lista de meses únicos de todas las tablas
    SELECT DISTINCT month FROM salaries
    UNION
    SELECT DISTINCT month FROM fixed_expenses
    UNION 
    SELECT DISTINCT DATE_FORMAT(STR_TO_DATE(date, '%Y-%m-%d'), '%Y-%m') FROM daily_expenses
    UNION
    SELECT DISTINCT month FROM daily_expenses_configs
) months

LEFT JOIN salaries s ON months.month = s.month

LEFT JOIN (
    SELECT 
        month,
        SUM(amount) as total_fixed_expenses,
        COUNT(*) as fixed_expenses_total,
        SUM(CASE WHEN is_paid = TRUE THEN 1 ELSE 0 END) as fixed_expenses_paid
    FROM fixed_expenses
    GROUP BY month
) fe_summary ON months.month = fe_summary.month

LEFT JOIN (
    SELECT 
        DATE_FORMAT(STR_TO_DATE(date, '%Y-%m-%d'), '%Y-%m') as month,
        SUM(amount) as total_daily_expenses,
        COUNT(*) as daily_expenses_count
    FROM daily_expenses
    GROUP BY DATE_FORMAT(STR_TO_DATE(date, '%Y-%m-%d'), '%Y-%m')
) de_summary ON months.month = de_summary.month

LEFT JOIN daily_expenses_configs dec ON months.month = dec.month;

-- =====================================================
-- VISTA: v_pockets_summary
-- Descripción: Resumen de bolsillos con estadísticas
-- =====================================================
CREATE OR REPLACE VIEW v_pockets_summary AS
SELECT 
    p.id,
    p.name,
    p.description,
    
    -- Estadísticas de gastos fijos
    COALESCE(fe_stats.fixed_expenses_count, 0) as fixed_expenses_count,
    COALESCE(fe_stats.total_fixed_amount, 0) as total_fixed_amount,
    COALESCE(fe_stats.paid_fixed_count, 0) as paid_fixed_count,
    
    -- Meses con actividad
    COALESCE(fe_stats.active_months_count, 0) as active_months_count
    
FROM pockets p
LEFT JOIN (
    SELECT 
        pocket_id,
        COUNT(*) as fixed_expenses_count,
        SUM(amount) as total_fixed_amount,
        SUM(CASE WHEN is_paid = TRUE THEN 1 ELSE 0 END) as paid_fixed_count,
        COUNT(DISTINCT month) as active_months_count
    FROM fixed_expenses
    GROUP BY pocket_id
) fe_stats ON p.id = fe_stats.pocket_id;

-- =====================================================
-- VISTA: v_expenses_calendar
-- Descripción: Vista calendario de gastos para el frontend
-- =====================================================
CREATE OR REPLACE VIEW v_expenses_calendar AS
SELECT 
    'fixed' as expense_type,
    fe.id,
    CONCAT(p.name, ': ', fe.concept_name) as title,
    fe.amount,
    CONCAT(fe.month, '-', LPAD(fe.payment_day, 2, '0')) as date,
    fe.is_paid as completed,
    fe.month,
    p.name as category
FROM fixed_expenses fe
JOIN pockets p ON fe.pocket_id = p.id

UNION ALL

SELECT 
    'daily' as expense_type,
    de.id,
    de.description as title,
    de.amount,
    de.date,
    TRUE as completed, -- Los gastos diarios siempre están "completados"
    DATE_FORMAT(STR_TO_DATE(de.date, '%Y-%m-%d'), '%Y-%m') as month,
    'Gasto Diario' as category
FROM daily_expenses de

ORDER BY date DESC;