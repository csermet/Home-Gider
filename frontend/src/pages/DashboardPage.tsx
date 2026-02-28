import { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { getSummary, getExpenses } from '../services/api';
import type { MonthlySummary, Expense } from '../types';
import { useAuth } from '../context/AuthContext';
import MonthSelector from '../components/MonthSelector';
import StatusBadge from '../components/StatusBadge';
import { TrendingUp, TrendingDown, Wallet, ArrowRight } from 'lucide-react';
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip, type PieLabelRenderProps } from 'recharts';

const COLORS = ['#3b82f6', '#ef4444', '#10b981', '#f59e0b', '#8b5cf6', '#ec4899', '#06b6d4', '#f97316', '#6b7280'];

export default function DashboardPage() {
  const { user } = useAuth();
  const now = new Date();
  const [searchParams, setSearchParams] = useSearchParams();
  const [month, setMonth] = useState(() => {
    const p = parseInt(searchParams.get('month') || '');
    return p >= 1 && p <= 12 ? p : now.getMonth() + 1;
  });
  const [year, setYear] = useState(() => {
    const p = parseInt(searchParams.get('year') || '');
    return p >= 2020 ? p : now.getFullYear();
  });
  const [summary, setSummary] = useState<MonthlySummary | null>(null);
  const [pendingExpenses, setPendingExpenses] = useState<Expense[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setSearchParams({ month: String(month), year: String(year) }, { replace: true });
    setLoading(true);
    Promise.all([
      getSummary(month, year),
      getExpenses(month, year),
    ]).then(([sumRes, expRes]) => {
      setSummary(sumRes.data);
      setPendingExpenses(expRes.data.filter((e) => e.status === 'pending'));
    }).finally(() => setLoading(false));
  }, [month, year]);

  const myBalance = summary?.user_summaries.find((u) => u.user_id === user?.id);
  const otherBalance = summary?.user_summaries.find((u) => u.user_id !== user?.id);

  const debtLabel = () => {
    if (!summary || summary.debt_amount === 0) return null;
    const debtor = summary.user_summaries.find((u) => u.user_id === summary.debtor_id);
    const creditor = summary.user_summaries.find((u) => u.user_id === summary.creditor_id);
    if (!debtor || !creditor) return null;
    return { debtor: debtor.display_name, creditor: creditor.display_name, amount: summary.debt_amount };
  };
  const debt = debtLabel();

  const chartData = summary?.category_breakdown.map((c, i) => ({
    name: c.category_name,
    value: c.total,
    color: COLORS[i % COLORS.length],
  })) || [];

  if (loading) {
    return <div className="flex justify-center py-20 text-slate-400">Yükleniyor...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-slate-800">Dashboard</h2>
        <MonthSelector month={month} year={year} onChange={(m, y) => { setMonth(m); setYear(y); }} />
      </div>

      {/* Özet Kartları */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-white rounded-xl p-5 shadow-sm border border-slate-200">
          <div className="flex items-center gap-3 mb-2">
            <div className="p-2 bg-blue-100 rounded-lg"><Wallet size={20} className="text-blue-600" /></div>
            <span className="text-sm text-slate-500">Toplam Harcama</span>
          </div>
          <p className="text-2xl font-bold text-slate-800">{summary?.total_expenses.toLocaleString('tr-TR')} TL</p>
          <p className="text-xs text-slate-400 mt-1">Ortak: {summary?.shared_expenses.toLocaleString('tr-TR')} TL</p>
        </div>

        <div className="bg-white rounded-xl p-5 shadow-sm border border-slate-200">
          <div className="flex items-center gap-3 mb-2">
            <div className="p-2 bg-green-100 rounded-lg"><TrendingUp size={20} className="text-green-600" /></div>
            <span className="text-sm text-slate-500">{myBalance?.display_name} Ödedi</span>
          </div>
          <p className="text-2xl font-bold text-slate-800">{myBalance?.total_paid.toLocaleString('tr-TR')} TL</p>
          <p className="text-xs text-slate-400 mt-1">Pay: {myBalance?.total_share.toLocaleString('tr-TR')} TL</p>
        </div>

        <div className="bg-white rounded-xl p-5 shadow-sm border border-slate-200">
          <div className="flex items-center gap-3 mb-2">
            <div className="p-2 bg-purple-100 rounded-lg"><TrendingDown size={20} className="text-purple-600" /></div>
            <span className="text-sm text-slate-500">{otherBalance?.display_name} Ödedi</span>
          </div>
          <p className="text-2xl font-bold text-slate-800">{otherBalance?.total_paid.toLocaleString('tr-TR')} TL</p>
          <p className="text-xs text-slate-400 mt-1">Pay: {otherBalance?.total_share.toLocaleString('tr-TR')} TL</p>
        </div>
      </div>

      {/* Borç Durumu */}
      {debt && (
        <div className={`bg-gradient-to-r ${summary!.remaining_debt === 0 ? 'from-green-500 to-green-600' : 'from-blue-500 to-blue-600'} rounded-xl p-6 text-white shadow-md`}>
          <div className="flex items-center justify-center gap-3 text-lg">
            <span className="font-semibold">{debt.debtor}</span>
            <ArrowRight size={24} />
            <span className="font-semibold">{debt.creditor}</span>
          </div>
          {summary!.remaining_debt > 0 ? (
            <>
              <p className="text-center text-3xl font-bold mt-2">{summary!.remaining_debt.toLocaleString('tr-TR')} TL</p>
              <p className="text-center text-blue-100 text-sm mt-1">
                kalan borç{summary!.total_payments > 0 ? ` (toplam ${debt.amount.toLocaleString('tr-TR')} TL, ${summary!.total_payments.toLocaleString('tr-TR')} TL ödendi)` : ''}
              </p>
            </>
          ) : (
            <>
              <p className="text-center text-3xl font-bold mt-2">Borç Ödendi!</p>
              <p className="text-center text-green-100 text-sm mt-1">{debt.amount.toLocaleString('tr-TR')} TL tamamen ödendi</p>
            </>
          )}
        </div>
      )}
      {summary && summary.debt_amount === 0 && summary.total_expenses > 0 && (
        <div className="bg-gradient-to-r from-green-500 to-green-600 rounded-xl p-6 text-white shadow-md text-center">
          <p className="text-xl font-bold">Eşit! Borç yok</p>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Kategori Dağılımı */}
        {chartData.length > 0 && (
          <div className="bg-white rounded-xl p-5 shadow-sm border border-slate-200">
            <h3 className="font-semibold text-slate-700 mb-4">Kategori Dağılımı</h3>
            <div className="h-64">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie data={chartData} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={80} label={(props: PieLabelRenderProps) => `${props.name ?? ''} %${((props.percent ?? 0) * 100).toFixed(0)}`}>
                    {chartData.map((entry, i) => (
                      <Cell key={i} fill={entry.color} />
                    ))}
                  </Pie>
                  <Tooltip formatter={(value) => `${Number(value).toLocaleString('tr-TR')} TL`} />
                </PieChart>
              </ResponsiveContainer>
            </div>
          </div>
        )}

        {/* Onay Bekleyenler */}
        <div className="bg-white rounded-xl p-5 shadow-sm border border-slate-200">
          <h3 className="font-semibold text-slate-700 mb-4">
            Onay Bekleyenler
            {pendingExpenses.length > 0 && (
              <span className="ml-2 bg-yellow-100 text-yellow-800 text-xs px-2 py-0.5 rounded-full">
                {pendingExpenses.length}
              </span>
            )}
          </h3>
          {pendingExpenses.length === 0 ? (
            <p className="text-slate-400 text-sm">Onay bekleyen gider yok</p>
          ) : (
            <div className="space-y-3">
              {pendingExpenses.slice(0, 5).map((e) => (
                <div key={e.id} className="flex items-center justify-between p-3 bg-slate-50 rounded-lg">
                  <div>
                    <p className="text-sm font-medium text-slate-700">{e.description}</p>
                    <p className="text-xs text-slate-400">{e.creator.display_name} - {e.category.name}</p>
                  </div>
                  <div className="text-right">
                    <p className="text-sm font-semibold text-slate-800">{e.amount.toLocaleString('tr-TR')} TL</p>
                    <StatusBadge status={e.status} />
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
