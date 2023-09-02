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
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useCookies } from "react-cookie";

export default function Register() {
	const [cookies, setCookies] = useCookies(["isLoggedIn", "username"]);
	const Router = useRouter();
	const [error, setError] = useState("");
	const [username, setUsername] = useState("");
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");

	const handleSubmit = async (e: any) => {
		e.preventDefault();
		try {
			const response = await fetch(
				"http://localhost:3000/api/user/signup",
				{
					method: "POST",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify({ username, email, password }),
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
			setCookies("username", username);
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
							Register
						</CardTitle>
						<CardDescription className="text-center">
							Enter your email, password and username to register
						</CardDescription>
					</CardHeader>
					<CardContent className="grid gap-4">
						<div className="grid gap-2">
							<Label htmlFor="username">Username</Label>
							<Input
								id="username"
								type="text"
								onChange={(e) => setUsername(e.target.value)}
							/>
						</div>
						<div className="grid gap-2">
							<Label htmlFor="email">Email</Label>
							<Input
								id="email"
								type="email"
								placeholder=""
								onChange={(e) => setEmail(e.target.value)}
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
					<div className="text-red-500 text-center">{error}</div>
					<CardFooter className="flex flex-col">
						<Button className="w-full" onClick={handleSubmit}>
							Register
						</Button>
					</CardFooter>
				</Card>
			</div>
		</div>
	);
}
