import { RouterProvider } from "@tanstack/react-router";
import ReactDOM from "react-dom/client";

import { Providers } from "./app/providers";
import { queryClient } from "./app/query-client";
import { buildRouter } from "./app/router";

const rootElement = document.getElementById("app");

if (!rootElement) {
  throw new Error("Root element not found");
}

const router = buildRouter(queryClient);

if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement);
  root.render(
    <Providers>
      <RouterProvider router={router} />
    </Providers>,
  );
}
