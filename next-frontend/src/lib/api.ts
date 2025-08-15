const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080"

type Account = {
  id: string
  name: string
}

import type { CreateInvoiceData, Invoice, InvoiceFilters } from "../../types"
import { logger } from "./logger"

class ApiClient {
  private apiKey: string | null = null

  setApiKey(apiKey: string) {
    this.apiKey = apiKey
    logger.api('API key set', { hasApiKey: !!apiKey });
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`

    const headers: HeadersInit = {
      "Content-Type": "application/json",
      ...options.headers,
    }

    if (this.apiKey) {
      (headers as Record<string, string>)["X-API-Key"] = this.apiKey
    }

    logger.api('API request initiated', {
      url,
      method: options.method || 'GET',
      hasApiKey: !!this.apiKey
    });

    const response = await fetch(url, {
      ...options,
      headers,
    })

    logger.api('API response received', {
      url: response.url,
      status: response.status,
      ok: response.ok
    });

    if (!response.ok) {
      logger.error('API request failed', {
        status: response.status,
        statusText: response.statusText,
        url
      });
      throw new Error(`API Error: ${response.status} ${response.statusText}`)
    }

    const data = await response.json();
    logger.debug('API response data', { 
      dataLength: Array.isArray(data) ? data.length : 'single object' 
    });

    return data;
  }

  async createAccount(data: { name: string }): Promise<Account> {
    logger.api('Creating account', { name: data.name });
    return this.request<Account>("/accounts", {
      method: "POST",
      body: JSON.stringify(data),
    })
  }

  async getAccount(apiKey: string): Promise<Account> {
    logger.api('Validating API key with account endpoint');
    const headers: Record<string, string> = { "X-API-Key": apiKey }
    return this.request<Account>("/accounts", { headers })
  }

  async createInvoice(data: CreateInvoiceData): Promise<Invoice> {
    logger.api('Creating invoice', { 
      amount: data.amount,
      paymentType: data.payment_type 
    });
    return this.request<Invoice>("/invoices", {
      method: "POST",
      body: JSON.stringify(data),
    })
  }

  async getInvoices(filters?: InvoiceFilters & { page?: number }): Promise<Invoice[]> {
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
    const endpoint = `/invoices${queryString ? `?${queryString}` : ""}`

    logger.api('Fetching invoices', { 
      hasFilters: !!filters,
      queryString 
    });

    return this.request(endpoint)
  }

  async getInvoice(id: string): Promise<Invoice> {
    logger.api('Getting invoice by ID', { id });
    // For demo purposes, return mock data
    // In real implementation, this would make an actual API call
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({
          id: 'ab387ffc-b52e-45fb-9e66-dd76652180dc',
          account_id: '6d75cfdf-f1aa-453d-aca3-7083b4d91cbd',
          amount: 100.5,
          status: 'pending',
          description: 'Test invoice for payment gateway',
          payment_type: 'credit_card',
          card_last_digits: '1234',
          created_at: '2025-08-13T02:38:48.112219Z',
          updated_at: '2025-08-13T02:38:48.112219Z'
        } as Invoice)
      }, 1000)
    })
  }
}

export const apiClient = new ApiClient()
