export interface User {
  id: number;
  username: string;
  display_name: string;
  is_admin: boolean;
  must_change_password: boolean;
  created_at: string;
}

export interface Category {
  id: number;
  name: string;
  icon: string;
}

export type ExpenseStatus = 'pending' | 'approved' | 'rejected';

export interface Expense {
  id: number;
  created_by: number;
  creator: User;
  category_id: number;
  category: Category;
  description: string;
  amount: number;
  expense_date: string;
  expense_month: number;
  expense_year: number;
  is_shared: boolean;
  split_ratio: number;
  is_installment: boolean;
  installment_no: number | null;
  installment_total: number | null;
  recurring_expense_id: number | null;
  status: ExpenseStatus;
  approved_by: number | null;
  approver: User | null;
  approved_at: string | null;
  created_at: string;
}

export type RecurringType = 'installment' | 'recurring';

export interface RecurringExpense {
  id: number;
  created_by: number;
  creator: User;
  category_id: number;
  category: Category;
  description: string;
  amount: number;
  total_amount: number | null;
  type: RecurringType;
  installment_count: number | null;
  installments_remaining: number | null;
  is_shared: boolean;
  split_ratio: number;
  is_active: boolean;
  status: ExpenseStatus;
  approved_by: number | null;
  approver: User | null;
  created_at: string;
}

export interface Payment {
  id: number;
  month: number;
  year: number;
  payer_id: number;
  payer: User;
  payee_id: number;
  payee: User;
  amount: number;
  created_at: string;
}

export interface UserSummary {
  user_id: number;
  display_name: string;
  total_paid: number;
  total_share: number;
  balance: number;
}

export interface CategorySum {
  category_id: number;
  category_name: string;
  category_icon: string;
  total: number;
}

export interface MonthlySummary {
  month: number;
  year: number;
  total_expenses: number;
  shared_expenses: number;
  user_summaries: UserSummary[];
  debtor_id: number | null;
  creditor_id: number | null;
  debt_amount: number;
  total_payments: number;
  remaining_debt: number;
  category_breakdown: CategorySum[];
}
