import { useLogout, useMe } from '../features/auth/authHooks'

export default function HomePage() {
    const { data, isLoading, isError } = useMe()
    const logoutMutation = useLogout()

    if (isLoading) return <div>Loading...</div>
    if (isError) return <div>Unauthorized</div>

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