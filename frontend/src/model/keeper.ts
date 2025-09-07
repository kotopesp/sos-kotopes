export interface Meta {
    current_page: number;
    per_page: number;
    total: number;
    total_pages: number;
  }
  
  export interface KeeperResponse {
    data: {
      meta: Meta
      payload: Keeper[];
    };
    status: string;   
  }
  
  export interface Keeper {
    user_id: number,
    id: number,
    location_id: number,
    name: string,
    photo: File,
    description: string,
    location_name: string,
    price: number,
    has_cage: boolean,
    boarding_duration: string,
    boarding_compensation: string,
    animal_acceptance: string,
    animal_category: string,
  }