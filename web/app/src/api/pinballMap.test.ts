import axios from "axios";
import {Api} from "./pinballMap";

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

beforeEach(() => {
  mockedAxios.get.mockReset()
})

describe('Api.getLocations()', () => {
  it('returns locations if api call was successful', async () => {
const locations = [{
      id: 1,
      name: 'Location 1',
      street: '123 Main St',
      city: 'Anytown',
      state: 'CA',
      zip: '12345',
      country: 'USA',
      num_machines: 5
    }]
    mockedAxios.get.mockResolvedValueOnce({data: {locations: locations}})

    const result = await new Api().locationsGet('Location 1')
    expect(mockedAxios.get).toBeCalledTimes(1)
    expect(mockedAxios.get).toBeCalledWith(
      'https://pinballmap.com/api/v1/locations.json?by_location_name=Location%201'
    )
    expect(result).toEqual(locations)
  })
})

describe('Api.parseError()', () => {
  it('returns error response if error is axios error', () => {
    const errorResponse = {
      errors: 'Error message'
    }
    const axiosError = {
      isAxiosError: true,
      response: {
        data: errorResponse
      }
    }
    const result = new Api().parseError(axiosError as any)
    expect(result).toEqual(errorResponse)
  })
  it('throws error if error is not axios error', () => {
    const errorResponse = {
      errors: 'Error message'
    }
    const axiosError = {
      isAxiosError: false,
      response: {
        data: errorResponse
      }
    }
    expect(() => {
      new Api().parseError(axiosError as any)
    }).toThrowError()
  })
})
