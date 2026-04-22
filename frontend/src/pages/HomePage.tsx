import { use } from 'react'
import { useLogout, useMe } from '../features/auth/authHooks'
import { useNavigate } from 'react-router-dom'

export default function HomePage() {
    const { data, isLoading, isError } = useMe()
    const logoutMutation = useLogout()
    const navigation = useNavigate()

    if (isLoading) return <div>Loading...</div>
    if (isError) {
        // navigation('/login')
        return <div>Error loading user info. Please <a href="/login">login</a> again.</div> 
    }

    return (
        <div>
            <h1>Dashboard</h1>
            <p>ID: {data?.id}</p>
            <p>Name: {data?.name}</p>
            <p>Email: {data?.email}</p>

            <button onClick={() => logoutMutation.mutate()}>
                Logout
            </button>
        </div>
    )
}