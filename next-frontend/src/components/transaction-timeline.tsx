import { Check } from "lucide-react"
import { formatDate } from "@/lib/utils"

interface TimelineItem {
  title: string
  timestamp: string
  completed: boolean
}

interface TransactionTimelineProps {
  items: TimelineItem[]
}

export function TransactionTimeline({ items }: TransactionTimelineProps) {
  return (
    <div className="space-y-4">
      {items.map((item, index) => (
        <div key={index} className="flex items-center gap-3">
          <div
            className={`flex-shrink-0 w-6 h-6 rounded-full flex items-center justify-center ${
              item.completed ? "bg-green-500 text-white" : "bg-slate-600 text-slate-400"
            }`}
          >
            {item.completed && <Check className="w-4 h-4" />}
          </div>
          <div className="flex-1">
            <p className="text-white font-medium font-dm-sans">{item.title}</p>
            <p className="text-slate-400 text-sm font-dm-sans">{formatDate(item.timestamp)}</p>
          </div>
        </div>
      ))}
    </div>
  )
}
