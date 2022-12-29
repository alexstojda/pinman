import {render} from "@testing-library/react";
import AuthTest from "./index";
import React from "react";

const mockedNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom') as any,
  useNavigate: () => mockedNavigate,
}));

beforeEach(() => {
  mockedNavigate.mockReset()
})

describe('AuthTest', () => {
  it('renders', async () => {
    const wrapper = render(<AuthTest/>)

    const header = await wrapper.findByText("User is authenticated")
    expect(header).toBeInTheDocument()
  })
  it('logs out when logout button is pressed', async () => {
    const wrapper = render(<AuthTest/>)

    const logout = await wrapper.findByText('Logout')
    logout.click()
    expect(mockedNavigate).toBeCalledWith('/login')
  })
})