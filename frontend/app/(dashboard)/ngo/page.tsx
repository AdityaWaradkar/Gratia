'use client'

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Clock, CheckCircle, Package, AlertCircle } from 'lucide-react'
import Link from 'next/link'

export default function NGODashboard() {
  const stats = [
    { label: 'Available Food', value: 24, icon: Package, color: 'bg-green-100 text-green-700' },
    { label: 'Pending Claims', value: 8, icon: Clock, color: 'bg-yellow-100 text-yellow-700' },
    { label: 'Delivered', value: 67, icon: CheckCircle, color: 'bg-blue-100 text-blue-700' },
    { label: 'Cancelled', value: 2, icon: AlertCircle, color: 'bg-red-100 text-red-700' },
  ]

  const availableFood = [
    { id: 1, title: 'Fresh Vegetables', donor: 'Green Harvest', quantity: 50, unit: 'kg', distance: '2.5 km' },
    { id: 2, title: 'Bread Pastry', donor: 'Bakery House', quantity: 30, unit: 'pcs', distance: '5.1 km' },
    { id: 3, title: 'Cooked Meals', donor: 'Taste of India', quantity: 20, unit: 'plates', distance: '3.8 km' },
  ]

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">NGO Dashboard</h1>
          <p className="text-gray-600 mt-1">Find and claim food for your community</p>
        </div>
        <Link href="/ngo/foods">
          <Button className="gap-2">
            <Package className="h-4 w-4" />
            Browse Food
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
          <CardTitle>Available Food Nearby</CardTitle>
          <CardDescription>Food donations ready to be claimed</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {availableFood.map((food) => (
              <div
                key={food.id}
                className="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition"
              >
                <div>
                  <h3 className="font-medium">{food.title}</h3>
                  <p className="text-sm text-gray-600">
                    {food.donor} • {food.quantity} {food.unit} • {food.distance}
                  </p>
                </div>
                <div className="flex items-center gap-4">
                  <Badge variant="success">Available</Badge>
                  <Button size="sm">Claim</Button>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}