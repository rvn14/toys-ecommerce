import { cache } from "react";
import type { ApiResponse, Brand, Category, Product, ProductListData } from "@/lib/types";

const API_BASE_URL =
  process.env.INTERNAL_API_BASE_URL || process.env.NEXT_PUBLIC_API_BASE_URL;

if (!API_BASE_URL) {
  throw new Error(
    "INTERNAL_API_BASE_URL or NEXT_PUBLIC_API_BASE_URL must be set"
  );
}

export async function getProducts(params?: {
  page?: number;
  limit?: number;
  search?: string;
  category?: string;
  brand?: string;
  sort?: string;
}): Promise<ProductListData> {
  const searchParams = new URLSearchParams();

  if (params?.page) searchParams.set("page", String(params.page));
  if (params?.limit) searchParams.set("limit", String(params.limit));
  if (params?.search) searchParams.set("search", params.search);
  if (params?.category) searchParams.set("category", params.category);
  if (params?.brand) searchParams.set("brand", params.brand);
  if (params?.sort) searchParams.set("sort", params.sort);

  const queryString = searchParams.toString();

  const response = await fetch(
    `${API_BASE_URL}/products${queryString ? `?${queryString}` : ""}`,
    {
      cache: "no-store",
    }
  );

  if (!response.ok) {
    throw new Error("Failed to fetch products");
  }

  const body = (await response.json()) as ApiResponse<ProductListData>;

  if (!body.success || !body.data) {
    throw new Error(body.message || "Failed to fetch products");
  }

  return body.data;
}

async function getCatalogCollection<T>(path: string): Promise<T[]> {
  const response = await fetch(`${API_BASE_URL}/${path}`, {
    cache: "no-store",
  });

  if (!response.ok) {
    throw new Error(`Failed to fetch ${path}`);
  }

  const body = (await response.json()) as ApiResponse<T[]>;

  if (!body.success || !body.data) {
    throw new Error(body.message || `Failed to fetch ${path}`);
  }

  return body.data;
}

export function getCategories() {
  return getCatalogCollection<Category>("categories");
}

export function getBrands() {
  return getCatalogCollection<Brand>("brands");
}

export const getProduct = cache(async (id: number): Promise<Product | null> => {
  const response = await fetch(`${API_BASE_URL}/products/${id}`, {
    cache: "no-store",
  });

  if (response.status === 404) {
    return null;
  }

  if (!response.ok) {
    throw new Error("Failed to fetch product");
  }

  const body = (await response.json()) as ApiResponse<Product>;

  if (!body.success || !body.data) {
    throw new Error(body.message || "Failed to fetch product");
  }

  return body.data;
});
