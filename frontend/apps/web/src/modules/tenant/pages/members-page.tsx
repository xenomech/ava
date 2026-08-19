import {
  Button,
  Chip,
  Field,
  Input,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@ava/ui";
import { TENANT_ROLES, type TenantRole } from "@ava/contracts";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { toast } from "sonner";

import { isApiError } from "@/config/http/request";
import { useSession } from "@/modules/auth";
import { Page, Section } from "@/shared/components/page";
import { inviteMember, removeMember } from "../api";
import { tenantQueries } from "../queries";

const INVITABLE = TENANT_ROLES.filter((role) => role !== "owner");

export function MembersPage() {
  const queryClient = useQueryClient();
  const { isAdmin } = useSession();

  const members = useQuery(tenantQueries.members());

  const [email, setEmail] = useState("");
  const [role, setRole] = useState<TenantRole>("member");

  const invalidateMembers = () =>
    queryClient.invalidateQueries({ queryKey: tenantQueries.members().queryKey });

  const invite = useMutation({
    mutationFn: () => inviteMember({ email, role }),
    onSuccess: async () => {
      toast.success(`Invitation sent to ${email}`);
      setEmail("");
      await invalidateMembers();
    },
    onError: (error) =>
      toast.error(isApiError(error) ? error.message : "Could not send the invitation"),
  });

  const remove = useMutation({
    mutationFn: removeMember,
    onSuccess: async () => {
      toast.success("Member removed");
      await invalidateMembers();
    },
    onError: (error) => toast.error(isApiError(error) ? error.message : "Could not remove member"),
  });

  return (
    <Page title="People" description="Everyone who can control this home.">
      {isAdmin ? (
        <Section title="Invite someone" description="They get an email with a link to join.">
          <form
            className="grid gap-4 p-5 sm:grid-cols-[minmax(0,1fr)_170px_auto] sm:items-end"
            onSubmit={(event) => {
              event.preventDefault();
              invite.mutate();
            }}
          >
            <Field
              label="Email"
              error={isApiError(invite.error) ? invite.error.details?.email : undefined}
            >
              {(props) => (
                <Input
                  {...props}
                  type="email"
                  required
                  placeholder="someone@example.com"
                  value={email}
                  onChange={(event) => setEmail(event.target.value)}
                />
              )}
            </Field>

            <Field label="Role">
              {({ id, ...props }) => (
                <Select value={role} onValueChange={(value) => setRole(value as TenantRole)}>
                  <SelectTrigger id={id} {...props}>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {INVITABLE.map((value) => (
                      <SelectItem key={value} value={value} className="capitalize">
                        {value}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            </Field>

            <Button type="submit" loading={invite.isPending}>
              Send invite
            </Button>
          </form>
        </Section>
      ) : null}

      <Section>
        {members.isPending ? (
          <p className="p-5 text-small text-muted">Loading people…</p>
        ) : members.isError ? (
          <p role="alert" className="p-5 text-small text-danger">
            {isApiError(members.error) ? members.error.message : "Could not load people"}
          </p>
        ) : (
          <ul className="divide-y divide-border">
            {members.data.map((member) => (
              <li key={member.user_id} className="flex items-center justify-between gap-4 p-5">
                <div className="flex min-w-0 items-center gap-3">
                  <span
                    aria-hidden
                    className="grid size-9 shrink-0 place-items-center rounded-full bg-raised text-small font-semibold uppercase"
                  >
                    {(member.name || member.email).charAt(0)}
                  </span>

                  <div className="min-w-0">
                    <p className="truncate text-body font-medium">{member.name || member.email}</p>
                    <p className="truncate text-caption text-muted">{member.email}</p>
                  </div>
                </div>

                <div className="flex shrink-0 items-center gap-2">
                  {member.status === "invited" ? <Chip tone="warning">Invited</Chip> : null}

                  <Chip tone="muted" className="capitalize">
                    {member.role}
                  </Chip>

                  {isAdmin && member.role !== "owner" ? (
                    <Button
                      variant="ghost"
                      size="sm"
                      loading={remove.isPending && remove.variables === member.user_id}
                      onClick={() => remove.mutate(member.user_id)}
                    >
                      Remove
                    </Button>
                  ) : null}
                </div>
              </li>
            ))}

            {members.data.length === 0 ? (
              <li className="p-8 text-center text-small text-muted">Just you, for now.</li>
            ) : null}
          </ul>
        )}
      </Section>
    </Page>
  );
}
