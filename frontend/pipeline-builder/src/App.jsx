import { Routes, Route, Navigate, useLocation } from 'react-router-dom';
import NavBar from './components/NavBar';
import ErrorBoundary from './components/ErrorBoundary';
import OverviewPage from './pages/OverviewPage';
import ConnectPage from './pages/ConnectPage';
import OpcuaConnectPage from './pages/OpcuaConnectPage';
import BuilderPage from './pages/BuilderPage';
import PipelinesPage from './pages/PipelinesPage';
import DashboardPage from './pages/DashboardPage';
import KnowledgeGraphPage from './pages/KnowledgeGraphPage';
import './App.css';

export default function App() {
  const location = useLocation();
  return (
    <div className="h-screen flex flex-col bg-dark-950 text-white">
      <NavBar />
      <div className="flex-1 min-h-0">
        <ErrorBoundary resetKey={location.pathname}>
          <Routes>
            <Route path="/" element={<Navigate to="/overview" replace />} />
            <Route path="/overview" element={<OverviewPage />} />
            <Route path="/connect" element={<ConnectPage />} />
            <Route path="/connect/opcua" element={<OpcuaConnectPage />} />
            <Route path="/compose" element={<BuilderPage />} />
            <Route path="/pipelines" element={<PipelinesPage />} />
            <Route path="/dashboards" element={<DashboardPage />} />
            <Route path="/kg" element={<KnowledgeGraphPage />} />
          </Routes>
        </ErrorBoundary>
      </div>
    </div>
  );
}
