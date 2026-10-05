import { BrowserRouter, Link, Navigate, Route, Routes } from 'react-router-dom'
import './App.css'
import { AuthProvider } from './auth/AuthContext'
import { useAuth } from './auth/context'
import DogsPage from './pages/DogsPage'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'

function Header() {
  const { user, loading, logout } = useAuth()
  return (
    <header>
      <Link to="/" className="brand">
        PetterHelp
      </Link>
      <nav>
        {loading ? null : user ? (
          <>
            <span className="muted">
              {user.name} · {user.role}
            </span>
            <button type="button" className="link" onClick={logout}>
              Sign out
            </button>
          </>
        ) : (
          <>
            <Link to="/login">Sign in</Link>
            <Link to="/register" className="button">
              Sign up
            </Link>
          </>
        )}
      </nav>
    </header>
  )
}

// Sends signed-in users away from the auth pages.
function GuestOnly({ children }: { children: React.ReactElement }) {
  const { user, loading } = useAuth()
  if (loading) return null
  return user ? <Navigate to="/" replace /> : children
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Header />
        <main>
          <Routes>
            <Route path="/" element={<DogsPage />} />
            <Route
              path="/login"
              element={
                <GuestOnly>
                  <LoginPage />
                </GuestOnly>
              }
            />
            <Route
              path="/register"
              element={
                <GuestOnly>
                  <RegisterPage />
                </GuestOnly>
              }
            />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </main>
      </AuthProvider>
    </BrowserRouter>
  )
}
