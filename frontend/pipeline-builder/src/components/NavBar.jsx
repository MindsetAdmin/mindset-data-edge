import { NavLink } from 'react-router-dom';
import {
    LayoutDashboard,
    Plug,
    Workflow,
    List,
    BarChart3,
    Network,
} from 'lucide-react';

const TABS = [
    { to: '/overview', label: 'Overview', Icon: LayoutDashboard },
    { to: '/connect', label: 'Connect', Icon: Plug },
    { to: '/compose', label: 'Compose', Icon: Workflow },
    { to: '/pipelines', label: 'Pipelines', Icon: List },
    { to: '/dashboards', label: 'Dashboards', Icon: BarChart3 },
    { to: '/kg', label: 'Knowledge Graph', Icon: Network },
];

export default function NavBar() {
    return (
        <header className="bg-panel border-b border-border-subtle px-4 py-2 flex items-center gap-6 shrink-0">
            <img
                src="/logo.png"
                alt="MindSet Data"
                className="h-7 w-auto shrink-0"
            />
            <nav className="flex gap-0.5">
                {TABS.map(({ to, label, Icon }) => (
                    <NavLink
                        key={to}
                        to={to}
                        className={({ isActive }) =>
                            `inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded text-13 transition-colors ${
                                isActive
                                    ? 'bg-elevated text-text-primary'
                                    : 'text-text-tertiary hover:text-text-primary hover:bg-panel-alt'
                            }`
                        }
                    >
                        <Icon size={14} strokeWidth={1.5} />
                        <span>{label}</span>
                    </NavLink>
                ))}
            </nav>
        </header>
    );
}
