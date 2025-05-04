import { useEffect, useState } from 'react'
import './App.css'
import apiClient from './lib/api/client'

type User = {
  id: number
  name: string
  email: string
}

function App() {
  const [user, setUser] = useState<User | null>(null)

  useEffect(() => {
    (async () => {
      const userData = await apiClient.get('/user/1');
      console.log(userData)
      setUser(userData)
    })()
  }, []);

  return (
    <>
      <div>
        {user ? (
          <pre>
            {JSON.stringify(user, null, 2)}
          </pre>) : (
            <h2>Loading...</h2>
          )}
      </div>
    </>
  )
}

export default App
