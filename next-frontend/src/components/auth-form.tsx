"use client"

import type React from "react"

import { useState } from "react"
import { ArrowRight, Info } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { LoadingSpinner } from "@/components/loading-spinner"
import { useAuth } from "@/hooks/use-auth"

export function AuthForm() {
  const [apiKey, setApiKey] = useState("")
  const [error, setError] = useState("")
  const [isSubmitting, setIsSubmitting] = useState(false)
  const { login } = useAuth()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!apiKey.trim()) {
      setError("Por favor, insira sua API Key")
      return
    }

    setIsSubmitting(true)
    setError("")

    console.log('----- apiKey -----');
    console.log(apiKey);

    const success = await login(apiKey.trim())

    console.log('----- success -----');
    console.log(success);

    if (!success) {
      setError("API Key inválida. Verifique suas credenciais e tente novamente.")
    }

    setIsSubmitting(false)
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <div className="w-full max-w-md space-y-6">
        <Card className="bg-slate-800 border-slate-700">
          <CardHeader className="text-center">
            <CardTitle className="text-2xl font-bold text-white font-space-grotesk">Autenticação Gateway</CardTitle>
            <CardDescription className="text-slate-400 font-dm-sans">
              Insira sua API Key para acessar o sistema
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <label htmlFor="apiKey" className="text-sm font-medium text-slate-200 font-dm-sans">
                  API Key
                </label>
                <Input
                  id="apiKey"
                  type="text"
                  value={apiKey}
                  onChange={(e) => setApiKey(e.target.value)}
                  placeholder="Digite sua API Key"
                  className="bg-slate-700 border-slate-600 text-white placeholder:text-slate-400"
                  disabled={isSubmitting}
                />
              </div>

              {error && (
                <Alert variant="destructive" className="bg-red-900/20 border-red-800 text-red-400">
                  <AlertDescription>{error}</AlertDescription>
                </Alert>
              )}

              <Button
                type="submit"
                className="w-full bg-blue-600 hover:bg-blue-700 text-white font-dm-sans"
                disabled={isSubmitting}
              >
                {isSubmitting ? <LoadingSpinner size="sm" className="mr-2" /> : <ArrowRight className="w-4 h-4 mr-2" />}
                {isSubmitting ? "Entrando..." : "Entrar"}
              </Button>
            </form>
          </CardContent>
        </Card>

        <Card className="bg-slate-800/50 border-slate-700">
          <CardContent className="pt-6">
            <div className="flex items-start gap-3">
              <Info className="w-5 h-5 text-blue-400 mt-0.5 flex-shrink-0" />
              <div>
                <h3 className="font-semibold text-white mb-2 font-space-grotesk">Como obter uma API Key?</h3>
                <p className="text-sm text-slate-400 font-dm-sans">
                  Para obter sua API Key, você precisa criar uma conta de comerciante. Entre em contato com nosso
                  suporte para mais informações.
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
