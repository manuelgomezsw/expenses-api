# 🗄️ Database Management

Este directorio contiene todos los scripts SQL para manejar la base de datos **manualmente**. No utilizamos migraciones automáticas de GORM.

## 📁 Estructura de Archivos

```
sql/database/
├── README.md                    # Esta documentación
├── setup_database.sql           # Script completo para setup inicial
├── 02_create_tables.sql         # Creación de tablas
├── 03_create_views.sql          # Vistas para consultas optimizadas
└── 04_insert_initial_data.sql   # Datos iniciales
```

## 🚀 Setup Inicial

### Opción 1: Script Completo (Recomendado)
```bash
# Ejecutar setup completo
mysql -u root -p < sql/database/setup_database.sql
```

### Opción 2: Scripts Individuales
```bash
# 1. Crear tablas
mysql -u root -p < sql/database/02_create_tables.sql

# 2. Crear vistas (opcional)
mysql -u root -p < sql/database/03_create_views.sql

# 3. Insertar datos iniciales (opcional)
mysql -u root -p < sql/database/04_insert_initial_data.sql
```

## 📊 Modelo de Datos

### Tablas Principales

1. **`salaries`** - Configuración de salarios mensuales
   ```sql
   CREATE TABLE salaries (
       id INT PRIMARY KEY AUTO_INCREMENT,
       monthly_amount DECIMAL(15,2) NOT NULL,
       month VARCHAR(7) NOT NULL UNIQUE, -- "2024-01"
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );
   ```

2. **`pockets`** - Bolsillos organizacionales
   ```sql
   CREATE TABLE pockets (
       id INT PRIMARY KEY AUTO_INCREMENT,
       name VARCHAR(255) NOT NULL UNIQUE,
       description TEXT NULL,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );
   ```

3. **`fixed_expenses`** - Gastos fijos mensuales
   ```sql
   CREATE TABLE fixed_expenses (
       id INT PRIMARY KEY AUTO_INCREMENT,
       pocket_id INT NOT NULL,
       concept_name VARCHAR(255) NOT NULL,
       amount DECIMAL(15,2) NOT NULL,
       payment_day INT NOT NULL, -- 1-31
       is_paid BOOLEAN DEFAULT FALSE,
       month VARCHAR(7) NOT NULL, -- "2024-01"
       paid_date VARCHAR(10) NULL, -- "2024-01-15"
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       FOREIGN KEY (pocket_id) REFERENCES pockets(id)
   );
   ```

4. **`daily_expenses`** - Gastos diarios
   ```sql
   CREATE TABLE daily_expenses (
       id INT PRIMARY KEY AUTO_INCREMENT,
       description VARCHAR(500) NOT NULL,
       amount DECIMAL(15,2) NOT NULL,
       date VARCHAR(10) NOT NULL, -- "2024-01-15"
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );
   ```

5. **`daily_expenses_configs`** - Configuración presupuesto mensual
   ```sql
   CREATE TABLE daily_expenses_configs (
       id INT PRIMARY KEY AUTO_INCREMENT,
       monthly_budget DECIMAL(15,2) NOT NULL,
       month VARCHAR(7) NOT NULL UNIQUE -- "2024-01"
   );
   ```

## 🔄 Migraciones

### Agregar Nueva Migración

1. **Crear archivo SQL:**
   ```bash
   # Ejemplo: 05_add_new_feature.sql
   touch sql/database/05_add_new_feature.sql
   ```

2. **Escribir SQL:**
   ```sql
   -- 05_add_new_feature.sql
   ALTER TABLE pockets ADD COLUMN budget DECIMAL(15,2) DEFAULT 0;
   ALTER TABLE pockets ADD COLUMN color VARCHAR(7) DEFAULT '#000000';
   ```

3. **Ejecutar migración:**
   ```bash
   mysql -u root -p expenses_db < sql/database/05_add_new_feature.sql
   ```

### Rollback

Para hacer rollback, crear script inverso:
```sql
-- 05_rollback_new_feature.sql
ALTER TABLE pockets DROP COLUMN budget;
ALTER TABLE pockets DROP COLUMN color;
```

## 🌱 Seeding de Datos

El backend incluye un **Seeder** para datos iniciales:

```go
// En el código Go
seeder := database.NewSeeder(db)
seeder.SeedInitialData() // Crea bolsillos y configuraciones básicas
```

### Datos que se crean automáticamente:
- ✅ Bolsillos básicos (Hogar, Alimentación, Transporte, etc.)
- ✅ Configuración de salario para mes actual (0.00)
- ✅ Configuración de presupuesto diario para mes actual (0.00)

## 🔍 Consultas Útiles

### Verificar estructura
```sql
SHOW TABLES;
DESCRIBE salaries;
DESCRIBE pockets;
DESCRIBE fixed_expenses;
DESCRIBE daily_expenses;
DESCRIBE daily_expenses_configs;
```

### Ver datos iniciales
```sql
SELECT * FROM pockets;
SELECT * FROM salaries;
SELECT * FROM daily_expenses_configs;
```

### Limpiar datos (desarrollo)
```sql
-- ⚠️ CUIDADO: Esto elimina todos los datos
TRUNCATE TABLE daily_expenses;
TRUNCATE TABLE fixed_expenses;
DELETE FROM daily_expenses_configs;
DELETE FROM salaries;
-- No eliminar pockets si tienen foreign keys
```

## 🚨 Importante

- ❌ **NO usar GORM AutoMigrate** - Todo se maneja manualmente
- ✅ **Siempre hacer backup** antes de ejecutar migraciones
- ✅ **Probar en desarrollo** antes de aplicar en producción
- ✅ **Documentar cambios** en este README

## 🔗 Conexión

El backend se conecta usando GORM pero **sin migraciones automáticas**:

```go
// internal/infrastructure/database/gorm_connection.go
// Solo conexión, sin AutoMigrate
db, err := gorm.Open(mysql.Open(dsn), config)
```
