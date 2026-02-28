import { useEffect, useState } from 'react';
import { getExpenses, getCategories, createExpense, approveExpense, rejectExpense, deleteExpense } from '../services/api';
import type { Expense, Category } from '../types';
import { useAuth } from '../context/AuthContext';
import MonthSelector from '../components/MonthSelector';
import StatusBadge from '../components/StatusBadge';
import { Plus, Check, X, Trash2 } from 'lucide-react';

export default function ExpensesPage() {
  const { user } = useAuth();
  const now = new Date();
  const [month, setMonth] = useState(now.getMonth() + 1);
  const [year, setYear] = useState(now.getFullYear());
  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [loading, setLoading] = useState(true);
  const [filterStatus, setFilterStatus] = useState<string>('all');

  // Form state
  const [form, setForm] = useState({
    category_id: 0,
    description: '',
    amount: '',
    expense_date: new Date().toISOString().split('T')[0],
    is_shared: true,
    split_ratio: '50',
  });

  const loadData = () => {
    setLoading(true);
    Promise.all([
      getExpenses(month, year),
      getCategories(),
    ]).then(([expRes, catRes]) => {
      setExpenses(expRes.data);
      setCategories(catRes.data);
      if (catRes.data.length > 0 && form.category_id === 0) {
        setForm((f) => ({ ...f, category_id: catRes.data[0].id }));
      }
    }).finally(() => setLoading(false));
  };

  useEffect(() => { loadData(); }, [month, year]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await createExpense({
      category_id: form.category_id,
      description: form.description,
      amount: parseFloat(form.amount),
      expense_date: form.expense_date,
      is_shared: form.is_shared,
      split_ratio: parseFloat(form.split_ratio),
    });
    setForm({ category_id: categories[0]?.id || 0, description: '', amount: '', expense_date: new Date().toISOString().split('T')[0], is_shared: true, split_ratio: '50' });
    setShowForm(false);
    loadData();
  };

  const handleApprove = async (id: number) => {
    await approveExpense(id);
    loadData();
  };

  const handleReject = async (id: number) => {
    await rejectExpense(id);
    loadData();
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Bu gideri silmek istediğinize emin misiniz?')) return;
    await deleteExpense(id);
    loadData();
  };

  const filtered = filterStatus === 'all'
    ? expenses
    : expenses.filter((e) => e.status === filterStatus);

  if (loading) {
    return <div className="flex justify-center py-20 text-slate-400">Yükleniyor...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <h2 className="text-2xl font-bold text-slate-800">Giderler</h2>
        <div className="flex items-center gap-3">
          <MonthSelector month={month} year={year} onChange={(m, y) => { setMonth(m); setYear(y); }} />
          <button
            onClick={() => setShowForm(!showForm)}
            className="flex items-center gap-1.5 px-4 py-2 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors cursor-pointer"
          >
            <Plus size={18} />
            Gider Ekle
          </button>
        </div>
      </div>

      {/* Gider Ekleme Formu */}
      {showForm && (
        <div className="bg-white rounded-xl p-6 shadow-sm border border-slate-200">
          <h3 className="font-semibold text-slate-700 mb-4">Yeni Gider</h3>
          <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-slate-600 mb-1">Kategori</label>
              <select
                value={form.category_id}
                onChange={(e) => setForm({ ...form, category_id: parseInt(e.target.value) })}
                className="w-full px-3 py-2 rounded-lg border border-slate-300 focus:border-blue-500 outline-none"
              >
                {categories.map((c) => (
                  <option key={c.id} value={c.id}>{c.name}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-600 mb-1">Açıklama</label>
              <input
                type="text"
                value={form.description}
                onChange={(e) => setForm({ ...form, description: e.target.value })}
                className="w-full px-3 py-2 rounded-lg border border-slate-300 focus:border-blue-500 outline-none"
                placeholder="Gider açıklaması"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-600 mb-1">Tutar (TL)</label>
              <input
                type="number"
                step="0.01"
                value={form.amount}
                onChange={(e) => setForm({ ...form, amount: e.target.value })}
                className="w-full px-3 py-2 rounded-lg border border-slate-300 focus:border-blue-500 outline-none"
                placeholder="0.00"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-600 mb-1">Tarih</label>
              <input
                type="date"
                value={form.expense_date}
                onChange={(e) => setForm({ ...form, expense_date: e.target.value })}
                className="w-full px-3 py-2 rounded-lg border border-slate-300 focus:border-blue-500 outline-none"
                required
              />
            </div>
            <div className="flex items-center gap-4">
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="checkbox"
                  checked={form.is_shared}
                  onChange={(e) => setForm({ ...form, is_shared: e.target.checked })}
                  className="w-4 h-4 rounded border-slate-300"
                />
                <span className="text-sm text-slate-600">Ortak gider</span>
              </label>
            </div>
            {form.is_shared && (
              <div>
                <label className="block text-sm font-medium text-slate-600 mb-1">Bölüşme Oranı (Sizin payınız %)</label>
                <input
                  type="number"
                  min="1"
                  max="99"
                  value={form.split_ratio}
                  onChange={(e) => setForm({ ...form, split_ratio: e.target.value })}
                  className="w-full px-3 py-2 rounded-lg border border-slate-300 focus:border-blue-500 outline-none"
                />
              </div>
            )}
            <div className="md:col-span-2 flex gap-3 justify-end">
              <button
                type="button"
                onClick={() => setShowForm(false)}
                className="px-4 py-2 text-slate-600 border border-slate-300 rounded-lg hover:bg-slate-50 cursor-pointer"
              >
                İptal
              </button>
              <button
                type="submit"
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 cursor-pointer"
              >
                Ekle
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Filtreler */}
      <div className="flex gap-2">
        {[
          { value: 'all', label: 'Tümü' },
          { value: 'pending', label: 'Bekleyen' },
          { value: 'approved', label: 'Onaylı' },
          { value: 'rejected', label: 'Reddedilen' },
        ].map((f) => (
          <button
            key={f.value}
            onClick={() => setFilterStatus(f.value)}
            className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors cursor-pointer ${
              filterStatus === f.value
                ? 'bg-blue-100 text-blue-700'
                : 'bg-slate-100 text-slate-500 hover:bg-slate-200'
            }`}
          >
            {f.label}
          </button>
        ))}
      </div>

      {/* Gider Listesi */}
      {filtered.length === 0 ? (
        <div className="bg-white rounded-xl p-12 shadow-sm border border-slate-200 text-center text-slate-400">
          Bu ay için gider bulunamadı
        </div>
      ) : (
        <div className="space-y-3">
          {filtered.map((expense) => (
            <div key={expense.id} className="bg-white rounded-xl p-4 shadow-sm border border-slate-200 flex items-center justify-between">
              <div className="flex items-center gap-4">
                <div className="w-10 h-10 bg-slate-100 rounded-lg flex items-center justify-center text-slate-500 text-sm font-medium">
                  {expense.category.name.charAt(0)}
                </div>
                <div>
                  <p className="font-medium text-slate-800">{expense.description}</p>
                  <div className="flex items-center gap-2 text-xs text-slate-400 mt-0.5">
                    <span>{expense.creator.display_name}</span>
                    <span>-</span>
                    <span>{expense.category.name}</span>
                    <span>-</span>
                    <span>{new Date(expense.expense_date).toLocaleDateString('tr-TR')}</span>
                    {!expense.is_shared && <span className="text-orange-500 font-medium">Kişisel</span>}
                    {expense.is_shared && expense.split_ratio !== 50 && (
                      <span className="text-blue-500">%{expense.split_ratio}/%{100 - expense.split_ratio}</span>
                    )}
                    {expense.is_installment && expense.installment_no && expense.installment_total && (
                      <span className="text-purple-500">Taksit {expense.installment_no}/{expense.installment_total}</span>
                    )}
                  </div>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <div className="text-right mr-2">
                  <p className="font-semibold text-slate-800">{expense.amount.toLocaleString('tr-TR')} TL</p>
                  <StatusBadge status={expense.status} />
                </div>
                {expense.status === 'pending' && (user?.is_admin || expense.created_by !== user?.id) && (
                  <div className="flex gap-1">
                    <button
                      onClick={() => handleApprove(expense.id)}
                      className="p-2 bg-green-100 text-green-600 rounded-lg hover:bg-green-200 cursor-pointer"
                      title="Onayla"
                    >
                      <Check size={16} />
                    </button>
                    <button
                      onClick={() => handleReject(expense.id)}
                      className="p-2 bg-red-100 text-red-600 rounded-lg hover:bg-red-200 cursor-pointer"
                      title="Reddet"
                    >
                      <X size={16} />
                    </button>
                  </div>
                )}
                {expense.status === 'pending' && expense.created_by === user?.id && (
                  <button
                    onClick={() => handleDelete(expense.id)}
                    className="p-2 bg-slate-100 text-slate-500 rounded-lg hover:bg-red-100 hover:text-red-600 cursor-pointer"
                    title="Sil"
                  >
                    <Trash2 size={16} />
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
