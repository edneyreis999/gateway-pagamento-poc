"use client";

import { useState } from "react";
import { getInvoicesAction } from "@/lib/actions/invoices";
import type { Invoice, InvoiceFilters } from "../../types";

interface UseInvoicesProps {
  initialData: {
    invoices: Invoice[];
    error: string | null;
  };
}

export function useInvoices({ initialData }: UseInvoicesProps) {
  const [invoices, setInvoices] = useState<Invoice[]>(initialData.invoices);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(initialData.error);
  const [filters, setFilters] = useState<InvoiceFilters>({
    status: "all",
    startDate: "",
    endDate: "",
    search: "",
  });

  const fetchInvoices = async (newFilters?: InvoiceFilters, page = 1) => {
    try {
      setIsLoading(true);
      setError(null);

      const filtersToUse = newFilters || filters;
      const response = await getInvoicesAction({
        ...filtersToUse,
        page,
      });

      setInvoices(response.invoices);
      setError(response.error);
    } catch (err) {
      setError("Erro ao carregar faturas. Tente novamente.");
      console.error("Error fetching invoices:", err);
    } finally {
      setIsLoading(false);
    }
  };

  const updateFilters = (newFilters: Partial<InvoiceFilters>) => {
    const updatedFilters = { ...filters, ...newFilters };
    setFilters(updatedFilters);
    fetchInvoices(updatedFilters, 1);
  };

  const goToPage = (page: number) => {
    fetchInvoices(filters, page);
  };

  return {
    invoices,
    isLoading,
    error,
    filters,
    updateFilters,
    goToPage,
  };
}
