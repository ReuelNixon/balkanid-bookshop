"use client";

import { Button } from "@/components/ui/button";
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useCookies } from "react-cookie";

export default function Login() {
	const Router = useRouter();
	const [error, setError] = useState("");
	const [identity, setIdentity] = useState("");
	const [password, setPassword] = useState("");
	const [cookies, setCookies] = useCookies(["isLoggedIn", "username"]);

	const handleSubmit = async (e: any) => {
		e.preventDefault();
		try {
			const response = await fetch(
				"http://localhost:3000/api/user/signin",
				{
					method: "POST",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify({ identity, password }),
					credentials: "include",
				}
			);
			const data = await response.json();
			if (data.error) {
				let message =
					data.username ||
					data.email ||
					data.password ||
					data.general;
				throw new Error(message);
			}
			setError("");
			setCookies("isLoggedIn", true);
			setCookies("username", identity);
			Router.push("/");
		} catch (error: any) {
			setError(error.message);
		}
	};
	return (
		<div className="relative flex flex-col justify-center items-center pt-24 overflow-hidden">
			<div className="w-full m-auto bg-white lg:max-w-lg">
				<Card>
					<CardHeader className="space-y-1">
						<CardTitle className="text-2xl text-center">
							Sign in
						</CardTitle>
						<CardDescription className="text-center">
							Enter your email and password to login
						</CardDescription>
					</CardHeader>
					<CardContent className="grid gap-4">
						<div className="grid gap-2">
							<Label htmlFor="identity">Username or Email</Label>
							<Input
								id="identity"
								type="text"
								placeholder=""
								onChange={(e) => setIdentity(e.target.value)}
							/>
						</div>
						<div className="grid gap-2">
							<Label htmlFor="password">Password</Label>
							<Input
								id="password"
								type="password"
								onChange={(e) => setPassword(e.target.value)}
							/>
						</div>
					</CardContent>
					<div className="justify-center text-red-500 text-center">
						{error}
					</div>
					<CardFooter className="flex flex-col">
						<Button className="w-full" onClick={handleSubmit}>
							Sign in
						</Button>
						<p className="mt-2 text-xs text-center text-gray-700">
							{" "}
							Don't have an account?{" "}
							<Link
								href="/register"
								className="text-blue-600 hover:underline"
							>
								Register
							</Link>
						</p>
					</CardFooter>
				</Card>
			</div>
		</div>
	);
}
