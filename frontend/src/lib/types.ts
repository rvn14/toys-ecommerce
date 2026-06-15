export type Category = {
  id: number;
  name: string;
  slug: string;
  description: string;
};

export type Brand = {
  id: number;
  name: string;
  slug: string;
  description: string;
};

export type ProductImage = {
  id: number;
  url: string;
  storagePath: string;
  altText: string;
  sortOrder: number;
};

export type Product = {
  id: number;
  name: string;
  slug: string;
  sku: string;
  description: string;
  priceCents: number;
  discountPriceCents: number | null;
  stockQuantity: number;
  status: string;
  category: Category;
  brand: Brand;
  specifications: Record<string, unknown>;
  images: ProductImage[];
};

export type PaginationMeta = {
  page: number;
  limit: number;
  totalItems: number;
  totalPages: number;
  hasNext: boolean;
  hasPrevious: boolean;
};

export type ProductListData = {
  items: Product[];
  pagination: PaginationMeta;
};

export type ApiResponse<T> = {
  success: boolean;
  message: string;
  data?: T;
  errors?: unknown;
};

export type User = {
  id: number;
  name: string;
  email: string;
  role: string;
};

export type AuthData = {
  accessToken: string;
  user: User;
};

export type LoginRequest = {
  email: string;
  password: string;
};

export type RegisterRequest = {
  name: string;
  email: string;
  password: string;
};

export type CatalogEntityRequest = {
  name: string;
  description: string;
};

export type ProductImageInput = {
  url: string;
  storagePath: string;
  altText: string;
  sortOrder: number;
};

export type ProductRequest = {
  name: string;
  sku: string;
  description: string;
  priceCents: number;
  discountPriceCents: number | null;
  stockQuantity: number;
  status: string;
  categoryId: number;
  brandId: number;
  specifications: Record<string, unknown>;
  images: ProductImageInput[];
};

export type CartItem = {
  id: number;
  productId: number;
  productName: string;
  productSlug: string;
  productImageUrl: string;
  priceCents: number;
  discountPriceCents: number | null;
  effectivePriceCents: number;
  stockQuantity: number;
  quantity: number;
  lineTotalCents: number;
};

export type Cart = {
  id: number;
  userId: number;
  items: CartItem[];
  totalItems: number;
  subtotalCents: number;
};
