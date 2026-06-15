import type { ApiResponse, Cart } from "@/lib/types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

if (!API_BASE_URL) {
  throw new Error("NEXT_PUBLIC_API_BASE_URL is not set");
}

export async function addToCart(token: string, productId: number, quantity: number): Promise<Cart> {
  const response = await fetch(`${API_BASE_URL}/cart/items`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ productId, quantity }),
  });
  const body = (await response.json()) as ApiResponse<Cart>;

  if (!response.ok || !body.success || !body.data) {
    throw new Error(body.message || "Could not add this product to your cart");
  }

  return body.data;
}

export async function getCart(token: string): Promise<Cart> {
  const response = await fetch(`${API_BASE_URL}/cart`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  const body = (await response.json()) as ApiResponse<Cart>;

  if (!response.ok || !body.success || !body.data) {
    throw new Error(body.message || "Could not load your cart");
  }

  return body.data;
}
