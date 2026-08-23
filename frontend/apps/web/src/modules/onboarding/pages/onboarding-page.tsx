import { Button } from "@ava/ui";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import { useCallback, useEffect, useState } from "react";

import { SESSION_QUERY_KEY } from "@/modules/auth";
import { hubQueries } from "@/modules/hub";
import { Loader } from "@/shared/components/loader";
import { Room, ROOM_STAGES, type RoomLight } from "@/shared/components/room";
import { Reveal } from "@/shared/components/reveal";
import { useFlow } from "@/shared/flow";
import { ONBOARDING_FLOW } from "../flow";
import { StepFields } from "../components/step-fields";
import { Eyebrow, StepFrame } from "../components/step-frame";
import { Welcome } from "../components/welcome";
import { toStepPayload } from "../payload";

// How lit the room is at each point in the sequence. Setup reads as the lights
// coming up: dark and cold at the door, warm and full by the time you are in.
const LIGHT_BY_STEP: Record<string, RoomLight> = {
  home: ROOM_STAGES.waking,
  hub: ROOM_STAGES.lit,
};

// Long enough to register as a payoff, short enough not to be in the way. The
// console is already mounting behind it.
const FINALE_MS = 1100;

export function OnboardingPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const hubs = useQuery(hubQueries.list());

  const [begun, setBegun] = useState(false);
  const [finishing, setFinishing] = useState(false);

  // The home step renames the tenant server-side, so the cached session is
  // stale the moment the flow completes.
  const onCompleted = useCallback(() => {
    void queryClient.invalidateQueries({ queryKey: SESSION_QUERY_KEY });
    setFinishing(true);
  }, [queryClient]);

  const flow = useFlow(ONBOARDING_FLOW, toStepPayload, onCompleted);

  // Hold on the fully lit room for a beat before handing over to the console.
  useEffect(() => {
    if (!finishing) return;

    const timer = setTimeout(() => void navigate({ to: "/", replace: true }), FINALE_MS);

    return () => clearTimeout(timer);
  }, [finishing, navigate]);

  if (finishing) {
    return (
      <Room light={ROOM_STAGES.full}>
        <p className="text-center text-balance font-semibold tracking-tighter text-[clamp(2rem,6vw,3rem)] animate-rise">
          Welcome home.
        </p>
      </Room>
    );
  }

  if (flow.isPending) return <Loader label="Getting things ready" />;

  if (flow.isError) {
    return (
      <Room light={ROOM_STAGES.asleep}>
        <div className="grid gap-5">
          <h1 className="text-display font-semibold">Could not start setup</h1>
          <p className="text-lead text-muted">{flow.loadError ?? "Please try again."}</p>
          <div>
            <Button onClick={flow.retry}>Retry</Button>
          </div>
        </div>
      </Room>
    );
  }

  if (!begun) {
    return (
      <Room light={ROOM_STAGES.asleep}>
        <Welcome onBegin={() => setBegun(true)} />
      </Room>
    );
  }

  const step = flow.step;
  const total = flow.flow?.steps.length ?? 0;

  if (!step) return <Loader label="Getting things ready" />;

  // Until a hub exists, submitting the hub step can only come back as an error
  // telling you to skip — so skipping *is* the way forward, and it is the only
  // button shown. Once a hub is paired, Continue takes over.
  const satisfied = step.id !== "hub" || (hubs.data ?? []).length > 0;

  const action = satisfied
    ? { label: "Continue", loading: flow.isAdvancing }
    : {
        label: "Skip for now",
        variant: "ghost" as const,
        loading: flow.isAdvancing,
        onClick: () => flow.advance("skip"),
      };

  return (
    <Room light={LIGHT_BY_STEP[step.id] ?? ROOM_STAGES.waking}>
      {/* Remounting per step restarts the cascade, so each screen arrives
          rather than swapping its text in place. */}
      <StepFrame
        key={step.id}
        eyebrow={<Eyebrow at={flow.at} total={total} />}
        title={step.title}
        description={step.description}
        error={flow.error}
        action={action}
        onSubmit={() => flow.advance("submit")}
      >
        <StepFields step={step} values={flow.values} onChange={flow.setValue} />
      </StepFrame>

      {flow.at > 0 ? (
        <Reveal at={6} className="mt-8">
          <button
            type="button"
            onClick={() => flow.advance("back")}
            className="min-h-11 text-small text-subtle transition-colors duration-150 ease-out hover:text-fg"
          >
            ← Back
          </button>
        </Reveal>
      ) : null}
    </Room>
  );
}
