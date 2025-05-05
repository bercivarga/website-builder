import { Link } from '@tanstack/react-router'

function App() {
  return (
    <div>
      <ul>
        <li>
          <Link to="/auth/login">Login</Link>
        </li>
        <li>
          <Link to="/auth/register">Register</Link>
        </li>
        <li>
          <Link to="/editor">Editor</Link>
        </li>
      </ul>
    </div>
  )
}

export default App
