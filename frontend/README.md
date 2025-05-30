# Stock Tracker Frontend

A Vue 3 application for tracking stock prices and recommendations.

## Features

- View real-time stock prices and recommendations
- Search stocks by company name or ticker
- Paginated stock list with sorting
- Clean and modern UI with Tailwind CSS
- Type-safe development with TypeScript

## Tech Stack

- Vue 3 with Composition API
- TypeScript
- Vite
- Pinia for state management
- Vue Router
- Tailwind CSS
- Axios for API calls

## Project Structure

The project follows a Domain-Driven Design (DDD) approach with the following structure:

```
src/
├── app/                      → App entrypoint, layout, router setup
├── modules/                  → Domain modules (stocks)
├── shared/                   → Shared utilities and components
└── assets/                   → Static assets
```

## Development Setup

1. Install dependencies:

   ```bash
   npm install
   ```

2. Create a `.env` file in the root directory:

   ```
   VITE_API_URL=http://localhost:3000
   ```

3. Start the development server:

   ```bash
   npm run dev
   ```

4. Build for production:
   ```bash
   npm run build
   ```

## API Integration

The application expects a REST API with the following endpoints:

- `GET /stocks` - List stocks with pagination
  - Query params: `page`, `limit`, `search`
  - Response: `{ items: Stock[], total: number, page: number, limit: number, totalPages: number }`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Submit a pull request
