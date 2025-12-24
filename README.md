📚 MangaHub - Multi-Protocol Manga Tracking System

Course Project: Netcentric Programming > Course Code: IT096IU

📖 Introduction

MangaHub is a comprehensive web-based system simulating an online manga reading and tracking platform. Unlike conventional websites, MangaHub is built upon a Netcentric architecture, concurrently integrating 5 core network protocols within a single Backend. This design addresses real-world challenges related to data synchronization and real-time interaction.

The project is more than just a manga website; it's a miniature "Telecommunications Hub" where HTTP, TCP, UDP, and WebSocket data streams operate in parallel.

🚀 Key Features (5 Protocols)

The system deploys 5 network protocols functioning simultaneously:

HTTP REST API (Port 8080):

User Management (Registration, Login, JWT Auth).

Library Management: Searching, retrieving detailed manga data.

WebSocket (Real-time):

Real-time data bridge for the Web interface.

Online Chat System.

Receive instant notifications and synchronize reading status.

TCP Protocol (Port 9090 - Sync):

Simulates high-reliability reading progress synchronization from IoT/Mobile devices.

UDP Protocol (Port 9091 - Notification):

"Public Address System" for broadcasting ultra-fast notifications about events (New Chapter, Maintenance) to all users.

Internal Service (Mock gRPC):

Models internal communication between services for complex data queries.

🛠️ Technology Stack

Backend: Go (Golang)

Framework: Gin-Gonic (HTTP), Gorilla WebSocket.

Database Driver: go-sqlite3.

Security: Bcrypt, JWT-Go.

Frontend:

HTML5, CSS3 (Custom Glassmorphism Design).

JavaScript (ES6+, Fetch API, WebSocket).

Database: SQLite (Embedded, Zero-configuration).

📂 Project Structure

The project follows the standard Golang project layout:

MangaHub/
├── cmd/
│   └── api-server/
│       └── main.go       # Entry point: Starts the entire system (HTTP, TCP, UDP, WS)
├── pkg/
│   └── database/
│       └── database.go   # Handles SQLite connection, Migration, and Data Seeding
├── web/                  # Frontend (User Interface)
│   ├── index.html        # Home & Dashboard
│   ├── login.html        # Login Page
│   ├── register.html     # Registration Page
│   ├── style.css         # Dark Mode UI & Animations
│   └── app.js            # API and WebSocket client logic
├── data/                 # Directory for real manga images (optional)
├── go.mod                # Dependency management
└── mangahub.db           # Database File (Auto-generated)


⚙️ Installation & Execution Guide

1. System Requirements

Go: Version 1.20 or higher.

GCC Compiler: Required for compiling SQLite libraries (Install MinGW on Windows).

2. Installation

Clone the project and download necessary dependencies:

git clone [https://github.com/username/MangaHub.git](https://github.com/username/MangaHub.git)
cd MangaHub
go mod tidy


3. Start the Server

Run the following command to start the "Super Server":

go run cmd/api-server/main.go


Success indication:

✅ Database initialized successfully...
✅ Đã tạo xong 50 bộ truyện mẫu!
🚀 Server đang chạy tại: http://localhost:8080


4. Usage

Open your browser and navigate to: http://localhost:8080

Default Admin Account:

Username: admin

Password: 123456

Alternatively, click Register to create a new account.

📸 Demo Features

1. Modern Dark Mode Interface

Designed with a dark theme optimized for reading, featuring smooth transitions.

2. Mock Data (Auto Seeding)

The system automatically generates 50 mock manga titles on the first run, allowing immediate testing without manual data entry.

3. Manga Reader

Supports 2 reading modes: Scroll (Webtoon) and Flip.

Automatically fetches images from the Server.

4. Admin Panel

Administrators have access to a dedicated management page to view user lists and remove violating accounts.
