import { NavLink } from 'react-router-dom';

const TABS = [
  ['/overview', '🏠 Overview'],
  ['/connect', '🔌 Connect'],
  ['/compose', '⚙️ Compose'],
  ['/pipelines', '📡 Pipelines'],
  ['/dashboards', '📊 Dashboards'],
  ['/kg', '🧠 KG'],
];

export default function NavBar() {
  return (
    <header className="bg-dark-900 border-b border-dark-700 px-4 py-2.5 flex items-center gap-6 shrink-0">
      <div className="font-bold text-lg bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent whitespace-nowrap">
        🧩 MindSet&nbsp;Data
      </div>
      <nav className="flex gap-1">
        {TABS.map(([to, label]) => (
          <NavLink
            key={to}
            to={to}
            className={({ isActive }) =>
              `px-3 py-1.5 rounded-md text-sm transition ${
                isActive ? 'bg-dark-700 text-white' : 'text-dark-400 hover:text-white hover:bg-dark-800'
              }`
            }
          >
            {label}
          </NavLink>
        ))}
      </nav>
    </header>
  );
}
