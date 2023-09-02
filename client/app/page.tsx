"use client";

import { BookGrid } from "@/components/book-grid";
import { Button } from "@/components/ui/button";
import { useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";

export default function Home() {
	const searchParams = useSearchParams();
	const pageSizeParam = searchParams.get("pageSize") || "12";
	const pageParam = searchParams.get("page") || "1";

	const [pageSize, setPageSize] = useState(pageSizeParam);
	const [page, setPage] = useState(pageParam);

	const [pages, setPages] = useState<number[]>([]);
	function setPagesArray(start: number) {
		setPages([]);
		let temp = [1];
		if (start > 5) {
			start -= 4;
		} else if (start > 2) {
			start -= 1;
		} else if (start == 1) {
			start += 1;
		}
		for (let i = start; i < start + 10; i++) {
			temp.push(i);
		}
		setPages(temp);
	}

	const [books, setBooks] = useState([]);
	async function getBooks(page: string, pageSize: string) {
		setBooks([]);
		const response = await fetch(
			`http://localhost:3000/api/book?page=${page}&pageSize=${pageSize}`,
			{
				method: "GET",
				credentials: "include",
			}
		);
		const data = await response.json();
		setBooks(data.data);
	}

	useEffect(() => {
		getBooks(page, pageSize);
		setPagesArray(parseInt(page));
	}, []);
	return (
		<>
			<BookGrid bookList={books} />
			<div className="flex justify-center space-x-5 pt-10">
				{pages.map((i) =>
					i == parseInt(page) ? (
						<Button
							key={i}
							variant="default"
							onClick={() => {
								getBooks(i.toString(), pageSize);
								setPagesArray(i);
							}}
							disabled
						>
							{i}
						</Button>
					) : (
						<Button
							key={i}
							variant="default"
							onClick={() => {
								setPage(i.toString());
								getBooks(i.toString(), pageSize);
								setPagesArray(i);
							}}
						>
							{i}
						</Button>
					)
				)}
			</div>
		</>
	);
}
