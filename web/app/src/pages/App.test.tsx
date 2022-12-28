import React from 'react';
import {render} from '@testing-library/react';
import {apiUrl} from "../helpers";

import axios from 'axios';
import App from './App';

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

test('renders with api hello data', async () => {
  mockedAxios.get.mockResolvedValue({
    status: 200,
    data: {
      hello: "world"
    }
  })

  const wrapper = render(<App/>)

  const hello = await wrapper.findByText('{"hello":"world"}')
  expect(hello).toBeInTheDocument()
  expect(mockedAxios.get).toBeCalledWith(apiUrl("/api/hello"), {
    headers: {
      Accept: "application/json",
    },
  })
});
