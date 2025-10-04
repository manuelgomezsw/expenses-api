#!/bin/bash

# Script de despliegue a GCP App Engine
# Uso: ./deploy.sh

set -e  # Salir si hay error

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  🚀 Expenses API - GCP Deployment     ║${NC}"
echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo ""

# Verificar que gcloud está instalado
if ! command -v gcloud &> /dev/null; then
    echo -e "${RED}❌ Error: gcloud no está instalado${NC}"
    echo "Instala Google Cloud SDK: https://cloud.google.com/sdk/docs/install"
    exit 1
fi

# Verificar proyecto actual
PROJECT=$(gcloud config get-value project 2>/dev/null)
if [ -z "$PROJECT" ]; then
    echo -e "${RED}❌ Error: No hay proyecto de GCP configurado${NC}"
    echo "Ejecuta: gcloud config set project YOUR_PROJECT_ID"
    exit 1
fi

echo -e "${GREEN}✓ Proyecto GCP: ${PROJECT}${NC}"
echo ""

# Obtener secretos de GCP Secret Manager
echo -e "${YELLOW}📦 Obteniendo secretos de Secret Manager...${NC}"

DB_USER=$(gcloud secrets versions access latest --secret="DB_USER" 2>/dev/null || echo "")
DB_PASS=$(gcloud secrets versions access latest --secret="DB_PASS" 2>/dev/null || echo "")
DB_NAME=$(gcloud secrets versions access latest --secret="DB_NAME_EXPENSES" 2>/dev/null || echo "")
INSTANCE_NAME=$(gcloud secrets versions access latest --secret="INSTANCE_CONNECTION_NAME" 2>/dev/null || echo "")

# Verificar secretos
if [ -z "$DB_USER" ] || [ -z "$DB_PASS" ] || [ -z "$DB_NAME" ]; then
    echo -e "${RED}❌ Error: Faltan secretos requeridos en Secret Manager${NC}"
    echo "Secretos requeridos: DB_USER, DB_PASS, DB_NAME_EXPENSES"
    exit 1
fi

echo -e "${GREEN}✓ Secretos obtenidos correctamente${NC}"
echo ""

# Crear DSN para Cloud SQL usando unix socket
# Formato para App Engine con Cloud SQL Proxy
DSN="${DB_USER}:${DB_PASS}@unix(/cloudsql/${INSTANCE_NAME})/${DB_NAME}?charset=utf8mb4&parseTime=True&loc=Local"

# Crear archivo temporal con variables de entorno
cat > .env.deploy <<EOF
DB_DSN=${DSN}
EOF

echo -e "${YELLOW}🔨 Construyendo aplicación...${NC}"
go build -o /dev/null . || {
    echo -e "${RED}❌ Error en la compilación${NC}"
    rm -f .env.deploy
    exit 1
}
echo -e "${GREEN}✓ Compilación exitosa${NC}"
echo ""

# Mostrar configuración de despliegue
echo -e "${BLUE}📋 Configuración de despliegue:${NC}"
echo "  - Proyecto: ${PROJECT}"
echo "  - Servicio: expenses-api"
echo "  - Runtime: go122"
echo "  - Base de datos: ${DB_NAME}"
echo "  - Instance: ${INSTANCE_NAME}"
echo ""

# Confirmar despliegue
read -p "¿Continuar con el despliegue? (y/n): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}⚠️  Despliegue cancelado${NC}"
    rm -f .env.deploy
    exit 0
fi

echo ""
echo -e "${YELLOW}🚀 Desplegando a App Engine...${NC}"
echo ""

# Crear app.yaml temporal con el DSN
cat > app.yaml.tmp <<EOF
runtime: go122

instance_class: F1
service: expenses-api

env_variables:
  ENV: "production"
  PORT: "8080"
  DEBUG: "false"
  LOG_LEVEL: "info"
  
  # Database connection string
  DB_DSN: "${DSN}"
  
  # CORS - ajustar según necesites
  CORS_ALLOWED_ORIGIN: "*"
  
  # JWT Secret temporal
  JWT_SECRET: "expenses-api-production-secret-2024-change-later"

# Conexión a Cloud SQL usando unix socket
beta_settings:
  cloud_sql_instances: "${INSTANCE_NAME}"

# Automatic scaling
automatic_scaling:
  min_instances: 0
  max_instances: 2
  target_cpu_utilization: 0.65
EOF

# Desplegar
gcloud app deploy app.yaml.tmp --quiet || {
    echo -e "${RED}❌ Error en el despliegue${NC}"
    rm -f .env.deploy app.yaml.tmp
    exit 1
}

# Limpiar archivos temporales
rm -f .env.deploy app.yaml.tmp

echo ""
echo -e "${GREEN}╔════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  ✅ Despliegue completado exitosamente ║${NC}"
echo -e "${GREEN}╔════════════════════════════════════════╗${NC}"
echo ""
echo -e "${BLUE}🌐 URL del servicio:${NC}"
gcloud app browse --service=expenses-api --no-launch-browser
echo ""
echo -e "${BLUE}📊 Ver logs:${NC}"
echo "  gcloud app logs tail -s expenses-api"
echo ""
echo -e "${BLUE}📈 Ver versiones:${NC}"
echo "  gcloud app versions list --service=expenses-api"
echo ""

