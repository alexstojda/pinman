import {act, renderHook, RenderHookResult} from '@testing-library/react'
import {AuthOptions, AuthState, useAuth} from "./useAuth";
import {faker} from "@faker-js/faker";
import {Api, UserResponse, UsersApi} from "./index";

jest.mock("./index")

const mockedNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom') as any,
  useNavigate: () => mockedNavigate,
}));

let hook: RenderHookResult<AuthState, AuthOptions> | null = null

beforeEach(() => {
  hook = null
})

test('renders when not authenticated and requireAuth false', async () => {
  const MockedApi = jest.mocked(Api)
  MockedApi.prototype.tryTokenRefresh.mockResolvedValue(false)

  await act(async () => {
    hook = renderHook<AuthState, AuthOptions>(() => useAuth({requireAuth: false}))
  })

  expect(hook!.result.current.user).toBeUndefined()
  expect(mockedNavigate).toBeCalledTimes(0)
})

test('redirects when not authenticated and requireAuth true', async () => {
  const MockedApi = jest.mocked(Api)
  MockedApi.prototype.tryTokenRefresh.mockResolvedValue(false)

  await act(async () => {
    hook = renderHook<AuthState, AuthOptions>(() => useAuth({requireAuth: true}))
  })

  expect(hook!.result.current.user).toBeUndefined()
  expect(mockedNavigate).toBeCalledWith("/login")
})

test('retrieves user information if isAuthenticated', async () => {
  const userResponse: UserResponse = {
    user: {
      id: faker.datatype.uuid(),
      name: faker.name.fullName(),
      email: faker.internet.email(),
      role: "user",
      updated_at: faker.date.recent(1).toISOString(),
      created_at: faker.date.recent(3).toISOString(),
    }
  }

  const MockedUserApi = jest.mocked(UsersApi)
  MockedUserApi.prototype.usersMeGet.mockResolvedValue({
    status: 200,
    statusText: "Ok",
    headers: {},
    data: userResponse,
    config: {}
  })

  const MockedApi = jest.mocked(Api)
  MockedApi.prototype.tryTokenRefresh.mockResolvedValue(true)
  MockedApi.prototype.userApi.mockReturnValue(new UsersApi())

  await act(async () => {
    hook = renderHook<AuthState, AuthOptions>(() => useAuth({requireAuth: false}))
  })

  expect(hook!.result.current.user).toEqual(userResponse.user)
  expect(mockedNavigate).toBeCalledTimes(0)
})


test('isAuthenticated but user info request faise', async () => {
  const MockedUserApi = jest.mocked(UsersApi)
  MockedUserApi.prototype.usersMeGet.mockResolvedValue({
    status: 404,
    statusText: "Not Found",
    headers: {},
    data: {
      code: 404,
      detail: "user does not exist"
    } as any,
    config: {}
  })

  const MockedApi = jest.mocked(Api)
  MockedApi.prototype.tryTokenRefresh.mockResolvedValue(true)
  MockedApi.prototype.userApi.mockReturnValue(new UsersApi())

  const error = jest.spyOn(console, "error").mockImplementation(() => {});

  await act(async () => {
    hook = renderHook<AuthState, AuthOptions>(() => useAuth({requireAuth: false}))
  })

  expect(hook!.result.current.user).toBeUndefined()
  expect(error).toBeCalledTimes(1)
  expect(mockedNavigate).toBeCalledTimes(0)
})
