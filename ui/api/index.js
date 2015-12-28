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

export function search(query) {
  return axios.get("/api/search/", { params: { q: query } });
}

export function addChannel(url) {
  return axios.post("/api/channels/", { url });
}

export function getChannels() {
  return axios.get("/api/channels/");
}

export function getChannel(id) {
  return axios.get(`/api/channels/${id}/`);
}

export function getLatestPodcasts(page=1) {
  return axios.get("/api/podcasts/latest/", { params: { page } });
}

export function subscribe(id) {
  return axios.post(`/api/subscriptions/${id}/`);
}

export function unsubscribe(id) {
  return axios.delete(`/api/subscriptions/${id}/`);
}

export function getBookmarks(page=1) {
  return axios.get("/api/bookmarks/", { params: { page } });
}

export function addBookmark(id) {
  return axios.post(`/api/bookmarks/${id}/`);
}

export function deleteBookmark(id) {
  return axios.delete(`/api/bookmarks/${id}/`);
}
