"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { ArrowRight, Boxes, LayoutDashboard, Package, RefreshCw, Tags } from "lucide-react";
import { useAuth } from "@/components/auth-provider";
import { CatalogEntityManager } from "@/components/admin/catalog-entity-manager";
import { ProductManager } from "@/components/admin/product-manager";
import { PageLoader } from "@/components/page-loader";
import { adminApi } from "@/lib/admin-api";
import type { Brand, CatalogEntityRequest, Category, Product, ProductRequest } from "@/lib/types";

type AdminSection = "dashboard" | "categories" | "brands" | "products";

const sections = [
  { id: "dashboard" as const, label: "Dashboard", icon: LayoutDashboard },
  { id: "categories" as const, label: "Categories", icon: Tags },
  { id: "brands" as const, label: "Brands", icon: Boxes },
  { id: "products" as const, label: "Products", icon: Package },
];

export function AdminPortal() {
  const { token, user } = useAuth();
  const [activeSection, setActiveSection] = useState<AdminSection>("dashboard");
  const [categories, setCategories] = useState<Category[]>([]);
  const [brands, setBrands] = useState<Brand[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");

  const loadCatalog = useCallback(async () => {
    if (!token) return;
    setIsLoading(true);
    setError("");
    try {
      const [categoryData, brandData, productData] = await Promise.all([
        adminApi.getCategories(token),
        adminApi.getBrands(token),
        adminApi.getProducts(token),
      ]);
      setCategories(categoryData);
      setBrands(brandData);
      setProducts(productData.items);
    } catch (loadError) {
      setError(loadError instanceof Error ? loadError.message : "Could not load the admin catalog");
    } finally {
      setIsLoading(false);
    }
  }, [token]);

  useEffect(() => {
    const timeoutId = window.setTimeout(() => {
      void loadCatalog();
    }, 0);

    return () => window.clearTimeout(timeoutId);
  }, [loadCatalog]);

  if (isLoading) {
    return <PageLoader label="Admin catalog" message="Loading inventory..." />;
  }

  const sortByName = <T extends { name: string }>(items: T[]) => [...items].sort((a, b) => a.name.localeCompare(b.name));

  const createCategory = async (data: CatalogEntityRequest) => {
    const created = await adminApi.createCategory(token!, data);
    setCategories((current) => sortByName([...current, created]));
  };
  const updateCategory = async (id: number, data: CatalogEntityRequest) => {
    const updated = await adminApi.updateCategory(token!, id, data);
    setCategories((current) => sortByName(current.map((item) => item.id === id ? updated : item)));
  };
  const deleteCategory = async (id: number) => {
    await adminApi.deleteCategory(token!, id);
    setCategories((current) => current.filter((item) => item.id !== id));
  };

  const createBrand = async (data: CatalogEntityRequest) => {
    const created = await adminApi.createBrand(token!, data);
    setBrands((current) => sortByName([...current, created]));
  };
  const updateBrand = async (id: number, data: CatalogEntityRequest) => {
    const updated = await adminApi.updateBrand(token!, id, data);
    setBrands((current) => sortByName(current.map((item) => item.id === id ? updated : item)));
  };
  const deleteBrand = async (id: number) => {
    await adminApi.deleteBrand(token!, id);
    setBrands((current) => current.filter((item) => item.id !== id));
  };

  const createProduct = async (data: ProductRequest) => {
    const created = await adminApi.createProduct(token!, data);
    setProducts((current) => [created, ...current]);
  };
  const updateProduct = async (id: number, data: ProductRequest) => {
    const updated = await adminApi.updateProduct(token!, id, data);
    setProducts((current) => current.map((item) => item.id === id ? updated : item));
  };
  const deleteProduct = async (id: number) => {
    await adminApi.deleteProduct(token!, id);
    setProducts((current) => current.filter((item) => item.id !== id));
  };

  return (
    <main className="min-h-screen bg-brand-cream">
      <section className="dot-pattern border-y border-brand-border bg-brand-purple py-12 text-white">
        <div className="site-shell flex flex-col justify-between gap-6 lg:flex-row lg:items-end">
          <div><p className="text-xs font-black uppercase tracking-[0.22em] text-brand-yellow">ToyJoy operations</p><h1 className="font-display mt-2 text-5xl sm:text-7xl">ADMIN PORTAL</h1><p className="mt-3 text-sm text-white/70">Welcome, {user?.name}. Manage the storefront catalog from one place.</p></div>
          <div className="flex flex-wrap gap-3"><button type="button" onClick={() => void loadCatalog()} className="inline-flex items-center gap-2 rounded-full border border-white/30 px-5 py-3 text-xs font-black uppercase tracking-wider"><RefreshCw className="h-4 w-4" /> Refresh data</button><Link href="/products" className="inline-flex items-center gap-2 rounded-full bg-white px-5 py-3 text-xs font-black uppercase tracking-wider text-brand-purple">View storefront <ArrowRight className="h-4 w-4" /></Link></div>
        </div>
      </section>

      <div className="site-shell py-8">
        <nav className="flex gap-2 overflow-x-auto rounded-2xl border border-brand-border bg-white p-2" aria-label="Admin sections">
          {sections.map(({ id, label, icon: Icon }) => <button key={id} type="button" onClick={() => setActiveSection(id)} className={`inline-flex min-w-fit flex-1 items-center justify-center gap-2 rounded-xl px-4 py-3 text-xs font-black uppercase tracking-wider transition ${activeSection === id ? "bg-brand-ink text-white" : "text-brand-muted hover:bg-brand-surface"}`}><Icon className="h-4 w-4" />{label}</button>)}
        </nav>

        {error && <div className="mt-6 flex items-center justify-between gap-4 rounded-2xl border border-destructive/20 bg-destructive/10 p-5 text-sm font-semibold text-destructive"><span>{error}</span><button type="button" onClick={() => void loadCatalog()} className="underline">Retry</button></div>}

        <div className="mt-8">
          {activeSection === "dashboard" && <AdminDashboard categories={categories} brands={brands} products={products} onOpen={setActiveSection} />}
          {activeSection === "categories" && <CatalogEntityManager title="Categories" singular="Category" items={categories} onCreate={createCategory} onUpdate={updateCategory} onDelete={deleteCategory} />}
          {activeSection === "brands" && <CatalogEntityManager title="Brands" singular="Brand" items={brands} onCreate={createBrand} onUpdate={updateBrand} onDelete={deleteBrand} />}
          {activeSection === "products" && <ProductManager products={products} categories={categories} brands={brands} onCreate={createProduct} onUpdate={updateProduct} onDelete={deleteProduct} />}
        </div>
      </div>
    </main>
  );
}

function AdminDashboard({ categories, brands, products, onOpen }: { categories: Category[]; brands: Brand[]; products: Product[]; onOpen: (section: AdminSection) => void }) {
  const activeProducts = products.filter((product) => product.status === "active").length;
  const lowStock = products.filter((product) => product.stockQuantity <= 5).length;
  const cards = [
    { label: "Categories", value: categories.length, section: "categories" as const, color: "bg-brand-yellow" },
    { label: "Brands", value: brands.length, section: "brands" as const, color: "bg-brand-aqua" },
    { label: "Products", value: products.length, section: "products" as const, color: "bg-brand-pink" },
    { label: "Active products", value: activeProducts, section: "products" as const, color: "bg-brand-lime" },
  ];

  return (
    <div>
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">{cards.map((card) => <button key={card.label} type="button" onClick={() => onOpen(card.section)} className={`rounded-[2rem] p-6 text-left transition hover:-translate-y-1 ${card.color}`}><p className="text-xs font-black uppercase tracking-[0.18em]">{card.label}</p><p className="font-display mt-3 text-5xl">{card.value}</p><span className="mt-5 inline-flex items-center gap-2 text-xs font-black uppercase">Manage <ArrowRight className="h-4 w-4" /></span></button>)}</div>
      <div className="mt-6 grid gap-6 lg:grid-cols-[1fr_360px]">
        <section className="rounded-[2rem] border border-brand-border bg-white p-7"><p className="text-xs font-black uppercase tracking-[0.18em] text-brand-purple">Catalog flow</p><h2 className="font-display mt-2 text-3xl">QUICK ACTIONS</h2><div className="mt-6 grid gap-3 sm:grid-cols-3"><button type="button" onClick={() => onOpen("categories")} className="rounded-2xl bg-brand-surface p-5 text-left text-sm font-black">Add a category<ArrowRight className="mt-4 h-4 w-4 text-brand-orange" /></button><button type="button" onClick={() => onOpen("brands")} className="rounded-2xl bg-brand-surface p-5 text-left text-sm font-black">Add a brand<ArrowRight className="mt-4 h-4 w-4 text-brand-orange" /></button><button type="button" onClick={() => onOpen("products")} className="rounded-2xl bg-brand-surface p-5 text-left text-sm font-black">Create a product<ArrowRight className="mt-4 h-4 w-4 text-brand-orange" /></button></div></section>
        <aside className="rounded-[2rem] bg-brand-ink p-7 text-white"><p className="text-xs font-black uppercase tracking-[0.18em] text-brand-yellow">Inventory watch</p><p className="font-display mt-3 text-5xl">{lowStock}</p><p className="mt-2 text-sm text-white/65">products have five or fewer units remaining.</p><button type="button" onClick={() => onOpen("products")} className="mt-6 inline-flex items-center gap-2 text-xs font-black uppercase text-brand-yellow">Review stock <ArrowRight className="h-4 w-4" /></button></aside>
      </div>
    </div>
  );
}
