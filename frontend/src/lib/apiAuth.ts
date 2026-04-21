
import { api } from './api'
import { refreshAccessToken } from './refreshToken'
import { authStore } from './authStore'

let isRedirecting = false

api.interceptors.request.use((config) => {
    const token = authStore.getAccessToken()
    if (token) {
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
})


let isRefreshing = false;
let failedQueue: Array<{
    resolve: (token: string) => void;
    reject: (error: unknown) => void;
}> = [];

function processQueue(error: unknown, token: string | null = null) {
    failedQueue.forEach((promise) => {
        if (error) {
            promise.reject(error);
        } else if (token) {
            promise.resolve(token);
        }
    });

    failedQueue = [];
}

api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config as any;
        if (error.response?.status === 401 && !originalRequest._retry) {
            if (isRefreshing) {
                return new Promise((resolve, reject) => {
                    failedQueue.push({ resolve, reject });
                }).then((token) => {
                    originalRequest.headers.Authorization = `Bearer ${token}`;
                    return api(originalRequest);
                }).catch((err) => {
                    return Promise.reject(err);
                });
            }

            originalRequest._retry = true;
            isRefreshing = true;
            
            return new Promise(async (resolve, reject) => {
                try {
                    const token = await refreshAccessToken();
                    if (token) {
                        originalRequest.headers.Authorization = `Bearer ${token}`;
                        processQueue(null, token);
                        resolve(api(originalRequest));
                    } else {
                        processQueue(new Error('Failed to refresh token'));
                        reject(new Error('Failed to refresh token'));
                    }
                } catch (err) {
                    processQueue(err);
                    reject(err);
                }
                finally {
                    isRefreshing = false;
                }
            });
        }
        return Promise.reject(error);
    }
)