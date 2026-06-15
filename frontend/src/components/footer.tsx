import Link from "next/link";
import { Camera, Check, Mail, MapPin, MessageCircle, Phone, Send, Truck } from "lucide-react";

export function Footer() {
  return (
    <footer className="border-t border-brand-border bg-white">
      <div className="grid w-full divide-y divide-white/20 overflow-hidden bg-brand-night text-white md:grid-cols-3 md:divide-x md:divide-y-0">
        <div className="flex items-center gap-4 px-6 py-5">
          <Truck className="h-7 w-7 text-brand-yellow" />
          <div><p className="font-bold">Free delivery</p><p className="text-xs text-white/60">On orders over $100</p></div>
        </div>
        <div className="flex items-center gap-4 px-6 py-5">
          <span className="font-display text-3xl text-brand-orange">30</span>
          <div><p className="font-bold">30 days return</p><p className="text-xs text-white/60">Easy, friendly returns</p></div>
        </div>
        <div className="flex items-center gap-4 px-6 py-5">
          <span className="grid h-9 w-9 place-items-center rounded-full bg-brand-green"><Check className="h-5 w-5" /></span>
          <div><p className="font-bold">Secure payment</p><p className="text-xs text-white/60">Protected checkout</p></div>
        </div>
      </div>

      <div className="site-shell grid gap-10 py-14 md:grid-cols-2 lg:grid-cols-[1.6fr_1fr_1fr_1.2fr]">
        <div>
          <Link href="/" className="font-display text-3xl text-brand-ink">TOY<span className="text-brand-orange">JOY</span></Link>
          <p className="mt-4 max-w-sm text-sm leading-6 text-brand-muted">Playful picks for curious kids, collectors, and everyone who still knows how to have fun.</p>
          <form className="mt-6 flex max-w-sm overflow-hidden rounded-full border border-brand-border bg-brand-cream" action="/products">
            <label htmlFor="footer-search" className="sr-only">Search the store</label>
            <input id="footer-search" name="search" placeholder="Find your next favorite toy" className="min-w-0 flex-1 bg-transparent px-5 py-3 text-sm outline-none" />
            <button className="m-1 grid h-10 w-10 place-items-center rounded-full bg-brand-orange text-white" aria-label="Search"><Send className="h-4 w-4" /></button>
          </form>
        </div>

        <div><h2 className="font-display text-xl">ABOUT</h2><div className="mt-4 grid gap-3 text-sm text-brand-muted"><Link href="/">Our story</Link><Link href="/products">New arrivals</Link><Link href="/products?sort=price_asc">Special offers</Link><Link href="/account">My account</Link></div></div>
        <div><h2 className="font-display text-xl">SHOP</h2><div className="mt-4 grid gap-3 text-sm text-brand-muted"><Link href="/products?category=building-sets">Building sets</Link><Link href="/products?category=plush-toys">Plush toys</Link><Link href="/products?category=action-figures">Action figures</Link><Link href="/products">All products</Link></div></div>
        <div><h2 className="font-display text-xl">CONTACT</h2><div className="mt-4 grid gap-3 text-sm text-brand-muted"><p className="flex gap-2"><MapPin className="mt-0.5 h-4 w-4 shrink-0 text-brand-orange" />Colombo, Sri Lanka</p><p className="flex gap-2"><Phone className="h-4 w-4 text-brand-orange" />+94 77 123 4567</p><p className="flex gap-2"><Mail className="h-4 w-4 text-brand-orange" />hello@toyjoy.store</p></div><div className="mt-5 flex gap-2"><span className="grid h-9 w-9 place-items-center rounded-full bg-brand-surface"><MessageCircle className="h-4 w-4" /></span><span className="grid h-9 w-9 place-items-center rounded-full bg-brand-surface"><Camera className="h-4 w-4" /></span></div></div>
      </div>

      <div className="border-t border-brand-border py-5 text-center text-xs text-brand-subtle">Copyright 2026 ToyJoy. Built for big imaginations.</div>
    </footer>
  );
}
