import axios from 'axios'
import { authStore } from './authStore'

export const api = axios.create({
    baseURL: 'http://localhost:8080/api',
    withCredentials: true,
})

api.interceptors.request.use((config) => {
    const token = authStore.getAccessToken()
    if (token) {
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
})