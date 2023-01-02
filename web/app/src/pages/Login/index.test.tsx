import {act, render, RenderResult} from "@testing-library/react";
import LoginPage from "./index";
import {MemoryRouter} from "react-router-dom";
import React from "react";
import {Api, useAuth, User} from "../../api";
import {faker} from "@faker-js/faker";
import {Simulate} from "react-dom/test-utils";

const mockNavigate = jest.fn()
jest.mock('react-router-dom', () => ({
  ...(jest.requireActual('react-router-dom')),
  useNavigate: () => mockNavigate
}))

jest.mock('../../api/useAuth')
const mockUseAuth = jest.mocked(useAuth)

jest.mock('../../api/api')
const mockApi = jest.mocked(Api)

beforeEach(() => {
  mockApi.prototype.login.mockReset()
  mockApi.prototype.parseError.mockReset()
  mockNavigate.mockReset()
  mockUseAuth.mockReset()
})

describe('LoginPage', () => {
  it('renders', async () => {
    mockUseAuth.mockReturnValue({user: undefined})
    const result = render(
      <MemoryRouter initialEntries={['/login']}>
        <LoginPage/>
      </MemoryRouter>
    )
    await result.findAllByText("Login")
  })

  it('shows success message if registered=true', async () => {
    mockUseAuth.mockReturnValue({user: undefined})
    const result = render(
      <MemoryRouter initialEntries={['/login?registered=true']}>
        <LoginPage/>
      </MemoryRouter>
    )
    await result.findAllByText("Account created, please log in")
  })

  it('redirects to /authenticated if user is signed-in', async () => {
    const mockUser: User = {
      name: faker.name.fullName(),
      email: faker.internet.email(),
      id: faker.datatype.uuid(),
      role: "user",
      created_at: faker.date.recent(5).toISOString(),
      updated_at: faker.date.recent(1).toISOString()
    }
    mockUseAuth.mockReturnValue({
      user: mockUser
    })

    await act(() => {
      render(
        <MemoryRouter>
          <LoginPage/>
        </MemoryRouter>
      )
    })

    expect(mockNavigate).toBeCalledWith("/authenticated")
  })

  async function typeLoginInfo(result: RenderResult) {
    const usernameField = await result.findByPlaceholderText("email address")
    usernameField.setAttribute("value", faker.internet.email())
    await act(() => {
      Simulate.change(usernameField)
    })

    const passwordField = await result.findByPlaceholderText('Password')
    passwordField.setAttribute("value", faker.internet.password())
    await act(() => {
      Simulate.change(passwordField)
    })
  }

  it('redirects to /navigate if user login was successful', async () => {
    mockUseAuth.mockReturnValue({user: undefined})
    mockApi.prototype.login.mockResolvedValue(true)

    const result = render(
      <MemoryRouter>
        <LoginPage/>
      </MemoryRouter>
    )

    await typeLoginInfo(result)

    await result.getByTestId("login").click()

    expect(mockNavigate).toBeCalledWith("/authenticated")
  })

  it('displays error message if login fails', async () => {
    mockUseAuth.mockReturnValue({user: undefined})
    mockApi.prototype.login.mockRejectedValue(false)
    mockApi.prototype.parseError.mockReturnValue({
      status: 401,
      title: "Unauthorized",
      detail: "Invalid username or password"
    })

    const result = render(
      <MemoryRouter>
        <LoginPage/>
      </MemoryRouter>
    )

    await typeLoginInfo(result)

    const logSpy = jest.spyOn(console, 'error').mockImplementation(jest.fn())

    await act(() => {
      result.getByTestId("login").click()
    })

    const errorMessage = await result.findByText("Invalid username or password")
    expect(errorMessage).toBeInTheDocument()
    expect(logSpy).toBeCalled()

  })
})