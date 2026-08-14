"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { claimAPI } from "@/lib/api/claim";
import { Claim, CreateClaimRequest } from "@/types/claim";

// Keys for React Query caching
export const claimKeys = {
  all: ["claims"] as const,
  lists: () => [...claimKeys.all, "list"] as const,
  listByNGO: (ngoUserId: string) => [...claimKeys.lists(), "ngo", ngoUserId] as const,
  listByDonor: (donorUserId: string) => [...claimKeys.lists(), "donor", donorUserId] as const,
  listByFood: (foodId: string) => [...claimKeys.lists(), "food", foodId] as const,
  details: () => [...claimKeys.all, "detail"] as const,
  detail: (id: string) => [...claimKeys.details(), id] as const,
};

export function useClaims() {
  const queryClient = useQueryClient();

  // Get all claims by NGO
  const useGetClaimsByNGO = (ngoUserId: string) => {
    return useQuery({
      queryKey: claimKeys.listByNGO(ngoUserId),
      queryFn: () => claimAPI.getClaimsByNGO(),
      enabled: !!ngoUserId,
    });
  };

  // Get all claims by Donor
  const useGetClaimsByDonor = (donorUserId: string) => {
    return useQuery({
      queryKey: claimKeys.listByDonor(donorUserId),
      queryFn: () => claimAPI.getClaimsByDonor(),
      enabled: !!donorUserId,
    });
  };

  // Get claims by food listing
  const useGetClaimsByFood = (foodId: string) => {
    return useQuery({
      queryKey: claimKeys.listByFood(foodId),
      queryFn: () => claimAPI.getClaimsByFood(foodId),
      enabled: !!foodId,
    });
  };

  // Get single claim
  const useGetClaim = (id: string) => {
    return useQuery({
      queryKey: claimKeys.detail(id),
      queryFn: () => claimAPI.getClaim(id),
      enabled: !!id,
    });
  };

  // Create claim
  const useCreateClaim = () => {
    return useMutation({
      mutationFn: (data: CreateClaimRequest) => claimAPI.createClaim(data),
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: claimKeys.lists() });
        toast.success("Claim created successfully!");
      },
      onError: (error: any) => {
        toast.error(error.response?.data?.error || "Failed to create claim");
      },
    });
  };

  // Approve claim (Donor action)
  const useApproveClaim = () => {
    return useMutation({
      mutationFn: (id: string) => claimAPI.approveClaim(id),
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: claimKeys.all });
        toast.success("Claim approved successfully!");
      },
      onError: (error: any) => {
        toast.error(error.response?.data?.error || "Failed to approve claim");
      },
    });
  };

  // Reject claim (Donor action)
  const useRejectClaim = () => {
    return useMutation({
      mutationFn: (id: string) => claimAPI.rejectClaim(id),
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: claimKeys.all });
        toast.success("Claim rejected successfully!");
      },
      onError: (error: any) => {
        toast.error(error.response?.data?.error || "Failed to reject claim");
      },
    });
  };

  // Cancel claim (NGO action)
  const useCancelClaim = () => {
    return useMutation({
      mutationFn: (id: string) => claimAPI.cancelClaim(id),
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: claimKeys.all });
        toast.success("Claim cancelled successfully!");
      },
      onError: (error: any) => {
        toast.error(error.response?.data?.error || "Failed to cancel claim");
      },
    });
  };

  // Mark picked up (NGO action)
  const useMarkPickedUp = () => {
    return useMutation({
      mutationFn: (id: string) => claimAPI.markPickedUp(id),
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: claimKeys.all });
        toast.success("Claim marked as picked up!");
      },
      onError: (error: any) => {
        toast.error(error.response?.data?.error || "Failed to mark picked up");
      },
    });
  };

  // Mark delivered (NGO action)
  const useMarkDelivered = () => {
    return useMutation({
      mutationFn: (id: string) => claimAPI.markDelivered(id),
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: claimKeys.all });
        toast.success("Claim marked as delivered!");
      },
      onError: (error: any) => {
        toast.error(error.response?.data?.error || "Failed to mark delivered");
      },
    });
  };

  return {
    useGetClaimsByNGO,
    useGetClaimsByDonor,
    useGetClaimsByFood,
    useGetClaim,
    useCreateClaim,
    useApproveClaim,
    useRejectClaim,
    useCancelClaim,
    useMarkPickedUp,
    useMarkDelivered,
  };
}