import { getUserClient } from "./client";
import { DonorProfile, NGOProfile, CreateDonorProfileRequest, CreateNGOProfileRequest } from "@/types/user";

export const userAPI = {
  // Donor Profile APIs
  createDonorProfile: async (data: CreateDonorProfileRequest): Promise<DonorProfile> => {
    const client = getUserClient();
    const response = await client.post("/donors/profile", data);
    return response.data;
  },

  getMyDonorProfile: async (): Promise<DonorProfile> => {
    const client = getUserClient();
    const response = await client.get("/donors/profile/me");
    return response.data;
  },

  updateMyDonorProfile: async (data: CreateDonorProfileRequest): Promise<void> => {
    const client = getUserClient();
    await client.put("/donors/profile/me", data);
  },

  // NGO Profile APIs
  createNGOProfile: async (data: CreateNGOProfileRequest): Promise<NGOProfile> => {
    const client = getUserClient();
    const response = await client.post("/ngos/profile", data);
    return response.data;
  },

  getMyNGOProfile: async (): Promise<NGOProfile> => {
    const client = getUserClient();
    const response = await client.get("/ngos/profile/me");
    return response.data;
  },

  verifyNGO: async (userId: string): Promise<void> => {
    const client = getUserClient();
    await client.put(`/admin/ngos/verify?userId=${userId}`);
  },

  // Internal: Get donor profile by user ID
  getDonorProfileByUserID: async (userId: string): Promise<DonorProfile> => {
    const client = getUserClient();
    const response = await client.get(`/internal/users/${userId}/donor`);
    return response.data;
  },

  // Internal: Get NGO profile by user ID
  getNGOProfileByUserID: async (userId: string): Promise<NGOProfile> => {
    const client = getUserClient();
    const response = await client.get(`/internal/users/${userId}/ngo`);
    return response.data;
  },
};