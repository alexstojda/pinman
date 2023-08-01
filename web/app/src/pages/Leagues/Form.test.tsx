import {act, render, RenderResult, waitFor} from "@testing-library/react";
import LeagueForm from "./Form";
import {faker} from "@faker-js/faker";
import {Simulate} from "react-dom/test-utils";
import {randomInt} from "crypto";
import * as helpers from '../../helpers'
import {Api, LeaguesApi, LocationsApi} from "../../api";
import {fake} from "../../test";
import {LocationOption} from "../../components/LocationSelect";

jest.mock('../../api')
const mockApi = jest.mocked(Api)

const mockedNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom') as any,
  useNavigate: () => mockedNavigate,
}));

let mockLocation: LocationOption
jest.mock('../../components/LocationSelect', () => ({
  __esModule: true,
  LocationSelect: (props: { onChange: (o: LocationOption) => void }) => <input
    type='text'
    onChange={(e) => {
      props.onChange({
        ...mockLocation,
        value: e.target.value,
      })
    }}
    placeholder={'test-location'}
  />,
}));

beforeEach(() => {
  mockLocation = {
    label: faker.lorem.words(2),
    value: '1234',
    type: 'pinman',
    pinballMapId: randomInt(1000).toString(),
  }
  mockedNavigate.mockReset()
  mockApi.prototype.leaguesApi.mockReset()
  mockApi.prototype.parseError.mockReset()
})

afterEach(() => {
  jest.clearAllMocks()
})

describe('LeagueForm', function () {
  it('should render', function () {
    render(
      <LeagueForm mode={'create'}/>
    )
  });
  it('redirects to /leagues if successful', async function () {
    jest.mocked(LeaguesApi).prototype.leaguesPost.mockResolvedValue({
      config: {},
      data: {
        league: fake.league()
      },
      headers: {},
      request: {},
      status: 201,
      statusText: "",
    })
    mockApi.prototype.leaguesApi.mockReturnValue(new LeaguesApi())

    const result = render(
      <LeagueForm mode={'create'}/>
    )
    await typeLeagueInfo(result)

    const submitButton = await result.findByText("Create")
    act(() => {
      Simulate.click(submitButton)
    })

    await waitFor(() => expect(mockedNavigate).toBeCalledWith('/leagues'))
  })
  it('renders error from API', async () => {
    const mockError = fake.errorResponse()
    mockApi.prototype.parseError.mockReturnValue(mockError)

    jest.mocked(LeaguesApi).prototype.leaguesPost.mockRejectedValue(false)
    mockApi.prototype.leaguesApi.mockReturnValue(new LeaguesApi())

    const result = render(
      <LeagueForm mode={'create'}/>
    )
    await typeLeagueInfo(result)

    const submitButton = await result.findByText("Create")
    act(() => {
      Simulate.click(submitButton)
    })

    await waitFor(() => {
      expect(mockApi.prototype.parseError).toBeCalled()
      expect(result.getByText(mockError.title)).toBeInTheDocument()
    })
  })
  it(
    'creates new location with existing location and redirects to /leagues if successful',
    async function () {
      jest.mocked(LeaguesApi).prototype.leaguesPost.mockResolvedValue({
        config: {},
        data: {
          league: fake.league()
        },
        headers: {},
        request: {},
        status: 201,
        statusText: "",
      })
      mockApi.prototype.leaguesApi.mockReturnValue(new LeaguesApi())

      mockLocation.type = 'pinmap'

      jest.mocked(LocationsApi).prototype.locationsPost.mockResolvedValue({
        config: {},
        data: {
          location: fake.location()
        },
        headers: {},
        request: {},
        status: 201,
        statusText: "",
      })
      mockApi.prototype.locationsApi.mockReturnValue(new LocationsApi())

      const result = render(
        <LeagueForm mode={'create'}/>
      )
      await typeLeagueInfo(result)

      const submitButton = await result.findByText("Create")
      act(() => {
        Simulate.click(submitButton)
      })

      await waitFor(() => expect(mockedNavigate).toBeCalledWith('/leagues'))
    })
  it('shows error when create new location fails', async function () {
    mockLocation.type = 'pinmap'

    jest.mocked(LocationsApi).prototype.locationsPost.mockRejectedValue(false)
    mockApi.prototype.locationsApi.mockReturnValue(new LocationsApi())
    const mockError = fake.errorResponse()
    mockApi.prototype.parseError.mockReturnValue(mockError)

    const result = render(
      <LeagueForm mode={'create'}/>
    )
    await typeLeagueInfo(result)

    const submitButton = await result.findByText("Create")
    act(() => {
      Simulate.click(submitButton)
    })

    const error = await result.findByText(mockError.detail)
    await waitFor(() => {
      expect(error).toBeInTheDocument()
      expect(mockedNavigate).not.toBeCalled()
    })
  })
  it('shows error if submitting and no location specified', async function () {
    const result = render(
      <LeagueForm mode={'create'}/>
    )
    await typeLeagueInfo(result, true)

    const submitButton = await result.findByText("Create")
    act(() => {
      Simulate.click(submitButton)
    })

    const error = await result.findByText('Please select a location for your league')
    await waitFor(() => {
      expect(error).toBeInTheDocument()
      expect(mockedNavigate).not.toBeCalled()
    })
  })
});

async function typeLeagueInfo(result: RenderResult, omitLocation?: boolean) {
  const name = faker.random.words(randomInt(3, 10))

  const nameField = await result.findByPlaceholderText("name")
  nameField.setAttribute("value", faker.name.fullName())
  act(() => {
    Simulate.change(nameField)
  })

  const slugField = await result.findByPlaceholderText("slug")
  slugField.setAttribute("value", helpers.slugify(name))
  act(() => {
    Simulate.change(slugField)
  })

  if (!omitLocation) {
    const locationSelectValue = await result.findByPlaceholderText('test-location')
    slugField.setAttribute("value", faker.datatype.uuid())
    act(() => {
      Simulate.change(locationSelectValue)
    })
  }
}
