# 🚀 Deployment Guide - GCP App Engine

## 🎯 Environment Configuration Strategy

### **Local Development** → `.env` file

### **GCP Production** → `app.yaml` + Secret Manager

---

## 🔧 Local Development Setup

### 1. **Create `.env` file:**

```bash
# Copy template
cp env.example .env

# Edit with your local credentials
nano .env
```

### 2. **Local `.env` example:**

```bash
ENV=development
PORT=8080

# Local MySQL
DB_HOST=localhost:3306
DB_USER=root
DB_PASSWORD=your_local_password
DB_NAME_EXPENSES=expenses_db

JWT_SECRET=your-dev-secret
CORS_ALLOWED_ORIGIN=http://localhost:4200
DEBUG=true
LOG_LEVEL=info
```

### 3. **Run locally:**

```bash
./dev.sh run
```

---

## 🌐 GCP Production Deployment

### **Option A: Direct Environment Variables (Simple)**

Update `app.yaml`:

```yaml
env_variables:
  ENV: "production"
  DB_DSN: "user:password@tcp(host:port)/database?parseTime=true"
  JWT_SECRET: "your-production-secret"
  CORS_ALLOWED_ORIGIN: "https://your-frontend-domain.com"
```

### **Option B: Secret Manager (Recommended for Production)**

#### 1. **Create secrets in GCP:**

```bash
# Database connection string
gcloud secrets create db-dsn --data-file=- <<< "quotes-user:M9@I*49e2b#9Ek}_@tcp(34.23.218.229:3306)/quotes?parseTime=true"

# JWT Secret
gcloud secrets create jwt-secret --data-file=- <<< "your-super-secure-jwt-secret"

# CORS Origin
gcloud secrets create cors-origin --data-file=- <<< "https://your-frontend-domain.com"
```

#### 2. **Update `app.yaml` to use secrets:**

```yaml
runtime: go122
instance_class: F1
service: expenses-api

env_variables:
  ENV: "production"
  PORT: "8080"
  DEBUG: "false"
  LOG_LEVEL: "error"

# Reference secrets from Secret Manager
beta_settings:
  cloud_sql_instances: your-project:region:instance-name

# Use Secret Manager
includes:
  - env_variables:
      DB_DSN:
        _secret: "db-dsn"
      JWT_SECRET:
        _secret: "jwt-secret"
      CORS_ALLOWED_ORIGIN:
        _secret: "cors-origin"
```

#### 3. **Grant permissions:**

```bash
# Allow App Engine to access secrets
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
    --member="serviceAccount:YOUR_PROJECT_ID@appspot.gserviceaccount.com" \
    --role="roles/secretmanager.secretAccessor"
```

---

## 🗄️ Database Configuration

### **Local Development:**

```bash
# Individual variables
DB_HOST=localhost:3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME_EXPENSES=expenses_db
```

### **GCP Production:**

```bash
# Single DSN (preferred)
DB_DSN=user:password@tcp(host:port)/database?parseTime=true

# Cloud SQL Proxy format
DB_DSN=user:password@unix(/cloudsql/project:region:instance)/database?parseTime=true
```

---

## 🔐 Security Best Practices

### **1. Environment-specific secrets:**

- **Development:** Use `.env` file (never commit)
- **Production:** Use GCP Secret Manager

### **2. JWT Secrets:**

- **Development:** Simple string in `.env`
- **Production:** Strong random secret in Secret Manager

### **3. Database Credentials:**

- **Development:** Local MySQL credentials
- **Production:** Cloud SQL with IAM or strong passwords

### **4. CORS Origins:**

- **Development:** `http://localhost:4200`
- **Production:** Your actual frontend domain

---

## 🚀 Deployment Commands

### **Deploy to App Engine:**

```bash
# Build and deploy
gcloud app deploy

# Deploy specific version
gcloud app deploy --version=v1

# Deploy with custom config
gcloud app deploy app.yaml
```

### **Environment Variables Check:**

```bash
# View current environment variables
gcloud app versions describe VERSION_ID --service=expenses-api

# View logs
gcloud app logs tail -s expenses-api
```

---

## 🔍 Configuration Validation

The application automatically validates configuration on startup:

### **Required Variables:**

- ✅ `PORT` - Server port
- ✅ `ENV` - Environment (development/production)
- ✅ Database connection (either `DB_DSN` or individual `DB_*` vars)
- ✅ `JWT_SECRET` - Must not be default in production

### **Startup Logs:**

```
Starting Expenses API in production mode on port 8080
Connecting to database with DSN: user:****@tcp(host:port)/db
GORM database connection established successfully
```

---

## 🧪 Testing Configuration

### **Local Test:**

```bash
# Test with development config
./dev.sh run

# Test with production-like config
ENV=production ./dev.sh run
```

### **GCP Test:**

```bash
# Deploy to staging
gcloud app deploy --version=staging

# Test staging version
curl https://staging-dot-expenses-api-dot-YOUR_PROJECT.appspot.com/health
```

---

## 🚨 Troubleshooting

### **Common Issues:**

1. **Database Connection Failed:**

   - Check DSN format
   - Verify Cloud SQL instance is running
   - Check firewall rules

2. **Secret Manager Access Denied:**

   - Verify IAM permissions
   - Check service account has `secretmanager.secretAccessor` role

3. **CORS Errors:**
   - Verify `CORS_ALLOWED_ORIGIN` matches frontend domain
   - Check protocol (http vs https)

### **Debug Commands:**

```bash
# Check environment variables
gcloud app versions describe VERSION_ID --service=expenses-api

# View application logs
gcloud app logs tail -s expenses-api

# Test database connection
gcloud sql connect INSTANCE_NAME --user=USERNAME
```

---

## ✅ Deployment Checklist

### **Before Deployment:**

- [ ] Update `DB_DSN` with production database
- [ ] Set strong `JWT_SECRET`
- [ ] Configure correct `CORS_ALLOWED_ORIGIN`
- [ ] Set `ENV=production`
- [ ] Set `DEBUG=false`
- [ ] Test locally with production-like config

### **After Deployment:**

- [ ] Verify application starts successfully
- [ ] Test database connectivity
- [ ] Test API endpoints
- [ ] Verify CORS configuration
- [ ] Check application logs

---

**🎉 Your application is now ready for both development and production!**
