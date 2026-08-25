import { useEffect } from "react";
import { toast } from "sonner";
import { useRegisterSW } from "virtual:pwa-register/react";

const CHECK_INTERVAL = 60 * 60 * 1000;

export function UpdatePrompt() {
  const {
    needRefresh: [needRefresh],
    updateServiceWorker,
  } = useRegisterSW({
    onRegisteredSW(_url, registration) {
      if (!registration) return;

      const check = () => {
        if (document.visibilityState === "hidden") return;

        void registration.update();
      };

      window.setInterval(check, CHECK_INTERVAL);
      window.addEventListener("online", check);
    },
  });

  useEffect(() => {
    if (!needRefresh) return;

    toast("A new version of Ava is ready", {
      id: "ava-update",
      duration: Infinity,
      dismissible: false,
      action: {
        label: "Reload",
        onClick: () => void updateServiceWorker(true),
      },
    });
  }, [needRefresh, updateServiceWorker]);

  return null;
}
