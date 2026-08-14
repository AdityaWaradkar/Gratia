export interface RegisterRequest {
  email: string;
  password: string;
  role: "DONOR" | "NGO" | "USER";
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  accessToken: string;
  refreshToken: string;
}

export interface User {
  id: string;
  email: string;
  role: "DONOR" | "NGO" | "ADMIN" | "USER";
  emailVerified: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Tokens {
  accessToken: string;
  refreshToken: string;
}