"use client";

import { useEffect } from "react";
import Link from "next/link";
import { MessageBox } from "@/components/message-box";

type ErrorPageProps = {
  error: Error & { digest?: string };
  unstable_retry: () => void;
};

export default function ErrorPage({ error, unstable_retry }: ErrorPageProps) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <main className="dot-pattern grid min-h-[65vh] place-items-center bg-brand-cream-soft px-5 py-16">
      <MessageBox
        eyebrow="Something went wrong"
        title="PLAYTIME HIT A BUMP"
        description="We could not finish loading this page. Try again, or head back to the shop while we put things in order."
        variant="error"
      >
        <button
          type="button"
          onClick={() => unstable_retry()}
          className="rounded-full bg-brand-orange px-7 py-3 text-xs font-black uppercase tracking-wider text-white transition hover:bg-brand-ink"
        >
          Try again
        </button>
        <Link
          href="/products"
          className="rounded-full border-2 border-brand-ink px-7 py-3 text-xs font-black uppercase tracking-wider transition hover:bg-brand-ink hover:text-white"
        >
          Browse toys
        </Link>
      </MessageBox>
    </main>
  );
}
