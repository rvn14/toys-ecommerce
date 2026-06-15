"use client";

import Link from "next/link";
import { FormEvent, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/components/auth-provider";

export default function LoginPage() {
  const router = useRouter();
  const { login, isAuthenticated, isLoading, user } = useAuth();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      router.replace(user?.role === "admin" ? "/admin" : "/account");
    }
  }, [isAuthenticated, isLoading, router, user?.role]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    setError("");
    setIsSubmitting(true);

    try {
      await login({
        email,
        password,
      });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <main className="dot-pattern flex min-h-[720px] items-center justify-center bg-brand-yellow px-4 py-14">
      <section className="toy-shadow w-full max-w-md rounded-[2rem] border border-white/70 bg-white p-8 sm:p-10">
        <div>
          <p className="text-xs font-black uppercase tracking-[0.2em] text-brand-orange">
            Welcome back
          </p>
          <h1 className="font-display mt-2 text-4xl text-brand-ink">
            Login to your account
          </h1>
          <p className="mt-2 text-sm text-brand-muted">
            Access your cart, orders, and account details.
          </p>
        </div>

        {error && (
          <div className="mt-6 rounded-2xl border border-destructive/20 bg-destructive/10 p-4 text-sm text-destructive">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="mt-8 space-y-5">
          <div>
            <label className="text-sm font-medium text-brand-muted">
              Email
            </label>
            <input
              type="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              className="mt-2 w-full rounded-xl border border-brand-border bg-brand-cream px-4 py-3 text-sm outline-none focus:border-brand-orange"
              placeholder="you@example.com"
              required
            />
          </div>

          <div>
            <label className="text-sm font-medium text-brand-muted">
              Password
            </label>
            <input
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              className="mt-2 w-full rounded-xl border border-brand-border bg-brand-cream px-4 py-3 text-sm outline-none focus:border-brand-orange"
              placeholder="Your password"
              required
            />
          </div>

          <button
            type="submit"
            disabled={isSubmitting}
            className="w-full rounded-full bg-brand-ink px-5 py-4 text-xs font-black uppercase tracking-wider text-white transition hover:bg-brand-orange disabled:cursor-not-allowed disabled:opacity-60"
          >
            {isSubmitting ? "Logging in..." : "Login"}
          </button>
        </form>

        <p className="mt-6 text-center text-sm text-brand-muted">
          Do not have an account?{" "}
          <Link href="/register" className="font-bold text-brand-orange">
            Register
          </Link>
        </p>
      </section>
    </main>
  );
}
