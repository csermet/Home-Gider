import { useEffect, useState } from 'react';
import { getUsers, adminResetPassword } from '../services/api';
import type { User } from '../types';
import { KeyRound, Shield } from 'lucide-react';

export default function AdminPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [resetUserId, setResetUserId] = useState<number | null>(null);
  const [newPassword, setNewPassword] = useState('');
  const [message, setMessage] = useState('');

  const loadUsers = () => {
    setLoading(true);
    getUsers().then((res) => setUsers(res.data)).finally(() => setLoading(false));
  };

  useEffect(() => { loadUsers(); }, []);

  const handleReset = async (userId: number) => {
    if (!newPassword || newPassword.length < 4) {
      setMessage('Şifre en az 4 karakter olmalı');
      return;
    }
    try {
      await adminResetPassword(userId, newPassword);
      setMessage('Şifre sıfırlandı! Kullanıcı bir sonraki girişte yeni şifre belirleyecek.');
      setResetUserId(null);
      setNewPassword('');
      loadUsers();
    } catch {
      setMessage('Şifre sıfırlanamadı');
    }
  };

  if (loading) {
    return <div className="flex justify-center py-20 text-slate-400">Yükleniyor...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Shield size={28} className="text-blue-600" />
        <h2 className="text-2xl font-bold text-slate-800">Admin Paneli</h2>
      </div>

      {message && (
        <div className="bg-blue-50 text-blue-700 p-3 rounded-lg text-sm">
          {message}
          <button onClick={() => setMessage('')} className="ml-2 underline cursor-pointer">Kapat</button>
        </div>
      )}

      <div className="bg-white rounded-xl shadow-sm border border-slate-200">
        <div className="p-5 border-b border-slate-100">
          <h3 className="font-semibold text-slate-700">Kullanıcı Yönetimi</h3>
        </div>
        <div className="divide-y divide-slate-100">
          {users.map((u) => (
            <div key={u.id} className="p-4 flex items-center justify-between">
              <div>
                <p className="font-medium text-slate-800">
                  {u.display_name}
                  {u.is_admin && <span className="ml-2 text-xs bg-blue-100 text-blue-700 px-2 py-0.5 rounded-full">Admin</span>}
                  {u.must_change_password && <span className="ml-2 text-xs bg-yellow-100 text-yellow-700 px-2 py-0.5 rounded-full">Şifre Değişmeli</span>}
                </p>
                <p className="text-sm text-slate-400">@{u.username}</p>
              </div>
              <div className="flex items-center gap-2">
                {!u.is_admin && (
                  <>
                    {resetUserId === u.id ? (
                      <div className="flex items-center gap-2">
                        <input
                          type="text"
                          value={newPassword}
                          onChange={(e) => setNewPassword(e.target.value)}
                          placeholder="Yeni geçici şifre"
                          className="px-3 py-1.5 rounded-lg border border-slate-300 text-sm focus:border-blue-500 outline-none"
                        />
                        <button
                          onClick={() => handleReset(u.id)}
                          className="px-3 py-1.5 bg-blue-600 text-white rounded-lg text-sm hover:bg-blue-700 cursor-pointer"
                        >
                          Sıfırla
                        </button>
                        <button
                          onClick={() => { setResetUserId(null); setNewPassword(''); }}
                          className="px-3 py-1.5 border border-slate-300 text-slate-600 rounded-lg text-sm hover:bg-slate-50 cursor-pointer"
                        >
                          İptal
                        </button>
                      </div>
                    ) : (
                      <button
                        onClick={() => setResetUserId(u.id)}
                        className="flex items-center gap-1 px-3 py-1.5 text-sm text-slate-600 border border-slate-300 rounded-lg hover:bg-slate-50 cursor-pointer"
                      >
                        <KeyRound size={14} />
                        Şifre Sıfırla
                      </button>
                    )}
                  </>
                )}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
