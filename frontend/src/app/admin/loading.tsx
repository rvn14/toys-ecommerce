import { PageLoader } from "@/components/page-loader";

export default function AdminLoading() {
  return (
    <PageLoader
      label="Admin workspace"
      message="Preparing the control room..."
    />
  );
}
