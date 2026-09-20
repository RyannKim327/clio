import BookPreview from "@/components/modals/book-preview";
import Book from "@/components/ui/book";
import Header from "@/components/widgets/header";
import Modal from "@/components/widgets/modal";
import { stories } from "@/data/stories";
import type { c_book } from "@/types/book";
import { useState } from "react";
import { Link } from "react-router";

export default function Reader() {
  const [action, setAction] = useState("")
  const [book, setBook] = useState<c_book | null>(null)

  function showModal(data: c_book) {
    setBook(data)
    setAction("book")
  }

  function closeModal() {
    setBook(null)
    setAction("")
  }

  return (
    <div className="flex flex-col w-full h-full overflow-hidden">
      <Header setAction={setAction} action={action} />
      <div className="grid grid-cols-5 gap-2 m-2 w-[calc(90%-1rem)]">
        {
          stories.map((story: c_book, i: number) => {
            return <Book
              key={`${i}. ${story.title}`}
              title={story.title}
              author={story.author}
              cover={story.cover}
              setAction={() => {
                showModal(({
                  ...story,
                  id: i
                }))
              }} />
          })
        }
      </div>
      <Modal closeModal={closeModal} show={action === "book" && book !== null}>
        <BookPreview book={book} />
      </Modal>
    </div>
  )
}
