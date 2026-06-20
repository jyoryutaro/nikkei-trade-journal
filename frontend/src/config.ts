// Base URL of the backend API. Override per environment with VITE_API_BASE
// (e.g. a deployed backend URL); defaults to the local dev server.
export const API_BASE = import.meta.env.VITE_API_BASE ?? 'http://localhost:8080'
