import { InvoiceDetailsPage } from "@/components/invoice-details-page"

interface InvoiceDetailsRouteProps {
  params: Promise<{
    id: string
  }>
}

export default async function InvoiceDetailsRoute({ params }: InvoiceDetailsRouteProps) {
  const { id } = await params
  return <InvoiceDetailsPage invoiceId={id} />
}
