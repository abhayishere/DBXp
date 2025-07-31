# DBXp - Database Explorer

A powerful Terminal User Interface (TUI) application for PostgreSQL database exploration and query execution, built with Go.

## 🌟 Features

### 📊 **Database Management**
- **Live Connection**: Connect to PostgreSQL databases using environment variables
- **Schema Explorer**: Browse tables in your database with an interactive sidebar
- **Real-time Updates**: Schema automatically refreshes after DDL operations

### 🔍 **Query Execution**
- **Interactive SQL Input**: Execute any SQL query with syntax support
- **Smart Query Handling**: Automatic detection of SELECT vs DDL/DML operations
- **Result Display**: Clean, tabular output with proper formatting
- **Error Handling**: Clear error messages for failed queries

### 🕹️ **Navigation & UX**
- **Keyboard Navigation**: 
  - `Tab`: Switch between schema explorer and query input
  - `Arrow Keys`: Navigate through table list
  - `Enter`: Select table or execute query
- **Click-to-Query**: Click any table name to auto-generate `SELECT * FROM table;`
- **Query History**: Navigate through last 20 executed queries
  - `↑ Arrow`: Previous query (newer to older)
  - `↓ Arrow`: Next query (older to newer)

### 💾 **Export Functionality**
- **CSV Export**: Export query results to `export.csv`
- **Keyboard Shortcut**: `Ctrl+E` to export current results
- **Smart Formatting**: Handles NULL values and different data types
- **Error Handling**: Clear feedback on export success/failure

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL database
- Docker (optional, for database setup)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/abhayishere/DBXp.git
   cd DBXp
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   Create a `.env` file:
   ```env
   DB_USER=postgres
   DB_PASSWORD=yourpass
   DB_HOST=localhost
   DB_PORT=5432
   DB_NAME=yourdb
   ```

4. **Run PostgreSQL** (if using Docker)
   ```bash
   docker run --name postgres-dbxp \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=yourpass \
     -e POSTGRES_DB=yourdb \
     -p 5432:5432 \
     -d postgres:latest
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

## 🎮 Usage

### Basic Navigation
1. **Start the app**: The cursor will be in the SQL input field
2. **Switch to table list**: Press `Tab` to focus on the schema explorer
3. **Select a table**: Use `↑/↓` arrows and press `Enter`
4. **Execute queries**: Type SQL and press `Enter`

### Sample Test Data
```sql
-- Create a test table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT,
    email TEXT
);

-- Insert sample data
INSERT INTO users (name, email) VALUES
('Alice', 'alice@example.com'),
('Bob', 'bob@example.com'),
('Charlie', 'charlie@example.com');

-- Query the data
SELECT * FROM users;
```

### Query History
- Press `↑` to navigate to previous queries
- Press `↓` to navigate to next queries
- History stores your last 20 executed queries

### Export Results
1. Execute a SELECT query
2. Press `Ctrl+E` to export results to CSV
3. Check `export.csv` in your project directory

## 🏗️ Architecture

### Project Structure
```
DBXp/
├── main.go                 # Application entry point
├── app/
│   └── app.go             # Main application logic
├── db/
│   └── connect.go         # Database connection
├── handlers/
│   ├── events.go          # UI event handling
│   ├── query.go           # Query execution
│   ├── history.go         # Query history management
│   └── export.go          # CSV export functionality
├── ui/
│   ├── layout.go          # UI layout management
│   └── schema.go          # Schema explorer
└── .env                   # Environment configuration
```

### Key Components

- **App**: Main application coordinator
- **QueryHandler**: Manages SQL execution and result formatting
- **EventHandler**: Handles keyboard inputs and user interactions
- **History**: Manages query history with navigation
- **Export**: Handles CSV export functionality
- **UI Components**: Schema explorer and layout management

## 🔧 Configuration

### Environment Variables
| Variable | Description | Default |
|----------|-------------|---------|
| `DB_USER` | PostgreSQL username | `postgres` |
| `DB_PASSWORD` | PostgreSQL password | `yourpass` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_NAME` | Database name | `yourdb` |

### Keyboard Shortcuts
| Key | Action |
|-----|--------|
| `Tab` | Switch focus between components |
| `↑` | Previous query in history |
| `↓` | Next query in history |
| `Enter` | Execute query / Select table |
| `Ctrl+E` | Export results to CSV |
| `Ctrl+C` | Exit application |

## 🛠️ Development

### Building
```bash
go build -o dbxp main.go
```

### Running Tests
```bash
go test ./...
```

### Dependencies
- [tview](https://github.com/rivo/tview) - Terminal UI framework
- [tcell](https://github.com/gdamore/tcell) - Terminal handling
- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [tview](https://github.com/rivo/tview) TUI framework
- Inspired by modern database administration tools
- PostgreSQL community for excellent Go drivers

---

**Happy Database Exploring!** 🚀


CREATE TABLE products (id SERIAL PRIMARY KEY,name VARCHAR(100) NOT NULL,price DECIMAL(10,2),category VARCHAR(50),created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);

INSERT INTO products (name, price, category) VALUES('Laptop', 999.99, 'Electronics'),('Coffee Mug', 12.50, 'Kitchen'),('Notebook', 5.99, 'Office'),('Headphones', 79.99, 'Electronics'),('Desk Chair', 199.99, 'Furniture');
