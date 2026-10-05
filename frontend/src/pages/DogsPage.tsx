import { useEffect, useState } from 'react'
import { dogs as dogsApi, type Dog } from '../api/client'
import { useAuth } from '../auth/context'

export default function DogsPage() {
  const { user } = useAuth()
  const [dogs, setDogs] = useState<Dog[]>([])
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    dogsApi
      .list()
      .then(setDogs)
      .catch((err) => setError(err instanceof Error ? err.message : String(err)))
  }, [])

  return (
    <>
      <p>
        {user
          ? `Welcome, ${user.name}. Ranked matches will land here once your lifestyle profile is in.`
          : 'Adoptable dogs in Cambridge, MA. Sign in to get matches ranked for you.'}
      </p>
      {error && <p className="error">Could not reach the API: {error}</p>}
      <ul className="dogs">
        {dogs.map((dog) => (
          <li key={dog.id}>
            <strong>{dog.name}</strong> · {dog.breed}, {dog.size}
            <br />
            <small>{dog.personality}</small>
            <br />
            <small>Needs: {dog.needs}</small>
            <br />
            <small>From {dog.purveyor}</small>
          </li>
        ))}
      </ul>
    </>
  )
}
