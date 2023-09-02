import Link from "next/link";

// harry, -> localhost:3000/harry

interface BookCardProps {
	name: string;
	id: string;
	imgUrl: string;
}

// <BookCard name="harry" />

export function BookCard({ name, id, imgUrl }: BookCardProps) {
	return (
		<Link
			href={`/${id}`}
			className="group rounded-lg border border-transparent w-[270px] m-3 px-5 py-4 dark:border-gray-500 hover:border-gray-300 hover:bg-gray-100 hover:dark:border-neutral-700 hover:dark:bg-neutral-800/30 transition-colors"
			key={name + "Card"}
		>
			<h2 className="text-center">{name}</h2>

			<img
				src={imgUrl}
				alt={"Image"}
				width={200}
				height={200}
				className="mx-auto my-auto p-4"
			/>
		</Link>
	);
}
