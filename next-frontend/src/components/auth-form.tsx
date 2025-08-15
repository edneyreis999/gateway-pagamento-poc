import type React from "react";

import { ArrowRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { apiClient } from "@/lib/api";

export async function loginAction(formData: FormData) {
  "use server";

  const apiKey = formData.get("apiKey") as string;

  await apiClient.getAccount(apiKey.trim());

  const cookiesStore = await cookies();
  cookiesStore.set("apiKey", apiKey.trim());

  redirect("/invoices");
}

export function AuthForm() {

  // const handleSubmit = async (e: React.FormEvent) => {
  //   e.preventDefault();

  //   if (!apiKey.trim()) {
  //     setError("Por favor, insira sua API Key");
  //     return;
  //   }

  //   setIsSubmitting(true);
  //   setError("");

  //   console.log("----- apiKey -----");
  //   console.log(apiKey);

  //   const success = await login(apiKey.trim());

  //   console.log("----- success -----");
  //   console.log(success);

  //   if (!success) {
  //     setError(
  //       "API Key inválida. Verifique suas credenciais e tente novamente."
  //     );
  //   }

  //   setIsSubmitting(false);
  // };

  return (
    <form action={loginAction} className="space-y-4">
      <div className="space-y-2">
        <label
          htmlFor="apiKey"
          className="text-sm font-medium text-slate-200 font-dm-sans"
        >
          API Key
        </label>
        <Input
          id="apiKey"
          placeholder="Digite sua API Key"
          className="bg-[#2a3749] border-gray-700 text-white placeholder-gray-400"
          name="apiKey"
        />
      </div>

      {/* {error && (
        <Alert
          variant="destructive"
          className="bg-red-900/20 border-red-800 text-red-400"
        >
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )} */}

      <Button
        type="submit"
        className="w-full bg-blue-600 hover:bg-blue-700 text-white font-dm-sans"
       
      >
        {/* {isSubmitting ? (
          <LoadingSpinner size="sm" className="mr-2" />
        ) : (
          <ArrowRight className="w-4 h-4 mr-2" />
        )}
        {isSubmitting ? "Entrando..." : "Entrar"} */}
        <ArrowRight className="w-4 h-4 mr-2" />
      </Button>
    </form>
  );
}
