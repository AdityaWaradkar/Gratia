import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  images: {
    domains: ["localhost", "res.cloudinary.com"],
  },
  reactStrictMode: true,
  env: {
    NEXT_PUBLIC_AUTH_URL: process.env.NEXT_PUBLIC_AUTH_URL || "http://localhost:8080",
    NEXT_PUBLIC_USER_URL: process.env.NEXT_PUBLIC_USER_URL || "http://localhost:8081",
    NEXT_PUBLIC_FOOD_URL: process.env.NEXT_PUBLIC_FOOD_URL || "http://localhost:8082",
    NEXT_PUBLIC_CLAIM_URL: process.env.NEXT_PUBLIC_CLAIM_URL || "http://localhost:8083",
  },
  // Ensure output directory is correct
  distDir: ".next",
};

export default nextConfig;