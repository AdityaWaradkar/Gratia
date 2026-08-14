import axios from "axios";

const AUTH_URL = process.env.NEXT_PUBLIC_AUTH_URL || "http://localhost:8080";
const USER_URL = process.env.NEXT_PUBLIC_USER_URL || "http://localhost:8081";
const FOOD_URL = process.env.NEXT_PUBLIC_FOOD_URL || "http://localhost:8082";
const CLAIM_URL = process.env.NEXT_PUBLIC_CLAIM_URL || "http://localhost:8083";

// Auth API client (doesn't require auth)
export const authClient = axios.create({
  baseURL: AUTH_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// API client for authenticated requests
export const apiClient = axios.create({
  baseURL: AUTH_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem("accessToken");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");
      localStorage.removeItem("userRole");
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);

// Service-specific clients with auth
export function getAuthClient() {
  return apiClient;
}

export function getUserClient() {
  const client = axios.create({
    baseURL: USER_URL,
    headers: { "Content-Type": "application/json" },
  });
  client.interceptors.request.use((config) => {
    const token = localStorage.getItem("accessToken");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  });
  return client;
}

export function getFoodClient() {
  const client = axios.create({
    baseURL: FOOD_URL,
    headers: { "Content-Type": "application/json" },
  });
  client.interceptors.request.use((config) => {
    const token = localStorage.getItem("accessToken");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  });
  return client;
}

export function getClaimClient() {
  const client = axios.create({
    baseURL: CLAIM_URL,
    headers: { "Content-Type": "application/json" },
  });
  client.interceptors.request.use((config) => {
    const token = localStorage.getItem("accessToken");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  });
  return client;
}