import axios from 'axios';
import type {
  User,
  Category,
  Expense,
  RecurringExpense,
  Payment,
  MonthlySummary,
} from '../types';

const api = axios.create({
  baseURL: '/api',
  withCredentials: true,
});

// Auth
export const login = (username: string, password: string) =>
  api.post<{ user: User }>('/auth/login', { username, password });

export const logout = () => api.post('/auth/logout');

export const getMe = () => api.get<User>('/auth/me');

export const changePassword = (new_password: string) =>
  api.post<{ user: User }>('/auth/change-password', { new_password });

// Admin
export const getUsers = () => api.get<User[]>('/admin/users');

export const adminResetPassword = (userId: number, new_password: string) =>
  api.post(`/admin/users/${userId}/reset-password`, { new_password });

// Categories
export const getCategories = () => api.get<Category[]>('/categories');

export const createCategory = (data: { name: string; icon: string }) =>
  api.post<Category>('/categories', data);

// Expenses
export const getExpenses = (month: number, year: number) =>
  api.get<Expense[]>('/expenses', { params: { month, year } });

export const createExpense = (data: {
  category_id: number;
  description: string;
  amount: number;
  expense_date: string;
  is_shared: boolean;
  split_ratio: number;
}) => api.post<Expense>('/expenses', data);

export const updateExpense = (id: number, data: Record<string, unknown>) =>
  api.put<Expense>(`/expenses/${id}`, data);

export const deleteExpense = (id: number) =>
  api.delete(`/expenses/${id}`);

export const approveExpense = (id: number) =>
  api.post<Expense>(`/expenses/${id}/approve`);

export const rejectExpense = (id: number) =>
  api.post<Expense>(`/expenses/${id}/reject`);

export const confirmDeleteExpense = (id: number) =>
  api.post(`/expenses/${id}/confirm-delete`);

export const cancelDeleteExpense = (id: number) =>
  api.post<Expense>(`/expenses/${id}/cancel-delete`);

// Recurring
export const getRecurring = () =>
  api.get<RecurringExpense[]>('/recurring');

export const createRecurring = (data: {
  category_id: number;
  description: string;
  amount: number;
  total_amount?: number;
  type: string;
  installment_count?: number;
  is_shared: boolean;
  split_ratio: number;
}) => api.post<RecurringExpense>('/recurring', data);

export const deleteRecurring = (id: number) =>
  api.delete(`/recurring/${id}`);

export const approveRecurring = (id: number) =>
  api.post(`/recurring/${id}/approve`);

export const rejectRecurring = (id: number) =>
  api.post(`/recurring/${id}/reject`);

// Summary & Payments
export const getSummary = (month: number, year: number) =>
  api.get<MonthlySummary>('/summary', { params: { month, year } });

export const getPayments = (month: number, year: number) =>
  api.get<Payment[]>('/payments', { params: { month, year } });

export const addPayment = (data: {
  month: number;
  year: number;
  payer_id: number;
  payee_id: number;
  amount: number;
}) => api.post('/payments', data);

export const deletePayment = (id: number) =>
  api.delete(`/payments/${id}`);
