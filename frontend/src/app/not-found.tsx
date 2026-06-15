import Link from "next/link";
import { MessageBox } from "@/components/message-box";

export default function NotFound() {
  return (
    <main className="dot-pattern grid min-h-[65vh] place-items-center bg-brand-yellow px-5 py-16">
      <MessageBox
        eyebrow="404 - Page not found"
        title="THIS TOY WANDERED OFF"
        description="The page you are looking for does not exist, has moved, or is taking an unscheduled play break."
        variant="empty"
      >
        <Link
          href="/"
          className="rounded-full bg-brand-ink px-7 py-3 text-xs font-black uppercase tracking-wider text-white transition hover:bg-brand-orange"
        >
          Back home
        </Link>
        <Link
          href="/products"
          className="rounded-full border-2 border-brand-ink px-7 py-3 text-xs font-black uppercase tracking-wider transition hover:bg-brand-ink hover:text-white"
        >
          Visit the shop
        </Link>
      </MessageBox>
    </main>
  );
}
