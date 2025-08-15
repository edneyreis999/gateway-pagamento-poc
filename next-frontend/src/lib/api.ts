const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080"

type Account = {
  id: string
  name: string
}

type CreateInvoiceData = {
  amount: number
  description: string
  paymentMethod: {
    type: string
    lastFourDigits: string
    cardholderName: string
  }
}

type Invoice = {
  id: string
  amount: number
  description: string
  status: string
  createdAt: string
  updatedAt: string
  paymentMethod: {
    type: string
    lastFourDigits: string
    cardholderName: string
  }
  accountId: string
  clientIp: string
  device: string
}

interface InvoiceFilters {
  status?: string
  startDate?: string
  endDate?: string
  search?: string
}

type PaginationInfo = {
  currentPage: number
  totalPages: number
}

class ApiClient {
  private apiKey: string | null = null

  setApiKey(apiKey: string) {
    this.apiKey = apiKey
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`

    const headers: HeadersInit = {
      "Content-Type": "application/json",
      ...options.headers,
    }

    if (this.apiKey) {
      headers["X-API-Key"] = this.apiKey
    }

    const response = await fetch(url, {
      ...options,
      headers,
    })

    if (!response.ok) {
      throw new Error(`API Error: ${response.status} ${response.statusText}`)
    }

    return response.json()
  }

  async createAccount(data: { name: string }): Promise<Account> {
    return this.request<Account>("/accounts", {
      method: "POST",
      body: JSON.stringify(data),
    })
  }

  async getAccount(apiKey: string): Promise<Account> {
    return this.request<Account>("/accounts", { headers: { 'X-API-Key': apiKey } })
  }

  async createInvoice(data: CreateInvoiceData): Promise<Invoice> {
    return this.request<Invoice>("/invoice", {
      method: "POST",
      body: JSON.stringify(data),
    })
  }

  async getInvoices(filters?: InvoiceFilters & { page?: number }): Promise<{
    invoices: Invoice[]
    pagination: PaginationInfo
  }> {
    const params = new URLSearchParams()

    if (filters?.status && filters.status !== "all") {
      params.append("status", filters.status)
    }
    if (filters?.startDate) {
      params.append("start_date", filters.startDate)
    }
    if (filters?.endDate) {
      params.append("end_date", filters.endDate)
    }
    if (filters?.search) {
      params.append("search", filters.search)
    }
    if (filters?.page) {
      params.append("page", filters.page.toString())
    }

    const queryString = params.toString()
    const endpoint = `/invoice${queryString ? `?${queryString}` : ""}`

    return this.request(endpoint)
  }

  async getInvoice(id: string): Promise<Invoice> {
    // For demo purposes, return mock data
    // In real implementation, this would make an actual API call
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({
          id: `#INV-${id.padStart(3, "0")}`,
          amount: 150000, // R$ 1.500,00 in cents
          description: "Compra Online #123",
          status: "approved",
          createdAt: "2025-03-30T14:30:00Z",
          updatedAt: "2025-03-30T14:35:00Z",
          paymentMethod: {
            type: "credit_card",
            lastFourDigits: "1234",
            cardholderName: "João da Silva",
          },
          accountId: "ACC-12345",
          clientIp: "192.168.1.1",
          device: "Desktop - Chrome",
        } as Invoice)
      }, 1000)
    })
  }
}

export const apiClient = new ApiClient()
