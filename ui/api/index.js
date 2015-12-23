import axios from 'axios';

axios.interceptors.request.use(config => {
    config.headers['X-CSRF-Token'] =  window.csrfToken;
    return config;
}, error => Promise.reject(error));


export function signup(name, email, password) {
  return axios.post("/api/auth/signup/", {
    name,
    email,
    password
  })
}
