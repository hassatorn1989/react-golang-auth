
export const authStore = {
    getAccessToken: () => {
        return localStorage.getItem('accessToken')
    },
    setAccessToken: (token: string | null) => {
        if (token) {
            localStorage.setItem('accessToken', token)
        } else {
            localStorage.removeItem('accessToken')
        }
    },
}