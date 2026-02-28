import type { ExpenseStatus } from '../types';

const config: Record<ExpenseStatus, { label: string; classes: string }> = {
  pending: { label: 'Onay Bekliyor', classes: 'bg-yellow-100 text-yellow-800' },
  approved: { label: 'Onaylandı', classes: 'bg-green-100 text-green-800' },
  rejected: { label: 'Reddedildi', classes: 'bg-red-100 text-red-800' },
};

export default function StatusBadge({ status }: { status: ExpenseStatus }) {
  const { label, classes } = config[status];
  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${classes}`}>
      {label}
    </span>
  );
}
