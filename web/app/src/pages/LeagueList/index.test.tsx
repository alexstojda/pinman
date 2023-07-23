import {render} from "@testing-library/react";
import React from "react";
import LeagueListPage from "./index";
import {Api, LeaguesApi} from "../../api";
import {fake} from "../../test";
import {MemoryRouter} from "react-router-dom";

const mockedNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom') as any,
  useNavigate: () => mockedNavigate,
}));

jest.mock('../../api')
const mockApi = jest.mocked(Api)

beforeEach(() => {
  mockedNavigate.mockReset()
  mockApi.prototype.leaguesApi.mockReset()
  mockApi.prototype.parseError.mockReset()
})


describe('LeagueListPage', () => {
  it('renders leagues from API', async () => {
    const mockLeagues = Array.from({length: 5}, () => fake.league())

    jest.mocked(LeaguesApi).prototype.leaguesGet.mockResolvedValue({
      config: {},
      data: {
        leagues: mockLeagues
      },
      headers: {},
      request: {},
      status: 200,
      statusText: "",
    })
    mockApi.prototype.leaguesApi.mockReturnValue(new LeaguesApi())

    const result = render(
      <MemoryRouter>
        <LeagueListPage/>
      </MemoryRouter>
    )

    const header = await result.findByText("Leagues")
    expect(header).toBeInTheDocument()

    const leagueCards = await result.findAllByTestId("league-card")
    expect(leagueCards).toHaveLength(mockLeagues.length)

    await mockLeagues.forEach((league) => {
      expect(result.getByText(league.name)).toBeInTheDocument()
      expect(result.getByText(league.location)).toBeInTheDocument()
    })
  })
  it('renders error from API', async () => {
    const mockError = fake.errorResponse()
    mockApi.prototype.parseError.mockReturnValue(mockError)

    jest.mocked(LeaguesApi).prototype.leaguesGet.mockRejectedValue(false)
    mockApi.prototype.leaguesApi.mockReturnValue(new LeaguesApi())

    const result = render(
      <MemoryRouter>
        <LeagueListPage/>
      </MemoryRouter>
    )

    const error = await result.findByText(mockError.detail)
    expect(error).toBeInTheDocument()
  })
  it('shows unknown error when parseError fails', async () => {
    mockApi.prototype.parseError.mockImplementation((e) => {
      throw e
    })

    jest.mocked(LeaguesApi).prototype.leaguesGet.mockRejectedValue(false)
    mockApi.prototype.leaguesApi.mockReturnValue(new LeaguesApi())

    const result = render(
      <MemoryRouter>
        <LeagueListPage/>
      </MemoryRouter>
    )

    const error = await result.findByText("Unknown error")
    expect(error).toBeInTheDocument()
  })
  it('navigates to create league page when button is pressed', async () => {
    jest.mocked(LeaguesApi).prototype.leaguesGet.mockResolvedValue({
      config: {},
      data: {
        leagues: []
      },
      headers: {},
      request: {},
      status: 200,
      statusText: "",
    })
    mockApi.prototype.leaguesApi.mockReturnValue(new LeaguesApi())

    const result = render(
      <MemoryRouter>
        <LeagueListPage/>
      </MemoryRouter>
    )

    const createButton = await result.findByText("Create League")
    expect(createButton).toBeInTheDocument()

    createButton.click()

    expect(mockedNavigate).toHaveBeenCalledWith('/leagues/create')
  })
})
