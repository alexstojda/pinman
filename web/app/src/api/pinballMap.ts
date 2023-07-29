import axios, {AxiosError} from 'axios'

export type Location = {
  id: number
  name: string
  street: string
  city: string
  state: string
  zip: string
  country: string
  num_machines: number
}

export type ErrorResponse = {
  errors: string
}

export class Api {
  async locationsGet(nameFilter: string): Promise<Location[]> {
    return axios.get(
      `https://pinballmap.com/api/v1/locations.json?by_location_name=${encodeURIComponent(nameFilter)}`
    ).then((response) => {
      return response.data.locations as Location[]
    })
  }

  parseError(e: AxiosError): ErrorResponse {
    if (e.isAxiosError && e.response && e.response.data) {
      return e.response.data as ErrorResponse
    }
    throw new Error("parseError: could not parse error")
  }
}
