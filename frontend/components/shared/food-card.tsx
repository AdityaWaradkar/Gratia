"use client";

import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { cn, formatDate, formatTimeLeft, getStatusColor } from "@/lib/utils";
import { FoodListing } from "@/types/food";
import { Clock, MapPin, Package, User } from "lucide-react";

interface FoodCardProps {
  listing: FoodListing;
  showClaim?: boolean;
  onClaim?: (listingId: string) => void;
  onView?: (listingId: string) => void;
  showActions?: boolean;
  role?: "donor" | "ngo";
}

const statusLabels: Record<string, string> = {
  AVAILABLE: "Available",
  CLAIMED: "Claimed",
  EXPIRED: "Expired",
  CANCELLED: "Cancelled",
};

export function FoodCard({
  listing,
  showClaim = false,
  onClaim,
  onView,
  showActions = false,
  role = "donor",
}: FoodCardProps) {
  const statusColor = getStatusColor(listing.status);
  const timeLeft = formatTimeLeft(listing.expiryTime);

  const handleClaim = () => {
    if (onClaim) {
      onClaim(listing.id);
    }
  };

  const handleView = () => {
    if (onView) {
      onView(listing.id);
    }
  };

  return (
    <Card className="hover:shadow-md transition-shadow overflow-hidden">
      {listing.imageUrl && (
        <div className="relative h-48 w-full bg-gray-100">
          <img
            src={listing.imageUrl}
            alt={listing.title}
            className="h-full w-full object-cover"
          />
        </div>
      )}

      <CardHeader className="pb-3">
        <div className="flex items-start justify-between gap-2">
          <CardTitle className="text-lg font-semibold line-clamp-1">
            {listing.title}
          </CardTitle>
          <Badge className={cn("shrink-0", statusColor)}>
            {statusLabels[listing.status] || listing.status}
          </Badge>
        </div>
      </CardHeader>

      <CardContent className="space-y-3">
        {listing.description && (
          <p className="text-sm text-gray-600 line-clamp-2">{listing.description}</p>
        )}

        <div className="space-y-2 text-sm">
          <div className="flex items-center gap-2 text-gray-600">
            <Package className="h-4 w-4 shrink-0" />
            <span>
              {listing.quantity} {listing.unit}
            </span>
          </div>

          <div className="flex items-center gap-2 text-gray-600">
            <MapPin className="h-4 w-4 shrink-0" />
            <span>{listing.location}</span>
          </div>

          <div className="flex items-center gap-2 text-gray-600">
            <Clock className="h-4 w-4 shrink-0" />
            <span className={cn(listing.status === "EXPIRED" && "text-red-600")}>
              {listing.status === "EXPIRED" ? "Expired" : timeLeft}
            </span>
          </div>

          {listing.donorUserId && (
            <div className="flex items-center gap-2 text-gray-600">
              <User className="h-4 w-4 shrink-0" />
              <span>Donor: {listing.donorUserId.slice(0, 8)}</span>
            </div>
          )}
        </div>
      </CardContent>

      <CardFooter className="pt-3 border-t border-gray-100 flex gap-2 flex-wrap">
        {showClaim && listing.status === "AVAILABLE" && (
          <Button
            size="sm"
            className="flex-1 bg-green-600 hover:bg-green-700"
            onClick={handleClaim}
          >
            Claim Food
          </Button>
        )}

        {showActions && (
          <Button
            size="sm"
            variant="outline"
            className="flex-1"
            onClick={handleView}
          >
            View Details
          </Button>
        )}

        {!showClaim && !showActions && (
          <Button
            size="sm"
            variant="outline"
            className="w-full"
            onClick={handleView}
          >
            View Details
          </Button>
        )}
      </CardFooter>
    </Card>
  );
}