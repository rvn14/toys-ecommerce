"use client";

import {
  useCallback,
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { useRouter } from "next/navigation";
import type {
  AuthData,
  LoginRequest,
  RegisterRequest,
  User,
} from "@/lib/types";
import {
  getCurrentUser,
  loginUser,
  logoutUser,
  registerUser,
} from "@/lib/auth-api";

const TOKEN_STORAGE_KEY = "toys_ecommerce_access_token";

type AuthContextValue = {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (request: LoginRequest) => Promise<void>;
  register: (request: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

type AuthProviderProps = {
  children: React.ReactNode;
};

export function AuthProvider({ children }: AuthProviderProps) {
  const router = useRouter();

  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const saveAuthData = useCallback(async (authData: AuthData) => {
    localStorage.setItem(TOKEN_STORAGE_KEY, authData.accessToken);
    setToken(authData.accessToken);
    setUser(authData.user);
  }, []);

  const login = useCallback(async (request: LoginRequest) => {
    const authData = await loginUser(request);
    await saveAuthData(authData);
    router.replace(authData.user.role === "admin" ? "/admin" : "/account");
    router.refresh();
  }, [router, saveAuthData]);

  const register = useCallback(async (request: RegisterRequest) => {
    const authData = await registerUser(request);
    await saveAuthData(authData);
    router.replace("/account");
    router.refresh();
  }, [router, saveAuthData]);

  const logout = useCallback(async () => {
    const currentToken = localStorage.getItem(TOKEN_STORAGE_KEY);

    if (currentToken) {
      try {
        await logoutUser(currentToken);
      } catch {
        // Even if backend logout fails, remove local token from browser.
      }
    }

    localStorage.removeItem(TOKEN_STORAGE_KEY);
    setToken(null);
    setUser(null);
    router.replace("/login");
  }, [router]);

  const refreshUser = useCallback(async () => {
    const storedToken = localStorage.getItem(TOKEN_STORAGE_KEY);

    if (!storedToken) {
      setToken(null);
      setUser(null);
      setIsLoading(false);
      return;
    }

    try {
      const currentUser = await getCurrentUser(storedToken);
      setToken(storedToken);
      setUser(currentUser);
    } catch {
      localStorage.removeItem(TOKEN_STORAGE_KEY);
      setToken(null);
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    const timeoutId = window.setTimeout(() => {
      void refreshUser();
    }, 0);

    return () => window.clearTimeout(timeoutId);
  }, [refreshUser]);

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      token,
      isLoading,
      isAuthenticated: Boolean(user && token),
      login,
      register,
      logout,
      refreshUser,
    }),
    [user, token, isLoading, login, register, logout, refreshUser]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error("useAuth must be used inside AuthProvider");
  }

  return context;
}
