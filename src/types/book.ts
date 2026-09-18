export interface c_chapter {
  content: string
  music?: string
  end?: boolean
}

export interface c_book {
  id?: number
  title: string
  author: string
  description: string
  cover?: string
  chapters: Record<string, c_chapter[]>
}
