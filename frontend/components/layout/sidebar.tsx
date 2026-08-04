'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { 
  LayoutDashboard, 
  Utensils, 
  ClipboardList, 
  Users, 
  LogOut,
  Leaf,
  Award,
  Home
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { useRouter } from 'next/navigation'
import { toast } from 'sonner'

interface SidebarProps {
  role?: 'user' | 'admin'
}

export function Sidebar({ role = 'user' }: SidebarProps) {
  const pathname = usePathname()
  const router = useRouter()

  const donorLinks = [
    { href: '/donor', label: 'Dashboard', icon: LayoutDashboard },
    { href: '/donor/listings', label: 'My Listings', icon: Utensils },
    { href: '/donor/create', label: 'Add Food', icon: ClipboardList },
  ]

  const ngoLinks = [
    { href: '/ngo', label: 'Dashboard', icon: LayoutDashboard },
    { href: '/ngo/foods', label: 'Available Food', icon: Utensils },
    { href: '/ngo/claims', label: 'My Claims', icon: ClipboardList },
  ]

  const adminLinks = [
    { href: '/admin', label: 'Dashboard', icon: LayoutDashboard },
    { href: '/admin/users', label: 'Users', icon: Users },
    { href: '/admin/ngos', label: 'NGOs', icon: Award },
  ]

  let links = donorLinks
  if (role === 'admin') links = adminLinks

  const handleLogout = () => {
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
    toast.success('Logged out successfully')
    router.push('/login')
  }

  return (
    <aside className="fixed left-0 top-0 h-full w-64 bg-white border-r border-gray-200 flex flex-col">
      <div className="p-6 border-b border-gray-200">
        <div className="flex items-center gap-2">
          <Leaf className="h-8 w-8 text-green-600" />
          <span className="text-xl font-bold text-green-700">Gratia</span>
        </div>
      </div>

      <nav className="flex-1 p-4 space-y-2">
        {links.map((link) => {
          const Icon = link.icon
          const isActive = pathname === link.href || pathname.startsWith(link.href + '/')
          
          return (
            <Link
              key={link.href}
              href={link.href}
              className={cn(
                'flex items-center gap-3 px-4 py-3 rounded-lg text-sm font-medium transition-colors',
                isActive 
                  ? 'bg-green-50 text-green-700' 
                  : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
              )}
            >
              <Icon className="h-5 w-5" />
              {link.label}
            </Link>
          )
        })}
      </nav>

      <div className="p-4 border-t border-gray-200 space-y-2">
        <Link
          href="/"
          className="flex items-center gap-3 px-4 py-3 rounded-lg text-sm font-medium text-gray-600 hover:bg-gray-50 transition-colors"
        >
          <Home className="h-5 w-5" />
          Home
        </Link>
        <Button
          variant="ghost"
          className="w-full justify-start gap-3 px-4 py-3 text-sm font-medium text-gray-600 hover:bg-gray-50 hover:text-red-600"
          onClick={handleLogout}
        >
          <LogOut className="h-5 w-5" />
          Logout
        </Button>
      </div>
    </aside>
  )
}