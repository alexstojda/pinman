import {render} from "@testing-library/react";
import AuthLayout from "./index";


describe('AuthLayout', () => {
  it('renders', () => {
    render(<AuthLayout
      title={"Auth Layout Page"}
      alert={{
        status: "error",
        detail: "An error occurred",
        title: "Error"
      }}
    >
      <h1>Hello World</h1>
    </AuthLayout>)
  })
})