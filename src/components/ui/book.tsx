interface book {
  title: string
  author: string
  cover?: string
}

export default function Book({
  title,
  author,
  cover
}: book) {
  return (
    <div className="relative aspect-9/16 border border-solid border-border rounded group overflow-hidden">
      {cover ?
        <img className="absolute w-full h-full top-0 left-0 right-0 bottom-0" src={cover} alt="Cover of Book" /> :
        <span className="absolute top-0 left-0 right-0 bottom-0 flex items-center justify-center text-center w-full h-full">No cover attached</span>}

      <div className="flex flex-col p-2 max-h-1/3 absolute bottom-0 z-10 bg-linear-to-b from-bg/25 to-bg w-full opacity-0 transition-all delay-75 group-hover:opacity-100">
        <span>{title}</span>
        <span className="text-xs">{author}</span>
      </div>
    </div>
  )
}
