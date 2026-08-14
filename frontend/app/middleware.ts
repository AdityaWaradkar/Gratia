import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

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
    // Check role from token and redirect accordingly
    try {
      const payload = JSON.parse(atob(token.split(".")[1]));
      const role = payload.role?.toLowerCase() || "donor";
      const validRoles = ["donor", "ngo", "admin"];
      const redirectPath = validRoles.includes(role) ? `/${role}` : "/donor";
      return NextResponse.redirect(new URL(redirectPath, request.url));
    } catch (_error) {
      // Invalid token - redirect to login
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