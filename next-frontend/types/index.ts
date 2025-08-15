export interface Invoice {
  id: string
  amount: number
  description: string
  status: "pending" | "approved" | "rejected"
  createdAt: string
  updatedAt: string
  paymentMethod: {
    type: "credit_card" | "debit_card" | "pix"
    lastFourDigits?: string
    cardholderName?: string
  }
  accountId: string
  clientIp?: string
  device?: string
}

export interface CreateInvoiceData {
  amount: number
  description: string
  paymentMethod: {
    type: "credit_card" | "debit_card" | "pix"
    cardNumber?: string
    cvv?: string
    expiryMonth?: string
    expiryYear?: string
    cardholderName?: string
  }
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
