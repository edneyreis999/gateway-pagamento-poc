"use client";

import { Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Navbar } from "@/components/navbar";
import { PageHeader } from "@/components/page-header";
import { InvoiceFiltersComponent } from "@/components/invoice-filters";
import { InvoiceTable } from "@/components/invoice-table";
import { ProtectedRoute } from "@/components/protected-route";
import { useInvoices } from "@/hooks/use-invoices";
import Link from "next/link";
import { Invoice } from "../../types";
import { logger } from "@/lib/logger";

interface InvoicesPageProps {
  initialData: {
    invoices: Invoice[];
    error: string | null;
  };
}

export function InvoicesPage({ initialData }: InvoicesPageProps) {
  logger.info('InvoicesPage render', {
    initialDataLength: initialData.invoices?.length,
    hasError: !!initialData.error
  });

  const { invoices, isLoading, error, filters, updateFilters } = useInvoices({
    initialData,
  });

  logger.state('InvoicesPage hook result', {
    invoicesLength: invoices?.length,
    isLoading,
    hasError: !!error
  });

  return (
    <ProtectedRoute>
      <div className="min-h-screen bg-slate-900">
        <Navbar />
        <main className="pt-16">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <PageHeader
              title="Faturas"
              subtitle="Gerencie suas faturas e acompanhe os pagamentos"
              action={
                <Button
                  asChild
                  className="bg-blue-600 hover:bg-blue-700 text-white"
                >
                  <Link href="/invoices/create">
                    <Plus className="w-4 h-4 mr-2" />
                    Nova Fatura
                  </Link>
                </Button>
              }
            />

            {error && (
              <div className="bg-red-900/20 border border-red-800 text-red-400 px-4 py-3 rounded-lg mb-6">
                {error}
              </div>
            )}

            <InvoiceFiltersComponent
              filters={filters}
              onFiltersChange={updateFilters}
            />

            <InvoiceTable invoices={invoices} isLoading={isLoading} />
          </div>
        </main>
      </div>
    </ProtectedRoute>
  );
}
