"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Check, Minus, Plus, ShoppingBag } from "lucide-react";
import { useAuth } from "@/components/auth-provider";
import { addToCart } from "@/lib/cart-api";

type AddToCartProps = {
  productId: number;
  stockQuantity: number;
};

export function AddToCart({ productId, stockQuantity }: AddToCartProps) {
  const router = useRouter();
  const { token } = useAuth();
  const [quantity, setQuantity] = useState(1);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [message, setMessage] = useState("");
  const [isSuccess, setIsSuccess] = useState(false);

  async function handleAddToCart() {
    if (!token) {
      router.push("/login");
      return;
    }

    setIsSubmitting(true);
    setMessage("");
    setIsSuccess(false);

    try {
      const cart = await addToCart(token, productId, quantity);
      setMessage(`${quantity} item${quantity === 1 ? "" : "s"} added to your cart.`);
      setIsSuccess(true);
      window.dispatchEvent(new CustomEvent("cart-updated", { detail: cart.totalItems }));
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Could not add this item to your cart");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <div className="mt-6">
      <div className="flex flex-col gap-3 sm:flex-row">
        <div className="flex h-14 items-center justify-between rounded-full border-2 border-brand-border bg-white px-2 sm:w-36">
          <button type="button" onClick={() => setQuantity((current) => Math.max(1, current - 1))} className="grid h-10 w-10 place-items-center rounded-full hover:bg-brand-surface" aria-label="Decrease quantity"><Minus className="h-4 w-4" /></button>
          <span className="font-bold" aria-live="polite">{quantity}</span>
          <button type="button" onClick={() => setQuantity((current) => Math.min(stockQuantity, current + 1))} className="grid h-10 w-10 place-items-center rounded-full hover:bg-brand-surface" aria-label="Increase quantity"><Plus className="h-4 w-4" /></button>
        </div>
        <button
          type="button"
          disabled={stockQuantity < 1 || isSubmitting}
          onClick={() => void handleAddToCart()}
          className="inline-flex min-h-14 flex-1 items-center justify-center gap-2 rounded-full bg-brand-ink px-8 py-4 text-sm font-black uppercase tracking-wider text-white transition hover:bg-brand-orange disabled:cursor-not-allowed disabled:opacity-50 sm:flex-none"
        >
          {isSuccess ? <Check className="h-5 w-5" /> : <ShoppingBag className="h-5 w-5" />}
          {isSubmitting ? "Adding..." : isSuccess ? "Added" : "Add to cart"}
        </button>
      </div>
      {message && <p className={`mt-3 text-sm font-semibold ${isSuccess ? "text-brand-success" : "text-destructive"}`} role="status">{message}</p>}
    </div>
  );
}
