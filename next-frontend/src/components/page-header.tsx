import type React from "react"
import { ArrowLeft } from "lucide-react"
import { Button } from "@/components/ui/button"
import Link from "next/link"

interface PageHeaderProps {
  title: string
  subtitle?: string
  showBackButton?: boolean
  backHref?: string
  action?: React.ReactNode
}

export function PageHeader({
  title,
  subtitle,
  showBackButton = false,
  backHref = "/invoices",
  action,
}: PageHeaderProps) {
  return (
    <div className="flex items-center justify-between mb-8">
      <div className="flex items-center gap-4">
        {showBackButton && (
          <Button variant="ghost" size="sm" asChild>
            <Link href={backHref}>
              <ArrowLeft className="w-4 h-4 mr-2" />
              Voltar
            </Link>
          </Button>
        )}
        <div>
          <h1 className="text-2xl font-bold text-white font-space-grotesk">{title}</h1>
          {subtitle && <p className="text-slate-400 mt-1 font-dm-sans">{subtitle}</p>}
        </div>
      </div>
      {action && <div>{action}</div>}
    </div>
  )
}
