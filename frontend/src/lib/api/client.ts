const fetchApi = async (url: string, options: RequestInit = {}) => {
  const baseUrl = import.meta.env.VITE_APP_API_URL;

  if (!baseUrl) {
    throw new Error('VITE_APP_API_URL is not defined');
  }

  if (!url.startsWith('/')) {
    throw new Error('URL must start with a forward slash');
  }

  const fullUrl = baseUrl + "/v1" + url;

  const response = await fetch(fullUrl, {
    ...options,
    headers: {
      ...options.headers,
    },
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  if (options?.headers?.['Content-Type' as keyof typeof options.headers] === 'application/json') {
    return response.json();
  }

  return response;
}

const get = async (url: string, options: RequestInit = {}) => {
  return fetchApi(url, { ...options, method: 'GET' });
};

const post = async (url: string, body: unknown, options: RequestInit = {}) => {
  return fetchApi(url, { ...options, method: 'POST', body: JSON.stringify(body) });
};

const put = async (url: string, body: unknown, options: RequestInit = {}) => {
  return fetchApi(url, { ...options, method: 'PUT', body: JSON.stringify(body) });
};

const del = async (url: string, options: RequestInit = {}) => {
  return fetchApi(url, { ...options, method: 'DELETE' });
};

const patch = async (url: string, body: unknown, options: RequestInit = {}) => {
  return fetchApi(url, { ...options, method: 'PATCH', body: JSON.stringify(body) });
};

const apiClient = {
  get,
  post,
  put,
  del,
  patch,
};

export default apiClient;