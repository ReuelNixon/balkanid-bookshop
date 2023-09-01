"use client";
import { useEffect, useState } from "react";

type ParamType = {
	params: {
		bookID: string;
	};
};

type BookType = {
	book_title: string;
	book_author: string;
	year_of_publication: string;
	ID: number;
	image_url: string;
	publisher: string;
	isbn: string;
};

export default function PokemonPage({ params }: ParamType) {
	const { bookID } = params;
	const [book, setBook] = useState<BookType>({} as BookType);
	const [reviews, setReviews] = useState([]);

	async function getBook(ID: string) {
		const response = await fetch("http://localhost:3000/api/book/" + ID);
		const data = await response.json();
		setBook(data.data);
	}

	async function getReviews(ID: string) {
		const response = await fetch(
			"http://localhost:3000/api/book/" + ID + "/reviews/"
		);
		const data = await response.json();
		setReviews(data.data);
	}

	useEffect(() => {
		getBook(bookID);
		getReviews(bookID);
	}, []);

	return (
		<>
			<div className="flex justify-around w-full max-w-3xl flex-wrap">
				<div className="py-8">
					<h1 className="text-4xl text-bold pt-4 text-center">
						{book.book_title}
					</h1>
					<div
						className="m-4"
						style={{
							position: "relative",
							width: "300px",
						}}
					>
						{book.image_url != null && (
							<img
								src={book.image_url}
								alt={"Picture"}
								style={{ objectFit: "contain" }}
							/>
						)}
					</div>
				</div>

				<div className="md:py-16">
					<p className="text-xl text-bold pt-4">
						Author: {book.book_author}
					</p>
					<p className="text-xl text-bold pt-4">
						Publisher: {book.publisher}
					</p>
					<p className="text-xl text-bold pt-4">
						Year of Publication: {book.year_of_publication}
					</p>
					<p className="text-xl text-bold pt-4">ISBN: {book.isbn}</p>
				</div>
			</div>

			{reviews.length > 0 && (
				<div className="flex flex-col justify-around w-full max-w-3xl flex-wrap">
					<h1 className="text-4xl text-bold pt-4 text-center">
						Reviews
					</h1>

					{reviews.map((review: any) => {
						return (
							<div className="py-8">
								<p className="text-xl text-bold pt-4">
									Posted by: {review.user_id}
								</p>
								<p className="text-xl text-bold pt-4">
									Review: {review.review}
								</p>
								<p className="text-xl text-bold pt-4">
									Rating: {review.rating}
								</p>
							</div>
						);
					})}
				</div>
			)}
		</>
	);
}
