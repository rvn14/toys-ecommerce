function ProductCardSkeleton() {
  return (
    <div className="animate-pulse">
      <div className="aspect-square rounded-[1.35rem] bg-brand-surface" />
      <div className="mt-4 h-3 w-20 rounded-full bg-brand-border" />
      <div className="mt-3 h-4 w-4/5 rounded-full bg-brand-border" />
      <div className="mt-2 h-4 w-3/5 rounded-full bg-brand-border" />
      <div className="mt-4 h-5 w-16 rounded-full bg-brand-border" />
    </div>
  );
}

export default function ProductsLoading() {
  return (
    <main className="fixed inset-0 z-[200] min-h-dvh overflow-y-auto bg-brand-paper" aria-busy="true" aria-label="Loading products">
      <section className="dot-pattern border-y border-brand-border bg-brand-cream-soft py-14 sm:py-20">
        <div className="site-shell animate-pulse">
          <div className="h-3 w-28 rounded-full bg-brand-border" />
          <div className="mt-5 h-12 max-w-2xl rounded-2xl bg-brand-border sm:h-16" />
          <div className="mt-5 h-4 max-w-lg rounded-full bg-brand-border" />
        </div>
      </section>
      <section className="site-shell py-10 sm:py-14">
        <div className="mb-10 h-16 animate-pulse rounded-[1.5rem] bg-brand-surface" />
        <div className="grid grid-cols-2 gap-x-4 gap-y-10 md:grid-cols-3 lg:grid-cols-4 lg:gap-x-6 xl:grid-cols-5 2xl:grid-cols-6">
          {Array.from({ length: 12 }).map((_, index) => (
            <ProductCardSkeleton key={index} />
          ))}
        </div>
      </section>
    </main>
  );
}
