"use client"

import type React from "react"

import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import { cn } from "@/lib/utils"

interface FormFieldProps {
  label: string
  error?: string
  required?: boolean
  className?: string
  children: React.ReactNode
}

export function FormField({ label, error, required, className, children }: FormFieldProps) {
  return (
    <div className={cn("space-y-2", className)}>
      <Label className="text-sm font-medium text-slate-200 font-dm-sans">
        {label}
        {required && <span className="text-red-400 ml-1">*</span>}
      </Label>
      {children}
      {error && <p className="text-sm text-red-400 font-dm-sans">{error}</p>}
    </div>
  )
}

interface CurrencyInputProps {
  value: string
  onChange: (value: string) => void
  placeholder?: string
  className?: string
}

export function CurrencyInput({ value, onChange, placeholder, className }: CurrencyInputProps) {
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const inputValue = e.target.value.replace(/[^\d]/g, "")
    onChange(inputValue)
  }

  const formatValue = (val: string) => {
    if (!val) return ""
    const numValue = Number.parseInt(val) / 100
    return new Intl.NumberFormat("pt-BR", {
      style: "currency",
      currency: "BRL",
    }).format(numValue)
  }

  return (
    <Input
      type="text"
      value={formatValue(value)}
      onChange={handleChange}
      placeholder={placeholder}
      className={cn("bg-slate-700 border-slate-600 text-white placeholder:text-slate-400", className)}
    />
  )
}

interface CardNumberInputProps {
  value: string
  onChange: (value: string) => void
  placeholder?: string
  className?: string
}

export function CardNumberInput({ value, onChange, placeholder, className }: CardNumberInputProps) {
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const inputValue = e.target.value.replace(/\s/g, "").replace(/[^\d]/g, "")
    if (inputValue.length <= 16) {
      const formatted = inputValue.replace(/(.{4})/g, "$1 ").trim()
      onChange(formatted)
    }
  }

  return (
    <Input
      type="text"
      value={value}
      onChange={handleChange}
      placeholder={placeholder}
      className={cn("bg-slate-700 border-slate-600 text-white placeholder:text-slate-400", className)}
      maxLength={19} // 16 digits + 3 spaces
    />
  )
}
