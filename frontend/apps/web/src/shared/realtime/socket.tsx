import { env } from "@ava/env/web";
import { createContext, use, useEffect, useMemo, useRef, useState, type ReactNode } from "react";

const RECONNECT_MIN_MS = 1_000;
const RECONNECT_MAX_MS = 15_000;

type Listener = (raw: string) => void;

type AvaSocket = {
  send: (frame: unknown) => boolean;
  subscribe: (listener: Listener) => () => void;
};

/* Two contexts on purpose. send/subscribe only touch refs and never change
   identity, while connected flips on every reconnect — folding them into one
   value made every consumer re-render (and ava-events re-subscribe its
   handlers) each time the socket flapped, for a field none of them read. */
const SocketContext = createContext<AvaSocket>({
  send: () => false,
  subscribe: () => () => undefined,
});

const SocketConnectedContext = createContext(false);

export const useAvaSocket = () => use(SocketContext);
export const useSocketConnected = () => use(SocketConnectedContext);

function socketURL() {
  /* The API base is absolute in development and a same-origin path in a
     deployed build, where the web server proxies /api. `new URL` throws on a
     bare path, so the page itself supplies the base — which is the right origin
     in exactly the case where the path is relative. */
  const url = new URL(`${env.VITE_API_URL}/socket`, window.location.origin);
  url.protocol = url.protocol === "https:" ? "wss:" : "ws:";

  return url.toString();
}

export function AvaSocketProvider({ children }: { children: ReactNode }) {
  const socket = useRef<WebSocket | null>(null);
  const listeners = useRef(new Set<Listener>());
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    let backoff = RECONNECT_MIN_MS;
    let retry: ReturnType<typeof setTimeout> | undefined;
    let stopped = false;

    const open = () => {
      const next = new WebSocket(socketURL());
      socket.current = next;

      next.onopen = () => {
        backoff = RECONNECT_MIN_MS;
        setConnected(true);
      };

      next.onmessage = (message) => {
        for (const listener of listeners.current) listener(message.data as string);
      };

      next.onclose = () => {
        setConnected(false);
        socket.current = null;

        if (stopped) return;

        retry = setTimeout(open, backoff);
        backoff = Math.min(backoff * 2, RECONNECT_MAX_MS);
      };
    };

    open();

    return () => {
      stopped = true;
      clearTimeout(retry);
      socket.current?.close();
      socket.current = null;
    };
  }, []);

  const value = useMemo<AvaSocket>(
    () => ({
      send: (frame) => {
        const live = socket.current;
        if (!live || live.readyState !== WebSocket.OPEN) return false;

        live.send(JSON.stringify(frame));

        return true;
      },
      subscribe: (listener) => {
        listeners.current.add(listener);

        return () => listeners.current.delete(listener);
      },
    }),
    [],
  );

  return (
    <SocketContext value={value}>
      <SocketConnectedContext value={connected}>{children}</SocketConnectedContext>
    </SocketContext>
  );
}
