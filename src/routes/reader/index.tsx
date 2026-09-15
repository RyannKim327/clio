import Book from "@/components/ui/book";
import Header from "@/components/widgets/header";
import { stories } from "@/data/stories";
import { useEffect, useState } from "react";

export default function Reader() {
  const [action, setAction] = useState("")

  return (
    <div className="flex flex-col bg-bg text-fg w-dvw h-dvh overflow-hidden">
      <Header setAction={setAction} />
      <div className="grid grid-cols-5 gap-2 m-2 w-[calc(90%-1rem)]">
        {
          stories.map((story, i: number) => {
            return <Book key={`${i}. ${story.title}`} title={story.title} author={story.author} cover={story.cover} />
          })
        }
      </div>
    </div>
  )
}
