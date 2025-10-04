# 🚀 Expenses API - Development Guide

## 🎯 Quick Start

### 1. **Setup Environment**

```bash
# Copy environment variables
cp env.example .env
# Edit .env with your database credentials

# Setup development environment
./dev.sh setup
```

### 2. **Database Setup**

```bash
# Setup database (creates tables, views, initial data)
./dev.sh db-setup
```

### 3. **Run Application**

```bash
# Run in development mode
./dev.sh run

# Or run with debug logging
./dev.sh debug
```

---

## 🛠️ Development Commands

### **Using dev.sh script:**

```bash
./dev.sh build      # Build application
./dev.sh run        # Run application
./dev.sh debug      # Run with debug logging
./dev.sh test       # Run tests
./dev.sh clean      # Clean build artifacts
./dev.sh setup      # Setup dev environment
./dev.sh db-setup   # Setup database
./dev.sh help       # Show help
```

### **Using Go directly:**

```bash
go run .            # Run application
go build -o main .  # Build application
go test ./...       # Run tests
go mod tidy         # Tidy dependencies
```

---

## 🐛 Debugging in Cursor

### **Method 1: Using Debug Configuration**

1. Open **Run and Debug** panel (`Cmd+Shift+D`)
2. Select configuration:

   - **🚀 Run Expenses API** - Normal execution
   - **🐛 Debug Expenses API** - Full debugging with breakpoints
   - **🧪 Debug with Test Database** - Debug with test DB
   - **🌐 Run with Production Settings** - Production simulation

3. Click **Start Debugging** or press `F5`

### **Method 2: Using Tasks**

1. Open **Command Palette** (`Cmd+Shift+P`)
2. Type "Tasks: Run Task"
3. Select:
   - **🔨 Build** - Build the application
   - **🚀 Run** - Run the application
   - **🧪 Test** - Run tests
   - **🗄️ Setup Database** - Setup database

---

## 🗄️ Database Configuration

### **Environment Variables:**

```bash
DB_HOST=localhost:3306
DB_USER=root
DB_PASSWORD=your_password_here
DB_NAME_EXPENSES=expenses_db
```

### **Database Setup:**

```bash
# Option 1: Using script
./dev.sh db-setup

# Option 2: Manual setup
mysql -u root -p < sql/database/setup_database.sql
```

### **Database Structure:**

- `salaries` - Monthly salary configuration
- `pockets` - Expense categories
- `fixed_expenses` - Monthly fixed expenses
- `daily_expenses` - Daily expense tracking
- `daily_expenses_configs` - Monthly budget configuration

---

## 🌐 API Endpoints

### **Configuration:**

- `GET /api/config/income` - Get salary configuration
- `PUT /api/config/income` - Update salary
- `GET /api/config/pockets` - List pockets
- `POST /api/config/pockets` - Create pocket
- `PUT /api/config/pockets/:id` - Update pocket
- `DELETE /api/config/pockets/:id` - Delete pocket

### **Summary:**

- `GET /api/summary/:month` - Monthly financial summary

### **Fixed Expenses:**

- `GET /api/fixed-expenses/:month` - Get fixed expenses
- `PUT /api/fixed-expenses/:id/status` - Update payment status

### **Daily Expenses:**

- `GET /api/daily-expenses/:month` - Get daily expenses
- `POST /api/daily-expenses` - Create daily expense
- `PUT /api/daily-expenses/:id` - Update daily expense
- `DELETE /api/daily-expenses/:id` - Delete daily expense

---

## 🧪 Testing

### **Run Tests:**

```bash
./dev.sh test
# or
go test ./... -v
```

### **Test with Coverage:**

```bash
go test ./... -cover
```

---

## 🔧 Troubleshooting

### **Common Issues:**

1. **Database Connection Error:**

   - Check if MySQL is running
   - Verify credentials in environment variables
   - Ensure database exists

2. **Port Already in Use:**

   - Change PORT in environment variables
   - Kill process using port: `lsof -ti:8080 | xargs kill`

3. **Build Errors:**
   - Run `go mod tidy`
   - Check Go version compatibility

### **Debug Tips:**

- Use `🐛 Debug Expenses API` configuration for breakpoints
- Check logs in integrated terminal
- Verify environment variables are set correctly

---

## 📁 Project Structure

```
expenses-api/
├── .vscode/                    # Cursor/VSCode configurations
├── internal/
│   ├── api/dto/               # API contracts
│   ├── application/           # Use cases & ports
│   ├── domain/               # Domain entities
│   └── infrastructure/       # Infrastructure layer
├── sql/database/             # Database scripts
├── dev.sh                    # Development script
├── env.example              # Environment variables template
└── main.go                  # Application entry point
```

---

## 🚀 Ready to Code!

1. **Start the API:** `./dev.sh run`
2. **Open browser:** `http://localhost:8080`
3. **Set breakpoints** in Cursor and use `🐛 Debug Expenses API`
4. **Test endpoints** with your Angular frontend or Postman

Happy coding! 🎉
