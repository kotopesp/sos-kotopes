export interface Meta {
  current_page: number;
  per_page: number;
  total: number;
  total_pages: number;
}

export interface PostResponse {
  data: {
    meta: Meta
    posts: Post[];
  };
  status: string;
}

export interface PostResponse {
  posts: Post[];
  meta: Meta;
}

export interface Post {
  id: number;
  animal_type: string
  author_username: string
  color: string
  comments: number
  created_at: string
  description: string
  gender: string
  is_favourite: boolean
  location: string
  photo: File
  status: string
}
