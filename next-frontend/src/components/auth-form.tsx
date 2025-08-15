"use client";

import type React from "react";

import { useState } from "react";
import { ArrowRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { LoadingSpinner } from "@/components/loading-spinner";
import { useAuth } from "@/hooks/use-auth";

export function AuthForm() {
  const [apiKey, setApiKey] = useState("");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { login } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!apiKey.trim()) {
      setError("Por favor, insira sua API Key");
      return;
    }

    setIsSubmitting(true);
    setError("");

    console.log("----- apiKey -----");
    console.log(apiKey);

    const success = await login(apiKey.trim());

    console.log("----- success -----");
    console.log(success);

    if (!success) {
      setError(
        "API Key inválida. Verifique suas credenciais e tente novamente."
      );
    }

    setIsSubmitting(false);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <label
          htmlFor="apiKey"
          className="text-sm font-medium text-slate-200 font-dm-sans"
        >
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
        <Alert
          variant="destructive"
          className="bg-red-900/20 border-red-800 text-red-400"
        >
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <Button
        type="submit"
        className="w-full bg-blue-600 hover:bg-blue-700 text-white font-dm-sans"
        disabled={isSubmitting}
      >
        {isSubmitting ? (
          <LoadingSpinner size="sm" className="mr-2" />
        ) : (
          <ArrowRight className="w-4 h-4 mr-2" />
        )}
        {isSubmitting ? "Entrando..." : "Entrar"}
      </Button>
    </form>
  );
}
