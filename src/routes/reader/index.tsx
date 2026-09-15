import Book from "@/components/ui/book";
import Header from "@/components/widgets/header";

export default function Reader() {
  return (
    <div className="flex flex-col bg-bg text-fg w-dvw h-dvh overflow-hidden">
      <Header />
      <div className="grid grid-cols-5 gap-2 w-[calc(90%-1rem)]">
        <Book title="sanaol santol" author="mang kaknorr" cover="https://wallpapers.com/images/high/iu-in-brown-coat-2b59nkxx8y8hpdes.webp" />
      </div>
    </div>
  )
}
