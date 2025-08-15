import { formatCurrency } from "@/lib/utils"

interface InvoiceSummaryProps {
  subtotal: number
  processingFee: number
  total: number
}

export function InvoiceSummary({ subtotal, processingFee, total }: InvoiceSummaryProps) {
  return (
    <div className="bg-slate-800 rounded-lg p-6 space-y-4">
      <div className="flex justify-between items-center text-slate-300">
        <span className="font-dm-sans">Subtotal</span>
        <span className="font-dm-sans">{formatCurrency(subtotal)}</span>
      </div>

      <div className="flex justify-between items-center text-slate-300">
        <span className="font-dm-sans">Taxa de Processamento (2%)</span>
        <span className="font-dm-sans">{formatCurrency(processingFee)}</span>
      </div>

      <div className="border-t border-slate-600 pt-4">
        <div className="flex justify-between items-center text-white text-lg font-semibold">
          <span className="font-space-grotesk">Total</span>
          <span className="font-space-grotesk">{formatCurrency(total)}</span>
        </div>
      </div>
    </div>
  )
}
