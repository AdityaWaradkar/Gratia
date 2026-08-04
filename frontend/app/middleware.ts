import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

export function middleware(request: NextRequest) {
  const token = request.cookies.get('accessToken')?.value
  const { pathname } = request.nextUrl

  // Public routes
  const publicRoutes = ['/', '/login', '/register', '/forgot-password']
  const isPublicRoute = publicRoutes.some(route => pathname === route || pathname.startsWith('/api'))

  // If no token and trying to access protected route, redirect to login
  if (!token && !isPublicRoute) {
    const loginUrl = new URL('/login', request.url)
    loginUrl.searchParams.set('redirect', pathname)
    return NextResponse.redirect(loginUrl)
  }

  // If token exists and trying to access auth pages, redirect to dashboard
  if (token && isPublicRoute && pathname !== '/') {
    return NextResponse.redirect(new URL('/donor', request.url))
  }

  return NextResponse.next()
}

export const config = {
  matcher: [
    '/((?!_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)',
  ],
}