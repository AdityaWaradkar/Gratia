"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Users, Building, CheckCircle, Clock, Award } from "lucide-react";

export default function AdminDashboard() {
  const stats = [
    { label: "Total Users", value: 256, icon: Users, color: "bg-blue-100 text-blue-700" },
    { label: "NGOs", value: 48, icon: Building, color: "bg-green-100 text-green-700" },
    { label: "Pending Verification", value: 12, icon: Clock, color: "bg-yellow-100 text-yellow-700" },
    { label: "Verified NGOs", value: 36, icon: Award, color: "bg-purple-100 text-purple-700" },
  ];

  const pendingNGOs = [
    { id: 1, organization: "Food for All", registrationNo: "NGO2024-001", requestedAt: "2024-01-15" },
    { id: 2, organization: "Community Kitchen", registrationNo: "NGO2024-002", requestedAt: "2024-01-14" },
    { id: 3, organization: "Hunger Relief", registrationNo: "NGO2024-003", requestedAt: "2024-01-13" },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Admin Dashboard</h1>
        <p className="text-gray-600 mt-1">Manage users and oversee the platform</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {stats.map((stat) => {
          const Icon = stat.icon;
          return (
            <Card key={stat.label}>
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm text-gray-600">{stat.label}</p>
                    <p className="text-2xl font-bold mt-1">{stat.value}</p>
                  </div>
                  <div className={`p-3 rounded-lg ${stat.color}`}>
                    <Icon className="h-6 w-6" />
                  </div>
                </div>
              </CardContent>
            </Card>
          );
        })}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Pending NGO Verifications</CardTitle>
          <CardDescription>Review and verify new NGO registrations</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {pendingNGOs.map((ngo) => (
              <div
                key={ngo.id}
                className="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition"
              >
                <div>
                  <h3 className="font-medium">{ngo.organization}</h3>
                  <p className="text-sm text-gray-600">
                    Registration: {ngo.registrationNo} • Requested: {ngo.requestedAt}
                  </p>
                </div>
                <div className="flex items-center gap-4">
                  <Badge variant="warning">Pending</Badge>
                  <Button size="sm" variant="outline" className="gap-2">
                    <CheckCircle className="h-3 w-3" />
                    Verify
                  </Button>
                  <Button size="sm" variant="destructive">Reject</Button>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}