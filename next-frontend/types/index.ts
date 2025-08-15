export interface Invoice {
  account_id: string
  amount: number
  created_at: string
  description: string
  id: string
  payment_type: "credit_card" | "debit_card" | "pix"
  status: "pending" | "approved" | "rejected"
  updated_at: string
}

export interface CreateInvoiceData {
  amount: number
  description: string
  payment_type: "credit_card" | "debit_card" | "pix"
  card_number: string
  cvv: string
  expiry_month: number
  expiry_year: number
  cardholder_name: string
}

export interface Account {
  id: string
  apiKey: string
  name: string
  createdAt: string
}

export interface AuthContextType {
  isAuthenticated: boolean
  apiKey: string | null
  login: (apiKey: string) => Promise<boolean>
  logout: () => void
}

export interface InvoiceFilters {
  status?: "all" | "pending" | "approved" | "rejected"
  startDate?: string
  endDate?: string
  search?: string
}

export interface PaginationInfo {
  currentPage: number
  totalPages: number
  totalItems: number
  itemsPerPage: number
}
