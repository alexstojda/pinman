import {render} from "@testing-library/react";
import React from "react";
import ColorToggle from "./index";

test('renders without error', async () => {
  render(<ColorToggle/>)
});
