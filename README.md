# Polling App

[Problem Description](https://docs.google.com/document/d/1mp8vuZsSAgvP7B7VnLVuTOHzirz5OopwwlCs81PhiHc/edit?tab=t.0#heading=h.85rgu5bm3g0)

A full-stack polling application built with a Go-based backend and a React-based frontend.

---

## Prerequisites

Ensure the following tools are installed on your system:
- **Go**: Version 1.16 or later
- **Node.js** and **npm**
- **PostgreSQL**

---

## Backend Setup

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/bhandari-ayush/polling-app.git
   cd polling-app/backend
   ```

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

3. Configure the PostgreSQL database:
   - Use the provided Docker Compose file for local testing.
   - Update the database credentials (e.g., host, port, username, password, and database name).
   - Pull the PostgreSQL image from the registry using the following details:
     - **Registry**: `registry.calm.nutanix.com`
     - **Credentials**: `calm/eye2eye`
   - Tag the image as `postgres:pollapp`:
     ```bash
     docker tag <postgres_img> postgres:pollapp
     ```
   - Access the database using:
     ```bash
     psql -U root -d pollapp
     ```

---

### Running the Backend

Start the backend server:
```bash
go run main.go
```

The backend server runs on port `8080`. The API is accessible at:
```
http://localhost:8080
```

- Backend configuration can be customized in:
   `backend/PollApp/service/application.go`
   - By default (for local testing), the backend assumes the database is accessible via the `DB_HOST` environment variable. If not set, you can configure it manually:
     ```bash
     export DB_HOST=<database_vm_ip>
     ```

---

## Frontend Setup

### Installation

1. Navigate to the frontend project directory:
   ```bash
   cd frontend
   ```

2. Install the required Node.js dependencies:
   ```bash
   npm install
   ```

---

### Running the Frontend

Start the development server:
```bash
npm start
```

Open your browser and navigate to:
```
http://localhost:3000
```

---

## Additional Notes

- Ensure the backend server is running before starting the frontend.
- Frontend configuration can be updated in the `.env` file located in the `frontend` directory.

---

## Deployment Using Calm

### Application Deployment
1. **App**: (Service Name => App)
   - Installation script: `app-install.sh` (uses `@@{DB.address}@@` where `DB` is the service name).  
   - Uninstallation script: `app-uninstall.sh`.

2. **Database**: (Service Name => DB)
   - Installation script: `postgress-install.sh`.
   - Uninstallation script: `postgress-uninstall.sh`.




