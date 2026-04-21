import { useState } from 'react'
import { useLogin } from '../features/auth/authHooks'
import { useNavigate } from 'react-router-dom'

export default function LoginPage() {
    const navigate = useNavigate()
    const loginMutation = useLogin()

    const [email, setEmail] = useState('admin@example.com')
    const [password, setPassword] = useState('123456')

    const onSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        try {
            await loginMutation.mutateAsync({ email, password })
            navigate('/')
        } catch {
            alert('login failed')
        }
    }

    return (
        <div style={{ maxWidth: 400, margin: '40px auto' }}>
            <h1>Login</h1>
            <form onSubmit={onSubmit}>
                <div>
                    <label>Email</label>
                    <input value={email} onChange={(e) => setEmail(e.target.value)} />
                </div>

                <div>
                    <label>Password</label>
                    <input
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                </div>

                <button type="submit" disabled={loginMutation.isPending}>
                    {loginMutation.isPending ? 'Logging in...' : 'Login'}
                </button>
            </form>
        </div>
    )
}