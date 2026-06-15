"use client";

import Link from "next/link";
import { Suspense, useEffect, useState } from "react";
import { usePathname, useSearchParams } from "next/navigation";
import { Heart, Search, ShoppingBag, UserRound } from "lucide-react";
import { useAuth } from "@/components/auth-provider";
import { cn } from "@/lib/utils";
import { getCart } from "@/lib/cart-api";

const links = [
  { href: "/", label: "Home" },
  { href: "/products", label: "Shop" },
  { href: "/products?category=building-sets", label: "Building Sets" },
  { href: "/products?category=plush-toys", label: "Plush Toys" },
];

function DesktopNav({ isAdmin }: { isAdmin: boolean }) {
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const navigationLinks = isAdmin ? [...links, { href: "/admin", label: "Admin" }] : links;

  const isActive = (href: string) => {
    const [linkPath, queryString] = href.split("?");
    const category = new URLSearchParams(queryString).get("category");

    if (category) {
      return pathname === linkPath && searchParams.get("category") === category;
    }

    if (linkPath === "/") {
      return pathname === "/";
    }

    return pathname.startsWith(linkPath) && !searchParams.get("category");
  };

  return (
    <nav className="hidden items-center gap-8 text-xs font-bold uppercase tracking-[0.11em] lg:flex" aria-label="Main navigation">
      {navigationLinks.map((link) => {
        const active = isActive(link.href);

        return (
          <Link
            key={link.href}
            href={link.href}
            aria-current={active ? "page" : undefined}
            className={cn(
              "relative py-2 transition hover:text-brand-orange after:absolute after:inset-x-0 after:bottom-0 after:h-0.5 after:origin-left after:scale-x-0 after:rounded-full after:bg-brand-orange after:transition-transform",
              active && "text-brand-orange after:scale-x-100"
            )}
          >
            {link.label}
          </Link>
        );
      })}
    </nav>
  );
}

export function Navbar() {
  const { isAuthenticated, token, user } = useAuth();
  const [cartCount, setCartCount] = useState(0);

  useEffect(() => {
    if (!token) {
      return;
    }

    void getCart(token).then((cart) => setCartCount(cart.totalItems)).catch(() => setCartCount(0));

    const handleCartUpdate = (event: Event) => {
      setCartCount((event as CustomEvent<number>).detail);
    };
    window.addEventListener("cart-updated", handleCartUpdate);
    return () => window.removeEventListener("cart-updated", handleCartUpdate);
  }, [token]);

  const displayedCartCount = token ? cartCount : 0;

  return (
    <header className="relative z-50 bg-white">
      <div className="bg-brand-surface px-4 py-2 text-[10px] font-semibold uppercase tracking-[0.14em] text-brand-muted">
        <div className="site-shell flex justify-between">
          <span>Free delivery over $100</span>
          <span className="hidden sm:inline">Joy delivered worldwide</span>
        </div>
      </div>

      <div className="site-shell flex h-[76px] items-center justify-between">
        <Link href="/" className="font-display text-3xl tracking-tight text-brand-ink">TOY<span className="text-brand-orange">JOY</span></Link>

        <Suspense fallback={<div className="hidden h-8 w-[430px] lg:block" />}>
          <DesktopNav isAdmin={user?.role === "admin"} />
        </Suspense>

        <div className="flex items-center gap-1 sm:gap-2">
          <Link href="/products" className="hidden h-10 w-10 place-items-center rounded-full transition hover:bg-brand-cream sm:grid" aria-label="Search"><Search className="h-[18px] w-[18px]" /></Link>
          <button className="hidden h-10 w-10 place-items-center rounded-full transition hover:bg-brand-cream sm:grid" aria-label="Wishlist"><Heart className="h-[18px] w-[18px]" /></button>
          <Link href={isAuthenticated ? "/account" : "/login"} className="hidden h-10 w-10 place-items-center rounded-full transition hover:bg-brand-cream sm:grid" aria-label="Account"><UserRound className="h-[18px] w-[18px]" /></Link>
          <button className="relative hidden h-10 w-10 place-items-center rounded-full transition hover:bg-brand-cream sm:grid" aria-label={`Cart with ${displayedCartCount} items`}><ShoppingBag className="h-[18px] w-[18px]" /><span className="absolute right-0 top-0 grid h-4 min-w-4 place-items-center rounded-full bg-brand-yellow px-1 text-[9px] font-black">{displayedCartCount}</span></button>
          <Link href="/products" className="fixed right-4 top-[46px] z-[100] rounded-full bg-brand-ink px-5 py-3 text-[10px] font-black uppercase tracking-wider text-white lg:hidden">Shop</Link>
        </div>
      </div>
    </header>
  );
}
