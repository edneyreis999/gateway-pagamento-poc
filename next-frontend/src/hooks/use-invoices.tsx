"use client"

import { useState, useEffect } from "react"
import { apiClient } from "@/lib/api"
import type { Invoice, InvoiceFilters, PaginationInfo } from "@/types"

export function useInvoices() {
  const [invoices, setInvoices] = useState<Invoice[]>([])
  const [pagination, setPagination] = useState<PaginationInfo>({
    currentPage: 1,
    totalPages: 1,
    totalItems: 0,
    itemsPerPage: 10,
  })
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [filters, setFilters] = useState<InvoiceFilters>({
    status: "all",
    startDate: "",
    endDate: "",
    search: "",
  })

  const fetchInvoices = async (newFilters?: InvoiceFilters, page = 1) => {
    try {
      setIsLoading(true)
      setError(null)

      const filtersToUse = newFilters || filters
      const response = await apiClient.getInvoices({
        ...filtersToUse,
        page,
      })

      setInvoices(response.invoices)
      setPagination(response.pagination)
    } catch (err) {
      setError("Erro ao carregar faturas. Tente novamente.")
      console.error("Error fetching invoices:", err)
    } finally {
      setIsLoading(false)
    }
  }

  const updateFilters = (newFilters: Partial<InvoiceFilters>) => {
    const updatedFilters = { ...filters, ...newFilters }
    setFilters(updatedFilters)
    fetchInvoices(updatedFilters, 1)
  }

  const goToPage = (page: number) => {
    fetchInvoices(filters, page)
  }

  const refreshInvoices = () => {
    fetchInvoices(filters, pagination.currentPage)
  }

  useEffect(() => {
    fetchInvoices()
  }, [])

  return {
    invoices,
    pagination,
    isLoading,
    error,
    filters,
    updateFilters,
    goToPage,
    refreshInvoices,
  }
}
