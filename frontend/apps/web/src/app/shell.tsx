import { Link } from "@tanstack/react-router";

import { ModeToggle } from "@/shared/components/mode-toggle";

export function AppShell({ children }: { children: React.ReactNode }) {
  return (
    <div className="grid h-svh grid-rows-[auto_1fr] bg-bg">
      <Header />
      <div className="overflow-y-auto">{children}</div>
    </div>
  );
}

function Header() {
  return (
    <header className="border-b border-border">
      <div className="mx-auto flex w-full max-w-5xl items-center justify-between px-4 py-3">
        <Link to="/" className="text-lead font-semibold">
          Ava
        </Link>
        <ModeToggle />
      </div>
    </header>
  );
}
