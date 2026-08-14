export interface DonorProfile {
  id: string;
  userId: string;
  name: string;
  phone?: string | null;
  address?: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface NGOProfile {
  id: string;
  userId: string;
  organization: string;
  registrationNo: string;
  verified: boolean;
  verifiedBy?: string | null;
  verifiedAt?: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface CreateDonorProfileRequest {
  name: string;
  phone?: string | null;
  address?: string | null;
}

export interface CreateNGOProfileRequest {
  organization: string;
  registrationNo: string;
}

export interface UpdateDonorProfileRequest {
  name?: string;
  phone?: string | null;
  address?: string | null;
}