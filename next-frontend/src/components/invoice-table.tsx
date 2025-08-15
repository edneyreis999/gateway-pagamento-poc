"use client"

import { Eye, Download } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { StatusBadge } from "@/components/status-badge"
import { LoadingSpinner } from "@/components/loading-spinner"
import { formatCurrency, formatDate } from "@/lib/utils"
import Link from "next/link"
import type { Invoice } from "@/types"

interface InvoiceTableProps {
  invoices: Invoice[]
  isLoading: boolean
}

export function InvoiceTable({ invoices, isLoading }: InvoiceTableProps) {
  if (isLoading) {
    return (
      <div className="bg-slate-800 rounded-lg p-8 flex items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    )
  }

  if (invoices.length === 0) {
    return (
      <div className="bg-slate-800 rounded-lg p-8 text-center">
        <p className="text-slate-400 font-dm-sans">Nenhuma fatura encontrada</p>
      </div>
    )
  }

  return (
    <div className="bg-slate-800 rounded-lg overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow className="border-slate-700 hover:bg-slate-700/50">
            <TableHead className="text-slate-300 font-dm-sans">ID</TableHead>
            <TableHead className="text-slate-300 font-dm-sans">DATA</TableHead>
            <TableHead className="text-slate-300 font-dm-sans">DESCRIÇÃO</TableHead>
            <TableHead className="text-slate-300 font-dm-sans">VALOR</TableHead>
            <TableHead className="text-slate-300 font-dm-sans">STATUS</TableHead>
            <TableHead className="text-slate-300 font-dm-sans">AÇÕES</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {invoices.map((invoice) => (
            <TableRow key={invoice.id} className="border-slate-700 hover:bg-slate-700/30">
              <TableCell className="font-mono text-slate-300">{invoice.id}</TableCell>
              <TableCell className="text-slate-300 font-dm-sans">
                {formatDate(invoice.createdAt).split(" ")[0]}
              </TableCell>
              <TableCell className="text-white font-dm-sans">{invoice.description}</TableCell>
              <TableCell className="text-white font-dm-sans font-semibold">{formatCurrency(invoice.amount)}</TableCell>
              <TableCell>
                <StatusBadge status={invoice.status} />
              </TableCell>
              <TableCell>
                <div className="flex items-center gap-2">
                  <Button variant="ghost" size="sm" asChild>
                    <Link href={`/invoices/${invoice.id}`}>
                      <Eye className="w-4 h-4 text-blue-400" />
                    </Link>
                  </Button>
                  <Button variant="ghost" size="sm">
                    <Download className="w-4 h-4 text-slate-400" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
