"use client";

import { useState } from "react";
import { Pencil, Plus, Trash2, X } from "lucide-react";
import { PaginationControls } from "@/components/pagination-controls";
import type { Brand, CatalogEntityRequest, Category } from "@/lib/types";

type Entity = Category | Brand;

type CatalogEntityManagerProps<T extends Entity> = {
  title: string;
  singular: string;
  items: T[];
  onCreate: (data: CatalogEntityRequest) => Promise<void>;
  onUpdate: (id: number, data: CatalogEntityRequest) => Promise<void>;
  onDelete: (id: number) => Promise<void>;
};

export function CatalogEntityManager<T extends Entity>({ title, singular, items, onCreate, onUpdate, onDelete }: CatalogEntityManagerProps<T>) {
  const itemsPerPage = 6;
  const [editingId, setEditingId] = useState<number | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [message, setMessage] = useState("");
  const totalPages = Math.max(1, Math.ceil(items.length / itemsPerPage));
  const safePage = Math.min(currentPage, totalPages);
  const paginatedItems = items.slice((safePage - 1) * itemsPerPage, safePage * itemsPerPage);

  function resetForm() {
    setEditingId(null);
    setName("");
    setDescription("");
  }

  function startEditing(item: T) {
    setEditingId(item.id);
    setName(item.name);
    setDescription(item.description);
    setMessage("");
  }

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setIsSubmitting(true);
    setMessage("");
    try {
      const data = { name: name.trim(), description: description.trim() };
      if (editingId) {
        await onUpdate(editingId, data);
        setMessage(`${singular} updated.`);
      } else {
        await onCreate(data);
        setMessage(`${singular} created.`);
      }
      resetForm();
    } catch (error) {
      setMessage(error instanceof Error ? error.message : `Could not save ${singular.toLowerCase()}`);
    } finally {
      setIsSubmitting(false);
    }
  }

  async function handleDelete(item: T) {
    if (!window.confirm(`Delete ${item.name}? This cannot be undone.`)) return;
    setMessage("");
    try {
      await onDelete(item.id);
      if (editingId === item.id) resetForm();
      setMessage(`${singular} deleted.`);
    } catch (error) {
      setMessage(error instanceof Error ? error.message : `Could not delete ${singular.toLowerCase()}`);
    }
  }

  return (
    <section className="grid gap-6 lg:grid-cols-[360px_1fr]">
      <form onSubmit={(event) => void handleSubmit(event)} className="h-fit rounded-[2rem] border border-brand-border bg-white p-6 toy-shadow">
        <div className="flex items-start justify-between gap-3">
          <div><p className="text-xs font-black uppercase tracking-[0.18em] text-brand-orange">{editingId ? "Edit" : "Create"}</p><h2 className="font-display mt-1 text-3xl">{editingId ? `UPDATE ${singular.toUpperCase()}` : `NEW ${singular.toUpperCase()}`}</h2></div>
          {editingId && <button type="button" onClick={resetForm} className="grid h-9 w-9 place-items-center rounded-full bg-brand-surface" aria-label="Cancel editing"><X className="h-4 w-4" /></button>}
        </div>
        <label className="mt-6 block text-xs font-black uppercase tracking-wider text-brand-muted">Name<input required minLength={2} maxLength={100} value={name} onChange={(event) => setName(event.target.value)} className="mt-2 w-full rounded-xl border border-brand-border bg-brand-cream px-4 py-3 text-sm font-medium normal-case tracking-normal outline-none focus:border-brand-orange" /></label>
        <label className="mt-4 block text-xs font-black uppercase tracking-wider text-brand-muted">Description<textarea value={description} onChange={(event) => setDescription(event.target.value)} rows={4} className="mt-2 w-full resize-y rounded-xl border border-brand-border bg-brand-cream px-4 py-3 text-sm font-medium normal-case tracking-normal outline-none focus:border-brand-orange" /></label>
        <button disabled={isSubmitting} className="mt-5 inline-flex w-full items-center justify-center gap-2 rounded-full bg-brand-ink px-6 py-3 text-xs font-black uppercase tracking-wider text-white transition hover:bg-brand-orange disabled:opacity-60"><Plus className="h-4 w-4" />{isSubmitting ? "Saving..." : editingId ? `Update ${singular}` : `Create ${singular}`}</button>
        {message && <p className="mt-4 text-sm font-semibold text-brand-muted" role="status">{message}</p>}
      </form>

      <div className="overflow-hidden rounded-[2rem] border border-brand-border bg-white">
        <div className="flex items-center justify-between border-b border-brand-border px-6 py-5"><div><p className="text-xs font-black uppercase tracking-[0.18em] text-brand-purple">Catalog</p><h2 className="font-display mt-1 text-3xl">{title.toUpperCase()}</h2></div><span className="rounded-full bg-brand-yellow px-3 py-1 text-xs font-black">{items.length}</span></div>
        <div className="divide-y divide-brand-border">
          {items.length === 0 && <p className="p-8 text-center text-sm text-brand-muted">No {title.toLowerCase()} yet.</p>}
          {paginatedItems.map((item) => (
            <article key={item.id} className="flex items-center justify-between gap-4 px-6 py-5">
              <div className="min-w-0"><h3 className="font-bold">{item.name}</h3><p className="mt-1 text-xs text-brand-subtle">/{item.slug}</p>{item.description && <p className="mt-2 line-clamp-2 text-sm text-brand-muted">{item.description}</p>}</div>
              <div className="flex shrink-0 gap-2"><button type="button" onClick={() => startEditing(item)} className="grid h-10 w-10 place-items-center rounded-full bg-brand-surface transition hover:bg-brand-yellow" aria-label={`Edit ${item.name}`}><Pencil className="h-4 w-4" /></button><button type="button" onClick={() => void handleDelete(item)} className="grid h-10 w-10 place-items-center rounded-full bg-brand-surface text-brand-orange transition hover:bg-brand-orange hover:text-white" aria-label={`Delete ${item.name}`}><Trash2 className="h-4 w-4" /></button></div>
            </article>
          ))}
        </div>
        <PaginationControls currentPage={safePage} totalPages={totalPages} onPageChange={setCurrentPage} className="border-t border-brand-border px-4 py-5" />
      </div>
    </section>
  );
}
