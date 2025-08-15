"use client";

import { useEffect, type ReactNode } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/hooks/use-auth";
import { LoadingSpinner } from "@/components/loading-spinner";

interface ProtectedRouteProps {
  children: ReactNode;
}

export function ProtectedRoute({ children }: { children: ReactNode }) {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  console.log("---- ProtectedRoute: Component render ----", {
    isAuthenticated,
    isLoading
  });

  useEffect(() => {
    console.log("---- ProtectedRoute: useEffect triggered ----", {
      isAuthenticated,
      isLoading,
      shouldRedirect: !isLoading && !isAuthenticated
    });
    
    if (!isLoading && !isAuthenticated) {
      console.log("---- ProtectedRoute: Redirecting to /auth ----");
      router.push("/auth");
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading) {
    console.log("---- ProtectedRoute: Showing loading spinner ----");
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  if (!isAuthenticated) {
    console.log("---- ProtectedRoute: Not authenticated, returning null ----");
    return null;
  }

  console.log("---- ProtectedRoute: Rendering children ----");
  return <>{children}</>;
}
