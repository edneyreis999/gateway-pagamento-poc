"use client"

import { CalendarIcon, Search } from "lucide-react"
import { Input } from "@/components/ui/input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import type { InvoiceFilters } from "@/types"

interface InvoiceFiltersProps {
  filters: InvoiceFilters
  onFiltersChange: (filters: Partial<InvoiceFilters>) => void
}

export function InvoiceFiltersComponent({ filters, onFiltersChange }: InvoiceFiltersProps) {
  return (
    <div className="bg-slate-800 rounded-lg p-6 mb-6">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Status Filter */}
        <div className="space-y-2">
          <label className="text-sm font-medium text-slate-200 font-dm-sans">Status</label>
          <Select value={filters.status || "all"} onValueChange={(value) => onFiltersChange({ status: value as any })}>
            <SelectTrigger className="bg-slate-700 border-slate-600 text-white">
              <SelectValue placeholder="Todos" />
            </SelectTrigger>
            <SelectContent className="bg-slate-700 border-slate-600">
              <SelectItem value="all">Todos</SelectItem>
              <SelectItem value="pending">Pendente</SelectItem>
              <SelectItem value="approved">Aprovado</SelectItem>
              <SelectItem value="rejected">Rejeitado</SelectItem>
            </SelectContent>
          </Select>
        </div>

        {/* Start Date */}
        <div className="space-y-2">
          <label className="text-sm font-medium text-slate-200 font-dm-sans">Data Inicial</label>
          <div className="relative">
            <Input
              type="date"
              value={filters.startDate || ""}
              onChange={(e) => onFiltersChange({ startDate: e.target.value })}
              className="bg-slate-700 border-slate-600 text-white placeholder:text-slate-400"
              placeholder="dd/mm/aaaa"
            />
            <CalendarIcon className="absolute right-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-slate-400" />
          </div>
        </div>

        {/* End Date */}
        <div className="space-y-2">
          <label className="text-sm font-medium text-slate-200 font-dm-sans">Data Final</label>
          <div className="relative">
            <Input
              type="date"
              value={filters.endDate || ""}
              onChange={(e) => onFiltersChange({ endDate: e.target.value })}
              className="bg-slate-700 border-slate-600 text-white placeholder:text-slate-400"
              placeholder="dd/mm/aaaa"
            />
            <CalendarIcon className="absolute right-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-slate-400" />
          </div>
        </div>

        {/* Search */}
        <div className="space-y-2">
          <label className="text-sm font-medium text-slate-200 font-dm-sans">Buscar</label>
          <div className="relative">
            <Input
              type="text"
              value={filters.search || ""}
              onChange={(e) => onFiltersChange({ search: e.target.value })}
              className="bg-slate-700 border-slate-600 text-white placeholder:text-slate-400 pl-10"
              placeholder="ID ou descrição"
            />
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-slate-400" />
          </div>
        </div>
      </div>
    </div>
  )
}
