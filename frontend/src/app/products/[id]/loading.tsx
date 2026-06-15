export default function ProductLoading() {
  return (
    <main className="fixed inset-0 z-[200] min-h-dvh overflow-y-auto bg-brand-paper" aria-label="Loading product" aria-busy="true">
      <div className="site-shell grid min-h-dvh animate-pulse items-center gap-10 py-10 lg:grid-cols-2 lg:gap-16 lg:py-16">
        <div className="aspect-square rounded-[2rem] bg-brand-surface" />
        <div>
          <div className="h-3 w-40 rounded-full bg-brand-border" />
          <div className="mt-6 h-14 w-4/5 rounded-2xl bg-brand-border" />
          <div className="mt-3 h-14 w-3/5 rounded-2xl bg-brand-border" />
          <div className="mt-7 h-8 w-28 rounded-full bg-brand-border" />
          <div className="mt-8 h-4 w-full rounded-full bg-brand-border" />
          <div className="mt-3 h-4 w-5/6 rounded-full bg-brand-border" />
          <div className="mt-8 h-14 w-44 rounded-full bg-brand-border" />
        </div>
      </div>
    </main>
  );
}
