import { Navigate, Route, Routes } from "react-router-dom";

import { LoginPage, PublicOnlyRoute, RequireAuth } from "../../features/auth";
import { DashboardPage } from "../../features/dashboard";
import { HealthPage } from "../pages/HealthPage";
import { NotFoundPage } from "../pages/NotFoundPage";

export function AppRouter() {
  return (
    <Routes>
      <Route path="/" element={<Navigate replace to="/app" />} />

      <Route element={<PublicOnlyRoute />}>
        <Route path="/login" element={<LoginPage />} />
      </Route>

      <Route element={<RequireAuth />}>
        <Route path="/app" element={<DashboardPage />} />
      </Route>

      <Route path="/health" element={<HealthPage />} />
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}
