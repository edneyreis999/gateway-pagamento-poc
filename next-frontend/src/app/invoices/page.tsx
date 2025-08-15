import { InvoicesPage } from "@/components/invoices-page"
import { getInvoicesAction } from "@/lib/actions/invoices"
import { logger } from "@/lib/logger"

export default async function InvoicesRoute() {
  logger.route('Invoices route accessed');
  
  const initialData = await getInvoicesAction()

  logger.info('Invoices page data loaded', {
    invoiceCount: initialData.invoices?.length,
    hasError: !!initialData.error
  });
  
  return <InvoicesPage initialData={initialData} />
}
