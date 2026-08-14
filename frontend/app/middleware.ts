import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

// Role mapping function
const getRoleRoute = (role: string): string => {
  const roleMap: Record<string, string> = {
    "DONOR": "donor",
    "NGO": "ngo",
    "ADMIN": "admin",
    "USER": "donor",
  };
  return roleMap[role] || "donor";
};

export function middleware(request: NextRequest) {
  const token = request.cookies.get("accessToken")?.value;
  const { pathname } = request.nextUrl;

  // Public routes
  const publicRoutes = ["/", "/login", "/register", "/forgot-password"];
  const isPublicRoute = publicRoutes.some(
    (route) => pathname === route || pathname.startsWith("/api")
  );

  // If no token and trying to access protected route, redirect to login
  if (!token && !isPublicRoute) {
    const loginUrl = new URL("/login", request.url);
    loginUrl.searchParams.set("redirect", pathname);
    return NextResponse.redirect(loginUrl);
  }

  // If token exists and trying to access auth pages, redirect to dashboard
  if (token && isPublicRoute && pathname !== "/") {
    try {
      const payload = JSON.parse(atob(token.split(".")[1]));
      const role = payload.role || "USER";
      const route = getRoleRoute(role);
      return NextResponse.redirect(new URL(`/${route}`, request.url));
    } catch (_error) {
      const response = NextResponse.redirect(new URL("/login", request.url));
      response.cookies.delete("accessToken");
      response.cookies.delete("userRole");
      return response;
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    "/((?!_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)",
  ],
};