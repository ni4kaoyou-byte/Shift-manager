import { useQuery } from "@tanstack/react-query";

import { fetchAuthPing } from "../api/ping";

export function useAuthPingQuery() {
  return useQuery({
    queryKey: ["auth", "ping"],
    queryFn: fetchAuthPing,
    retry: false,
  });
}
