import apiClient from './lib/api/client'
import './App.css'

function App() {
  const onLogin = async (email: string, password: string) => {
    try {
      const response = await apiClient.post('/auth/login', {
        email,
        password
      });

      console.log('Login response:', response);
    } catch (error) {
      console.error('Error during login:', error);
    }
  }

  const onRegister = async (username: string, password: string, email: string) => {
    try {
      const response = await apiClient.post('/auth/register', {
        username,
        password,
        email
      });
      console.log('Register response:', response);
    } catch (error) {
      console.error('Error during registration:', error);
    }
  }

  const onLoginAction = async (formData: FormData) => {
    const email = formData.get('email') as string;
    const password = formData.get('password') as string;

    if (!email || !password) {
      console.error('email and password are required');
      return;
    }

    await onLogin(email, password);
  }

  const onRegisterAction = async (formData: FormData) => {
    const username = formData.get('username') as string;
    const password = formData.get('password') as string;
    const email = formData.get('email') as string;

    if (!username || !password) {
      console.error('Username and password are required');
      return;
    }

    await onRegister(username, password, email);
  }

  return (
    <>
      <div>
        <h2>Login</h2>
        <form action={onLoginAction}>
          <input name="email" type="email" placeholder="Enter your email" />
          <input name="password" type="password" placeholder="Enter your password" />
          <button type="submit">Submit</button>
        </form>
        <br/>
        <h2>Register</h2>
        <form action={onRegisterAction}>
          <input name="username" type="text" placeholder="Enter your username" />
          <input name="email" type="email" placeholder="Enter your email" />
          <input name="password" type="password" placeholder="Enter your password" />
          <button type="submit">Submit</button>
        </form>
      </div>
    </>
  )
}

export default App
