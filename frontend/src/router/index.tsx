import { createBrowserRouter, Navigate } from 'react-router-dom'
import Layout from '@/pages/Layout'
import Login from '@/pages/Login'
import Dashboard from '@/pages/Dashboard'
import IncidentManage from '@/pages/IncidentManage'
import InspectionManage from '@/pages/InspectionManage'
import TrainingManage from '@/pages/TrainingManage'
import CertReview from '@/pages/CertReview'
import Profile from '@/pages/Profile'
import AuditLogs from '@/pages/AuditLogs'
import { RequireAuth, RequireRole } from './guards'

const router = createBrowserRouter([
  { path: '/login', element: <Login /> },
  {
    path: '/',
    element: (
      <RequireAuth>
        <Layout />
      </RequireAuth>
    ),
    children: [
      { index: true, element: <Navigate to="/dashboard" replace /> },
      { path: 'dashboard', element: <Dashboard /> },
      { path: 'incidents', element: <IncidentManage /> },
      { path: 'inspections', element: <InspectionManage /> },
      { path: 'trainings', element: <TrainingManage /> },
      { path: 'certifications', element: <CertReview /> },
      { path: 'profile', element: <Profile /> },
      {
        path: 'audit-logs',
        element: (
          <RequireRole roles={['admin']}>
            <AuditLogs />
          </RequireRole>
        ),
      },
    ],
  },
])

export default router
