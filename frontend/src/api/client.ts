// Thin fetch wrapper for the PetterHelp API. Attaches the bearer token when
// one is set and turns non-2xx responses into thrown ApiErrors carrying the
// server's message, so forms can show it directly.

const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

export class ApiError extends Error {
  status: number
  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

let token: string | null = null

export function setToken(t: string | null) {
  token = t
}

export async function api<T>(path: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers)
  if (init.body) headers.set('Content-Type', 'application/json')
  if (token) headers.set('Authorization', `Bearer ${token}`)

  const res = await fetch(`${API_URL}${path}`, { ...init, headers })
  if (!res.ok) {
    let message = res.statusText
    try {
      const body = await res.json()
      if (typeof body.message === 'string') message = body.message
    } catch {
      // not JSON; keep statusText
    }
    throw new ApiError(res.status, message)
  }
  return res.json() as Promise<T>
}

export type Role = 'seeker' | 'purveyor'

export type User = {
  id: number
  email: string
  name: string
  role: Role
  created_at: string
}

export type Session = { user: User; token: string }

export type Dog = {
  id: number
  name: string
  breed: string
  size: string
  personality: string
  needs: string
  purveyor: string
}

export const auth = {
  register: (body: { email: string; password: string; name: string; role: Role }) =>
    api<Session>('/api/auth/register', { method: 'POST', body: JSON.stringify(body) }),
  login: (body: { email: string; password: string }) =>
    api<Session>('/api/auth/login', { method: 'POST', body: JSON.stringify(body) }),
  me: () => api<User>('/api/auth/me'),
}

export const dogs = {
  list: () => api<Dog[]>('/api/dogs'),
}
