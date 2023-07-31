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

jest.mock('../../components/LocationSelect', () => ({
  __esModule: true,
  LocationSelect: (props: { onChange: (o: LocationOption) => void }) => <input
    type='text'
    onChange={(e) => {
      props.onChange({
        label: 'test-location',
        type: 'pinman',
        value: e.target.value,
        pinballMapId: '1',
      })
    }}
    placeholder={'test-location'}
  />,
}));

beforeEach(() => {
  mockedNavigate.mockReset()
  mockApi.prototype.leaguesApi.mockReset()
  mockApi.prototype.parseError.mockReset()
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
    'creates new location with pinball map location and redirects to /leagues if successful',
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
});


async function typeLeagueInfo(result: RenderResult) {
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

  const locationSelectValue = await result.findByPlaceholderText('test-location')
  slugField.setAttribute("value", faker.datatype.uuid())
  act(() => {
    Simulate.change(locationSelectValue)
  })
}
