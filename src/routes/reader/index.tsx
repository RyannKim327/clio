import Book from "@/components/ui/book";
import Header from "@/components/widgets/header";
import Modal from "@/components/widgets/modal";
import { stories } from "@/data/stories";
import { useEffect, useState } from "react";
import { Link } from "react-router";

interface _chapter {
  content: string
  music?: string
  end?: boolean
}

interface _book {
  id?: number
  title: string
  author: string
  description: string
  cover?: string
  chapters: _chapter[][]
}

export default function Reader() {
  const [action, setAction] = useState("")
  const [book, setBook] = useState<_book | null>(null)

  function showModal(data: _book) {
    setBook(data)
    setAction("book")
  }

  function closeModal() {
    setBook(null)
    setAction("")
  }

  return (
    <div className="flex flex-col w-full h-full overflow-hidden">
      <Header setAction={setAction} />
      <div className="grid grid-cols-5 gap-2 m-2 w-[calc(90%-1rem)]">
        {
          stories.map((story: _book, i: number) => {
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
        <div className="flex flex-col items-center">
          <span>{book?.title}</span>
          <span className="text-xs">{book?.author}</span>
          <span className="border-l-2 border-l-solid border-l-border pl-3">{book?.description}</span>
          <Link
            className="bg-border/50 rounded-full w-full text-center p-2 mt-5 text-xs"

            to={`read/?id=${book?.id}`}>Read story</Link>
        </div>
      </Modal>
    </div>
  )
}
