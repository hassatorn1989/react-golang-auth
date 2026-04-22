
import { api } from './api'
import { refreshAccessToken } from './refreshToken'
import { authStore } from './authStore'
import type { AxiosError, InternalAxiosRequestConfig } from 'axios'

// interface RetryableRequest extends InternalAxiosRequestConfig {
//     _retry?: boolean;
// }

// let isRedirecting = false

// api.interceptors.request.use((config) => {
//     const token = authStore.getAccessToken()
//     if (token) {
//         config.headers.Authorization = `Bearer ${token}`
//     }
//     return config
// })




// api.interceptors.response.use(
//     (response) => response,
//     async (error) => {
//         const originalRequest = error.config as any;
//         if (error.response?.status === 401 && !originalRequest._retry) {
//             if (isRefreshing) {
//                 return new Promise((resolve, reject) => {
//                     failedQueue.push({ resolve, reject });
//                 }).then((token) => {
//                     originalRequest.headers.Authorization = `Bearer ${token}`;
//                     return api(originalRequest);
//                 }).catch((err) => {
//                     return Promise.reject(err);
//                 });
//             }

//             originalRequest._retry = true;
//             isRefreshing = true;
            
//             return new Promise(async (resolve, reject) => {
//                 try {
//                     const token = await refreshAccessToken();
//                     if (token) {
//                         originalRequest.headers.Authorization = `Bearer ${token}`;
//                         processQueue(null, token);
//                         resolve(api(originalRequest));
//                     } else {
//                         processQueue(new Error('Failed to refresh token'));
//                         reject(new Error('Failed to refresh token'));
//                     }
//                 } catch (err) {
//                     processQueue(err);
//                     reject(err);
//                 }
//                 finally {
//                     isRefreshing = false;
//                 }
//             });
//         }
//         return Promise.reject(error);
//     }
// )
// let isRefreshing = false;
// let failedQueue: Array<{
//     resolve: (token: string) => void;
//     reject: (error: unknown) => void;
// }> = [];

// function processQueue(error: unknown, token: string | null = null) {
//     failedQueue.forEach((promise) => {
//         if (error) {
//             promise.reject(error);
//         } else if (token) {
//             promise.resolve(token);
//         }
//     });

//     failedQueue = [];
// }

// // response interceptor
// api.interceptors.response.use(
//     (response) => response,
//     async (error: AxiosError) => {
//         const originalRequest = error.config as RetryableRequest | undefined;

//         if (!originalRequest) {
//             return Promise.reject(error);
//         }

//         const status = error.response?.status;

//         const isRefreshEndpoint =
//             originalRequest.url?.includes("/refresh-token");

//         // ถ้า 401 และไม่ใช่ request refresh เอง
//         if (status === 401 && !originalRequest._retry && !isRefreshEndpoint) {
//             if (isRefreshing) {
//                 return new Promise((resolve, reject) => {
//                     failedQueue.push({
//                         resolve: (token: string) => {
//                             if (originalRequest.headers) {
//                                 originalRequest.headers.Authorization = `Bearer ${token}`;
//                             }
//                             resolve(api(originalRequest));
//                         },
//                         reject,
//                     });
//                 });
//             }

//             originalRequest._retry = true;
//             isRefreshing = true;

//             try {
//                 const token = await refreshAccessToken();
//                 processQueue(null, token);
                
//                 if (originalRequest.headers) {
//                     originalRequest.headers.Authorization = `Bearer ${token}`;
//                 }

//                 return api(originalRequest);
//             } catch (refreshError) {
//                 processQueue(refreshError, null);

//                 window.location.href = "/login";
//                 return Promise.reject(refreshError);
//             } finally {
//                 isRefreshing = false;
//             }
//         }

//         return Promise.reject(error);
//     }
// );