import { Button, Field, Input } from "@ava/ui";
import { useMutation } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { useState } from "react";

import { AuthCard } from "../components/auth-card";
import { FormError } from "../components/form-error";
import { register } from "../api";
import { isApiError } from "@/config/http/request";

function slugify(value: string): string {
  return value
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
}

export function RegisterPage() {
  const [form, setForm] = useState({
    name: "",
    username: "",
    email: "",
    password: "",
    tenant_name: "",
    tenant_slug: "",
  });
  const [slugTouched, setSlugTouched] = useState(false);

  const attempt = useMutation({
    mutationFn: () => register(form),
  });

  const fieldErrors = isApiError(attempt.error) ? attempt.error.details : undefined;

  if (attempt.isSuccess) {
    return (
      <AuthCard
        title="Check your email"
        description={`We sent a verification link to ${form.email}. Confirm it to activate your account.`}
        footer={
          <Link to="/auth/login" className="font-semibold text-fg">
            Back to sign in
          </Link>
        }
      />
    );
  }

  return (
    <AuthCard
      title="Create your workspace"
      description="One account, one workspace to start with."
      footer={
        <div className="flex items-center justify-between">
          <span>Already have an account?</span>
          <Link to="/auth/login" className="font-semibold text-fg">
            Sign in
          </Link>
        </div>
      }
    >
      <form
        className="grid gap-4"
        onSubmit={(event) => {
          event.preventDefault();
          attempt.mutate();
        }}
      >
        <Field label="Name" error={fieldErrors?.name}>
          {(props) => (
            <Input
              {...props}
              autoComplete="name"
              required
              value={form.name}
              onChange={(event) => setForm({ ...form, name: event.target.value })}
            />
          )}
        </Field>

        <Field label="Username" error={fieldErrors?.username}>
          {(props) => (
            <Input
              {...props}
              autoComplete="username"
              required
              value={form.username}
              onChange={(event) => setForm({ ...form, username: event.target.value })}
            />
          )}
        </Field>

        <Field label="Email" error={fieldErrors?.email}>
          {(props) => (
            <Input
              {...props}
              type="email"
              autoComplete="email"
              required
              value={form.email}
              onChange={(event) => setForm({ ...form, email: event.target.value })}
            />
          )}
        </Field>

        <Field label="Password" error={fieldErrors?.password} hint="At least 8 characters">
          {(props) => (
            <Input
              {...props}
              type="password"
              autoComplete="new-password"
              required
              minLength={8}
              value={form.password}
              onChange={(event) => setForm({ ...form, password: event.target.value })}
            />
          )}
        </Field>

        <Field label="Workspace" error={fieldErrors?.tenant_name}>
          {(props) => (
            <Input
              {...props}
              required
              placeholder="Acme"
              value={form.tenant_name}
              onChange={(event) =>
                setForm({
                  ...form,
                  tenant_name: event.target.value,
                  tenant_slug: slugTouched ? form.tenant_slug : slugify(event.target.value),
                })
              }
            />
          )}
        </Field>

        <Field
          label="Workspace URL"
          error={fieldErrors?.tenant_slug}
          hint="Lowercase letters, numbers and hyphens"
        >
          {(props) => (
            <Input
              {...props}
              required
              placeholder="acme"
              value={form.tenant_slug}
              onChange={(event) => {
                setSlugTouched(true);
                setForm({ ...form, tenant_slug: event.target.value });
              }}
            />
          )}
        </Field>

        {fieldErrors ? null : <FormError error={attempt.error} fallback="Could not register" />}

        <Button type="submit" className="mt-1 w-full" disabled={attempt.isPending}>
          {attempt.isPending ? "Creating…" : "Create account"}
        </Button>
      </form>
    </AuthCard>
  );
}
