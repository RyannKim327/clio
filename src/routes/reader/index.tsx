import BookPreview from "@/components/modals/book-preview";
import Book from "@/components/ui/book";
import Header from "@/components/widgets/header";
import Modal from "@/components/widgets/modal";
import { stories } from "@/data/stories";
import type { c_book } from "@/types/book";
import { useState } from "react";

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
      <div className="flex flex-col h-full w-full px-5">
        <h1 className="text-center text-lg border-b-border border-b-solid border-b mb-10 pb-2 pt-5">Recent Stories</h1>
        <div className="grid grid-cols-5 gap-8 m-2 w-[calc(90%-1rem)] overflow-hidden overflow-y-auto">
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
      </div>
      <Modal closeModal={closeModal} show={action === "book" && book !== null}>
        <BookPreview book={book} />
      </Modal>
    </div>
  )
}
