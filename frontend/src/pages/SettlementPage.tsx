import { useEffect, useState } from 'react';
import { getSummary, getPayments, addPayment, deletePayment } from '../services/api';
import type { MonthlySummary, Payment } from '../types';
import { useAuth } from '../context/AuthContext';
import MonthSelector from '../components/MonthSelector';
import { ArrowRight, CheckCircle, CircleDollarSign, Plus, Trash2 } from 'lucide-react';

export default function SettlementPage() {
  const { user } = useAuth();
  const now = new Date();
  const [month, setMonth] = useState(now.getMonth() + 1);
  const [year, setYear] = useState(now.getFullYear());
  const [summary, setSummary] = useState<MonthlySummary | null>(null);
  const [payments, setPayments] = useState<Payment[]>([]);
  const [loading, setLoading] = useState(true);
  const [showPaymentForm, setShowPaymentForm] = useState(false);
  const [paymentAmount, setPaymentAmount] = useState('');

  const loadData = () => {
    setLoading(true);
    Promise.all([
      getSummary(month, year),
      getPayments(month, year),
    ]).then(([sumRes, payRes]) => {
      setSummary(sumRes.data);
      setPayments(payRes.data);
    }).finally(() => setLoading(false));
  };

  useEffect(() => { loadData(); }, [month, year]);

  const handleAddPayment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!summary || !summary.debtor_id || !summary.creditor_id) return;
    const amount = parseFloat(paymentAmount);
    if (amount <= 0) return;

    await addPayment({
      month,
      year,
      payer_id: summary.debtor_id,
      payee_id: summary.creditor_id,
      amount,
    });
    setPaymentAmount('');
    setShowPaymentForm(false);
    loadData();
  };

  const handleDeletePayment = async (id: number) => {
    if (!confirm('Bu ödemeyi silmek istediğinize emin misiniz?')) return;
    await deletePayment(id);
    loadData();
  };

  const debtor = summary?.user_summaries.find((u) => u.user_id === summary?.debtor_id);
  const creditor = summary?.user_summaries.find((u) => u.user_id === summary?.creditor_id);

  if (loading) {
    return <div className="flex justify-center py-20 text-slate-400">Yükleniyor...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-slate-800">Hesaplaşma</h2>
        <MonthSelector month={month} year={year} onChange={(m, y) => { setMonth(m); setYear(y); }} />
      </div>

      {/* Kişi bazlı özet */}
      {summary && summary.user_summaries.length > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {summary.user_summaries.map((us) => (
            <div key={us.user_id} className="bg-white rounded-xl p-5 shadow-sm border border-slate-200">
              <h3 className="font-semibold text-slate-700 mb-3">{us.display_name}</h3>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-slate-500">Toplam Ödeme</span>
                  <span className="font-medium text-slate-800">{us.total_paid.toLocaleString('tr-TR')} TL</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-slate-500">Toplam Pay</span>
                  <span className="font-medium text-slate-800">{us.total_share.toLocaleString('tr-TR')} TL</span>
                </div>
                <hr className="border-slate-100" />
                <div className="flex justify-between">
                  <span className="text-slate-500">Bakiye</span>
                  <span className={`font-bold ${us.balance > 0 ? 'text-green-600' : us.balance < 0 ? 'text-red-600' : 'text-slate-600'}`}>
                    {us.balance > 0 ? '+' : ''}{us.balance.toLocaleString('tr-TR')} TL
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Borç durumu */}
      {summary && summary.debt_amount > 0 && debtor && creditor && (
        <div className="bg-gradient-to-r from-blue-500 to-blue-600 rounded-xl p-6 text-white shadow-md">
          <div className="flex items-center justify-center gap-3 text-lg">
            <span className="font-semibold">{debtor.display_name}</span>
            <ArrowRight size={24} />
            <span className="font-semibold">{creditor.display_name}</span>
          </div>
          <p className="text-center text-3xl font-bold mt-2">{summary.debt_amount.toLocaleString('tr-TR')} TL</p>
          <p className="text-center text-blue-100 text-sm mt-1">toplam borç</p>

          {/* Ödeme durumu */}
          {summary.total_payments > 0 && (
            <div className="mt-4 bg-white/15 rounded-lg p-3">
              <div className="flex justify-between text-sm">
                <span>Yapılan Ödemeler</span>
                <span className="font-semibold">{summary.total_payments.toLocaleString('tr-TR')} TL</span>
              </div>
              <div className="flex justify-between text-sm mt-1">
                <span>Kalan Borç</span>
                <span className="font-bold text-lg">{summary.remaining_debt.toLocaleString('tr-TR')} TL</span>
              </div>
              <div className="mt-2 w-full bg-white/20 rounded-full h-2">
                <div
                  className="bg-white h-2 rounded-full transition-all"
                  style={{ width: `${Math.min((summary.total_payments / summary.debt_amount) * 100, 100)}%` }}
                />
              </div>
            </div>
          )}

          {summary.remaining_debt > 0 && summary.debtor_id === user?.id && (
            <div className="text-center mt-4">
              <button
                onClick={() => setShowPaymentForm(!showPaymentForm)}
                className="px-6 py-2 bg-white text-blue-600 rounded-lg font-medium hover:bg-blue-50 transition-colors cursor-pointer inline-flex items-center gap-2"
              >
                <Plus size={18} />
                Ödeme Ekle
              </button>
            </div>
          )}
        </div>
      )}

      {/* Borç yok */}
      {summary && summary.debt_amount === 0 && summary.total_expenses > 0 && (
        <div className="bg-gradient-to-r from-green-500 to-green-600 rounded-xl p-6 text-white shadow-md text-center">
          <CheckCircle size={48} className="mx-auto mb-2" />
          <p className="text-xl font-bold">Bu ay eşit! Borç yok.</p>
        </div>
      )}

      {/* Borç tamamen ödendi */}
      {summary && summary.debt_amount > 0 && summary.remaining_debt === 0 && (
        <div className="bg-gradient-to-r from-green-500 to-green-600 rounded-xl p-6 text-white shadow-md text-center">
          <CheckCircle size={48} className="mx-auto mb-2" />
          <p className="text-xl font-bold">Borç tamamen ödendi!</p>
        </div>
      )}

      {summary && summary.total_expenses === 0 && (
        <div className="bg-white rounded-xl p-12 shadow-sm border border-slate-200 text-center text-slate-400">
          Bu ay için gider bulunamadı
        </div>
      )}

      {/* Ödeme Ekleme Formu */}
      {showPaymentForm && summary && summary.debtor_id && summary.creditor_id && (
        <div className="bg-white rounded-xl p-6 shadow-sm border border-slate-200">
          <h3 className="font-semibold text-slate-700 mb-4">Ödeme Ekle</h3>
          <form onSubmit={handleAddPayment} className="flex items-end gap-4">
            <div className="flex-1">
              <label className="block text-sm font-medium text-slate-600 mb-1">
                Tutar (TL) — Kalan: {summary.remaining_debt.toLocaleString('tr-TR')} TL
              </label>
              <input
                type="number"
                step="0.01"
                min="0.01"
                max={summary.remaining_debt}
                value={paymentAmount}
                onChange={(e) => setPaymentAmount(e.target.value)}
                className="w-full px-3 py-2 rounded-lg border border-slate-300 focus:border-blue-500 outline-none"
                placeholder="0.00"
                required
              />
            </div>
            <button
              type="button"
              onClick={() => setPaymentAmount(summary.remaining_debt.toString())}
              className="px-4 py-2 text-blue-600 border border-blue-300 rounded-lg hover:bg-blue-50 cursor-pointer text-sm"
            >
              Tamamını Öde
            </button>
            <button type="button" onClick={() => setShowPaymentForm(false)} className="px-4 py-2 text-slate-600 border border-slate-300 rounded-lg hover:bg-slate-50 cursor-pointer">
              İptal
            </button>
            <button type="submit" className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 cursor-pointer">
              Kaydet
            </button>
          </form>
        </div>
      )}

      {/* Ödeme Geçmişi */}
      {payments.length > 0 && (
        <div className="bg-white rounded-xl p-5 shadow-sm border border-slate-200">
          <h3 className="font-semibold text-slate-700 mb-4">Ödeme Geçmişi</h3>
          <div className="space-y-3">
            {payments.map((p) => (
              <div key={p.id} className="flex items-center justify-between p-4 bg-slate-50 rounded-lg">
                <div className="flex items-center gap-3">
                  <CircleDollarSign size={24} className="text-green-500" />
                  <div>
                    <p className="text-sm font-medium text-slate-700">
                      {p.payer.display_name} → {p.payee.display_name}
                    </p>
                    <p className="text-xs text-slate-400">
                      {new Date(p.created_at).toLocaleDateString('tr-TR')} — {new Date(p.created_at).toLocaleTimeString('tr-TR', { hour: '2-digit', minute: '2-digit' })}
                    </p>
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <span className="font-semibold text-green-600">{p.amount.toLocaleString('tr-TR')} TL</span>
                  {(user?.is_admin || p.payer_id === user?.id) && (
                    <button
                      onClick={() => handleDeletePayment(p.id)}
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
        </div>
      )}
    </div>
  );
}
