import { Routes, Route, Navigate } from 'react-router-dom'
import LoginPage from './pages/LoginPage'
import HomePage from './pages/HomePage'
import { useMe } from './features/auth/authHooks'
import { authStore } from './lib/authStore'

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const token = authStore.getAccessToken()
  
  const { isPending, data, isError } = useMe()

  if (!token) {
    return <Navigate to="/login" replace />
  }

  if (isPending) return <div>Checking auth...</div>

  if (isError || !data) return <Navigate to="/login" replace />
  return <>{children}</>
}

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route
        path="/"
        element={
          <ProtectedRoute>
            <HomePage />
          </ProtectedRoute>
        }
      />
    </Routes>
  )
}