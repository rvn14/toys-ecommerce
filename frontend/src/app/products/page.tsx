import type { Metadata } from "next";
import Link from "next/link";
import { MessageBox } from "@/components/message-box";
import { ProductCard } from "@/components/product-card";
import { ProductFilters } from "@/components/product-filters";
import { ProductPagination } from "@/components/product-pagination";
import { getBrands, getCategories, getProducts } from "@/lib/api-client";

type ProductsPageProps = {
  searchParams: Promise<{
    page?: string;
    search?: string;
    category?: string;
    brand?: string;
    sort?: string;
  }>;
};

export const metadata: Metadata = {
  title: "Shop Toys, Games & Collectibles",
  description: "Browse ToyJoy's newest toys, building sets, plush friends, action figures, games, and collectibles.",
  alternates: { canonical: "/products" },
  openGraph: {
    title: "Shop Toys, Games & Collectibles",
    description: "Find colorful toys and collectible favorites for every kind of play.",
    url: "/products",
  },
};

export default async function ProductsPage({ searchParams }: ProductsPageProps) {
  const params = await searchParams;
  const page = Math.max(1, Number(params.page || "1") || 1);
  const [products, categories, brands] = await Promise.all([
    getProducts({
      page,
      limit: 12,
      search: params.search,
      category: params.category,
      brand: params.brand,
      sort: params.sort || "newest",
    }),
    getCategories(),
    getBrands(),
  ]);

  return (
    <main>
      <section className="dot-pattern border-y border-brand-border bg-brand-cream-soft px-5 py-14 sm:py-20">
        <div className="site-shell">
          <p className="text-xs font-black uppercase tracking-[0.22em] text-brand-orange">The toy box</p>
          <h1 className="font-display mt-2 text-5xl sm:text-7xl">FIND YOUR NEW FAVORITE</h1>
          <p className="mt-4 max-w-2xl text-sm leading-6 text-brand-muted sm:text-base">
            Toys, games, collectibles, plush friends, and building sets for every kind of play.
          </p>
        </div>
      </section>

      <section className="site-shell py-10 sm:py-14">
        <ProductFilters categories={categories} brands={brands} values={params} />

        <div className="mb-8 flex flex-col justify-between gap-3 border-b border-brand-border pb-5 sm:flex-row sm:items-end">
          <div>
            <p className="text-xs font-black uppercase tracking-[0.18em] text-brand-subtle">{products.pagination.totalItems} products</p>
            <h2 className="font-display mt-1 text-3xl">
              {params.search ? `Results for "${params.search}"` : params.category ? params.category.replaceAll("-", " ") : "All products"}
            </h2>
          </div>
          {(params.search || params.category || params.brand) && <Link href="/products" className="text-xs font-black uppercase text-brand-orange">Clear filters</Link>}
        </div>

        {products.items.length === 0 ? (
          <MessageBox eyebrow="Nothing matched" title="NO TOYS FOUND" description="Try a different search or browse the complete collection." variant="empty" headingLevel="h2">
            <Link href="/products" className="rounded-full bg-brand-orange px-6 py-3 text-xs font-black uppercase text-white transition hover:bg-brand-ink">Browse all products</Link>
          </MessageBox>
        ) : (
          <div className="grid grid-cols-2 gap-x-4 gap-y-10 md:grid-cols-3 lg:grid-cols-4 lg:gap-x-6 xl:grid-cols-5 2xl:grid-cols-6">
            {products.items.map((product) => <ProductCard key={product.id} product={product} />)}
          </div>
        )}

        <ProductPagination pagination={products.pagination} filters={{ search: params.search, category: params.category, brand: params.brand, sort: params.sort }} />
      </section>
    </main>
  );
}
