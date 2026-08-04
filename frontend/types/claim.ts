export interface Claim {
  id: string
  foodListingId: string
  ngoUserId: string
  donorUserId: string
  status: 'CREATED' | 'ACCEPTED' | 'REJECTED' | 'PICKED_UP' | 'DELIVERED' | 'CANCELLED'
  createdAt: string
  updatedAt: string
  acceptedAt?: string
  rejectedAt?: string
  pickedUpAt?: string
  deliveredAt?: string
  cancelledAt?: string
}