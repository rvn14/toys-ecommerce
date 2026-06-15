"use client";

import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";

type PaginationControlsProps = {
  currentPage: number;
  totalPages: number;
  className?: string;
  basePath?: string;
  query?: Record<string, string | undefined>;
  onPageChange?: (page: number) => void;
};

type PageItem = number | "ellipsis-start" | "ellipsis-end";

function getPageItems(currentPage: number, totalPages: number): PageItem[] {
  if (totalPages <= 7) {
    return Array.from({ length: totalPages }, (_, index) => index + 1);
  }

  if (currentPage <= 4) {
    return [1, 2, 3, 4, 5, "ellipsis-end", totalPages];
  }

  if (currentPage >= totalPages - 3) {
    return [1, "ellipsis-start", totalPages - 4, totalPages - 3, totalPages - 2, totalPages - 1, totalPages];
  }

  return [1, "ellipsis-start", currentPage - 1, currentPage, currentPage + 1, "ellipsis-end", totalPages];
}

export function PaginationControls({
  currentPage,
  totalPages,
  className,
  basePath,
  query = {},
  onPageChange,
}: PaginationControlsProps) {
  if (totalPages <= 1) {
    return null;
  }

  const safePage = Math.min(Math.max(currentPage, 1), totalPages);
  const pageItems = getPageItems(safePage, totalPages);

  const pageHref = (page: number) => {
    if (!basePath) return "#";
    const searchParams = new URLSearchParams();
    searchParams.set("page", String(page));
    Object.entries(query).forEach(([key, value]) => {
      if (value) searchParams.set(key, value);
    });
    return `${basePath}?${searchParams.toString()}`;
  };

  const handlePageClick = (event: React.MouseEvent<HTMLAnchorElement>, page: number) => {
    if (!onPageChange) return;
    event.preventDefault();
    onPageChange(page);
  };

  const linkClassName = "font-bold text-brand-ink hover:bg-brand-surface data-[active=true]:border-brand-border data-[active=true]:bg-brand-surface";

  return (
    <Pagination className={className}>
      <PaginationContent className="gap-1 sm:gap-2">
        <PaginationItem>
          <PaginationPrevious
            href={safePage > 1 ? pageHref(safePage - 1) : undefined}
            onClick={safePage > 1 ? (event) => handlePageClick(event, safePage - 1) : undefined}
            aria-disabled={safePage === 1}
            tabIndex={safePage === 1 ? -1 : undefined}
            className={`${linkClassName} ${safePage === 1 ? "pointer-events-none opacity-40" : ""}`}
          />
        </PaginationItem>

        {pageItems.map((item) => (
          <PaginationItem key={item}>
            {typeof item === "number" ? (
              <PaginationLink
                href={pageHref(item)}
                onClick={(event) => handlePageClick(event, item)}
                isActive={item === safePage}
                className={linkClassName}
              >
                {item}
              </PaginationLink>
            ) : (
              <PaginationEllipsis />
            )}
          </PaginationItem>
        ))}

        <PaginationItem>
          <PaginationNext
            href={safePage < totalPages ? pageHref(safePage + 1) : undefined}
            onClick={safePage < totalPages ? (event) => handlePageClick(event, safePage + 1) : undefined}
            aria-disabled={safePage === totalPages}
            tabIndex={safePage === totalPages ? -1 : undefined}
            className={`${linkClassName} ${safePage === totalPages ? "pointer-events-none opacity-40" : ""}`}
          />
        </PaginationItem>
      </PaginationContent>
    </Pagination>
  );
}
