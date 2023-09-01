"use client";

import { BookGrid } from "@/components/book-grid";
import { useEffect, useState } from "react";

export default function Home() {
	const [books, setBooks] = useState([]);
	async function getBooks() {
		const response = await fetch("http://localhost:3000/api/book/");
		const data = await response.json();
		console.log(data);
		setBooks(data.data);
	}
	useEffect(() => {
		getBooks();
	}, []);
	return (
		<>
			<BookGrid bookList={books} />
			<div className="mb-32 grid text-center lg:max-w-5xl lg:w-full lg:mb-0 lg:grid-cols-4 lg:text-left mt-10"></div>
		</>
	);
}
