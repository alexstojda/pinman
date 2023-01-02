import SignUpPage from "./index";
import {act, render, RenderResult} from "@testing-library/react";
import {MemoryRouter} from "react-router-dom";
import {Api, useAuth, User, UsersApi} from "../../api";
import {faker} from "@faker-js/faker";
import {Simulate} from "react-dom/test-utils";
import React from "react";

const mockNavigate = jest.fn()
jest.mock('react-router-dom', () => ({
  ...(jest.requireActual('react-router-dom')),
  useNavigate: () => mockNavigate
}))

jest.mock('../../api')
const mockUseAuth = jest.mocked(useAuth)

const mockApi = jest.mocked(Api)

describe('SignUpPage', () => {
  it('renders', async () => {
    mockUseAuth.mockReturnValue({user: undefined})
    const result = render(
      <MemoryRouter>
        <SignUpPage/>
      </MemoryRouter>
    )
    await result.findByText("Register")
  })
  it('redirects to /authenticated if user is logged in', async () => {
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
          <SignUpPage/>
        </MemoryRouter>
      )
    })

    expect(mockNavigate).toBeCalledWith("/authenticated")
  })
  it('redirects to /login?registered=true for successful login', async () => {
    mockUseAuth.mockReturnValue({user: undefined})
    const mockUserApi = jest.mocked(UsersApi)
    mockUserApi.prototype.usersRegisterPost.mockResolvedValue({
      config: {},
      data: {},
      request: {},
      headers: {},
      status: 201,
      statusText: "",
    })

    mockApi.prototype.userApi.mockReturnValue(new UsersApi())

    const result = render(
      <MemoryRouter>
        <SignUpPage/>
      </MemoryRouter>
    )

    await typeRegisterInfo(result, true)

    await act(() => {
      result.getByTestId("create account").click()
    })

    expect(mockNavigate).toBeCalledWith("/login?registered=true")
  })
  it("shows error if passwords don't match", async () => {
    mockUseAuth.mockReturnValue({user: undefined})

    const result = render(
      <MemoryRouter>
        <SignUpPage/>
      </MemoryRouter>
    )

    await typeRegisterInfo(result, false)

    await act(() => {
      result.getByTestId("create account").click()
    })

    const error = await result.getByText("Password confirmation does not match")

    expect(mockNavigate).not.toBeCalled()
    expect(error).toBeInTheDocument()
  })
  it('shows error if registration fails', async () => {
    mockUseAuth.mockReturnValue({user: undefined})
    jest.mocked(UsersApi).prototype.usersRegisterPost.mockRejectedValue(false)
    mockApi.prototype.parseError.mockReturnValue({
      status: 400,
      title: "Bad Request",
      detail: "username is already taken"
    })
    mockApi.prototype.userApi.mockReturnValue(new UsersApi())

    const result = render(
      <MemoryRouter>
        <SignUpPage/>
      </MemoryRouter>
    )

    await typeRegisterInfo(result, true)

    const logSpy = jest.spyOn(console, 'error').mockImplementation(jest.fn())

    await act(() => {
      result.getByTestId("create account").click()
    })

    const error = await result.findByText("username is already taken")
    expect(error).toBeInTheDocument()
    expect(logSpy).toBeCalled()
  })
})

async function typeRegisterInfo(result: RenderResult, samePassword: boolean) {
  const nameField = await result.findByPlaceholderText("name")
  nameField.setAttribute("value", faker.name.fullName())
  await act(() => {
    Simulate.change(nameField)
  })

  // noinspection DuplicatedCode
  const usernameField = await result.findByPlaceholderText("email address")
  usernameField.setAttribute("value", faker.internet.email())
  await act(() => {
    Simulate.change(usernameField)
  })

  const password = faker.internet.password()

  const passwordField = await result.findByPlaceholderText('password')
  passwordField.setAttribute("value", password)
  await act(() => {
    Simulate.change(passwordField)
  })

  const confirmPasswordField = await result.findByPlaceholderText('confirm password')
  confirmPasswordField.setAttribute("value", samePassword ? password : faker.internet.password())
  await act(() => {
    Simulate.change(confirmPasswordField)
  })
}