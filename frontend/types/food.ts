export interface FoodListing {
  id: string
  donorUserId: string
  title: string
  description?: string
  quantity: number
  unit: string
  expiryTime: string
  location: string
  imageUrl?: string
  status: 'AVAILABLE' | 'CLAIMED' | 'EXPIRED' | 'CANCELLED'
  createdAt: string
  updatedAt: string
}

export interface CreateFoodListingRequest {
  title: string
  description?: string
  quantity: number
  unit: string
  expiryTime: string
  location: string
  imageUrl?: string
}