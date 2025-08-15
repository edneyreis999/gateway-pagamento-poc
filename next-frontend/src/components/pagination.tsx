"use client"

import { ChevronLeft, ChevronRight } from "lucide-react"
import { Button } from "@/components/ui/button"
import type { PaginationInfo } from "../../types"

interface PaginationProps {
  pagination: PaginationInfo
  onPageChange: (page: number) => void
}

export function Pagination({ pagination, onPageChange }: PaginationProps) {
  const { currentPage, totalPages, totalItems, itemsPerPage } = pagination

  const startItem = (currentPage - 1) * itemsPerPage + 1
  const endItem = Math.min(currentPage * itemsPerPage, totalItems)

  const getVisiblePages = () => {
    const pages = []
    const maxVisible = 5

    let start = Math.max(1, currentPage - Math.floor(maxVisible / 2))
    const end = Math.min(totalPages, start + maxVisible - 1)

    if (end - start + 1 < maxVisible) {
      start = Math.max(1, end - maxVisible + 1)
    }

    for (let i = start; i <= end; i++) {
      pages.push(i)
    }

    return pages
  }

  if (totalPages <= 1) return null

  return (
    <div className="flex items-center justify-between mt-6">
      <p className="text-sm text-slate-400 font-dm-sans">
        Mostrando {startItem} - {endItem} de {totalItems} resultados
      </p>

      <div className="flex items-center gap-2">
        <Button
          variant="outline"
          size="sm"
          onClick={() => onPageChange(currentPage - 1)}
          disabled={currentPage === 1}
          className="bg-slate-700 border-slate-600 text-white hover:bg-slate-600"
        >
          <ChevronLeft className="w-4 h-4" />
        </Button>

        {getVisiblePages().map((page) => (
          <Button
            key={page}
            variant={page === currentPage ? "default" : "outline"}
            size="sm"
            onClick={() => onPageChange(page)}
            className={
              page === currentPage
                ? "bg-blue-600 text-white hover:bg-blue-700"
                : "bg-slate-700 border-slate-600 text-white hover:bg-slate-600"
            }
          >
            {page}
          </Button>
        ))}

        <Button
          variant="outline"
          size="sm"
          onClick={() => onPageChange(currentPage + 1)}
          disabled={currentPage === totalPages}
          className="bg-slate-700 border-slate-600 text-white hover:bg-slate-600"
        >
          <ChevronRight className="w-4 h-4" />
        </Button>
      </div>
    </div>
  )
}
