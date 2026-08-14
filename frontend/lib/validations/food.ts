import * as z from "zod";

export const createFoodListingSchema = z.object({
  title: z.string().min(1, "Title is required"),
  description: z.string().optional().nullable(),
  quantity: z.number().min(1, "Quantity must be at least 1"),
  unit: z.string().min(1, "Unit is required"),
  expiryTime: z.string().refine((val) => {
    const date = new Date(val);
    return !isNaN(date.getTime()) && date > new Date();
  }, "Expiry time must be in the future"),
  location: z.string().min(1, "Location is required"),
  imageUrl: z.string().url("Must be a valid URL").optional().nullable(),
});

export const updateFoodListingSchema = z.object({
  title: z.string().min(1, "Title is required").optional(),
  description: z.string().optional().nullable(),
  quantity: z.number().min(1, "Quantity must be at least 1").optional(),
  unit: z.string().min(1, "Unit is required").optional(),
  expiryTime: z
    .string()
    .refine((val) => {
      const date = new Date(val);
      return !isNaN(date.getTime()) && date > new Date();
    }, "Expiry time must be in the future")
    .optional(),
  location: z.string().min(1, "Location is required").optional(),
  imageUrl: z.string().url("Must be a valid URL").optional().nullable(),
});

export type CreateFoodListingFormData = z.infer<typeof createFoodListingSchema>;
export type UpdateFoodListingFormData = z.infer<typeof updateFoodListingSchema>;