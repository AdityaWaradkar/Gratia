export interface RegisterRequest {
  email: string
  password: string
  role: 'user' | 'admin'
}

export interface LoginRequest {
  email: string
  password: string
}

export interface AuthResponse {
  accessToken: string
  refreshToken: string
}

export interface User {
  id: string
  email: string
  role: 'user' | 'admin'
}

export interface Tokens {
  accessToken: string
  refreshToken: string
}