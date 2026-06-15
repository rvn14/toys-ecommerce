import Link from "next/link";
import { Search, SlidersHorizontal, Tags } from "lucide-react";
import type { Brand, Category } from "@/lib/types";

type ProductFiltersProps = {
  categories: Category[];
  brands: Brand[];
  values: {
    search?: string;
    category?: string;
    brand?: string;
    sort?: string;
  };
};

export function ProductFilters({ categories, brands, values }: ProductFiltersProps) {
  const hasFilters = Boolean(values.search || values.category || values.brand || (values.sort && values.sort !== "newest"));

  return (
    <div className="mb-10">
      <form className="grid gap-3 rounded-[1.5rem] border border-brand-border bg-white p-3 shadow-sm lg:grid-cols-[minmax(240px,1fr)_190px_190px_190px_auto]" action="/products">
        <label className="flex items-center gap-3 rounded-xl bg-brand-surface px-4">
          <Search className="h-4 w-4 shrink-0 text-brand-subtle" />
          <span className="sr-only">Search products</span>
          <input name="search" defaultValue={values.search} placeholder="Search products or SKU" className="min-w-0 flex-1 bg-transparent py-3 text-sm outline-none" />
        </label>
        <label className="flex items-center gap-3 rounded-xl bg-brand-surface px-4">
          <Tags className="h-4 w-4 shrink-0 text-brand-subtle" />
          <span className="sr-only">Filter by category</span>
          <select name="category" defaultValue={values.category || ""} className="w-full bg-transparent py-3 text-sm font-semibold outline-none">
            <option value="">All categories</option>
            {categories.map((category) => <option key={category.id} value={category.slug}>{category.name}</option>)}
          </select>
        </label>
        <label className="flex items-center gap-3 rounded-xl bg-brand-surface px-4">
          <Tags className="h-4 w-4 shrink-0 text-brand-subtle" />
          <span className="sr-only">Filter by brand</span>
          <select name="brand" defaultValue={values.brand || ""} className="w-full bg-transparent py-3 text-sm font-semibold outline-none">
            <option value="">All brands</option>
            {brands.map((brand) => <option key={brand.id} value={brand.slug}>{brand.name}</option>)}
          </select>
        </label>
        <label className="flex items-center gap-3 rounded-xl bg-brand-surface px-4">
          <SlidersHorizontal className="h-4 w-4 shrink-0 text-brand-subtle" />
          <span className="sr-only">Sort products</span>
          <select name="sort" defaultValue={values.sort || "newest"} className="w-full bg-transparent py-3 text-sm font-semibold outline-none">
            <option value="newest">Newest first</option>
            <option value="price_asc">Price: low to high</option>
            <option value="price_desc">Price: high to low</option>
            <option value="name_asc">Name: A to Z</option>
            <option value="name_desc">Name: Z to A</option>
          </select>
        </label>
        <button className="rounded-xl bg-brand-ink px-7 py-3 text-xs font-black uppercase tracking-wider text-white transition hover:bg-brand-orange">Apply filters</button>
      </form>
      {hasFilters && (
        <div className="mt-3 flex justify-end">
          <Link href="/products" className="text-xs font-black uppercase tracking-wider text-brand-orange">Clear all filters</Link>
        </div>
      )}
    </div>
  );
}
