import { NavLink } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import {
    LayoutDashboard,
    Plug,
    Workflow,
    List,
    BarChart3,
    Network,
} from 'lucide-react';

const TABS = [
    { to: '/overview', key: 'nav.overview', Icon: LayoutDashboard },
    { to: '/connectors', key: 'nav.connectors', Icon: Plug },
    { to: '/compose', key: 'nav.compose', Icon: Workflow },
    { to: '/pipelines', key: 'nav.pipelines', Icon: List },
    { to: '/dashboards', key: 'nav.dashboards', Icon: BarChart3 },
    { to: '/kg', key: 'nav.kg', Icon: Network },
];

function LanguageToggle() {
    const { i18n } = useTranslation();
    const setLang = (lng) => {
        i18n.changeLanguage(lng);
        localStorage.setItem('mindset_lang', lng);
    };
    return (
        <div className="ml-auto flex items-center gap-0.5 shrink-0">
            {['fr', 'en'].map((lng) => (
                <button
                    key={lng}
                    onClick={() => setLang(lng)}
                    className={`px-2 py-1 rounded text-[11px] font-mono uppercase tracking-wider transition-colors ${
                        i18n.resolvedLanguage === lng
                            ? 'bg-elevated text-text-primary'
                            : 'text-text-tertiary hover:text-text-primary hover:bg-panel-alt'
                    }`}
                >
                    {lng}
                </button>
            ))}
        </div>
    );
}

export default function NavBar() {
    const { t } = useTranslation();
    return (
        <header className="bg-panel border-b border-border-subtle px-4 py-2 flex items-center gap-6 shrink-0">
            <img
                src="/logo.png"
                alt="MindSet Data"
                className="h-7 w-auto shrink-0"
            />
            <nav className="flex gap-0.5">
                {TABS.map(({ to, key, Icon }) => (
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
                        <span>{t(key)}</span>
                    </NavLink>
                ))}
            </nav>
            <LanguageToggle />
        </header>
    );
}
