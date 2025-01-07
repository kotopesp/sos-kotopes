export interface Meta {
  current_page: number;
  per_page: number;
  total: number;
  total_pages: number;
}

export interface TrapperResponse {
  data: {
    meta: Meta
    trappers: Trapper[];
  };
  status: string;   
}

export interface TrapperResponse {
  trappers: Trapper[];
  meta: Meta;
}

export interface Trapper {
  name: string,
  photo: File,
  animal_type: string,
  description: string,
  equipment_rental: number,
  have_car: boolean,
  have_ladder: boolean,
  have_metal_cage: boolean,
  have_net: boolean,
  have_other: string,
  have_plastic_cage: boolean,
  id: number,
  location: string,
  price: number,
  user_id: number,
  willingness_carry: string,
}