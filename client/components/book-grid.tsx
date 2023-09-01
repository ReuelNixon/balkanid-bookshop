"use client";
import { useState } from "react";
import { BookCard } from "./book-card";
import { Input } from "./ui/input";
import { Label } from "./ui/label";

// <BookGrid bookList={data}/>

interface BookGridProps {
	bookList: any;
}

export function BookGrid({ bookList }: BookGridProps) {
	const [searchText, setSearchText] = useState("");

	const searchFilter = (bookList: any) => {
		return bookList.filter((book: any) =>
			book.book_title.toLowerCase().includes(searchText.toLowerCase())
		);
	};
	const filteredbookList = searchFilter(bookList);
	console.log(filteredbookList);

	return (
		<>
			<div>
				<h3 className="text-2xl py-6 text-center">
					Search For The Book!
				</h3>
				<div className="grid w-full max-w-sm items-center gap-1.5">
					<Label htmlFor="bookName">Book Name</Label>
					<Input
						type="text"
						value={searchText}
						autoComplete="off"
						id="bookName"
						placeholder="Type and press Enter"
						onChange={(e) => setSearchText(e.target.value)}
					/>
				</div>
				<h3 className="text-3xl pt-12 pb-6 text-center">
					Book Collection
				</h3>
			</div>

			<div className="mb-32 grid text-center md:grid-cols-2 lg:mb-0 lg:grid-cols-3 lg:text-left">
				{filteredbookList.map((book: any) => {
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
