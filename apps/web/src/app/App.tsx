import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import { AuthProvider } from "../features/auth";
import { AppRouter } from "./router/AppRouter";

const queryClient = new QueryClient();

export function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <AppRouter />
      </AuthProvider>
    </QueryClientProvider>
  );
}
