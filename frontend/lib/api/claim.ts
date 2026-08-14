import { getClaimClient } from "./client";
import { Claim, CreateClaimRequest, ClaimStatus } from "@/types/claim";

export const claimAPI = {
  // Create a new claim
  createClaim: async (data: CreateClaimRequest): Promise<Claim> => {
    const client = getClaimClient();
    const response = await client.post("/claims", data);
    return response.data;
  },

  // Get claim by ID
  getClaim: async (id: string): Promise<Claim> => {
    const client = getClaimClient();
    const response = await client.get(`/claims/${id}`);
    return response.data;
  },

  // Get all claims by NGO
  getClaimsByNGO: async (): Promise<Claim[]> => {
    const client = getClaimClient();
    const response = await client.get("/claims/ngo");
    return response.data;
  },

  // Get all claims by Donor
  getClaimsByDonor: async (): Promise<Claim[]> => {
    const client = getClaimClient();
    const response = await client.get("/claims/donor");
    return response.data;
  },

  // Get all claims for a food listing
  getClaimsByFood: async (foodId: string): Promise<Claim[]> => {
    const client = getClaimClient();
    const response = await client.get(`/foods/${foodId}/claims`);
    return response.data;
  },

  // Approve claim (Donor action)
  approveClaim: async (id: string): Promise<void> => {
    const client = getClaimClient();
    await client.post(`/claims/${id}/approve`);
  },

  // Reject claim (Donor action)
  rejectClaim: async (id: string): Promise<void> => {
    const client = getClaimClient();
    await client.post(`/claims/${id}/reject`);
  },

  // Cancel claim (NGO action)
  cancelClaim: async (id: string): Promise<void> => {
    const client = getClaimClient();
    await client.post(`/claims/${id}/cancel`);
  },

  // Mark claim as picked up (NGO action)
  markPickedUp: async (id: string): Promise<void> => {
    const client = getClaimClient();
    await client.post(`/claims/${id}/pickup`);
  },

  // Mark claim as delivered (NGO action)
  markDelivered: async (id: string): Promise<void> => {
    const client = getClaimClient();
    await client.post(`/claims/${id}/deliver`);
  },
};