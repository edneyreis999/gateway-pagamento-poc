"use client"

import { LogOut } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useAuth } from "@/hooks/use-auth"

export function Navbar() {
  const { logout } = useAuth()

  const handleLogout = () => {
    logout()
  }

  return (
    <nav className="fixed top-0 left-0 right-0 z-50 bg-slate-800/95 backdrop-blur-sm border-b border-slate-700">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <div className="flex-shrink-0">
            <h1 className="text-xl font-bold text-white font-space-grotesk">Full Cycle Gateway</h1>
          </div>

          {/* User area */}
          <div className="flex items-center gap-4">
            <span className="text-sm text-slate-300 font-dm-sans">Olá, usuário</span>
            <Button
              onClick={handleLogout}
              variant="destructive"
              size="sm"
              className="bg-red-600 hover:bg-red-700 text-white"
            >
              <LogOut className="w-4 h-4 mr-2" />
              Logout
            </Button>
          </div>
        </div>
      </div>
    </nav>
  )
}
