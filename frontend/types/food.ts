export type FoodStatus = "AVAILABLE" | "CLAIMED" | "EXPIRED" | "CANCELLED";

export interface FoodListing {
  id: string;
  donorUserId: string;
  title: string;
  description?: string | null;
  quantity: number;
  unit: string;
  expiryTime: string;
  location: string;
  imageUrl?: string | null;
  status: FoodStatus;
  createdAt: string;
  updatedAt: string;
}

export interface CreateFoodListingRequest {
  title: string;
  description?: string | null;
  quantity: number;
  unit: string;
  expiryTime: string;
  location: string;
  imageUrl?: string | null;
}

export interface UpdateFoodListingRequest extends Partial<CreateFoodListingRequest> {
  id?: string;
}