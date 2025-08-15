import { AuthForm } from "@/components/auth-form";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Info } from "lucide-react";

export default function AuthPage() {
  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <div className="w-full max-w-md space-y-6">
        <Card className="bg-slate-800 border-slate-700">
          <CardHeader className="text-center">
            <CardTitle className="text-2xl font-bold text-white font-space-grotesk">
              Autenticação Gateway
            </CardTitle>
            <CardDescription className="text-slate-400 font-dm-sans">
              Insira sua API Key para acessar o sistema
            </CardDescription>
          </CardHeader>
          <CardContent>
            <AuthForm />
          </CardContent>
        </Card>

        <Card className="bg-slate-800/50 border-slate-700">
          <CardContent className="pt-6">
            <div className="flex items-start gap-3">
              <Info className="w-5 h-5 text-blue-400 mt-0.5 flex-shrink-0" />
              <div>
                <h3 className="font-semibold text-white mb-2 font-space-grotesk">
                  Como obter uma API Key?
                </h3>
                <p className="text-sm text-slate-400 font-dm-sans">
                  Para obter sua API Key, você precisa criar uma conta de
                  comerciante. Entre em contato com nosso suporte para mais
                  informações.
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
