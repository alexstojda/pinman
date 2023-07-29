import {act, render, RenderResult} from "@testing-library/react";
import LeagueForm from "./Form";
import {faker} from "@faker-js/faker";
import {Simulate} from "react-dom/test-utils";
import {randomInt} from "crypto";
import * as helpers from '../../helpers'
import {Api, LeaguesApi} from "../../api";
import {fake} from "../../test";

jest.mock('../../api')
const mockApi = jest.mocked(Api)

const mockedNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom') as any,
  useNavigate: () => mockedNavigate,
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

    expect(mockedNavigate).toHaveBeenCalledWith('/leagues')
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

    expect(mockApi.prototype.parseError).toBeCalled()
    expect(result.getByText(mockError.title)).toBeInTheDocument()
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

  const locationField = await result.findByPlaceholderText("location")
  locationField.setAttribute("value", faker.random.words(4))
  act(() => {
    Simulate.change(locationField)
  })
}
