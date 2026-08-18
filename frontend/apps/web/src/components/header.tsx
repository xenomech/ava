import { Link } from "@tanstack/react-router";

import { ModeToggle } from "./mode-toggle";

const links = [{ to: "/", label: "Home" }] as const;

export default function Header() {
  return (
    <header className="border-b">
      <div className="mx-auto flex w-full max-w-5xl items-center justify-between px-4 py-3">
        <div className="flex items-center gap-6">
          <Link to="/" className="font-semibold tracking-tight">
            Ava
          </Link>
          <nav className="flex items-center gap-4 text-sm text-muted-foreground">
            {links.map((link) => (
              <Link
                key={link.to}
                to={link.to}
                activeProps={{ className: "text-foreground" }}
              >
                {link.label}
              </Link>
            ))}
          </nav>
        </div>
        <ModeToggle />
      </div>
    </header>
  );
}
