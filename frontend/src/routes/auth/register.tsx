import { createFileRoute, useRouter } from '@tanstack/react-router'
import apiClient from '../../lib/api/client';
import { useState } from 'react';

export const Route = createFileRoute('/auth/register')({
  component: RouteComponent,
})

function RouteComponent() {
  const [hasError, setHasError] = useState(false);

  const router = useRouter();

  const onRegister = async (username: string, password: string, email: string) => {
    try {
      const response = await apiClient.post('/auth/register', {
        username,
        password,
        email
      });

      if (!response.ok) {
        throw new Error('Network response was not ok');
      }

      router.navigate({to: '/auth/login'});
    } catch (error) {
      console.error('Error during registration:', error);
      setHasError(true);
    }
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
        <h2>Register</h2>
        <form action={onRegisterAction}>
          <input name="username" type="text" placeholder="Enter your username" />
          <input name="email" type="email" placeholder="Enter your email" />
          <input name="password" type="password" placeholder="Enter your password" />
          <button type="submit">Submit</button>
        </form>
        {hasError && <p className="text-red-500">Registration failed. Please try again.</p>}
      </div>
    </>
  )
}
