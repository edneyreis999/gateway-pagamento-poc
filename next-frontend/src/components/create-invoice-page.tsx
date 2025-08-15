"use client"

import { Navbar } from "@/components/navbar"
import { PageHeader } from "@/components/page-header"
import { InvoiceForm } from "@/components/invoice-form"
import { ProtectedRoute } from "@/components/protected-route"

export function CreateInvoicePage() {
  return (
    <ProtectedRoute>
      <div className="min-h-screen bg-slate-900">
        <Navbar />
        <main className="pt-16">
          <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <PageHeader
              title="Criar Nova Fatura"
              subtitle="Preencha os dados abaixo para processar um novo pagamento"
              showBackButton
              backHref="/invoices"
            />

            <div className="bg-slate-800 rounded-lg p-6">
              <InvoiceForm />
            </div>
          </div>
        </main>

        {/* Footer */}
        <footer className="bg-slate-800 border-t border-slate-700 py-4 mt-16">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <p className="text-center text-slate-400 text-sm font-dm-sans">
              © 2025 Full Cycle Gateway. Todos os direitos reservados.
            </p>
          </div>
        </footer>
      </div>
    </ProtectedRoute>
  )
}
