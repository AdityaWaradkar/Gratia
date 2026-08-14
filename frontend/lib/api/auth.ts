import { authClient } from "./client";

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterCredentials {
  email: string;
  password: string;
  role?: "DONOR" | "NGO" | "USER";
}

export interface Tokens {
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

export interface AuthResponse {
  user: User;
  tokens: Tokens;
}

export const authAPI = {
  login: async (credentials: LoginCredentials): Promise<Tokens> => {
    const response = await authClient.post("/login", credentials);
    return response.data;
  },

  register: async (credentials: RegisterCredentials): Promise<User> => {
    const response = await authClient.post("/register", credentials);
    return response.data;
  },

  refresh: async (refreshToken: string): Promise<Tokens> => {
    const response = await authClient.post("/refresh", { refreshToken });
    return response.data;
  },

  forgotPassword: async (email: string): Promise<void> => {
    await authClient.post("/forgot-password", { email });
  },

  resetPassword: async (token: string, newPassword: string): Promise<void> => {
    await authClient.post("/reset-password", { token, newPassword });
  },

  logout: async (refreshToken: string): Promise<void> => {
    await authClient.post("/logout", { refreshToken });
  },

  getCurrentUser: async (): Promise<User> => {
    const response = await authClient.get("/me");
    return response.data;
  },
};