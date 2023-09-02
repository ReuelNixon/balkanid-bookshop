"use client";
import { Button } from "@/components/ui/button";
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

export default function BookPage({ params }: ParamType) {
	const { bookID } = params;
	const [book, setBook] = useState<BookType>({} as BookType);
	const [reviews, setReviews] = useState([]);
	const [isInCart, setIsInCart] = useState(false);
	const [isPurchased, setIsPurchased] = useState(false);
	const [error, setError] = useState("");
	const [review, setReview] = useState("");
	const [rating, setRating] = useState(0);

	async function checkCart(ID: string) {
		const IDnum = parseInt(ID);
		const response = await fetch(
			"http://localhost:3000/api/book/private/cart",
			{
				method: "GET",
				credentials: "include",
			}
		);
		const data = await response.json();
		if (data.data == null) {
			setIsInCart(false);
		} else {
			data.data.forEach((book: any) => {
				if (book.ID === IDnum) {
					setIsInCart(true);
				}
			});
		}
	}

	async function checkPurchases(ID: string) {
		const IDnum = parseInt(ID);
		const response = await fetch(
			"http://localhost:3000/api/book/private/purchases",
			{
				method: "GET",
				credentials: "include",
			}
		);
		const data = await response.json();
		if (data.data == null) {
			setIsPurchased(false);
		} else {
			data.data.forEach((book: any) => {
				if (book.ID === IDnum) {
					setIsPurchased(true);
				}
			});
		}
	}

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

	async function addTocart(ID: string) {
		const response = await fetch(
			"http://localhost:3000/api/book/private/addToCart",
			{
				method: "POST",
				credentials: "include",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify({ book_id: parseInt(ID) }),
			}
		);
		const data = await response.json();
		if (data.error) {
			setError(data.error);
		} else {
			setError("");
			setIsInCart(true);
		}
	}

	async function checkout(ID: string) {
		const response = await fetch(
			"http://localhost:3000/api/book/private/" + ID + "/checkout",
			{
				method: "POST",
				credentials: "include",
			}
		);
		const data = await response.json();
		if (data.error) {
			setError(data.message);
		} else {
			setError("");
			setIsPurchased(true);
		}
	}

	async function postReview(ID: number, review: string, rating: number) {
		const response = await fetch(
			"http://localhost:3000/api/book/private/postReview",
			{
				method: "POST",
				credentials: "include",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify({ book_id: ID, review, rating }),
			}
		);
		const data = await response.json();
		if (data.error) {
			setError(data.error);
		} else {
			setError("");
			getReviews(bookID);
		}
	}

	useEffect(() => {
		getBook(bookID);
		getReviews(bookID);
		checkCart(bookID);
		checkPurchases(bookID);
	}, []);

	return (
		<>
			<div className="flex justify-around w-full max-w-3xl flex-wrap">
				<div className="py-8 max-w-sm justify-center">
					<h1 className="text-2xl text-bold pt-4 text-center">
						{book.book_title}
					</h1>
					<div className="m-4 relative w-[300px] h-[400px] flex justify-center">
						{book.image_url != null && (
							<img
								src={book.image_url}
								alt={"Picture"}
								style={{ objectFit: "contain" }}
							/>
						)}
					</div>
				</div>

				<div className="md:pt-28">
					<p className="text-lg text-bold pt-4">
						Author: {book.book_author}
					</p>
					<p className="text-lg text-bold pt-4">
						Publisher: {book.publisher}
					</p>
					<p className="text-lg text-bold pt-4">
						Year of Publication: {book.year_of_publication}
					</p>
					<p className="text-lg text-bold pt-4">ISBN: {book.isbn}</p>

					<div className="flex flex-col justify-around w-full max-w-3xl flex-wrap py-24">
						<div className="text-red-500 text-lg text-bold pt-4 text-center">
							{error}
						</div>
						{isPurchased ? (
							<Button variant="secondary" disabled>
								Already purchased
							</Button>
						) : (
							<div className="flex flex-col">
								{isInCart ? (
									<Button variant="secondary" disabled>
										Already in cart
									</Button>
								) : (
									<Button
										variant="secondary"
										onClick={() => addTocart(bookID)}
									>
										Add to cart
									</Button>
								)}
								<br />
								<Button
									variant="secondary"
									onClick={() => checkout(bookID)}
								>
									Buy Now
								</Button>
							</div>
						)}
					</div>
				</div>
			</div>

			{reviews.length > 0 && (
				<div className="flex flex-col justify-around w-full max-w-3xl flex-wrap">
					<h1 className="text-2xl text-bold pt-4 text-center">
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
			<div className="flex flex-col justify-around w-full max-w-3xl flex-wrap">
				<h3 className="text-2xl text-bold pt-4 text-center">
					Post a review
				</h3>
				<div className="flex flex-col justify-around w-full max-w-3xl flex-wrap">
					<div className="text-red-500 text-lg text-bold pt-4 text-center">
						{error}
					</div>
					<div className="flex flex-col justify-around w-full max-w-3xl flex-wrap">
						<label
							className="text-xl text-bold pt-4"
							htmlFor="review"
						>
							Review:
						</label>
						<textarea
							className="border-2 border-black rounded-lg p-4"
							name="review"
							id="review"
							rows={4}
							cols={50}
							onChange={(e) => {
								setError("");
								setReview(e.target.value);
							}}
						/>
					</div>
					<div className="flex flex-col justify-around w-full max-w-3xl flex-wrap">
						<label
							className="text-xl text-bold pt-4"
							htmlFor="rating"
						>
							Rating:
						</label>
						<input
							className="border-2 border-black rounded-lg"
							type="range"
							min="0"
							max="5"
							step="1"
							name="rating"
							id="rating"
							onChange={(e) => {
								setError("");
								setRating(parseInt(e.target.value));
							}}
						/>
					</div>
					<br />
					<div className="flex flex-col justify-around w-full max-w-3xl flex-wrap">
						<Button
							variant="secondary"
							onClick={() =>
								postReview(parseInt(bookID), review, rating)
							}
						>
							Post Review
						</Button>
					</div>
				</div>
			</div>
		</>
	);
}
