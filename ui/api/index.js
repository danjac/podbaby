import axios from 'axios';

axios.interceptors.request.use(config => {
    config.headers['X-CSRF-Token'] =  window.csrfToken;
    return config;
}, error => Promise.reject(error));

export function getCurrentUser() {
  return axios.get("/api/auth/currentuser/");
}

export function logout() {
  return axios.delete("/api/auth/logout/");
}

export function login(identifier, password) {
  return axios.post("/api/auth/login/", {
    identifier,
    password
  });
}

export function signup(name, email, password) {
  return axios.post("/api/auth/signup/", {
    name,
    email,
    password
  });
}

export function addChannel(url) {
  return axios.post("/api/channels/", { url });
}

export function getLatestPodcasts() {
  return axios.get("/api/podcasts/latest/");
}
