"use client";

import {
  createContext,
  useContext,
  useState,
  useEffect,
  type ReactNode,
} from "react";
import { useRouter } from "next/navigation";
import { apiClient } from "@/lib/api";

interface AuthContextType {
  isAuthenticated: boolean;
  apiKey: string | null;
  login: (apiKey: string) => Promise<boolean>;
  logout: () => void;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [apiKey, setApiKey] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();

  console.log("---- AuthProvider: Component mounted ----");
  console.log("---- AuthProvider: Current state ----", {
    isAuthenticated,
    apiKey,
    isLoading
  });

  useEffect(() => {
    console.log("---- AuthProvider: useEffect triggered ----");
    console.log("---- AuthProvider: Client-side useEffect ----");
    console.log("---- AuthProvider: Window available? ----", typeof window !== 'undefined');
    
    // Função para ler cookies no client-side
    const getCookie = (name: string) => {
      const value = `; ${document.cookie}`;
      const parts = value.split(`; ${name}=`);
      if (parts.length === 2) return parts.pop()?.split(";").shift();
      return null;
    };

    const storedApiKey = getCookie("apiKey");
    console.log("---- AuthProvider: Cookie value ----", storedApiKey);
    console.log("---- AuthProvider: All cookies ----", document.cookie);

    if (storedApiKey) {
      console.log("---- AuthProvider: Setting authenticated state ----");
      setApiKey(storedApiKey);
      setIsAuthenticated(true);
      apiClient.setApiKey(storedApiKey);
    } else {
      console.log("---- AuthProvider: No API key found in cookie ----");
    }
    
    console.log("---- AuthProvider: Setting isLoading to false ----");
    setIsLoading(false);
  }, []);

  const login = async (newApiKey: string): Promise<boolean> => {
    try {
      console.log("---- AuthProvider: Login function called ----", { newApiKey });
      setIsLoading(true);
      apiClient.setApiKey(newApiKey);

      // Try to fetch account to validate API key
      await apiClient.getAccount(newApiKey);

      localStorage.setItem("apiKey", newApiKey);
      setApiKey(newApiKey);
      setIsAuthenticated(true);
      console.log("---- AuthProvider: Login successful, redirecting to /invoices ----");
      router.push("/invoices");
      return true;
    } catch (error) {
      console.error("Authentication failed:", error);
      return false;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = () => {
    console.log("---- AuthProvider: Logout function called ----");
    localStorage.removeItem("apiKey");
    setApiKey(null);
    setIsAuthenticated(false);
    apiClient.setApiKey("");
    router.push("/auth");
  };

  console.log("---- AuthProvider: Before return ----", {
    isAuthenticated,
    apiKey,
    isLoading
  });

  return (
    <AuthContext.Provider
      value={{ isAuthenticated, apiKey, login, logout, isLoading }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
