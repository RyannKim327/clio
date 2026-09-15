import Book from "@/components/ui/book";
import Header from "@/components/widgets/header";
import { useEffect, useState } from "react";

export default function Reader() {
  const [action, setAction] = useState("")

  useEffect(() => { }, [action])

  return (
    <div className="flex flex-col bg-bg text-fg w-dvw h-dvh overflow-hidden">
      <Header setAction={setAction} />
      <div className="grid grid-cols-5 gap-2 m-2 w-[calc(90%-1rem)]">
        <Book title="sanaol santol" author="mang kaknorr" cover="https://wallpapers.com/images/high/iu-in-brown-coat-2b59nkxx8y8hpdes.webp" />
      </div>
    </div>
  )
}
