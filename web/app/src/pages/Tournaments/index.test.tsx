import {Api, TournamentsApi} from "../../api";
import {fake} from '../../test';
import {render} from "@testing-library/react";
import {MemoryRouter} from "react-router-dom";
import TournamentsListPage from "./index";

jest.mock('../../api')
const mockApi = jest.mocked(Api)

beforeEach(() => {
  mockApi.prototype.leaguesApi.mockReset()
  mockApi.prototype.parseError.mockReset()
})

describe('TournamentsListPage', () => {
  it('renders tournaments from API', async () => {
    const mockTournaments = Array.from({length: 5}, () => fake.tournament())

    jest.mocked(TournamentsApi).prototype.tournamentsGet.mockResolvedValue({
      config: {},
      data: {
        tournaments: mockTournaments
      },
      headers: {},
      request: {},
      status: 200,
      statusText: "",
    })
    mockApi.prototype.tournamentsApi.mockReturnValue(new TournamentsApi())

    const result = render(
      <MemoryRouter>
        <TournamentsListPage/>
      </MemoryRouter>
    )

    const header = await result.findByText("Tournaments")
    expect(header).toBeInTheDocument()

    const tournamentCards = await result.findAllByTestId("tournament-card")
    expect(tournamentCards).toHaveLength(mockTournaments.length)

    mockTournaments.forEach((tournament) => {
      expect(result.getByText(tournament.name)).toBeInTheDocument()
      expect(result.getByText(tournament.location!.name)).toBeInTheDocument()
    })
  })
  it('renders error from API', async () => {
    const mockError = fake.errorResponse()
    mockApi.prototype.parseError.mockReturnValue(mockError)

    jest.mocked(TournamentsApi).prototype.tournamentsGet.mockRejectedValue(false)
    mockApi.prototype.tournamentsApi.mockReturnValue(new TournamentsApi())

    const result = render(
      <MemoryRouter>
        <TournamentsListPage/>
      </MemoryRouter>
    )

    const header = await result.findByText("Tournaments")
    expect(header).toBeInTheDocument()

    const error = await result.findByText(mockError.detail)
    expect(error).toBeInTheDocument()
  })
  it('shows unknown error when parseError fails', async () => {
    mockApi.prototype.parseError.mockImplementation((e) => {
      throw e
    })

    mockApi.prototype.tournamentsApi.mockReturnValue(new TournamentsApi())
    jest.mocked(TournamentsApi).prototype.tournamentsGet.mockRejectedValue(false)

    const result = render(
      <MemoryRouter>
        <TournamentsListPage/>
      </MemoryRouter>
    )

    const header = await result.findByText("Tournaments")
    expect(header).toBeInTheDocument()

    const error = await result.findByText("Unknown error")
    expect(error).toBeInTheDocument()
  })
})
