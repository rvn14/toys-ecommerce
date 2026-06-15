import { PaginationControls } from "@/components/pagination-controls";
import type { PaginationMeta } from "@/lib/types";

type ProductPaginationProps = {
  pagination: PaginationMeta;
  filters: Record<string, string | undefined>;
};

export function ProductPagination({ pagination, filters }: ProductPaginationProps) {
  return (
    <PaginationControls
      currentPage={pagination.page}
      totalPages={pagination.totalPages}
      basePath="/products"
      query={filters}
      className="mt-14"
    />
  );
}
