export type ClaimStatus = "CREATED" | "ACCEPTED" | "REJECTED" | "PICKED_UP" | "DELIVERED" | "CANCELLED";

export interface Claim {
  id: string;
  foodListingId: string;
  ngoUserId: string;
  donorUserId: string;
  status: ClaimStatus;
  createdAt: string;
  updatedAt: string;
  acceptedAt?: string | null;
  rejectedAt?: string | null;
  pickedUpAt?: string | null;
  deliveredAt?: string | null;
  cancelledAt?: string | null;
}

export interface CreateClaimRequest {
  foodListingId: string;
}