"use client";

import { useEffect, type ReactNode } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/hooks/use-auth";
import { LoadingSpinner } from "@/components/loading-spinner";
import { logger } from "@/lib/logger";

interface ProtectedRouteProps {
  children: ReactNode;
}

export function ProtectedRoute({ children }: { children: ReactNode }) {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  logger.route('ProtectedRoute render', { isAuthenticated, isLoading });

  useEffect(() => {
    logger.route('ProtectedRoute effect triggered', {
      isAuthenticated,
      isLoading,
      shouldRedirect: !isLoading && !isAuthenticated
    });
    
    if (!isLoading && !isAuthenticated) {
      logger.warn('Unauthorized access, redirecting to auth');
      router.push("/auth");
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading) {
    logger.debug('Showing loading spinner');
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  if (!isAuthenticated) {
    logger.warn('Not authenticated, blocking access');
    return null;
  }

  logger.route('Access granted, rendering protected content');
  return <>{children}</>;
}
