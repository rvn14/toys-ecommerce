"use client";

import Link from "next/link";
import { FormEvent, useState } from "react";
import { useAuth } from "@/components/auth-provider";

export default function RegisterPage() {
  const { register } = useAuth();

  const [name, setName] = useState("Test User");
  const [email, setEmail] = useState("newuser@example.com");
  const [password, setPassword] = useState("Password123!");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    setError("");
    setIsSubmitting(true);

    try {
      await register({
        name,
        email,
        password,
      });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Registration failed");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <main className="dot-pattern flex min-h-[820px] items-center justify-center bg-brand-pink px-4 py-14">
      <section className="toy-shadow w-full max-w-md rounded-[2rem] border border-white/70 bg-white p-8 sm:p-10">
        <div>
          <p className="text-xs font-black uppercase tracking-[0.2em] text-brand-orange">
            Create account
          </p>
          <h1 className="font-display mt-2 text-4xl text-brand-ink">
            Join Global Toys Store
          </h1>
          <p className="mt-2 text-sm text-brand-muted">
            Create an account to save your cart and place orders.
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
              Name
            </label>
            <input
              type="text"
              value={name}
              onChange={(event) => setName(event.target.value)}
              className="mt-2 w-full rounded-xl border border-brand-border bg-brand-cream px-4 py-3 text-sm outline-none focus:border-brand-orange"
              placeholder="Your name"
              required
            />
          </div>

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
              placeholder="Minimum 8 characters"
              required
            />
          </div>

          <button
            type="submit"
            disabled={isSubmitting}
            className="w-full rounded-full bg-brand-ink px-5 py-4 text-xs font-black uppercase tracking-wider text-white transition hover:bg-brand-orange disabled:cursor-not-allowed disabled:opacity-60"
          >
            {isSubmitting ? "Creating account..." : "Register"}
          </button>
        </form>

        <p className="mt-6 text-center text-sm text-brand-muted">
          Already have an account?{" "}
          <Link href="/login" className="font-bold text-brand-orange">
            Login
          </Link>
        </p>
      </section>
    </main>
  );
}
