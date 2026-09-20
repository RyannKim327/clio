export default function Search() {
  return (
    <div
      className="flex flex-col w-md">
      <h1
        className="text-xl border-b-border border-b border-b-solid w-full text-center mb-5">Search</h1>
      <form
        action=""
        className="flex gap-1 border-fg border-solid border px-2 py-1 rounded w-full">

        <input
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
