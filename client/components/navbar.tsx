"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { useCookies } from "react-cookie";
import { Button } from "./ui/button";

export default function Navbar() {
	const [cookies, setCookie] = useCookies(["isLoggedIn", "username"]);
	// const isLoggedIn = cookies.isLoggedIn == true;
	const [isLoggedIn, setIsLoggedIn] = useState(false);

	async function logout() {
		setCookie("isLoggedIn", undefined, {
			path: "/",
		});
		setCookie("username", undefined, {
			path: "/",
		});

		const response = await fetch(
			"http://localhost:3000/api/user/private/logout",
			{
				method: "POST",
				headers: { "Content-Type": "application/json" },
				credentials: "include",
			}
		);
		const data = response.json();
	}

	useEffect(() => {
		if (cookies.isLoggedIn == true) {
			setIsLoggedIn(true);
		}
	}, []);
	return (
		<div className="z-10 p-8 max-w-5xl w-full items-center justify-between text-sm lg:p-0 lg:flex">
			<div className="items-center min-w-full justify-between text-xl fixed left-0 top-0 flex w-full border-b border-gray-300 bg-gradient-to-b from-zinc-200 pb-6 pt-8 backdrop-blur-2xl dark:border-neutral-800 dark:bg-zinc-800/30 dark:from-inherit pl-10  lg:static lg:w-auto lg:p-4 lg:rounded-xl lg:border lg:bg-gray-200 lg:dark:bg-zinc-800/30">
				<Link href={"/"}>
					<div>BalkanID Bookshop</div>
				</Link>

				{isLoggedIn ? (
					<div className="flex justify-between pr-6 lg:pr-0">
						<Link href="/cart">
							<Button variant="secondary">Cart</Button>
						</Link>
						<p>.</p>
						<Link href="/purchases">
							<Button variant="secondary">Purchases</Button>
						</Link>
						<p>.</p>
						<Button variant="secondary" onClick={logout}>
							Logout
						</Button>
					</div>
				) : (
					<div className="flex justify-between pr-6 lg:pr-0">
						<Link href="/login">
							<Button variant="secondary">Login</Button>
						</Link>
						<p>.</p>
						<Link href="/register">
							<Button variant="secondary">Register</Button>
						</Link>
					</div>
				)}
			</div>
		</div>
	);
}
