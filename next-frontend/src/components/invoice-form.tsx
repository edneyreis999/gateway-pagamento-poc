"use client"

import type React from "react"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { CreditCard, Lock } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Textarea } from "@/components/ui/textarea"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { FormField, CurrencyInput, CardNumberInput } from "@/components/form-field"
import { InvoiceSummary } from "@/components/invoice-summary"
import { LoadingSpinner } from "@/components/loading-spinner"
import { apiClient } from "@/lib/api"
import { Input } from "@/components/ui/input"
import type { CreateInvoiceData } from "@/types"

interface FormData {
  amount: string
  description: string
  paymentType: "credit_card" | "debit_card" | "pix"
  cardNumber: string
  expiryMonth: string
  expiryYear: string
  cvv: string
  cardholderName: string
}

interface FormErrors {
  [key: string]: string
}

export function InvoiceForm() {
  const router = useRouter()
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [formData, setFormData] = useState<FormData>({
    amount: "",
    description: "",
    paymentType: "credit_card",
    cardNumber: "",
    expiryMonth: "",
    expiryYear: "",
    cvv: "",
    cardholderName: "",
  })
  const [errors, setErrors] = useState<FormErrors>({})

  const updateField = (field: keyof FormData, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }))
    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: "" }))
    }
  }

  const validateForm = (): boolean => {
    const newErrors: FormErrors = {}

    if (!formData.amount || Number.parseInt(formData.amount) <= 0) {
      newErrors.amount = "Valor é obrigatório e deve ser maior que zero"
    }

    if (!formData.description.trim()) {
      newErrors.description = "Descrição é obrigatória"
    }

    if (formData.paymentType !== "pix") {
      if (!formData.cardNumber || formData.cardNumber.replace(/\s/g, "").length !== 16) {
        newErrors.cardNumber = "Número do cartão deve ter 16 dígitos"
      }

      if (!formData.expiryMonth) {
        newErrors.expiryMonth = "Mês de expiração é obrigatório"
      }

      if (!formData.expiryYear) {
        newErrors.expiryYear = "Ano de expiração é obrigatório"
      }

      if (!formData.cvv || formData.cvv.length !== 3) {
        newErrors.cvv = "CVV deve ter 3 dígitos"
      }

      if (!formData.cardholderName.trim()) {
        newErrors.cardholderName = "Nome no cartão é obrigatório"
      }
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validateForm()) return

    setIsSubmitting(true)

    try {
      const invoiceData: CreateInvoiceData = {
        amount: Number.parseInt(formData.amount),
        description: formData.description,
        paymentMethod: {
          type: formData.paymentType,
          ...(formData.paymentType !== "pix" && {
            cardNumber: formData.cardNumber.replace(/\s/g, ""),
            cvv: formData.cvv,
            expiryMonth: formData.expiryMonth,
            expiryYear: formData.expiryYear,
            cardholderName: formData.cardholderName,
          }),
        },
      }

      const invoice = await apiClient.createInvoice(invoiceData)
      router.push(`/invoices/${invoice.id}`)
    } catch (error) {
      console.error("Error creating invoice:", error)
      setErrors({ submit: "Erro ao criar fatura. Tente novamente." })
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleCancel = () => {
    router.push("/invoices")
  }

  const subtotal = Number.parseInt(formData.amount || "0")
  const processingFee = Math.round(subtotal * 0.02)
  const total = subtotal + processingFee

  const currentYear = new Date().getFullYear()
  const years = Array.from({ length: 10 }, (_, i) => currentYear + i)
  const months = [
    { value: "01", label: "01" },
    { value: "02", label: "02" },
    { value: "03", label: "03" },
    { value: "04", label: "04" },
    { value: "05", label: "05" },
    { value: "06", label: "06" },
    { value: "07", label: "07" },
    { value: "08", label: "08" },
    { value: "09", label: "09" },
    { value: "10", label: "10" },
    { value: "11", label: "11" },
    { value: "12", label: "12" },
  ]

  return (
    <form onSubmit={handleSubmit} className="space-y-8">
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Left Column - Basic Info */}
        <div className="space-y-6">
          <FormField label="Valor" required error={errors.amount}>
            <CurrencyInput
              value={formData.amount}
              onChange={(value) => updateField("amount", value)}
              placeholder="R$ 0,00"
            />
          </FormField>

          <FormField label="Descrição" required error={errors.description}>
            <Textarea
              value={formData.description}
              onChange={(e) => updateField("description", e.target.value)}
              placeholder="Descreva o motivo do pagamento"
              className="bg-slate-700 border-slate-600 text-white placeholder:text-slate-400 min-h-[120px]"
            />
          </FormField>
        </div>

        {/* Right Column - Payment Info */}
        <div className="space-y-6">
          <h3 className="text-lg font-semibold text-white font-space-grotesk">Dados do Cartão</h3>

          <FormField label="Número do Cartão" required error={errors.cardNumber}>
            <div className="relative">
              <CardNumberInput
                value={formData.cardNumber}
                onChange={(value) => updateField("cardNumber", value)}
                placeholder="0000 0000 0000 0000"
              />
              <CreditCard className="absolute right-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-slate-400" />
            </div>
          </FormField>

          <div className="grid grid-cols-2 gap-4">
            <FormField label="Data de Expiração" required error={errors.expiryMonth}>
              <Select value={formData.expiryMonth} onValueChange={(value) => updateField("expiryMonth", value)}>
                <SelectTrigger className="bg-slate-700 border-slate-600 text-white">
                  <SelectValue placeholder="MM/AA" />
                </SelectTrigger>
                <SelectContent className="bg-slate-700 border-slate-600">
                  {months.map((month) => (
                    <SelectItem key={month.value} value={month.value}>
                      {month.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </FormField>

            <FormField label="CVV" required error={errors.cvv}>
              <Input
                type="text"
                value={formData.cvv}
                onChange={(e) => {
                  const value = e.target.value.replace(/\D/g, "")
                  if (value.length <= 3) {
                    updateField("cvv", value)
                  }
                }}
                placeholder="123"
                className="bg-slate-700 border-slate-600 text-white placeholder:text-slate-400"
                maxLength={3}
              />
            </FormField>
          </div>

          <FormField label="Ano de Expiração" required error={errors.expiryYear}>
            <Select value={formData.expiryYear} onValueChange={(value) => updateField("expiryYear", value)}>
              <SelectTrigger className="bg-slate-700 border-slate-600 text-white">
                <SelectValue placeholder="AAAA" />
              </SelectTrigger>
              <SelectContent className="bg-slate-700 border-slate-600">
                {years.map((year) => (
                  <SelectItem key={year} value={year.toString()}>
                    {year}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </FormField>

          <FormField label="Nome no Cartão" required error={errors.cardholderName}>
            <Input
              type="text"
              value={formData.cardholderName}
              onChange={(e) => updateField("cardholderName", e.target.value)}
              placeholder="Como aparece no cartão"
              className="bg-slate-700 border-slate-600 text-white placeholder:text-slate-400"
            />
          </FormField>
        </div>
      </div>

      {/* Invoice Summary */}
      <InvoiceSummary subtotal={subtotal} processingFee={processingFee} total={total} />

      {/* Error Message */}
      {errors.submit && (
        <div className="bg-red-900/20 border border-red-800 text-red-400 px-4 py-3 rounded-lg">{errors.submit}</div>
      )}

      {/* Action Buttons */}
      <div className="flex justify-end gap-4">
        <Button type="button" variant="outline" onClick={handleCancel} disabled={isSubmitting}>
          Cancelar
        </Button>
        <Button
          type="submit"
          disabled={isSubmitting}
          className="bg-blue-600 hover:bg-blue-700 text-white min-w-[180px]"
        >
          {isSubmitting ? (
            <>
              <LoadingSpinner size="sm" className="mr-2" />
              Processando...
            </>
          ) : (
            <>
              <Lock className="w-4 h-4 mr-2" />
              Processar Pagamento
            </>
          )}
        </Button>
      </div>
    </form>
  )
}
