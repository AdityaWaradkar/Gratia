"use client";

import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { cn, formatDate, getStatusColor } from "@/lib/utils";
import { Claim } from "@/types/claim";
import { FoodListing } from "@/types/food";
import { CheckCircle, XCircle, Clock, Truck, Package, User } from "lucide-react";

interface ClaimCardProps {
  claim: Claim;
  food?: FoodListing;
  onAction?: (claimId: string, action: string) => void;
  showActions?: boolean;
  role?: "donor" | "ngo";
}

const statusIcons: Record<string, React.ReactNode> = {
  CREATED: <Clock className="h-4 w-4" />,
  ACCEPTED: <CheckCircle className="h-4 w-4" />,
  REJECTED: <XCircle className="h-4 w-4" />,
  PICKED_UP: <Truck className="h-4 w-4" />,
  DELIVERED: <Package className="h-4 w-4" />,
  CANCELLED: <XCircle className="h-4 w-4" />,
};

const statusLabels: Record<string, string> = {
  CREATED: "Pending",
  ACCEPTED: "Accepted",
  REJECTED: "Rejected",
  PICKED_UP: "Picked Up",
  DELIVERED: "Delivered",
  CANCELLED: "Cancelled",
};

export function ClaimCard({ claim, food, onAction, showActions = false, role = "donor" }: ClaimCardProps) {
  const statusColor = getStatusColor(claim.status);
  const StatusIcon = statusIcons[claim.status] || <Clock className="h-4 w-4" />;

  const handleAction = (action: string) => {
    if (onAction) {
      onAction(claim.id, action);
    }
  };

  const getActions = () => {
    if (!showActions) return null;

    const actions: React.ReactNode[] = [];

    if (role === "donor") {
      if (claim.status === "CREATED") {
        actions.push(
          <Button
            key="approve"
            size="sm"
            className="bg-green-600 hover:bg-green-700"
            onClick={() => handleAction("approve")}
          >
            Approve
          </Button>
        );
        actions.push(
          <Button
            key="reject"
            size="sm"
            variant="outline"
            className="border-red-300 text-red-600 hover:bg-red-50"
            onClick={() => handleAction("reject")}
          >
            Reject
          </Button>
        );
      }
    }

    if (role === "ngo") {
      if (claim.status === "CREATED") {
        actions.push(
          <Button
            key="cancel"
            size="sm"
            variant="outline"
            className="border-red-300 text-red-600 hover:bg-red-50"
            onClick={() => handleAction("cancel")}
          >
            Cancel
          </Button>
        );
      }
      if (claim.status === "ACCEPTED") {
        actions.push(
          <Button
            key="pickup"
            size="sm"
            className="bg-blue-600 hover:bg-blue-700"
            onClick={() => handleAction("pickup")}
          >
            Mark Picked Up
          </Button>
        );
      }
      if (claim.status === "PICKED_UP") {
        actions.push(
          <Button
            key="deliver"
            size="sm"
            className="bg-purple-600 hover:bg-purple-700"
            onClick={() => handleAction("deliver")}
          >
            Mark Delivered
          </Button>
        );
      }
    }

    return actions.length > 0 ? <div className="flex gap-2 flex-wrap">{actions}</div> : null;
  };

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between">
          <CardTitle className="text-lg font-semibold">
            Claim #{claim.id.slice(0, 8)}
          </CardTitle>
          <Badge className={cn("flex items-center gap-1", statusColor)}>
            {StatusIcon}
            {statusLabels[claim.status] || claim.status}
          </Badge>
        </div>
      </CardHeader>

      <CardContent className="space-y-3">
        {food && (
          <div className="space-y-1">
            <p className="font-medium text-gray-900">{food.title}</p>
            <p className="text-sm text-gray-600">
              {food.quantity} {food.unit} • {food.location}
            </p>
          </div>
        )}

        <div className="grid grid-cols-2 gap-2 text-sm">
          <div className="flex items-center gap-2 text-gray-600">
            <User className="h-4 w-4" />
            <span>{role === "donor" ? "NGO" : "Donor"}: {claim.ngoUserId.slice(0, 8)}</span>
          </div>
          <div className="flex items-center gap-2 text-gray-600">
            <Clock className="h-4 w-4" />
            <span>{formatDate(claim.createdAt)}</span>
          </div>
        </div>

        {claim.acceptedAt && (
          <div className="text-sm text-green-600">
            Accepted: {formatDate(claim.acceptedAt)}
          </div>
        )}

        {claim.deliveredAt && (
          <div className="text-sm text-purple-600">
            Delivered: {formatDate(claim.deliveredAt)}
          </div>
        )}
      </CardContent>

      {(showActions || getActions()) && (
        <CardFooter className="pt-3 border-t border-gray-100">
          {getActions()}
        </CardFooter>
      )}
    </Card>
  );
}