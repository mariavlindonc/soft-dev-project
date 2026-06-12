import client from './client'
import type { AuthResponse, LoginRequest, RegisterRequest } from '../types/user'

export async function login(data: LoginRequest): Promise<AuthResponse> {
  const res = await client.post('/auth/login', data)
  return res.data
}

export async function register(data: RegisterRequest): Promise<AuthResponse> {
  const res = await client.post('/auth/register', data)
  return res.data
}
