import { getCookie, setCookie } from "./cookies";
import { JWT_COOKIE_NAME, JWT_REFRESH_COOKIE_NAME } from "./config";
import apiClient from "../lib/api/client";

export type JWTResponse = {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type: string;
};

export const setJWT = (response: JWTResponse) => {
  const { access_token, refresh_token, expires_in } = response;
  setCookie(JWT_COOKIE_NAME, access_token, expires_in);

  const refreshExpiration = 7 * 24 * 60 * 60;
  setCookie(JWT_REFRESH_COOKIE_NAME, refresh_token, refreshExpiration);
}

export const getJWT = () => {
  return getCookie(JWT_COOKIE_NAME);
}

export const getRefreshJWT = () => {
  return getCookie(JWT_REFRESH_COOKIE_NAME);
}

export const deleteJWT = () => {
  setCookie(JWT_COOKIE_NAME, '', -1);
  setCookie(JWT_REFRESH_COOKIE_NAME, '', -1);
}

export const isJWTExpired = () => {
  const token = getJWT();
  if (!token) return true;

  const payload = JSON.parse(atob(token.split('.')[1]));
  const exp = payload.exp;
  const now = Math.floor(Date.now() / 1000);

  return exp < now;
}

export const refreshToken = async () => {
  const refreshToken = getRefreshJWT();
  if (!refreshToken) return null;

  const response = await apiClient.post('/auth/refresh', null, {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${refreshToken}`
    }
  });

  if (!response.ok) {
    deleteJWT();
    return null;
  }

  const data = await response.json();
  setJWT(data);
  return data;
}