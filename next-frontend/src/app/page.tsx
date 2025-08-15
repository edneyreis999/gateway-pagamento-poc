"use client"

import { useEffect } from "react"
import { useRouter } from "next/navigation"
import { useAuth } from "@/hooks/use-auth"
import { LoadingSpinner } from "@/components/loading-spinner"
import { logger } from "@/lib/logger"

export default function HomePage() {
  const { isAuthenticated, isLoading } = useAuth()
  const router = useRouter()

  logger.route('HomePage render', { isAuthenticated, isLoading });

  useEffect(() => {
    logger.route('HomePage effect triggered', {
      isAuthenticated,
      isLoading
    });
    
    if (!isLoading) {
      if (isAuthenticated) {
        logger.route('User authenticated, redirecting to invoices');
        router.push("/invoices")
      } else {
        logger.route('User not authenticated, redirecting to auth');
        router.push("/auth")
      }
    }
  }, [isAuthenticated, isLoading, router])

  return (
    <div className="min-h-screen flex items-center justify-center">
      <LoadingSpinner size="lg" />
    </div>
  )
}
