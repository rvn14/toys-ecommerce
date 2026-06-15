import type { Metadata } from "next";
import Image from "next/image";
import Link from "next/link";
import { notFound } from "next/navigation";
import { ArrowLeft, PackageCheck, ShieldCheck, Star, Truck } from "lucide-react";
import { AddToCart } from "@/components/add-to-cart";
import { getProduct } from "@/lib/api-client";
import { formatPriceFromCents } from "@/lib/format";

type ProductPageProps = {
  params: Promise<{ id: string }>;
};

function parseProductId(value: string) {
  const id = Number(value);
  return Number.isInteger(id) && id > 0 ? id : null;
}

export async function generateMetadata({ params }: ProductPageProps): Promise<Metadata> {
  const { id: idParam } = await params;
  const id = parseProductId(idParam);

  if (!id) {
    return { title: "Product not found" };
  }

  const product = await getProduct(id);

  if (!product) {
    return { title: "Product not found" };
  }

  const description = product.description || `Shop ${product.name} from ToyJoy.`;
  const image = product.images[0]?.url;

  return {
    title: product.name,
    description,
    alternates: { canonical: `/products/${product.id}` },
    openGraph: {
      title: product.name,
      description,
      type: "website",
      url: `/products/${product.id}`,
      images: image ? [{ url: image, alt: product.images[0]?.altText || product.name }] : undefined,
    },
  };
}

export default async function ProductPage({ params }: ProductPageProps) {
  const { id: idParam } = await params;
  const id = parseProductId(idParam);

  if (!id) {
    notFound();
  }

  const product = await getProduct(id);

  if (!product) {
    notFound();
  }

  const firstImage = product.images[0];
  const hasDiscount =
    product.discountPriceCents !== null &&
    product.discountPriceCents > 0 &&
    product.discountPriceCents < product.priceCents;
  const effectivePrice = hasDiscount ? product.discountPriceCents! : product.priceCents;
  const specifications = Object.entries(product.specifications || {});
  const formatSpecification = (value: unknown) => {
    if (typeof value === "string" || typeof value === "number" || typeof value === "boolean") {
      return String(value);
    }

    return JSON.stringify(value, null, 2);
  };

  return (
    <main>
      <section className="border-y border-brand-border bg-brand-cream-soft py-6">
        <div className="site-shell">
          <Link href="/products" className="inline-flex items-center gap-2 text-xs font-black uppercase tracking-wider text-brand-muted transition hover:text-brand-orange">
            <ArrowLeft className="h-4 w-4" /> Back to all toys
          </Link>
        </div>
      </section>

      <section className="site-shell grid gap-10 py-10 lg:grid-cols-2 lg:gap-16 lg:py-16">
        <div className="relative aspect-square overflow-hidden rounded-[2rem] bg-brand-surface toy-shadow">
          {hasDiscount && (
            <span className="absolute left-5 top-5 z-10 rounded-full bg-brand-orange px-4 py-2 text-xs font-black uppercase tracking-wider text-white">
              Sale
            </span>
          )}
          {firstImage?.url ? (
            <Image
              src={firstImage.url}
              alt={firstImage.altText || product.name}
              fill
              priority
              className="object-cover"
              sizes="(max-width: 1024px) 100vw, 50vw"
            />
          ) : (
            <div className="grid h-full place-items-center text-sm text-brand-subtle">No image available</div>
          )}
        </div>

        <div className="self-center">
          <p className="text-xs font-black uppercase tracking-[0.2em] text-brand-orange">
            {product.category.name} | {product.brand.name}
          </p>
          <h1 className="font-display mt-3 text-5xl leading-[0.95] sm:text-6xl">{product.name}</h1>
          <div className="mt-5 flex items-center gap-2 text-brand-gold" aria-label="5 out of 5 stars">
            {Array.from({ length: 5 }).map((_, index) => (
              <Star key={index} className="h-4 w-4 fill-current" />
            ))}
            <span className="ml-2 text-xs font-bold text-brand-muted">A ToyJoy favorite</span>
          </div>
          <div className="mt-7 flex items-baseline gap-3">
            <span className="font-display text-4xl text-brand-orange">{formatPriceFromCents(effectivePrice)}</span>
            {hasDiscount && <span className="text-base text-brand-subtle line-through">{formatPriceFromCents(product.priceCents)}</span>}
          </div>
          <p className="mt-6 text-sm leading-7 text-brand-muted sm:text-base">
            {product.description || "A playful pick made for curious minds and big imaginations."}
          </p>

          <div className="mt-7 rounded-2xl bg-brand-surface p-5">
            <p className="text-xs font-black uppercase tracking-wider text-brand-subtle">Availability</p>
            <p className={`mt-2 font-bold ${product.stockQuantity > 0 ? "text-brand-success" : "text-brand-orange"}`}>
              {product.stockQuantity > 0 ? `${product.stockQuantity} ready to ship` : "Currently out of stock"}
            </p>
          </div>

          <AddToCart productId={product.id} stockQuantity={product.stockQuantity} />

          <div className="mt-8 grid gap-3 border-t border-brand-border pt-7 sm:grid-cols-3">
            <div className="flex items-center gap-3 text-xs font-bold"><Truck className="h-5 w-5 text-brand-orange" /> Fast delivery</div>
            <div className="flex items-center gap-3 text-xs font-bold"><ShieldCheck className="h-5 w-5 text-brand-orange" /> Secure checkout</div>
            <div className="flex items-center gap-3 text-xs font-bold"><PackageCheck className="h-5 w-5 text-brand-orange" /> Easy returns</div>
          </div>
        </div>
      </section>

      {specifications.length > 0 && (
        <section className="site-shell pb-16">
          <div className="rounded-[2rem] border border-brand-border bg-white p-7 sm:p-10">
            <p className="text-xs font-black uppercase tracking-[0.2em] text-brand-purple">Product details</p>
            <h2 className="font-display mt-2 text-3xl">GOOD TO KNOW</h2>
            <dl className="mt-7 grid gap-3 sm:grid-cols-2">
              {specifications.map(([key, value]) => (
                <div key={key} className="rounded-2xl bg-brand-surface p-5">
                  <dt className="text-[10px] font-black uppercase tracking-wider text-brand-subtle">{key}</dt>
                  <dd className="mt-2 whitespace-pre-wrap text-sm font-bold">{formatSpecification(value)}</dd>
                </div>
              ))}
            </dl>
          </div>
        </section>
      )}
    </main>
  );
}
