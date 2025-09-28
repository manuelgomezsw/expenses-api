# 🚀 API Guide para Frontend Angular

## 📋 Nuevos Endpoints Implementados

### **Resumen Mensual**
```http
GET /api/summary/{month}
```
**Parámetros:**
- `month`: Formato YYYY-MM (ej: 2024-01)

**Respuesta:**
```json
{
  "month": "2024-01",
  "total_income": 5000000,
  "total_fixed_expenses": 2500000,
  "total_daily_expenses": 800000,
  "remaining_budget": 1700000,
  "fixed_expenses_paid": 8,
  "fixed_expenses_total": 12,
  "daily_budget_used": 800000,
  "daily_budget_total": 1500000
}
```

---

### **Configuración de Ingresos**

#### Obtener configuración actual
```http
GET /api/config/income
```
**Respuesta:**
```json
{
  "monthly_amount": 5000000,
  "currency": "COP"
}
```

#### Actualizar configuración
```http
PUT /api/config/income
```
**Body:**
```json
{
  "monthly_amount": 5500000,
  "currency": "COP"
}
```

---

### **Gestión de Bolsillos**

#### Obtener todos los bolsillos
```http
GET /api/config/pockets
```
**Respuesta:**
```json
[
  {
    "id": 1,
    "name": "Comida",
    "budget": 500000,
    "color": "#FF6B6B"
  },
  {
    "id": 2,
    "name": "Transporte",
    "budget": 200000,
    "color": "#4ECDC4"
  }
]
```

#### Crear nuevo bolsillo
```http
POST /api/config/pockets
```
**Body:**
```json
{
  "name": "Entretenimiento",
  "budget": 300000,
  "color": "#45B7D1"
}
```

#### Actualizar bolsillo
```http
PUT /api/config/pockets/{id}
```

#### Eliminar bolsillo
```http
DELETE /api/config/pockets/{id}
```

---

### **Gastos Fijos**

#### Obtener gastos fijos del mes
```http
GET /api/fixed-expenses/{month}
```
**Parámetros:**
- `month`: Formato YYYY-MM

**Respuesta:**
```json
[
  {
    "id": 1,
    "name": "Arriendo",
    "amount": 1200000,
    "due_date": 5,
    "is_paid": true,
    "paid_date": "2024-01-05T10:30:00Z",
    "created_at": "2024-01-01T00:00:00Z"
  },
  {
    "id": 2,
    "name": "Internet",
    "amount": 80000,
    "due_date": 15,
    "is_paid": false,
    "paid_date": null,
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

#### Marcar gasto fijo como pagado
```http
PUT /api/fixed-expenses/{id}/status
```
**Body:**
```json
{
  "is_paid": true
}
```

---

### **Gastos Diarios**

#### Obtener gastos diarios del mes
```http
GET /api/daily-expenses/{month}
```
**Parámetros:**
- `month`: Formato YYYY-MM

**Respuesta:**
```json
[
  {
    "id": 1,
    "amount": 25000,
    "description": "Almuerzo",
    "date": "2024-01-15",
    "pocket_id": 1
  },
  {
    "id": 2,
    "amount": 8000,
    "description": "Transporte público",
    "date": "2024-01-15",
    "pocket_id": 2
  }
]
```

#### Crear gasto diario
```http
POST /api/daily-expenses
```
**Body:**
```json
{
  "amount": 35000,
  "description": "Cena restaurante",
  "date": "2024-01-16",
  "pocket_id": 1
}
```

#### Actualizar gasto diario
```http
PUT /api/daily-expenses/{id}
```

#### Eliminar gasto diario
```http
DELETE /api/daily-expenses/{id}
```

---

## 🔄 Mapeo de Modelos

### Frontend → Backend
| **Frontend Model** | **Backend Equivalent** | **Notas** |
|-------------------|----------------------|-----------|
| `Salary` | `config.Salary` | Nuevo modelo |
| `FixedExpense` | `concepts.Concept` | Adaptado con estados |
| `DailyExpense` | `expenses.Expense` | Simplificado |
| `Pocket` | `pockets.Pocket` | Agregado budget y color |
| `DailyExpensesConfig` | `config.DailyExpensesConfig` | Nuevo modelo |

---

## 🚧 Estado Actual vs Próximos Pasos

### ✅ **Implementado:**
- [x] Estructura de DTOs para frontend
- [x] Controladores con endpoints específicos
- [x] Rutas configuradas
- [x] Validación básica de parámetros
- [x] Respuestas con datos de ejemplo

### 🔄 **En Progreso:**
- [ ] Conectar con servicios reales del dominio
- [ ] Implementar lógica de filtrado por mes
- [ ] Mapear concepts → fixed expenses
- [ ] Adaptar expenses → daily expenses

### ⏳ **Pendiente:**
- [ ] Implementar repositories para nuevos modelos
- [ ] Crear tablas de base de datos necesarias
- [ ] Implementar lógica de cálculo de resumen
- [ ] Agregar validaciones robustas
- [ ] Implementar manejo de errores consistente
- [ ] Agregar tests unitarios

---

## 🎯 Próximos Pasos Recomendados

### **Fase 2: Implementación Real (Próxima semana)**
1. **Crear tablas de BD:**
   - `salary_config`
   - `daily_expenses_config`
   - Modificar `pockets` (agregar budget, color)

2. **Implementar repositories:**
   - `SalaryRepository`
   - `DailyExpensesConfigRepository`
   - Actualizar `PocketsRepository`

3. **Conectar servicios reales:**
   - Implementar lógica de cálculo en `SummaryController`
   - Mapear `concepts` a `fixed expenses`
   - Adaptar `expenses` a `daily expenses`

### **Fase 3: Optimización (Semana siguiente)**
1. **Filtrado por mes dinámico**
2. **Estados de gastos fijos (pagado/pendiente/vencido)**
3. **Cálculos de presupuesto en tiempo real**
4. **Validaciones y manejo de errores**

---

## 🧪 Testing de Endpoints

Para probar los endpoints, puedes usar:

```bash
# Obtener resumen mensual
curl -X GET "http://localhost:8080/api/summary/2024-01"

# Obtener configuración de ingresos
curl -X GET "http://localhost:8080/api/config/income"

# Obtener gastos fijos del mes
curl -X GET "http://localhost:8080/api/fixed-expenses/2024-01"

# Obtener gastos diarios del mes
curl -X GET "http://localhost:8080/api/daily-expenses/2024-01"
```

---

¿Te parece bien esta estructura? ¿Hay algún endpoint específico que quieras que implemente completamente primero?
