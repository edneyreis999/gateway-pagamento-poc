"use server"

import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import { apiClient } from "@/lib/api"
import type { Invoice, InvoiceFilters } from "../../../types"

export async function getInvoicesAction(filters?: InvoiceFilters & { page?: number }) {
  const cookiesStore = await cookies()
  const apiKey = cookiesStore.get("apiKey")?.value

  console.log('----- apiKey -----')
  console.log(apiKey)

  if (!apiKey) {
    redirect("/auth")
  }

  try {
    // Configurar a API key no cliente
    apiClient.setApiKey(apiKey)
    
    // Buscar as invoices
    const response = await apiClient.getInvoices(filters)

    console.log('----- response ------')
    console.log(response)
    
    // Fazer cast dos tipos para serem compatíveis
    const invoices: Invoice[] = response.map(invoice => ({
      ...invoice,
      status: invoice.status as "pending" | "approved" | "rejected",
      payment_type: invoice.payment_type as "credit_card" | "debit_card" | "pix"
    }))
    
    return {
      invoices,
      error: null
    }
  } catch (error) {
    console.error("Error fetching invoices:", error)
    return {
      invoices: [],
      pagination: {
        currentPage: 1,
        totalPages: 1,
        totalItems: 0,
        itemsPerPage: 10,
      },
      error: "Erro ao carregar faturas. Tente novamente."
    }
  }
}
