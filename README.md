# Go Dashboard

A modern, stylish dashboard application built with Go, HTMX, Alpine.js, and Tailwind CSS.

## Features

- 🏠 **Home Page** - Welcome dashboard with API data loader
- 💰 **Sales Dashboard** - Track sales performance and revenue
- 📦 **Inventory Dashboard** - Monitor and manage product inventory
- 🎨 **Modern UI** - Beautiful interface using Tailwind CSS
- ⚡ **Interactive** - Dynamic content loading with HTMX and Alpine.js

## Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerized deployment)

## Running Locally

### Without Docker

1. Clone the repository:
```bash
git clone <repository-url>
cd go-dashboard
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run main.go
```

4. Open your browser and navigate to:
```
http://localhost:8080
```

### With Docker

#### Using Docker CLI

1. Build the Docker image:
```bash
docker build -t go-dashboard .
```

2. Run the container:
```bash
docker run -p 8080:8080 go-dashboard
```

3. Open your browser and navigate to:
```
http://localhost:8080
```

#### Using Docker Compose

1. Start the application:
```bash
docker-compose up -d
```

2. View logs:
```bash
docker-compose logs -f
```

3. Stop the application:
```bash
docker-compose down
```

## Project Structure

```
go-dashboard/
├── api/
│   └── data.go           # API handlers and data structures
├── templates/
│   ├── layout.html       # Base layout with sidebar
│   ├── index.html        # Home page
│   ├── sales.html        # Sales dashboard
│   └── inventory.html    # Inventory dashboard
├── static/
│   └── index.html        # (Legacy static file)
├── main.go               # Application entry point
├── go.mod                # Go module dependencies
├── Dockerfile            # Docker build configuration
├── docker-compose.yml    # Docker Compose configuration
└── README.md             # This file
```

## API Endpoints

- `GET /` - Home page
- `GET /sales` - Sales dashboard
- `GET /inventory` - Inventory dashboard
- `GET /api/data` - Get sample data
- `GET /api/sales/json` - Get sales data (JSON)
- `GET /api/inventory/json` - Get inventory data (JSON)

## Technologies Used

- **Backend**: Go with Gorilla Mux router
- **Frontend**: 
  - HTMX for dynamic content
  - Alpine.js for interactivity
  - Tailwind CSS for styling
- **Deployment**: Docker & Docker Compose

## Development

To modify the templates or add new features:

1. Edit the HTML templates in `templates/`
2. Modify API handlers in `api/data.go`
3. Update routes in `main.go`
4. Rebuild and restart the application

For live development with Docker Compose, uncomment the volumes section in `docker-compose.yml` to mount local directories.

## License

MIT License

## Test