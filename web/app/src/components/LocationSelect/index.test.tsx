import {act, fireEvent, render, waitFor} from "@testing-library/react";
import React from "react";
import {LocationSelect} from "./index";
import {ChakraProvider} from "@chakra-ui/react";
import {Api, LocationsApi, pinballMap} from "../../api";
import {fake} from "../../test";

window.matchMedia = window.matchMedia || function () {
  return {
    matches: false,
    addListener: function () {
    },
    removeListener: function () {
    }
  };
};
jest.useFakeTimers();

jest.mock('../../api')
const mockApi = jest.mocked(Api)
const mockPinballMapApi = jest.mocked(pinballMap.Api)

beforeEach(() => {
  mockApi.prototype.parseError.mockReset()
  mockApi.prototype.locationsApi.mockReset()

  mockPinballMapApi.prototype.locationsGet.mockReset()
  mockPinballMapApi.prototype.parseError.mockReset()
  jest.resetAllMocks()
})

describe('LocationSelect', () => {
  it('renders without error', async () => {
    render(
      <ChakraProvider>
        <LocationSelect/>
      </ChakraProvider>
    )
  });
  it('invokes the onError callback if pinball map api call fails', async () => {
    mockPinballMapApi.prototype.locationsGet.mockRejectedValue(false)
    const mockError: pinballMap.ErrorResponse = {
      errors: "test error",
    }
    mockPinballMapApi.prototype.parseError.mockReturnValue(mockError)

    jest.mocked(LocationsApi).prototype.locationsGet.mockResolvedValue({
      config: {},
      data: {
        locations: [fake.location()],
      },
      headers: {},
      request: {},
      status: 200,
      statusText: "",
    })
    mockApi.prototype.locationsApi.mockReturnValue(new LocationsApi())

    const mockOnError = jest.fn()

    const wrapper = render(
      <ChakraProvider>
        <LocationSelect onError={mockOnError}/>
      </ChakraProvider>
    )

    await act(async () => {
      fireEvent.change(wrapper.getByRole('combobox'), {target: {value: 'valid search'}})
      jest.runAllTimers()
    })

    await waitFor(() => {
      expect(mockOnError).toBeCalled()
    })
  })
  it('invokes the onError callback when pinman api call fails', async () => {
    const mockError = fake.errorResponse()
    mockApi.prototype.parseError.mockReturnValue(mockError)

    mockPinballMapApi.prototype.locationsGet.mockResolvedValue([
      fake.pinballMapLocation(),
    ])
    jest.mocked(LocationsApi).prototype.locationsGet.mockRejectedValue(false)
    mockApi.prototype.locationsApi.mockReturnValue(new LocationsApi())
    const mockOnError = jest.fn()

    const wrapper = render(
      <ChakraProvider>
        <LocationSelect onError={mockOnError}/>
      </ChakraProvider>
    )

    await act(async () => {
      fireEvent.change(wrapper.getByRole('combobox'), {target: {value: 'te'}})
      jest.advanceTimersByTime(500)
      fireEvent.change(wrapper.getByRole('combobox'), {target: {value: 'test l'}})
      jest.runAllTimers()
    })

    await waitFor(() => {
      expect(mockOnError).toBeCalled()
    })
  })
  it('renders with options when filter provided', async () => {
    const mockPinballMapLocation = {
      ...fake.pinballMapLocation(),
      name: "test location",
      id: 456,
    }
    const mockKnownPinballMapLocation = {
      ...fake.pinballMapLocation(),
      name: "should not be in results",
      id: 123,
    }
    const mockPinmanLocation = {
      ...fake.location(),
      name: "test location 2",
      pinball_map_id: 123,
    }
    mockPinballMapApi.prototype.locationsGet.mockResolvedValue([
      mockPinballMapLocation,
      mockKnownPinballMapLocation
    ])
    jest.mocked(LocationsApi).prototype.locationsGet.mockResolvedValue({
      config: {},
      data: {
        locations: [
          mockPinmanLocation,
        ]
      },
      headers: {},
      request: {},
      status: 200,
      statusText: "",
    })
    mockApi.prototype.locationsApi.mockReturnValue(new LocationsApi())

    const wrapper = render(
      <ChakraProvider>
        <LocationSelect/>
      </ChakraProvider>
    )

    jest.spyOn(global, 'setTimeout')

    await act(async () => {
      fireEvent.change(wrapper.getByRole('combobox'), {target: {value: 'te'}})
      expect(setTimeout).not.toBeCalled()
    })

    await act(async () => {
      fireEvent.change(wrapper.getByRole('combobox'), {target: {value: 'test'}})
      jest.advanceTimersByTime(500)
      fireEvent.change(wrapper.getByRole('combobox'), {target: {value: 'test l'}})
      jest.runAllTimers()
    })

    await waitFor(async () => {
      expect(await wrapper.findByText(mockPinballMapLocation.name)).toBeInTheDocument()
      expect(await wrapper.findByText(mockPinmanLocation.name)).toBeInTheDocument()
      expect(wrapper.baseElement).not.toHaveTextContent(mockKnownPinballMapLocation.name)
    })
  });
});
