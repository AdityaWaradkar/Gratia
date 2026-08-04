import axios from 'axios'

const AUTH_API_URL = process.env.NEXT_PUBLIC_AUTH_API_URL || 'http://localhost:8080'
const USER_API_URL = process.env.NEXT_PUBLIC_USER_API_URL || 'http://localhost:8081'
const FOOD_API_URL = process.env.NEXT_PUBLIC_FOOD_API_URL || 'http://localhost:8082'
const CLAIM_API_URL = process.env.NEXT_PUBLIC_CLAIM_API_URL || 'http://localhost:8083'

export const authApi = axios.create({
  baseURL: AUTH_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const userApi = axios.create({
  baseURL: USER_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const foodApi = axios.create({
  baseURL: FOOD_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const claimApi = axios.create({
  baseURL: CLAIM_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add auth token interceptor
userApi.interceptors.request.use((config) => {
  const token = localStorage.getItem('accessToken')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

foodApi.interceptors.request.use((config) => {
  const token = localStorage.getItem('accessToken')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

claimApi.interceptors.request.use((config) => {
  const token = localStorage.getItem('accessToken')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})