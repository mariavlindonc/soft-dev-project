/* eslint-disable react-refresh/only-export-components */

import { createContext, useContext, useState, useCallback, type ReactNode } from 'react'
import type { User, LoginRequest, RegisterRequest } from '../types'
import * as authApi from '../api/auth'

interface AuthContextValue {
  user: User | null
  token: string | null
  login: (data: LoginRequest) => Promise<void>
  register: (data: RegisterRequest) => Promise<void>
  logout: () => void
  isAuthenticated: boolean
  isAdmin: boolean
}

const AuthContext = createContext<AuthContextValue | null>(null)

function getStoredUser(): User | null {
  const raw = localStorage.getItem('user')
  return raw ? (JSON.parse(raw) as User) : null
}

function getStoredToken(): string | null {
  return localStorage.getItem('token')
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(getStoredUser)
  const [token, setToken] = useState<string | null>(getStoredToken)

  const login = useCallback(async (data: LoginRequest) => {
    const res = await authApi.login(data)
    localStorage.setItem('token', res.token)
    localStorage.setItem('user', JSON.stringify(res.user))
    setToken(res.token)
    setUser(res.user)
  }, [])

  const register = useCallback(async (data: RegisterRequest) => {
    const res = await authApi.register(data)
    localStorage.setItem('token', res.token)
    localStorage.setItem('user', JSON.stringify(res.user))
    setToken(res.token)
    setUser(res.user)
  }, [])

  const logout = useCallback(() => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    setToken(null)
    setUser(null)
  }, [])

  const value: AuthContextValue = {
    user,
    token,
    login,
    register,
    logout,
    isAuthenticated: !!token && !!user,
    isAdmin: user?.role === 'admin',
  }

  return <AuthContext value={value}>{children}</AuthContext>
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
