'use client'

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Plus, Package, Clock, CheckCircle, XCircle } from 'lucide-react'
import Link from 'next/link'

export default function DonorDashboard() {
  const stats = [
    { label: 'Active Listings', value: 12, icon: Package, color: 'bg-blue-100 text-blue-700' },
    { label: 'Pending Claims', value: 5, icon: Clock, color: 'bg-yellow-100 text-yellow-700' },
    { label: 'Completed', value: 45, icon: CheckCircle, color: 'bg-green-100 text-green-700' },
    { label: 'Cancelled', value: 3, icon: XCircle, color: 'bg-red-100 text-red-700' },
  ]

  const recentListings = [
    { id: 1, title: 'Fresh Vegetables', quantity: 50, unit: 'kg', status: 'AVAILABLE', createdAt: '2024-01-15' },
    { id: 2, title: 'Bread Pastry', quantity: 30, unit: 'pcs', status: 'CLAIMED', createdAt: '2024-01-14' },
    { id: 3, title: 'Cooked Meals', quantity: 20, unit: 'plates', status: 'AVAILABLE', createdAt: '2024-01-14' },
  ]

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Donor Dashboard</h1>
          <p className="text-gray-600 mt-1">Manage your food donations</p>
        </div>
        <Link href="/donor/create">
          <Button className="gap-2">
            <Plus className="h-4 w-4" />
            Add New Listing
          </Button>
        </Link>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {stats.map((stat) => {
          const Icon = stat.icon
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
          )
        })}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Recent Listings</CardTitle>
          <CardDescription>Your latest food donations</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {recentListings.map((listing) => (
              <div
                key={listing.id}
                className="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition"
              >
                <div>
                  <h3 className="font-medium">{listing.title}</h3>
                  <p className="text-sm text-gray-600">
                    {listing.quantity} {listing.unit} • {listing.createdAt}
                  </p>
                </div>
                <div className="flex items-center gap-4">
                  <Badge
                    variant={
                      listing.status === 'AVAILABLE' ? 'success' :
                      listing.status === 'CLAIMED' ? 'warning' : 'secondary'
                    }
                  >
                    {listing.status}
                  </Badge>
                  <Button variant="outline" size="sm">
                    View
                  </Button>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}