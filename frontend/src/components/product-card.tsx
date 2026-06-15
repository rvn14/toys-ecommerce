import Image from "next/image";
import Link from "next/link";
import { Heart, ShoppingBag, Star } from "lucide-react";
import type { Product } from "@/lib/types";
import { formatPriceFromCents } from "@/lib/format";

type ProductCardProps = { product: Product };

export function ProductCard({ product }: ProductCardProps) {
  const firstImage = product.images[0];
  const hasDiscount = product.discountPriceCents !== null && product.discountPriceCents > 0 && product.discountPriceCents < product.priceCents;
  const effectivePrice = hasDiscount ? product.discountPriceCents! : product.priceCents;
  const productHref = `/products/${product.id}`;

  return (
    <article className="group min-w-0">
      <Link href={productHref} className="relative block aspect-square overflow-hidden rounded-[1.35rem] bg-brand-surface" aria-label={`View ${product.name}`}>
        {hasDiscount && <span className="absolute left-3 top-3 z-10 rounded-full bg-brand-orange px-3 py-1 text-[10px] font-black uppercase tracking-wider text-white">Sale</span>}
        <span className="absolute right-3 top-3 z-10 grid h-9 w-9 place-items-center rounded-full bg-white/90 opacity-0 shadow-sm transition group-hover:opacity-100"><Heart className="h-4 w-4" /></span>
        {firstImage?.url ? <Image src={firstImage.url} alt={firstImage.altText || product.name} fill className="object-cover transition duration-500 group-hover:scale-105" sizes="(max-width: 640px) 50vw, (max-width: 1024px) 33vw, 20vw" /> : <div className="grid h-full place-items-center text-sm text-brand-subtle">No image</div>}
        <span className="absolute bottom-3 left-3 right-3 flex translate-y-4 items-center justify-center gap-2 rounded-full bg-brand-ink px-4 py-3 text-xs font-bold text-white opacity-0 transition group-hover:translate-y-0 group-hover:opacity-100"><ShoppingBag className="h-4 w-4" /> View product</span>
      </Link>
      <div className="px-1 pt-4">
        <div className="flex items-center gap-1 text-brand-gold" aria-label="5 out of 5 stars">{Array.from({ length: 5 }).map((_, index) => <Star key={index} className="h-3 w-3 fill-current" />)}</div>
        <p className="mt-2 text-[10px] font-bold uppercase tracking-[0.12em] text-brand-subtle">{product.category.name}</p>
        <h2 className="mt-1 line-clamp-2 min-h-10 text-sm font-bold leading-5 text-brand-ink"><Link href={productHref} className="hover:text-brand-orange">{product.name}</Link></h2>
        <div className="mt-2 flex items-baseline gap-2"><span className="font-display text-lg text-brand-orange">{formatPriceFromCents(effectivePrice)}</span>{hasDiscount && <span className="text-xs text-brand-subtle line-through">{formatPriceFromCents(product.priceCents)}</span>}</div>
      </div>
    </article>
  );
}
