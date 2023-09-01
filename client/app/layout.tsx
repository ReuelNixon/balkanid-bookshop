"use client";

import Navbar from "@/components/navbar";
import { ThemeProvider } from "@/components/theme-provider";
import { Inter } from "next/font/google";
import { CookiesProvider } from "react-cookie";
import "./globals.css";
const inter = Inter({ subsets: ["latin"] });

export default function RootLayout({
	children,
}: {
	children: React.ReactNode;
}) {
	return (
		<html lang="en">
			<body className={inter.className}>
				<CookiesProvider>
					<ThemeProvider attribute="class" defaultTheme="dark">
						<main className="flex min-h-screen flex-col items-center p-10">
							<Navbar />
							{children}
						</main>
					</ThemeProvider>
				</CookiesProvider>
			</body>
		</html>
	);
}
