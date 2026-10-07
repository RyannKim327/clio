export type t_chapter = Record<string, c_chapter[]>

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
  chapters: t_chapter
}
