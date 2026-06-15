import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { Navbar } from "@/components/navbar";
import { AuthProvider } from "@/components/auth-provider";
import { Footer } from "@/components/footer";

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-inter",
  display: "swap",
});

export const metadata: Metadata = {
  metadataBase: new URL(process.env.NEXT_PUBLIC_SITE_URL || "https://toyjoy.store"),
  title: {
    default: "ToyJoy | Toys, Games & Collectibles",
    template: "%s | ToyJoy",
  },
  description: "Discover colorful toys, clever building sets, plush friends, games, gifts, and collectibles for every kind of imagination.",
  applicationName: "ToyJoy",
  keywords: ["toys", "games", "collectibles", "building sets", "plush toys", "gifts for kids"],
  authors: [{ name: "ToyJoy" }],
  creator: "ToyJoy",
  publisher: "ToyJoy",
  openGraph: {
    type: "website",
    siteName: "ToyJoy",
    title: "ToyJoy | Toys, Games & Collectibles",
    description: "Colorful toys, games, gifts, and collectibles for every kind of imagination.",
    url: "/",
    images: [{ url: "/images/toy-hero.png", alt: "ToyJoy robot and dinosaur toy mascots" }],
  },
  twitter: {
    card: "summary_large_image",
    title: "ToyJoy | Toys, Games & Collectibles",
    description: "Colorful toys, games, gifts, and collectibles for every kind of imagination.",
    images: ["/images/toy-hero.png"],
  },
  robots: { index: true, follow: true },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        suppressHydrationWarning
        className={`${inter.variable} font-sans antialiased`}
      >
        <AuthProvider>
          <Navbar />
          {children}
          <Footer />
        </AuthProvider>
      </body>
    </html>
  );
}
