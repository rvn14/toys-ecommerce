"use client";

import { useMemo, useState } from "react";
import Image from "next/image";
import { Pencil, Plus, Search, Trash2, X } from "lucide-react";
import { PaginationControls } from "@/components/pagination-controls";
import type { Brand, Category, Product, ProductRequest } from "@/lib/types";
import { formatPriceFromCents } from "@/lib/format";

type ProductManagerProps = {
  products: Product[];
  categories: Category[];
  brands: Brand[];
  onCreate: (data: ProductRequest) => Promise<void>;
  onUpdate: (id: number, data: ProductRequest) => Promise<void>;
  onDelete: (id: number) => Promise<void>;
};

const emptyForm = {
  name: "",
  sku: "",
  description: "",
  price: "",
  discountPrice: "",
  stockQuantity: "0",
  status: "active",
  categoryId: "",
  brandId: "",
  imageUrls: "",
  specifications: "{}",
};

export function ProductManager({ products, categories, brands, onCreate, onUpdate, onDelete }: ProductManagerProps) {
  const itemsPerPage = 8;
  const [form, setForm] = useState(emptyForm);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [search, setSearch] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [message, setMessage] = useState("");

  const filteredProducts = useMemo(() => {
    const query = search.trim().toLowerCase();
    if (!query) return products;
    return products.filter((product) => [product.name, product.sku, product.category.name, product.brand.name].some((value) => value.toLowerCase().includes(query)));
  }, [products, search]);
  const totalPages = Math.max(1, Math.ceil(filteredProducts.length / itemsPerPage));
  const safePage = Math.min(currentPage, totalPages);
  const paginatedProducts = filteredProducts.slice((safePage - 1) * itemsPerPage, safePage * itemsPerPage);

  function updateField(field: keyof typeof form, value: string) {
    setForm((current) => ({ ...current, [field]: value }));
  }

  function resetForm() {
    setEditingId(null);
    setForm(emptyForm);
  }

  function startEditing(product: Product) {
    setEditingId(product.id);
    setForm({
      name: product.name,
      sku: product.sku,
      description: product.description,
      price: (product.priceCents / 100).toFixed(2),
      discountPrice: product.discountPriceCents ? (product.discountPriceCents / 100).toFixed(2) : "",
      stockQuantity: String(product.stockQuantity),
      status: product.status,
      categoryId: String(product.category.id),
      brandId: String(product.brand.id),
      imageUrls: product.images.map((image) => image.url).join("\n"),
      specifications: JSON.stringify(product.specifications || {}, null, 2),
    });
    setMessage("");
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  function buildRequest(): ProductRequest {
    const specifications = JSON.parse(form.specifications || "{}") as unknown;
    if (!specifications || Array.isArray(specifications) || typeof specifications !== "object") {
      throw new Error("Specifications must be a JSON object");
    }

    const imageUrls = form.imageUrls.split(/\r?\n/).map((url) => url.trim()).filter(Boolean);
    const priceCents = Math.round(Number(form.price) * 100);
    const discountPriceCents = form.discountPrice ? Math.round(Number(form.discountPrice) * 100) : null;

    if (!Number.isFinite(priceCents) || priceCents < 1) throw new Error("Enter a valid product price");
    if (discountPriceCents !== null && (!Number.isFinite(discountPriceCents) || discountPriceCents >= priceCents)) throw new Error("Discount price must be lower than the regular price");

    return {
      name: form.name.trim(),
      sku: form.sku.trim(),
      description: form.description.trim(),
      priceCents,
      discountPriceCents,
      stockQuantity: Number(form.stockQuantity),
      status: form.status,
      categoryId: Number(form.categoryId),
      brandId: Number(form.brandId),
      specifications: specifications as Record<string, unknown>,
      images: imageUrls.map((url, index) => ({ url, storagePath: "", altText: form.name.trim(), sortOrder: index })),
    };
  }

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setIsSubmitting(true);
    setMessage("");
    try {
      const request = buildRequest();
      if (editingId) {
        await onUpdate(editingId, request);
        setMessage("Product updated.");
      } else {
        await onCreate(request);
        setCurrentPage(1);
        setMessage("Product created.");
      }
      resetForm();
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Could not save product");
    } finally {
      setIsSubmitting(false);
    }
  }

  async function handleDelete(product: Product) {
    if (!window.confirm(`Delete ${product.name}? This cannot be undone.`)) return;
    setMessage("");
    try {
      await onDelete(product.id);
      if (editingId === product.id) resetForm();
      setMessage("Product deleted.");
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Could not delete product");
    }
  }

  return (
    <section className="grid gap-6 xl:grid-cols-[440px_1fr]">
      <form onSubmit={(event) => void handleSubmit(event)} className="h-fit rounded-[2rem] border border-brand-border bg-white p-6 toy-shadow">
        <div className="flex items-start justify-between gap-3"><div><p className="text-xs font-black uppercase tracking-[0.18em] text-brand-orange">{editingId ? "Edit product" : "Add inventory"}</p><h2 className="font-display mt-1 text-3xl">{editingId ? "UPDATE PRODUCT" : "NEW PRODUCT"}</h2></div>{editingId && <button type="button" onClick={resetForm} className="grid h-9 w-9 place-items-center rounded-full bg-brand-surface" aria-label="Cancel editing"><X className="h-4 w-4" /></button>}</div>

        <div className="mt-6 grid gap-4 sm:grid-cols-2">
          <AdminField label="Name"><input required minLength={2} maxLength={150} value={form.name} onChange={(event) => updateField("name", event.target.value)} className="admin-input" /></AdminField>
          <AdminField label="SKU"><input required minLength={2} maxLength={80} value={form.sku} onChange={(event) => updateField("sku", event.target.value)} className="admin-input" /></AdminField>
          <AdminField label="Price"><input required min="0.01" step="0.01" type="number" value={form.price} onChange={(event) => updateField("price", event.target.value)} className="admin-input" placeholder="29.99" /></AdminField>
          <AdminField label="Discount price"><input min="0.01" step="0.01" type="number" value={form.discountPrice} onChange={(event) => updateField("discountPrice", event.target.value)} className="admin-input" placeholder="Optional" /></AdminField>
          <AdminField label="Stock"><input required min="0" max="100000" type="number" value={form.stockQuantity} onChange={(event) => updateField("stockQuantity", event.target.value)} className="admin-input" /></AdminField>
          <AdminField label="Status"><select value={form.status} onChange={(event) => updateField("status", event.target.value)} className="admin-input"><option value="active">Active</option><option value="draft">Draft</option><option value="out_of_stock">Out of stock</option><option value="archived">Archived</option></select></AdminField>
          <AdminField label="Category"><select required value={form.categoryId} onChange={(event) => updateField("categoryId", event.target.value)} className="admin-input"><option value="">Select category</option>{categories.map((category) => <option key={category.id} value={category.id}>{category.name}</option>)}</select></AdminField>
          <AdminField label="Brand"><select required value={form.brandId} onChange={(event) => updateField("brandId", event.target.value)} className="admin-input"><option value="">Select brand</option>{brands.map((brand) => <option key={brand.id} value={brand.id}>{brand.name}</option>)}</select></AdminField>
        </div>
        <AdminField label="Description" className="mt-4"><textarea rows={4} maxLength={5000} value={form.description} onChange={(event) => updateField("description", event.target.value)} className="admin-input resize-y" /></AdminField>
        <AdminField label="Image URLs - one per line" className="mt-4"><textarea rows={3} value={form.imageUrls} onChange={(event) => updateField("imageUrls", event.target.value)} className="admin-input resize-y font-mono text-xs" placeholder="https://example.com/toy.jpg" /></AdminField>
        <AdminField label="Specifications JSON" className="mt-4"><textarea rows={7} value={form.specifications} onChange={(event) => updateField("specifications", event.target.value)} className="admin-input resize-y font-mono text-xs" placeholder={'{"age": "6+", "pieces": 250}'} /></AdminField>
        <button disabled={isSubmitting || categories.length === 0 || brands.length === 0} className="mt-5 inline-flex w-full items-center justify-center gap-2 rounded-full bg-brand-ink px-6 py-3 text-xs font-black uppercase tracking-wider text-white transition hover:bg-brand-orange disabled:opacity-50"><Plus className="h-4 w-4" />{isSubmitting ? "Saving..." : editingId ? "Update product" : "Create product"}</button>
        {(categories.length === 0 || brands.length === 0) && <p className="mt-3 text-xs font-semibold text-brand-orange">Create at least one category and brand first.</p>}
        {message && <p className="mt-4 text-sm font-semibold text-brand-muted" role="status">{message}</p>}
      </form>

      <div className="overflow-hidden rounded-[2rem] border border-brand-border bg-white">
        <div className="border-b border-brand-border p-6"><div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-end"><div><p className="text-xs font-black uppercase tracking-[0.18em] text-brand-purple">Inventory</p><h2 className="font-display mt-1 text-3xl">PRODUCTS</h2></div><label className="flex items-center gap-3 rounded-xl bg-brand-surface px-4"><Search className="h-4 w-4 text-brand-subtle" /><span className="sr-only">Search admin products</span><input value={search} onChange={(event) => { setSearch(event.target.value); setCurrentPage(1); }} placeholder="Search inventory" className="w-full bg-transparent py-3 text-sm outline-none" /></label></div></div>
        <div className="divide-y divide-brand-border">
          {filteredProducts.length === 0 && <p className="p-8 text-center text-sm text-brand-muted">No products match this search.</p>}
          {paginatedProducts.map((product) => (
            <article key={product.id} className="flex items-center gap-4 px-5 py-5">
              <div className="relative h-16 w-16 shrink-0 overflow-hidden rounded-xl bg-brand-surface">{product.images[0]?.url ? <Image src={product.images[0].url} alt={product.name} fill className="object-cover" sizes="64px" /> : <div className="grid h-full place-items-center text-[10px] text-brand-subtle">No image</div>}</div>
              <div className="min-w-0 flex-1"><div className="flex flex-wrap items-center gap-2"><h3 className="truncate font-bold">{product.name}</h3><span className={`rounded-full px-2 py-1 text-[9px] font-black uppercase ${product.status === "active" ? "bg-brand-lime text-brand-lime-ink" : "bg-brand-surface text-brand-muted"}`}>{product.status.replaceAll("_", " ")}</span></div><p className="mt-1 text-xs text-brand-subtle">{product.sku} | {product.category.name} | {product.brand.name}</p><p className="mt-2 font-display text-lg text-brand-orange">{formatPriceFromCents(product.discountPriceCents || product.priceCents)} <span className="font-sans text-xs text-brand-muted">| Stock {product.stockQuantity}</span></p></div>
              <div className="flex shrink-0 gap-2"><button type="button" onClick={() => startEditing(product)} className="grid h-10 w-10 place-items-center rounded-full bg-brand-surface transition hover:bg-brand-yellow" aria-label={`Edit ${product.name}`}><Pencil className="h-4 w-4" /></button><button type="button" onClick={() => void handleDelete(product)} className="grid h-10 w-10 place-items-center rounded-full bg-brand-surface text-brand-orange transition hover:bg-brand-orange hover:text-white" aria-label={`Delete ${product.name}`}><Trash2 className="h-4 w-4" /></button></div>
            </article>
          ))}
        </div>
        <PaginationControls currentPage={safePage} totalPages={totalPages} onPageChange={setCurrentPage} className="border-t border-brand-border px-4 py-5" />
      </div>
    </section>
  );
}

function AdminField({ label, className = "", children }: { label: string; className?: string; children: React.ReactNode }) {
  return <label className={`block text-xs font-black uppercase tracking-wider text-brand-muted ${className}`}>{label}{children}</label>;
}
