import type { c_book } from "@/types/book"
import { Link } from "react-router"

interface b {
  book: c_book
}

export default function BookPreview({ book }: b) {
  return (
    <div className="flex flex-col items-center">
      <div className="flex flex-col gap-2 border-b border-b-solid border-b-border pb-2 mb-5 w-full items-center">
        <span>{book?.title}</span>
        <span className="text-xs">{book?.author}</span>
      </div>
      <span className="border-l-2 border-l-solid border-l-border pl-3">{book?.description}</span>
      <Link
        className="bg-border/50 rounded-sm w-full text-center p-2 mt-5 text-xs"
        to={`read/${book?.id}`}>Read story</Link>
    </div>
  )
}
