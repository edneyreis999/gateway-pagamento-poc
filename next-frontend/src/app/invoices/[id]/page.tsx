import { InvoiceDetailsPage } from "@/components/invoice-details-page"

interface InvoiceDetailsRouteProps {
  params: {
    id: string
  }
}

export default function InvoiceDetailsRoute({ params }: InvoiceDetailsRouteProps) {
  return <InvoiceDetailsPage invoiceId={params.id} />
}
