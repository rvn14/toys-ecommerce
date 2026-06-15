import { AdminPortal } from "@/components/admin/admin-portal";
import { AdminRoute } from "@/components/admin-route";

export default function AdminPage() {
  return (
    <AdminRoute>
      <AdminPortal />
    </AdminRoute>
  );
}
