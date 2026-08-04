export interface DonorProfile {
  id: string
  userId: string
  name: string
  phone?: string
  address?: string
  createdAt: string
  updatedAt: string
}

export interface NGOProfile {
  id: string
  userId: string
  organization: string
  registrationNo: string
  verified: boolean
  verifiedBy?: string
  verifiedAt?: string
  createdAt: string
  updatedAt: string
}