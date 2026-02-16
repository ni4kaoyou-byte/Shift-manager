import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";

import { clearAccessToken, setAccessToken } from "../shared/auth/token-storage";
import { App } from "./App";

describe("App routing", () => {
  beforeEach(() => {
    clearAccessToken();
  });

  it("redirects unauthenticated /app to login page", () => {
    render(
      <MemoryRouter initialEntries={["/app"]}>
        <App />
      </MemoryRouter>,
    );

    expect(screen.getByRole("heading", { name: "ログイン" })).toBeInTheDocument();
  });

  it("allows authenticated users to open dashboard", () => {
    setAccessToken("token-123");

    render(
      <MemoryRouter initialEntries={["/app"]}>
        <App />
      </MemoryRouter>,
    );

    expect(screen.getByRole("heading", { name: "Shift Manager" })).toBeInTheDocument();
  });
});
