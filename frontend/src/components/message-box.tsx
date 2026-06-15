import type { ReactNode } from "react";
import { CircleAlert, PackageOpen, SearchX } from "lucide-react";

type MessageBoxProps = {
  eyebrow?: string;
  title: string;
  description: string;
  variant?: "default" | "empty" | "error";
  headingLevel?: "h1" | "h2";
  children?: ReactNode;
};

const variants = {
  default: {
    icon: CircleAlert,
    iconClassName: "bg-brand-purple text-white",
  },
  empty: {
    icon: SearchX,
    iconClassName: "bg-brand-yellow text-brand-ink",
  },
  error: {
    icon: PackageOpen,
    iconClassName: "bg-brand-orange text-white",
  },
};

export function MessageBox({
  eyebrow,
  title,
  description,
  variant = "default",
  headingLevel = "h1",
  children,
}: MessageBoxProps) {
  const { icon: Icon, iconClassName } = variants[variant];
  const Heading = headingLevel;

  return (
    <section className="toy-shadow mx-auto w-full max-w-2xl rounded-[2rem] border-2 border-dashed border-brand-border-strong bg-brand-cream p-8 text-center sm:p-12">
      <span className={`mx-auto grid h-16 w-16 place-items-center rounded-2xl ${iconClassName}`}>
        <Icon className="h-8 w-8" aria-hidden="true" />
      </span>
      {eyebrow && (
        <p className="mt-6 text-xs font-black uppercase tracking-[0.2em] text-brand-orange">
          {eyebrow}
        </p>
      )}
      <Heading className="font-display mt-2 text-4xl leading-none sm:text-5xl">{title}</Heading>
      <p className="mx-auto mt-4 max-w-lg text-sm leading-6 text-brand-muted sm:text-base">
        {description}
      </p>
      {children && <div className="mt-7 flex flex-wrap justify-center gap-3">{children}</div>}
    </section>
  );
}
