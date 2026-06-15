import Link from "next/link";
import { MessageBox } from "@/components/message-box";

export default function ProductNotFound() {
  return (
    <main className="dot-pattern grid min-h-[65vh] place-items-center bg-brand-cream-soft px-5 py-16">
      <MessageBox
        eyebrow="Product not found"
        title="THAT TOY LEFT THE SHELF"
        description="This product may no longer be available, or the link may be incorrect. There are plenty more favorites waiting in the toy box."
        variant="empty"
      >
        <Link
          href="/products"
          className="rounded-full bg-brand-orange px-7 py-3 text-xs font-black uppercase tracking-wider text-white transition hover:bg-brand-ink"
        >
          Browse all products
        </Link>
      </MessageBox>
    </main>
  );
}
