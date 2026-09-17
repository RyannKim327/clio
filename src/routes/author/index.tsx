import Header from "@/components/widgets/header"
import { useState } from "react"

export default function AuthorPage() {
  const [action, setAction] = useState("")

  return (
    <div className="flex flex-col gap-2 w-full h-full">
      <Header setAction={setAction} />
      <span>For authors</span>
    </div>
  )
}
