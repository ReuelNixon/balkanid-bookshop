"use client";

import { BookCard } from "@/components/book-card";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useCookies } from "react-cookie";

export default function Purchases() {
	const [books, setBooks] = useState([]);
	const [message, setMessage] = useState("");
	const [cookies, setCookie] = useCookies(["isLoggedIn"]);

	const router = useRouter();

	if (cookies.isLoggedIn != true) {
		router.push("/login");
	}
	async function getBooks() {
		const response = await fetch(
			"http://localhost:3000/api/book/private/purchases",
			{
				method: "GET",
				credentials: "include",
			}
		);
		const data = await response.json();
		if (data.data == null) {
			setMessage("You haven't purchased any books yet");
			setBooks([]);
		} else {
			setBooks(data.data);
		}
	}

	useEffect(() => {
		getBooks();
	}, []);

	return (
		<>
			{books.length > 0 ? (
				<div>
					<h1 className="text-2xl text-bold text-center pt-4">
						Your Purchases
					</h1>
					<div className="min-w-[200px] flex justify-between m-10">
						<Button variant="secondary" onClick={getBooks}>
							Refresh
						</Button>
						<Button variant="secondary" onClick={getBooks}>
							Download
						</Button>
					</div>
				</div>
			) : (
				<h1 className="text-2xl text-bold text-center pt-24">
					{message}
				</h1>
			)}
			<div className="mb-32 grid text-center md:grid-cols-2 lg:mb-0 lg:grid-cols-3 lg:text-left">
				{books.map((book: any) => {
					return (
						<BookCard
							name={book.book_title}
							id={book.ID.toString()}
							imgUrl={book.image_url}
							key={book.book_title + "Card"}
						/>
					);
				})}
			</div>
		</>
	);
}
