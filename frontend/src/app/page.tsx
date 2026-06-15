import Image from "next/image";
import Link from "next/link";
import { ArrowRight, Clock3, Gift, ShieldCheck, Sparkles, Truck } from "lucide-react";
import { getProducts } from "@/lib/api-client";
import { ProductCard } from "@/components/product-card";
import { formatPriceFromCents } from "@/lib/format";

const ageGroups = [
  { title: "Birth to 24 months", subtitle: "Soft, safe & sensory", color: "bg-brand-pink", index: 4 },
  { title: "2 to 4 years", subtitle: "Learn through play", color: "bg-brand-aqua", index: 6 },
  { title: "5 to 7 years", subtitle: "Big imagination", color: "bg-brand-lime", index: 8 },
  { title: "8 to 13 years", subtitle: "Build & discover", color: "bg-brand-coral", index: 1 },
];

export default async function HomePage() {
  const products = await getProducts({ page: 1, limit: 15, sort: "newest" });
  const items = products.items;
  const featured = items.slice(0, 10);
  const latest = items.slice(10, 14).length ? items.slice(10, 14) : items.slice(0, 4);
  const saleItems = items.filter((product) => product.discountPriceCents).slice(0, 3);
  const imageAt = (index: number) => items[index % items.length]?.images[0];

  return (
    <main>
      <section className="relative min-h-[560px] overflow-hidden bg-brand-orange lg:min-h-[650px]">
        <Image src="/images/toy-hero.png" alt="Original robot and dinosaur toy mascots" fill priority className="object-cover object-[68%_center] lg:object-center" sizes="100vw" />
        <div className="absolute inset-0 bg-gradient-to-r from-brand-orange via-brand-orange/80 to-transparent lg:via-brand-orange/28" />
        <div className="site-shell relative flex min-h-[560px] items-center py-20 lg:min-h-[650px]">
          <div className="max-w-xl text-white">
            <div className="mb-5 inline-flex items-center gap-2 rounded-full bg-white/15 px-4 py-2 text-xs font-bold uppercase tracking-[0.16em] backdrop-blur-sm"><Sparkles className="h-4 w-4 text-brand-yellow" /> New adventures inside</div>
            <h1 className="font-display text-5xl leading-[0.88] sm:text-7xl lg:text-[6.8rem]">PLAY BIG.<br /><span className="text-brand-yellow">DREAM<span className="sm:hidden"><br /></span><span className="hidden sm:inline"> </span>LOUD.</span></h1>
            <p className="mt-6 max-w-[310px] text-sm leading-7 text-white/88 sm:max-w-md sm:text-lg">Discover colorful toys, clever building sets, and collectible friends made for every kind of imagination.</p>
            <div className="mt-8 flex flex-wrap gap-3"><Link href="/products" className="inline-flex items-center gap-2 rounded-full bg-white px-7 py-4 text-sm font-black text-brand-ink transition hover:-translate-y-0.5">Shop now <ArrowRight className="h-4 w-4" /></Link><Link href="#featured" className="rounded-full border border-white/50 px-7 py-4 text-sm font-black text-white transition hover:bg-white/10">Explore favorites</Link></div>
          </div>
        </div>
      </section>

      <section className="border-b border-brand-border bg-white">
        <div className="site-shell grid grid-cols-2 divide-x divide-y divide-brand-border md:grid-cols-4 md:divide-y-0">
          {[{ icon: Truck, title: "Free delivery", text: "Orders over $100" }, { icon: Gift, title: "Gift-ready", text: "Made to delight" }, { icon: Clock3, title: "Easy returns", text: "Within 30 days" }, { icon: ShieldCheck, title: "Secure pay", text: "Protected checkout" }].map(({ icon: Icon, title, text }) => <div key={title} className="flex min-w-0 items-center gap-3 px-3 py-5 sm:px-4"><Icon className="h-6 w-6 shrink-0 text-brand-orange" /><div className="min-w-0"><p className="text-[10px] font-black uppercase tracking-wide sm:text-xs">{title}</p><p className="mt-1 truncate text-xs text-brand-subtle">{text}</p></div></div>)}
        </div>
      </section>

      <section id="featured" className="site-shell py-20">
        <div className="mb-10 flex flex-col justify-between gap-5 sm:flex-row sm:items-end"><div><p className="text-xs font-black uppercase tracking-[0.2em] text-brand-orange">Hand-picked fun</p><h2 className="font-display mt-2 text-4xl sm:text-5xl">FEATURED PRODUCTS</h2></div><div className="flex gap-5 text-xs font-bold"><Link href="/products?category=action-figures" className="border-b-2 border-brand-orange pb-2">Action figures</Link><Link href="/products?category=building-sets" className="pb-2 text-brand-subtle">Building sets</Link><Link href="/products?category=plush-toys" className="pb-2 text-brand-subtle">Plush toys</Link></div></div>
        <div className="grid grid-cols-2 gap-x-4 gap-y-10 md:grid-cols-3 lg:grid-cols-5 lg:gap-x-5">{featured.map((product) => <ProductCard key={product.id} product={product} />)}</div>
        <div className="mt-12 text-center"><Link href="/products" className="inline-flex items-center gap-2 rounded-full border-2 border-brand-ink px-7 py-3 text-xs font-black uppercase tracking-wider transition hover:bg-brand-ink hover:text-white">View all toys <ArrowRight className="h-4 w-4" /></Link></div>
      </section>

      <section className="squiggle-pattern overflow-hidden py-20">
        <div className="site-shell">
          <div className="mb-10 text-center"><p className="text-xs font-black uppercase tracking-[0.2em]">Choose your adventure</p><h2 className="font-display mt-2 text-4xl sm:text-5xl">SHOP TOYS BY AGE</h2></div>
          <div className="grid gap-5 md:grid-cols-2 lg:grid-cols-4">{ageGroups.map((group) => { const image = imageAt(group.index); return <Link key={group.title} href="/products" className="group relative min-h-64 overflow-hidden rounded-[2rem] bg-white p-7 toy-shadow"><div className={`absolute -right-16 -top-16 h-52 w-52 rounded-full ${group.color} opacity-65 transition duration-500 group-hover:scale-110`} />{image?.url && <Image src={image.url} alt={image.altText || group.title} width={210} height={210} className="absolute -bottom-5 -right-5 h-44 w-44 rounded-full object-cover mix-blend-multiply transition duration-500 group-hover:scale-105" />}<p className="relative text-xs font-bold uppercase tracking-wider text-brand-muted">Explore ages</p><h3 className="font-display relative mt-4 max-w-[150px] text-3xl leading-[0.95]">{group.title}</h3><p className="relative mt-3 text-xs font-semibold text-brand-muted">{group.subtitle}</p><span className="relative mt-8 inline-flex h-10 w-10 items-center justify-center rounded-full bg-brand-ink text-white"><ArrowRight className="h-4 w-4" /></span></Link>; })}</div>
        </div>
      </section>

      <section className="site-shell grid gap-10 py-20 lg:grid-cols-[1fr_380px]">
        <div><div className="mb-8 flex items-end justify-between"><div><p className="text-xs font-black uppercase tracking-[0.2em] text-brand-purple">Fresh from the toy box</p><h2 className="font-display mt-2 text-4xl">LATEST PRODUCTS</h2></div><Link href="/products" className="text-xs font-black uppercase">View all</Link></div><div className="grid grid-cols-2 gap-5 md:grid-cols-4">{latest.map((product) => <ProductCard key={product.id} product={product} />)}</div></div>
        <aside className="rounded-[2rem] border-2 border-dashed border-brand-sale-border bg-brand-sale p-7"><p className="text-xs font-black uppercase tracking-[0.2em] text-brand-orange">Limited offers</p><h2 className="font-display mt-2 text-4xl">ON SALE</h2><div className="mt-6 grid gap-5">{saleItems.map((product) => <Link href={`/products/${product.id}`} key={product.id} className="flex items-center gap-4 border-b border-brand-sale-divider pb-5 last:border-0 last:pb-0">{product.images[0]?.url && <Image src={product.images[0].url} alt={product.name} width={72} height={72} className="h-16 w-16 rounded-xl bg-white object-cover" />}<div className="min-w-0"><h3 className="line-clamp-2 text-xs font-bold leading-4">{product.name}</h3><p className="mt-2 font-display text-lg text-brand-orange">{formatPriceFromCents(product.discountPriceCents || product.priceCents)}</p></div></Link>)}</div></aside>
      </section>

      <section className="overflow-hidden bg-brand-purple text-white">
        <div className="site-shell grid min-h-[360px] items-center gap-10 py-14 md:grid-cols-2">
          <div className="relative order-2 h-72 md:order-1">{imageAt(5)?.url && <><div className="absolute left-1/2 top-1/2 h-64 w-64 -translate-x-1/2 -translate-y-1/2 rounded-full bg-brand-yellow" /><Image src={imageAt(5)!.url} alt="Weekly toy pick" fill className="relative object-contain mix-blend-lighten" sizes="50vw" /></>}</div>
          <div className="order-1 text-center md:order-2 md:text-left"><p className="text-xs font-black uppercase tracking-[0.24em] text-brand-yellow">This week only</p><h2 className="font-display mt-3 text-5xl leading-[0.92] sm:text-6xl">WE CAN&apos;T KEEP<br /><span className="text-brand-yellow">THESE IN STOCK!</span></h2><p className="mx-auto mt-5 max-w-md text-sm leading-6 text-white/70 md:mx-0">Popular picks are moving fast. Grab a new favorite before playtime starts without you.</p><Link href="/products?sort=newest" className="mt-7 inline-flex items-center gap-2 rounded-full bg-white px-7 py-3 text-xs font-black uppercase text-brand-purple-deep">Shop the drop <ArrowRight className="h-4 w-4" /></Link></div>
        </div>
      </section>

      <section className="site-shell py-20">
        <div className="text-center"><p className="text-xs font-black uppercase tracking-[0.2em] text-brand-orange">Ideas for brighter play</p><h2 className="font-display mt-2 text-4xl sm:text-5xl">FROM THE PLAYROOM</h2></div>
        <div className="mt-10 grid gap-6 md:grid-cols-3">{[{ title: "5 ways building toys grow big ideas", index: 2, tag: "Play & learn" }, { title: "The perfect gift for every little explorer", index: 7, tag: "Gift guide" }, { title: "How to start a collection they will love", index: 11, tag: "Collectors" }].map((post) => { const image = imageAt(post.index); return <article key={post.title} className="group"><div className="relative aspect-[16/9] overflow-hidden rounded-[1.6rem] bg-brand-surface">{image?.url && <Image src={image.url} alt={post.title} fill className="object-cover transition duration-500 group-hover:scale-105" sizes="(max-width: 768px) 100vw, 33vw" />}</div><p className="mt-5 text-[10px] font-black uppercase tracking-[0.18em] text-brand-orange">{post.tag}</p><h3 className="font-display mt-2 text-2xl leading-tight">{post.title}</h3><Link href="/products" className="mt-3 inline-flex items-center gap-2 text-xs font-black uppercase">Read story <ArrowRight className="h-3 w-3" /></Link></article>; })}</div>
      </section>
    </main>
  );
}
