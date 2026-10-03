import { useEffect, useState } from 'react'
import './App.css'

const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

type Dog = {
  id: number
  name: string
  breed: string
  size: string
  personality: string
  needs: string
  purveyor: string
}

function App() {
  const [dogs, setDogs] = useState<Dog[]>([])
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetch(`${API_URL}/api/dogs`)
      .then((res) => (res.ok ? res.json() : Promise.reject(res.statusText)))
      .then(setDogs)
      .catch((err) => setError(String(err)))
  }, [])

  return (
    <main>
      <h1>PetterHelp</h1>
      <p>Adoptable dogs in Cambridge, MA. Matching comes next.</p>
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
    </main>
  )
}

export default App
