"use client"

import { useState, useEffect } from "react"
import { Download } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { StatusBadge } from "@/components/status-badge"
import { TransactionTimeline } from "@/components/transaction-timeline"
import { LoadingSpinner } from "@/components/loading-spinner"
import { formatCurrency, formatDate } from "@/lib/utils"
import { apiClient } from "@/lib/api"
import type { Invoice } from "@/types"

interface InvoiceDetailsProps {
  invoiceId: string
}

export function InvoiceDetails({ invoiceId }: InvoiceDetailsProps) {
  const [invoice, setInvoice] = useState<Invoice | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchInvoice = async () => {
      try {
        setIsLoading(true)
        setError(null)
        const data = await apiClient.getInvoice(invoiceId)
        setInvoice(data)
      } catch (err) {
        setError("Erro ao carregar detalhes da fatura")
        console.error("Error fetching invoice:", err)
      } finally {
        setIsLoading(false)
      }
    }

    fetchInvoice()
  }, [invoiceId])

  const handleDownloadPDF = () => {
    // Simulate PDF download
    console.log("Downloading PDF for invoice:", invoiceId)
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <LoadingSpinner size="lg" />
      </div>
    )
  }

  if (error || !invoice) {
    return (
      <div className="bg-red-900/20 border border-red-800 text-red-400 px-4 py-3 rounded-lg">
        {error || "Fatura não encontrada"}
      </div>
    )
  }

  const timelineItems = [
    {
      title: "Fatura Criada",
      timestamp: invoice.createdAt,
      completed: true,
    },
    {
      title: "Pagamento Processado",
      timestamp: invoice.createdAt, // In real app, this would be different
      completed: invoice.status !== "pending",
    },
    {
      title: "Transação Aprovada",
      timestamp: invoice.updatedAt,
      completed: invoice.status === "approved",
    },
  ]

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <div>
            <div className="flex items-center gap-3 mb-2">
              <h1 className="text-2xl font-bold text-white font-space-grotesk">Fatura {invoice.id}</h1>
              <StatusBadge status={invoice.status} />
            </div>
            <p className="text-slate-400 font-dm-sans">Criada em {formatDate(invoice.createdAt)}</p>
          </div>
        </div>
        <Button onClick={handleDownloadPDF} className="bg-slate-700 hover:bg-slate-600 text-white">
          <Download className="w-4 h-4 mr-2" />
          Download PDF
        </Button>
      </div>

      {/* Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Invoice Information */}
        <Card className="bg-slate-800 border-slate-700">
          <CardHeader>
            <CardTitle className="text-white font-space-grotesk">Informações da Fatura</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex justify-between">
              <span className="text-slate-400 font-dm-sans">ID da Fatura</span>
              <span className="text-white font-mono">{invoice.id}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-slate-400 font-dm-sans">Valor</span>
              <span className="text-white font-semibold font-dm-sans">{formatCurrency(invoice.amount)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-slate-400 font-dm-sans">Data de Criação</span>
              <span className="text-white font-dm-sans">{formatDate(invoice.createdAt)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-slate-400 font-dm-sans">Última Atualização</span>
              <span className="text-white font-dm-sans">{formatDate(invoice.updatedAt)}</span>
            </div>
            <div className="pt-2 border-t border-slate-600">
              <span className="text-slate-400 font-dm-sans">Descrição</span>
              <p className="text-white mt-1 font-dm-sans">{invoice.description}</p>
            </div>
          </CardContent>
        </Card>

        {/* Transaction Status */}
        <Card className="bg-slate-800 border-slate-700">
          <CardHeader>
            <CardTitle className="text-white font-space-grotesk">Status da Transação</CardTitle>
          </CardHeader>
          <CardContent>
            <TransactionTimeline items={timelineItems} />
          </CardContent>
        </Card>

        {/* Payment Method */}
        <Card className="bg-slate-800 border-slate-700">
          <CardHeader>
            <CardTitle className="text-white font-space-grotesk">Método de Pagamento</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex justify-between">
              <span className="text-slate-400 font-dm-sans">Tipo</span>
              <span className="text-white font-dm-sans">
                {invoice.paymentMethod.type === "credit_card"
                  ? "Cartão de Crédito"
                  : invoice.paymentMethod.type === "debit_card"
                    ? "Cartão de Débito"
                    : "PIX"}
              </span>
            </div>
            {invoice.paymentMethod.lastFourDigits && (
              <div className="flex justify-between">
                <span className="text-slate-400 font-dm-sans">Últimos Dígitos</span>
                <span className="text-white font-mono">**** **** **** {invoice.paymentMethod.lastFourDigits}</span>
              </div>
            )}
            {invoice.paymentMethod.cardholderName && (
              <div className="flex justify-between">
                <span className="text-slate-400 font-dm-sans">Titular</span>
                <span className="text-white font-dm-sans">{invoice.paymentMethod.cardholderName}</span>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Additional Data */}
        <Card className="bg-slate-800 border-slate-700">
          <CardHeader>
            <CardTitle className="text-white font-space-grotesk">Dados Adicionais</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex justify-between">
              <span className="text-slate-400 font-dm-sans">ID da Conta</span>
              <span className="text-white font-mono">{invoice.accountId}</span>
            </div>
            {invoice.clientIp && (
              <div className="flex justify-between">
                <span className="text-slate-400 font-dm-sans">IP do Cliente</span>
                <span className="text-white font-mono">{invoice.clientIp}</span>
              </div>
            )}
            {invoice.device && (
              <div className="flex justify-between">
                <span className="text-slate-400 font-dm-sans">Dispositivo</span>
                <span className="text-white font-dm-sans">{invoice.device}</span>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
