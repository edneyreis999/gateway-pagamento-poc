"use server"

import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import { apiClient } from "@/lib/api"
import { logger } from "@/lib/logger"
import type { Invoice, InvoiceFilters } from "../../../types"

export async function getInvoicesAction(filters?: InvoiceFilters & { page?: number }) {
  const cookiesStore = await cookies()
  const apiKey = cookiesStore.get("apiKey")?.value

  logger.api('getInvoicesAction called', { 
    hasApiKey: !!apiKey,
    hasFilters: !!filters 
  });

  if (!apiKey) {
    logger.warn('No API key found, redirecting to auth');
    redirect("/auth")
  }

  try {
    // Configurar a API key no cliente
    apiClient.setApiKey(apiKey)
    
    // Buscar as invoices
    logger.api('Fetching invoices from API');
    const response = await apiClient.getInvoices(filters)

    logger.api('Invoices fetched successfully', { 
      count: response.length 
    });
    
    // Fazer cast dos tipos para serem compatíveis
    const invoices: Invoice[] = response.map(invoice => ({
      ...invoice,
      status: invoice.status as "pending" | "approved" | "rejected",
      payment_type: invoice.payment_type as "credit_card" | "debit_card" | "pix"
    }))
    
    logger.info('getInvoicesAction completed successfully', {
      invoiceCount: invoices.length
    });

    return {
      invoices,
      error: null
    }
  } catch (error) {
    logger.error('Error fetching invoices', error);
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
