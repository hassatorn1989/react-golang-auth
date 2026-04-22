import { Routes, Route, Navigate } from 'react-router-dom'
import LoginPage from './pages/LoginPage'
import HomePage from './pages/HomePage'
import { useMe } from './features/auth/authHooks'

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isLoading, data, isError } = useMe()

  if (isLoading) return <div>Checking auth...</div>
  // if (isError || !data) return <Navigate to="/login" replace />
if (isError || !data) return <div>Error loading user info. Please <a href="/login">login</a> again.</div>
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