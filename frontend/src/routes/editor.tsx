import { createFileRoute } from '@tanstack/react-router'
import { useEffect, useState } from 'react'
import apiClient from '../lib/api/client'

export const Route = createFileRoute('/editor')({
  component: RouteComponent,
})

function RouteComponent() {
  const [userData, setUserData] = useState(null)

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const response = await apiClient.get('/user/me')
        if (!response.ok) {
          throw new Error('Network response was not ok')
        }
        const data = await response.json()
        setUserData(data)
      } catch (error) {
        console.error('Error fetching user data:', error)
      }
    }
    fetchUserData()
  }, [])

  if (!userData) {
    return <div>Loading...</div>
  }

  return (
    <div>
      <pre>
        {JSON.stringify(userData, null, 2)}
      </pre>
    </div>
  )
}
