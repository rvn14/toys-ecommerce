import type {
  ApiResponse,
  Brand,
  CatalogEntityRequest,
  Category,
  Product,
  ProductListData,
  ProductRequest,
} from "@/lib/types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

if (!API_BASE_URL) {
  throw new Error("NEXT_PUBLIC_API_BASE_URL is not set");
}

async function request<T>(path: string, token: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers: {
      Authorization: `Bearer ${token}`,
      ...(init?.body ? { "Content-Type": "application/json" } : {}),
      ...init?.headers,
    },
  });
  const body = (await response.json()) as ApiResponse<T>;

  if (!response.ok || !body.success) {
    throw new Error(body.message || "Admin request failed");
  }

  return body.data as T;
}

export const adminApi = {
  getCategories: (token: string) => request<Category[]>("/categories", token),
  createCategory: (token: string, data: CatalogEntityRequest) => request<Category>("/admin/categories", token, { method: "POST", body: JSON.stringify(data) }),
  updateCategory: (token: string, id: number, data: CatalogEntityRequest) => request<Category>(`/admin/categories/${id}`, token, { method: "PUT", body: JSON.stringify(data) }),
  deleteCategory: (token: string, id: number) => request<void>(`/admin/categories/${id}`, token, { method: "DELETE" }),

  getBrands: (token: string) => request<Brand[]>("/brands", token),
  createBrand: (token: string, data: CatalogEntityRequest) => request<Brand>("/admin/brands", token, { method: "POST", body: JSON.stringify(data) }),
  updateBrand: (token: string, id: number, data: CatalogEntityRequest) => request<Brand>(`/admin/brands/${id}`, token, { method: "PUT", body: JSON.stringify(data) }),
  deleteBrand: (token: string, id: number) => request<void>(`/admin/brands/${id}`, token, { method: "DELETE" }),

  getProducts: (token: string) => request<ProductListData>("/admin/products?limit=100&sort=newest", token),
  createProduct: (token: string, data: ProductRequest) => request<Product>("/admin/products", token, { method: "POST", body: JSON.stringify(data) }),
  updateProduct: (token: string, id: number, data: ProductRequest) => request<Product>(`/admin/products/${id}`, token, { method: "PUT", body: JSON.stringify(data) }),
  deleteProduct: (token: string, id: number) => request<void>(`/admin/products/${id}`, token, { method: "DELETE" }),
};
