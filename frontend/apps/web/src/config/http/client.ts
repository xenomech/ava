import { env } from "@ava/env/web";
import axios from "axios";

import { attachInterceptors } from "./interceptors";

export const apiClient = axios.create({
  baseURL: env.VITE_API_URL,
  timeout: 15_000,
  withCredentials: true,
  headers: { "Content-Type": "application/json" },
});

attachInterceptors(apiClient);
