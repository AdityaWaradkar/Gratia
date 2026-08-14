"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { authAPI } from "@/lib/api/auth";
import { User } from "@/types/auth";

export function useAuth() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem("accessToken");
    if (token) {
      fetchUser();
    } else {
      setIsLoading(false);
      setIsAuthenticated(false);
    }
  }, []);

  const fetchUser = async () => {
    try {
      setIsLoading(true);
      const userData = await authAPI.getCurrentUser();
      setUser(userData);
      setIsAuthenticated(true);
    } catch (error) {
      console.error("Failed to fetch user:", error);
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");
      localStorage.removeItem("userRole");
      setIsAuthenticated(false);
    } finally {
      setIsLoading(false);
    }
  };

  const login = async (email: string, password: string) => {
    try {
      setIsLoading(true);
      const tokens = await authAPI.login({ email, password });

      localStorage.setItem("accessToken", tokens.accessToken);
      localStorage.setItem("refreshToken", tokens.refreshToken);

      // Decode token to get user role
      try {
        const payload = JSON.parse(atob(tokens.accessToken.split(".")[1]));
        const role = payload.role?.toLowerCase() || "donor";
        localStorage.setItem("userRole", role);
      } catch (_error) {
        // Invalid token payload
      }

      await fetchUser();
      toast.success("Login successful!");
      return true;
    } catch (error: any) {
      toast.error(error.response?.data?.error || "Login failed. Please try again.");
      return false;
    } finally {
      setIsLoading(false);
    }
  };

  const register = async (email: string, password: string, role?: "DONOR" | "NGO" | "USER") => {
    try {
      setIsLoading(true);
      const userData = await authAPI.register({ email, password, role });
      toast.success("Registration successful! Please login.");
      router.push("/login");
      return true;
    } catch (error: any) {
      toast.error(error.response?.data?.error || "Registration failed. Please try again.");
      return false;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async () => {
    try {
      const refreshToken = localStorage.getItem("refreshToken");
      if (refreshToken) {
        await authAPI.logout(refreshToken);
      }
    } catch (error) {
      // Ignore logout errors
      console.error("Logout error:", error);
    } finally {
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");
      localStorage.removeItem("userRole");
      setUser(null);
      setIsAuthenticated(false);
      toast.success("Logged out successfully");
      router.push("/login");
    }
  };

  const forgotPassword = async (email: string) => {
    try {
      await authAPI.forgotPassword(email);
      toast.success("Password reset email sent!");
      return true;
    } catch (error: any) {
      toast.error(error.response?.data?.error || "Failed to send reset email.");
      return false;
    }
  };

  const resetPassword = async (token: string, newPassword: string) => {
    try {
      await authAPI.resetPassword(token, newPassword);
      toast.success("Password reset successful! Please login.");
      router.push("/login");
      return true;
    } catch (error: any) {
      toast.error(error.response?.data?.error || "Password reset failed.");
      return false;
    }
  };

  const getUserRole = () => {
    const role = localStorage.getItem("userRole");
    return role || "user";
  };

  return {
    user,
    isLoading,
    isAuthenticated,
    login,
    register,
    logout,
    forgotPassword,
    resetPassword,
    getUserRole,
    fetchUser,
  };
}