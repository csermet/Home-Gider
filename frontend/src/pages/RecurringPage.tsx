import { useEffect, useState } from 'react';
import { getRecurring, getCategories, createRecurring, deleteRecurring, approveRecurring, rejectRecurring } from '../services/api';
import type { RecurringExpense, Category } from '../types';
import { useAuth } from '../context/AuthContext';
import StatusBadge from '../components/StatusBadge';
import { Plus, Check, X, Trash2 } from 'lucide-react';

export default function RecurringPage() {
  const { user } = useAuth();
  const [items, setItems] = useState<RecurringExpense[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [loading, setLoading] = useState(true);

  const [form, setForm] = useState({
    category_id: 0,
    description: '',
    amount: '',
    total_amount: '',
    type: 'recurring' as 'recurring' | 'installment',
    installment_count: '',
    is_shared: true,
    split_ratio: '50',
  });

  const loadData = () => {
    setLoading(true);
    Promise.all([getRecurring(), getCategories()])
      .then(([recRes, catRes]) => {
        setItems(recRes.data);
        setCategories(catRes.data);
        if (catRes.data.length > 0 && form.category_id === 0) {
          setForm((f) => ({ ...f, category_id: catRes.data[0].id }));
        }
      })
      .finally(() => setLoading(false));
  };

  useEffect(() => { loadData(); }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const data: Parameters<typeof createRecurring>[0] = {
      category_id: form.category_id,
      description: form.description,
      amount: parseFloat(form.amount),
      type: form.type,
      is_shared: form.is_shared,
      split_ratio: parseFloat(form.split_ratio),
    };
    if (form.type === 'installment') {
      data.installment_count = parseInt(form.installment_count);
      if (form.total_amount) data.total_amount = parseFloat(form.total_amount);
    }
    await createRecurring(data);
    setForm({ category_id: categories[0]?.id || 0, description: '', amount: '', total_amount: '', type: 'recurring', installment_count: '', is_shared: true, split_ratio: '50' });
    setShowForm(false);
    loadData();
  };

  const canApprove = (item: RecurringExpense) =>
    item.status === 'pending' && (user?.is_admin || item.created_by !== user?.id);

  const handleApprove = async (id: number) => { await approveRecurring(id); loadData(); };
  const handleReject = async (id: number) => { await rejectRecurring(id); loadData(); };
  const handleDelete = async (id: number) => {
    if (!confirm('Bu şablonu deaktif etmek istediğinize emin misiniz?')) return;
    await deleteRecurring(id);
    loadData();
  };

  if (loading) {
    return <div className="flex justify-center py-20 text-slate-400">Yükleniyor...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-slate-800">Sabit & Taksitli Giderler</h2>
        <button
          onClick={() => setShowForm(!showForm)}
          className="flex items-center gap-1.5 px-4 py-2 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors cursor-pointer"
        >
          <Plus size={18} />
          Yeni Şablon
        </button>
      </div>

      {showForm && (
        <div className="bg-white rounded-xl p-6 shadow-sm border border-slate-200">
          <h3 className="font-semibold text-slate-700 mb-4">Yeni Şablon</h3>
          <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-slate-600 mb-1">Tip</label>
              <select
                value={form.type}
                onChange={(e) => setForm({ ...form, type: e.target.value as 'recurring' | 'installment' })}
                className="w-full px-3 py-2 rounded-lg border border-slate-300 focus:border-blue-500 outline-none"
              >
                <option value="recurring">Sabit Gider (Her ay tekrarlar)</option>
                <option value="installment">Taksitli Gider</option>
              </select>
            </div>
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
                placeholder="ör: İnternet faturası"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-600 mb-1">Aylık Tutar (TL)</label>
              <input
                type="number"
                step="0.01"
                value={form.amount}
                onChange={(e) => setForm({ ...form, amount: e.target.value })}
                className="w-full px-3 py-2 rounded-lg border border-slate-300 focus:border-blue-500 outline-none"
                required
              />
            </div>
            {form.type === 'installment' && (
              <>
                <div>
                  <label className="block text-sm font-medium text-slate-600 mb-1">Toplam Taksit Sayısı</label>
                  <input
                    type="number"
                    min="2"
                    value={form.installment_count}
                    onChange={(e) => setForm({ ...form, installment_count: e.target.value })}
                    className="w-full px-3 py-2 rounded-lg border border-slate-300 focus:border-blue-500 outline-none"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-600 mb-1">Toplam Tutar (TL, opsiyonel)</label>
                  <input
                    type="number"
                    step="0.01"
                    value={form.total_amount}
                    onChange={(e) => setForm({ ...form, total_amount: e.target.value })}
                    className="w-full px-3 py-2 rounded-lg border border-slate-300 focus:border-blue-500 outline-none"
                    placeholder="ör: 12000"
                  />
                </div>
              </>
            )}
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
              <button type="button" onClick={() => setShowForm(false)} className="px-4 py-2 text-slate-600 border border-slate-300 rounded-lg hover:bg-slate-50 cursor-pointer">
                İptal
              </button>
              <button type="submit" className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 cursor-pointer">
                Oluştur
              </button>
            </div>
          </form>
        </div>
      )}

      {items.length === 0 ? (
        <div className="bg-white rounded-xl p-12 shadow-sm border border-slate-200 text-center text-slate-400">
          Henüz sabit veya taksitli gider tanımlanmamış
        </div>
      ) : (
        <div className="space-y-3">
          {items.map((item) => (
            <div key={item.id} className={`bg-white rounded-xl p-4 shadow-sm border ${item.is_active ? 'border-slate-200' : 'border-slate-100 opacity-60'}`}>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-4">
                  <div className={`w-10 h-10 rounded-lg flex items-center justify-center text-sm font-bold ${
                    item.type === 'installment' ? 'bg-purple-100 text-purple-600' : 'bg-blue-100 text-blue-600'
                  }`}>
                    {item.type === 'installment' ? 'T' : 'S'}
                  </div>
                  <div>
                    <p className="font-medium text-slate-800">{item.description}</p>
                    <div className="flex items-center gap-2 text-xs text-slate-400 mt-0.5">
                      <span>{item.creator.display_name}</span>
                      <span>-</span>
                      <span>{item.category.name}</span>
                      {item.type === 'installment' && item.installment_count && item.installments_remaining !== null && (
                        <>
                          <span>-</span>
                          <span className="text-purple-500">
                            Kalan: {item.installments_remaining}/{item.installment_count} taksit
                          </span>
                        </>
                      )}
                      {!item.is_active && <span className="text-red-500 font-medium">Tamamlandı</span>}
                    </div>
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <div className="text-right mr-2">
                    <p className="font-semibold text-slate-800">{item.amount.toLocaleString('tr-TR')} TL/ay</p>
                    <StatusBadge status={item.status} />
                  </div>
                  {canApprove(item) && (
                    <div className="flex gap-1">
                      <button onClick={() => handleApprove(item.id)} className="p-2 bg-green-100 text-green-600 rounded-lg hover:bg-green-200 cursor-pointer" title="Onayla">
                        <Check size={16} />
                      </button>
                      <button onClick={() => handleReject(item.id)} className="p-2 bg-red-100 text-red-600 rounded-lg hover:bg-red-200 cursor-pointer" title="Reddet">
                        <X size={16} />
                      </button>
                    </div>
                  )}
                  {item.is_active && item.created_by === user?.id && (
                    <button onClick={() => handleDelete(item.id)} className="p-2 bg-slate-100 text-slate-500 rounded-lg hover:bg-red-100 hover:text-red-600 cursor-pointer" title="Deaktif Et">
                      <Trash2 size={16} />
                    </button>
                  )}
                </div>
              </div>
              {item.type === 'installment' && item.installment_count && item.installments_remaining !== null && (
                <div className="mt-3">
                  <div className="w-full bg-slate-100 rounded-full h-2">
                    <div
                      className="bg-purple-500 h-2 rounded-full transition-all"
                      style={{ width: `${((item.installment_count - item.installments_remaining) / item.installment_count) * 100}%` }}
                    />
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
