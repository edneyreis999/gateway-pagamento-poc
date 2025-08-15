"use client"

import { createContext, useContext, useState, useEffect, type ReactNode } from "react"
import { useRouter } from "next/navigation"
import { apiClient } from "@/lib/api"

interface AuthContextType {
  isAuthenticated: boolean
  apiKey: string | null
  login: (apiKey: string) => Promise<boolean>
  logout: () => void
  isLoading: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [apiKey, setApiKey] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const router = useRouter()

  useEffect(() => {
    // Check for stored API key on mount
    const storedApiKey = localStorage.getItem("gateway_api_key")
    if (storedApiKey) {
      setApiKey(storedApiKey)
      setIsAuthenticated(true)
      apiClient.setApiKey(storedApiKey)
    }
    setIsLoading(false)
  }, [])

  const login = async (newApiKey: string): Promise<boolean> => {
    try {
      setIsLoading(true)
      apiClient.setApiKey(newApiKey)

      // Try to fetch account to validate API key
      await apiClient.getAccount(newApiKey)

      localStorage.setItem("gateway_api_key", newApiKey)
      setApiKey(newApiKey)
      setIsAuthenticated(true)
      router.push("/invoices")
      return true
    } catch (error) {
      console.error("Authentication failed:", error)
      return false
    } finally {
      setIsLoading(false)
    }
  }

  const logout = () => {
    localStorage.removeItem("gateway_api_key")
    setApiKey(null)
    setIsAuthenticated(false)
    apiClient.setApiKey("")
    router.push("/auth")
  }

  return (
    <AuthContext.Provider value={{ isAuthenticated, apiKey, login, logout, isLoading }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider")
  }
  return context
}
