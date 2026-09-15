export default function Header() {
  return (
    <div className="flex justify-between items-center w-full h-10 bg-secondary-bg px-2 sticky top-0 z-10">
      <span>Istorya</span>
      <nav className="flex gap-2">
        <li className="list-none">Home</li>
        <li className="list-none">About</li>
      </nav>
    </div>
  )
}
