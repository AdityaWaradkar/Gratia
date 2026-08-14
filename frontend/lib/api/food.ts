import { getFoodClient } from "./client";
import { FoodListing, CreateFoodListingRequest, UpdateFoodListingRequest } from "@/types/food";

export const foodAPI = {
  // Create a new food listing
  createListing: async (data: CreateFoodListingRequest): Promise<FoodListing> => {
    const client = getFoodClient();
    const response = await client.post("/foods", data);
    return response.data;
  },

  // Get all available food listings
  getAvailableListings: async (): Promise<FoodListing[]> => {
    const client = getFoodClient();
    const response = await client.get("/foods");
    return response.data;
  },

  // Get food listing by ID
  getListing: async (id: string): Promise<FoodListing> => {
    const client = getFoodClient();
    const response = await client.get(`/foods/${id}`);
    return response.data;
  },

  // Update food listing
  updateListing: async (id: string, data: UpdateFoodListingRequest): Promise<void> => {
    const client = getFoodClient();
    await client.put(`/foods/${id}`, data);
  },

  // Cancel food listing
  cancelListing: async (id: string): Promise<void> => {
    const client = getFoodClient();
    await client.delete(`/foods/${id}`);
  },

  // Internal: Validate food for claim
  validateFoodForClaim: async (id: string): Promise<{ claimable: boolean; reason?: string }> => {
    const client = getFoodClient();
    const response = await client.get(`/internal/foods/${id}/validate`);
    return response.data;
  },

  // Internal: Mark food as claimed
  markFoodClaimed: async (id: string): Promise<void> => {
    const client = getFoodClient();
    await client.patch(`/internal/foods/${id}/claim`);
  },
};