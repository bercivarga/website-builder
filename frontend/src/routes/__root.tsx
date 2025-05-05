import { createRootRoute, Link, Outlet } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'
import { useEffect } from 'react'

import { refreshToken } from '../utils/jwt'

import '../App.css'

export const Route = createRootRoute({
  component: () => (
    <MainContainer />
  ),
})

function MainContainer() {
  // TODO: extract this to a custom hook
  useEffect(() => {
    const onWindowFocus = async () => {
      await refreshToken()
    }

    window.addEventListener('focus', onWindowFocus)

    return () => {
      window.removeEventListener('focus', onWindowFocus)
    }
  }, []);

  return (
    <>
      <div className="p-2 flex gap-2">
        <Link to="/" className="[&.active]:font-bold">
          Home
        </Link>{' '}
        <Link to="/about" className="[&.active]:font-bold">
          About
        </Link>
      </div>
      <hr />
      <Outlet />
      <TanStackRouterDevtools />
    </>
  )
}