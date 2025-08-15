"use client"

import { useEffect } from "react"
import { useRouter } from "next/navigation"
import { useAuth } from "@/hooks/use-auth"
import { LoadingSpinner } from "@/components/loading-spinner"

export default function HomePage() {
  const { isAuthenticated, isLoading } = useAuth()
  const router = useRouter()

  console.log("---- HomePage: Component render ----", {
    isAuthenticated,
    isLoading
  });

  useEffect(() => {
    console.log("---- HomePage: useEffect triggered ----", {
      isAuthenticated,
      isLoading
    });
    
    if (!isLoading) {
      if (isAuthenticated) {
        console.log("---- HomePage: router.push('/invoices') ----");
        router.push("/invoices")
      } else {
        console.log("---- HomePage: router.push('/auth') ----");
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
