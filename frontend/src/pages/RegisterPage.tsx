import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import type { Role } from '../api/client'
import { useAuth } from '../auth/context'

export default function RegisterPage() {
  const { register } = useAuth()
  const navigate = useNavigate()
  const [role, setRole] = useState<Role>('seeker')
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  async function onSubmit(e: FormEvent) {
    e.preventDefault()
    setError(null)
    setBusy(true)
    try {
      await register({ email, password, name, role })
      navigate('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not create account')
    } finally {
      setBusy(false)
    }
  }

  return (
    <section className="card">
      <h2>Create an account</h2>
      <form onSubmit={onSubmit}>
        <fieldset className="roles">
          <legend>I am a…</legend>
          <label className={role === 'seeker' ? 'selected' : ''}>
            <input
              type="radio"
              name="role"
              value="seeker"
              checked={role === 'seeker'}
              onChange={() => setRole('seeker')}
            />
            <span>
              <strong>Seeker</strong>
              <small>Looking for a dog to adopt</small>
            </span>
          </label>
          <label className={role === 'purveyor' ? 'selected' : ''}>
            <input
              type="radio"
              name="role"
              value="purveyor"
              checked={role === 'purveyor'}
              onChange={() => setRole('purveyor')}
            />
            <span>
              <strong>Purveyor</strong>
              <small>Shelter, rescue, breeder or rehoming a dog</small>
            </span>
          </label>
        </fieldset>
        <label>
          Name
          <input
            type="text"
            autoComplete="name"
            required
            maxLength={100}
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </label>
        <label>
          Email
          <input
            type="email"
            autoComplete="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </label>
        <label>
          Password
          <input
            type="password"
            autoComplete="new-password"
            required
            minLength={8}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <small className="muted">At least 8 characters</small>
        </label>
        {error && <p className="error">{error}</p>}
        <button type="submit" disabled={busy}>
          {busy ? 'Creating…' : 'Create account'}
        </button>
      </form>
      <p className="muted">
        Already have an account? <Link to="/login">Sign in</Link>
      </p>
    </section>
  )
}
