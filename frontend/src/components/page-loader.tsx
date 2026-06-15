type PageLoaderProps = {
  label?: string;
  message?: string;
};

export function PageLoader({
  label = "Loading",
  message = "Opening the toy box...",
}: PageLoaderProps) {
  return (
    <div
      className="fixed inset-0 z-[200] grid min-h-dvh place-items-center overflow-hidden bg-brand-cream-soft px-5"
      role="status"
      aria-live="polite"
    >
      <div className="text-center">
        <div className="relative mx-auto h-20 w-24" aria-hidden="true">
          <span className="absolute bottom-1 left-1/2 h-12 w-12 -translate-x-1/2 animate-bounce rounded-2xl bg-brand-orange" />
          <span className="absolute left-1 top-2 h-7 w-7 animate-pulse rounded-full bg-brand-yellow" />
          <span className="absolute right-1 top-1 h-8 w-8 animate-pulse rounded-xl bg-brand-purple [animation-delay:180ms]" />
        </div>
        <p className="mt-4 text-xs font-black uppercase tracking-[0.22em] text-brand-orange">{label}</p>
        <p className="font-display mt-2 text-3xl">{message}</p>
      </div>
    </div>
  );
}
