import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { getMe, login, logout, type LoginPayload } from './authApi'
import { authStore } from '../../lib/authStore'
import { useNavigate } from 'react-router-dom'

export function useMe() {
    return useQuery({
        queryKey: ['me'],
        queryFn: getMe,
        retry: false,
    })
}

export function useLogin() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (payload: LoginPayload) => login(payload),
        onSuccess: (data) => {
            queryClient.setQueryData(['me'], data.user)
            authStore.setAccessToken(data.accessToken)
        },
        onError: () => {
            authStore.setAccessToken(null)
        }
    })
}

export function useLogout() {
    const navigation = useNavigate()
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: logout,
        onSuccess: () => {
            queryClient.removeQueries({ queryKey: ['me'] })
            authStore.setAccessToken(null)
            navigation('/login')
        },
    })
}