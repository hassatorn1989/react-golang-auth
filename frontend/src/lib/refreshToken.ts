import { api } from './api'
import { authStore } from './authStore'

let refreshPromise: Promise<string | null> | null = null

export async function refreshAccessToken(): Promise<string | null> {
    if (!refreshPromise) {
        refreshPromise = api
            .post('/auth/refresh')
            .then((res) => {
                const token = res.data.accessToken as string
                authStore.setAccessToken(token)
                return token
            })
            .catch(() => {
                authStore.setAccessToken(null)
                return null
            })
            .finally(() => {
                refreshPromise = null
            })
    }

    return refreshPromise
}