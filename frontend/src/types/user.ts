export interface User {
  id: number
  name: string
  email: string
  role: 'client' | 'admin'
  created_at: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  name: string
  email: string
  password: string
}

export interface AuthResponse {
  token: string
  user: User
}
