import { InvoicesPage } from "@/components/invoices-page"
import { getInvoicesAction } from "@/lib/actions/invoices"

export default async function InvoicesRoute() {
  const initialData = await getInvoicesAction()

  console.log('---- initialData ----')
  console.log(initialData.invoices[0])
  
  return <InvoicesPage initialData={initialData} />
}
