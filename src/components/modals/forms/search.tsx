import { useState, type ChangeEvent, type SubmitEvent } from "react";

export default function SearchForm() {

  const [search, setSearch] = useState("")

  function submitForm(event: SubmitEvent<HTMLFormElement>) {

  }

  return (
    <div
      className="flex flex-col w-md">
      <h1
        className="text-xl border-b-border border-b border-b-solid w-full text-center mb-5">Search</h1>
      <form
        onSubmit={submitForm}
        className="flex gap-1 border-fg border-solid border px-2 py-1 rounded w-full">

        <input
          value={search}
          onChange={(e: ChangeEvent<HTMLInputElement, HTMLInputElement>) => {
            setSearch(e.target.value)
          }}
          className="outline-none border-none w-full"
          type="search"
          placeholder="e.g: Title, Author, Genre" />

        <span>|</span>
        <input
          className="hover:text-border"
          type="submit"
          value="Search" />

      </form>
    </div>
  )
}
