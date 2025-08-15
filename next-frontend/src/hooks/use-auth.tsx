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
import { logger } from "@/lib/logger";

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

  logger.auth('Provider mounted', { isAuthenticated, isLoading });

  useEffect(() => {
    logger.auth('Initializing authentication state');
    
    // Função para ler cookies no client-side
    const getCookie = (name: string) => {
      const value = `; ${document.cookie}`;
      const parts = value.split(`; ${name}=`);
      if (parts.length === 2) return parts.pop()?.split(";").shift();
      return null;
    };

    const storedApiKey = getCookie("apiKey");
    logger.debug('Cookie retrieved', { storedApiKey: storedApiKey ? '***' : null });

    if (storedApiKey) {
      logger.auth('API key found in cookie, setting authenticated state');
      setApiKey(storedApiKey);
      setIsAuthenticated(true);
      apiClient.setApiKey(storedApiKey);
    } else {
      logger.warn('No API key found in cookie');
    }
    
    logger.auth('Setting loading state to false');
    setIsLoading(false);
  }, []);

  const login = async (newApiKey: string): Promise<boolean> => {
    try {
      logger.auth('Login attempt initiated', { apiKey: '***' });
      setIsLoading(true);
      apiClient.setApiKey(newApiKey);

      // Try to fetch account to validate API key
      await apiClient.getAccount(newApiKey);

      localStorage.setItem("apiKey", newApiKey);
      setApiKey(newApiKey);
      setIsAuthenticated(true);
      logger.auth('Login successful, redirecting to invoices');
      router.push("/invoices");
      return true;
    } catch (error) {
      logger.error('Authentication failed', error);
      return false;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = () => {
    logger.auth('Logout initiated');
    localStorage.removeItem("apiKey");
    setApiKey(null);
    setIsAuthenticated(false);
    apiClient.setApiKey("");
    router.push("/auth");
  };

  logger.state('AuthProvider state update', {
    isAuthenticated,
    hasApiKey: !!apiKey,
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
