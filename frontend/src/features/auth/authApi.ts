import { replace } from 'react-router-dom'
import { api } from '../../lib/api'
import { authStore } from '../../lib/authStore'

export type LoginPayload = {
    email: string
    password: string
}

export type User = {
    id: number
    name: string
    email: string
}

export type LoginResponse = {
    accessToken: string
    user: User
}

export async function login(payload: LoginPayload) {
    const res = await api.post('/auth/login', payload)
    authStore.setAccessToken(res.data.accessToken)
    return res.data as LoginResponse
}

export async function logout() {
    await api.post('/auth/logout')
    authStore.setAccessToken(null)
}

export async function getMe() {
    const res = await api.get('/me')
    return res.data as User
}