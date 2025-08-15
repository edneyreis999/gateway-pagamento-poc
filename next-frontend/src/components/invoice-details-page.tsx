"use client"

import { Navbar } from "@/components/navbar"
import { PageHeader } from "@/components/page-header"
import { InvoiceDetails } from "@/components/invoice-details"
import { ProtectedRoute } from "@/components/protected-route"

interface InvoiceDetailsPageProps {
  invoiceId: string
}

export function InvoiceDetailsPage({ invoiceId }: InvoiceDetailsPageProps) {
  return (
    <ProtectedRoute>
      <div className="min-h-screen bg-slate-900">
        <Navbar />
        <main className="pt-16">
          <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <PageHeader showBackButton backHref="/invoices" title="" />
            <InvoiceDetails invoiceId={invoiceId} />
          </div>
        </main>
      </div>
    </ProtectedRoute>
  )
}
