import { NavLink, Outlet } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { LayoutDashboard, BarChart3, Receipt, Repeat, Handshake, LogOut, Shield } from 'lucide-react';

const navItems = [
  { to: '/', label: 'Dashboard', icon: LayoutDashboard },
  { to: '/overview', label: 'Genel Bakış', icon: BarChart3 },
  { to: '/expenses', label: 'Giderler', icon: Receipt },
  { to: '/recurring', label: 'Sabit/Taksit', icon: Repeat },
  { to: '/settlement', label: 'Hesaplaşma', icon: Handshake },
];

export default function Layout() {
  const { user, logoutUser } = useAuth();

  return (
    <div className="min-h-screen bg-slate-50">
      <header className="bg-white border-b border-slate-200 shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <h1 className="text-xl font-bold text-slate-800">Ev Giderleri</h1>
            <div className="flex items-center gap-4">
              <span className="text-sm text-slate-600">
                Merhaba, <strong>{user?.display_name}</strong>
                {user?.is_admin && <span className="ml-1 text-xs text-blue-600">(Admin)</span>}
              </span>
              <button
                onClick={logoutUser}
                className="flex items-center gap-1 text-sm text-slate-500 hover:text-red-600 transition-colors cursor-pointer"
              >
                <LogOut size={16} />
                Çıkış
              </button>
            </div>
          </div>
        </div>
      </header>

      <nav className="bg-white border-b border-slate-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex gap-1">
            {navItems.map((item) => (
              <NavLink
                key={item.to}
                to={item.to}
                end={item.to === '/'}
                className={({ isActive }) =>
                  `flex items-center gap-2 px-4 py-3 text-sm font-medium border-b-2 transition-colors ${
                    isActive
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'
                  }`
                }
              >
                <item.icon size={18} />
                {item.label}
              </NavLink>
            ))}
            {user?.is_admin && (
              <NavLink
                to="/admin"
                className={({ isActive }) =>
                  `flex items-center gap-2 px-4 py-3 text-sm font-medium border-b-2 transition-colors ${
                    isActive
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'
                  }`
                }
              >
                <Shield size={18} />
                Admin
              </NavLink>
            )}
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <Outlet />
      </main>
    </div>
  );
}
