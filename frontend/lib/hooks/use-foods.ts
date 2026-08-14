"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { foodAPI } from "@/lib/api/food";
import { FoodListing, CreateFoodListingRequest, UpdateFoodListingRequest } from "@/types/food";

// Keys for React Query caching
export const foodKeys = {
  all: ["foods"] as const,
  lists: () => [...foodKeys.all, "list"] as const,
  available: () => [...foodKeys.lists(), "available"] as const,
  details: () => [...foodKeys.all, "detail"] as const,
  detail: (id: string) => [...foodKeys.details(), id] as const,
};

export function useFoods() {
  const queryClient = useQueryClient();

  // Get all available food listings
  const useGetAvailableListings = () => {
    return useQuery({
      queryKey: foodKeys.available(),
      queryFn: () => foodAPI.getAvailableListings(),
      staleTime: 60000, // 1 minute
    });
  };

  // Get single food listing
  const useGetListing = (id: string) => {
    return useQuery({
      queryKey: foodKeys.detail(id),
      queryFn: () => foodAPI.getListing(id),
      enabled: !!id,
    });
  };

  // Create food listing
  const useCreateListing = () => {
    return useMutation({
      mutationFn: (data: CreateFoodListingRequest) => foodAPI.createListing(data),
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: foodKeys.lists() });
        toast.success("Food listing created successfully!");
      },
      onError: (error: any) => {
        toast.error(error.response?.data?.error || "Failed to create listing");
      },
    });
  };

  // Update food listing
  const useUpdateListing = () => {
    return useMutation({
      mutationFn: ({ id, data }: { id: string; data: UpdateFoodListingRequest }) =>
        foodAPI.updateListing(id, data),
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: foodKeys.all });
        toast.success("Food listing updated successfully!");
      },
      onError: (error: any) => {
        toast.error(error.response?.data?.error || "Failed to update listing");
      },
    });
  };

  // Cancel food listing
  const useCancelListing = () => {
    return useMutation({
      mutationFn: (id: string) => foodAPI.cancelListing(id),
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: foodKeys.all });
        toast.success("Food listing cancelled successfully!");
      },
      onError: (error: any) => {
        toast.error(error.response?.data?.error || "Failed to cancel listing");
      },
    });
  };

  // Internal: Validate food for claim
  const useValidateFoodForClaim = (id: string) => {
    return useQuery({
      queryKey: [...foodKeys.detail(id), "validate"] as const,
      queryFn: () => foodAPI.validateFoodForClaim(id),
      enabled: !!id,
    });
  };

  return {
    useGetAvailableListings,
    useGetListing,
    useCreateListing,
    useUpdateListing,
    useCancelListing,
    useValidateFoodForClaim,
  };
}