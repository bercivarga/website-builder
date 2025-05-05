import { createFileRoute, useRouter } from '@tanstack/react-router'
import apiClient from '../../lib/api/client';
import { JWTResponse, setJWT } from '../../utils/jwt';
import { useState } from 'react';

export const Route = createFileRoute('/auth/login')({
  component: RouteComponent,
})

function RouteComponent() {
  const [hasError, setHasError] = useState(false);

  const router = useRouter();

  const onLogin = async (email: string, password: string) => {
    try {
      const response = await apiClient.post('/auth/login', {
        email,
        password
      });

      if (!response.ok) {
        throw new Error('Network response was not ok');
      }

      const data = (await response.json() as JWTResponse);
      setJWT(data);
      router.navigate({to: '/editor'});
    } catch (error) {
      console.error('Error during login:', error);
      setHasError(true);
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

return (
  <>
    <div>
      <h2>Login</h2>
      <form action={onLoginAction}>
        <input name="email" type="email" placeholder="Enter your email" />
        <input name="password" type="password" placeholder="Enter your password" />
        <button type="submit">Submit</button>
      </form>
      {hasError && <p className="text-red-500">Login failed. Please try again.</p>}
    </div>
  </>
)
}
