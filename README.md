# MovieStream

A full-stack movie streaming application with AI-powered recommendations.

## Features

- Browse and discover movies with genre filtering
- User authentication (register, login, logout)
- Stream movies with embedded video player
- AI-powered movie recommendations
- Admin review management
- JWT-based authorization

## Tech Stack

**Backend:**
- Go with Gin-Gonic framework
- MongoDB database
- JWT authentication
- CORS enabled

**Frontend:**
- React 19 with Vite
- React Router for navigation
- Bootstrap & React Bootstrap UI
- Axios for HTTP requests
- React Player for video streaming
- FontAwesome icons

## Quick Start

**Backend:**
```bash
cd Server/appServer
go mod tidy
go run main.go
```
**Frontend:**
```bash
cd Client/movie-stream-client
npm install
npm run dev
```
