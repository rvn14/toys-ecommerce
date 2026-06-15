"use client";

import Link from "next/link";
import { ArrowRight, LogOut, Package, ShoppingBag, UserRound } from "lucide-react";
import { ProtectedRoute } from "@/components/protected-route";
import { useAuth } from "@/components/auth-provider";

export default function AccountPage() {
  const { user, logout } = useAuth();

  return (
    <ProtectedRoute>
      <main className="bg-brand-cream">
        <section className="dot-pattern border-y border-brand-border bg-brand-lime py-14"><div className="site-shell"><p className="text-xs font-black uppercase tracking-[0.2em]">Welcome back</p><h1 className="font-display mt-2 text-5xl sm:text-7xl">HEY, {user?.name?.split(" ")[0]?.toUpperCase()}!</h1><p className="mt-3 text-sm font-semibold text-brand-lime-ink">Your ToyJoy account is ready for another adventure.</p></div></section>

        <section className="site-shell grid gap-6 py-12 lg:grid-cols-[1fr_380px]">
          <div className="rounded-[2rem] border border-brand-border bg-white p-7 toy-shadow sm:p-9"><div className="flex items-center gap-4"><span className="grid h-14 w-14 place-items-center rounded-2xl bg-brand-cream text-brand-orange"><UserRound className="h-7 w-7" /></span><div><p className="text-xs font-black uppercase tracking-[0.18em] text-brand-subtle">Profile details</p><h2 className="font-display mt-1 text-3xl">YOUR ACCOUNT</h2></div></div><dl className="mt-8 grid gap-4 sm:grid-cols-2"><div className="rounded-2xl bg-brand-surface p-5"><dt className="text-[10px] font-black uppercase tracking-wider text-brand-subtle">Name</dt><dd className="mt-2 font-bold">{user?.name}</dd></div><div className="rounded-2xl bg-brand-surface p-5"><dt className="text-[10px] font-black uppercase tracking-wider text-brand-subtle">Email</dt><dd className="mt-2 break-all font-bold">{user?.email}</dd></div><div className="rounded-2xl bg-brand-surface p-5"><dt className="text-[10px] font-black uppercase tracking-wider text-brand-subtle">Membership</dt><dd className="mt-2 font-bold capitalize">{user?.role}</dd></div><div className="rounded-2xl bg-brand-surface p-5"><dt className="text-[10px] font-black uppercase tracking-wider text-brand-subtle">Status</dt><dd className="mt-2 font-bold text-brand-success">Ready to play</dd></div></dl></div>

          <aside className="rounded-[2rem] bg-brand-purple p-7 text-white toy-shadow"><p className="text-xs font-black uppercase tracking-[0.18em] text-brand-yellow">Quick actions</p><h2 className="font-display mt-2 text-3xl">WHAT&apos;S NEXT?</h2><div className="mt-7 grid gap-3"><Link href="/products" className="flex items-center justify-between rounded-2xl bg-white px-5 py-4 text-sm font-black text-brand-ink"><span className="flex items-center gap-3"><ShoppingBag className="h-5 w-5 text-brand-orange" /> Continue shopping</span><ArrowRight className="h-4 w-4" /></Link><div className="flex items-center gap-3 rounded-2xl bg-white/10 px-5 py-4 text-sm font-bold"><Package className="h-5 w-5 text-brand-yellow" /> Order history coming soon</div><button type="button" onClick={() => void logout()} className="flex items-center gap-3 rounded-2xl border border-white/30 px-5 py-4 text-left text-sm font-bold"><LogOut className="h-5 w-5" /> Logout</button></div></aside>
        </section>
      </main>
    </ProtectedRoute>
  );
}
