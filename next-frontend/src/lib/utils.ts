import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatCurrency(amount: number): string {
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
  }).format(amount / 100) // Assumindo que o valor vem em centavos
}

export function formatDate(dateString: string): string {
  return new Intl.DateTimeFormat("pt-BR", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(dateString))
}

export function formatCardNumber(cardNumber: string): string {
  return cardNumber
    .replace(/\s/g, "")
    .replace(/(.{4})/g, "$1 ")
    .trim()
}

export function maskCardNumber(cardNumber: string): string {
  const cleaned = cardNumber.replace(/\s/g, "")
  const lastFour = cleaned.slice(-4)
  return `**** **** **** ${lastFour}`
}

export function getStatusColor(status: string): string {
  switch (status) {
    case "approved":
      return "status-approved"
    case "pending":
      return "status-pending"
    case "rejected":
      return "status-rejected"
    default:
      return "bg-gray-500/20 text-gray-400 border-gray-500/30"
  }
}

export function getStatusLabel(status: string): string {
  switch (status) {
    case "approved":
      return "Aprovado"
    case "pending":
      return "Pendente"
    case "rejected":
      return "Rejeitado"
    default:
      return status
  }
}
