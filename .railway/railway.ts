import { defineRailway, github, postgres, preserve, project, service, volume } from "railway/iac";

/**
 * The whole Ava project, one environment at a time.
 *
 * Anything omitted here is deleted on apply, so this file has to describe the
 * environment completely. Always read `railway config plan` before applying.
 *
 * Four services: Postgres and the broker hold state, the API and the web server
 * are stateless and rebuild from the repo.
 *
 * The target environment has to be passed in. The CLI evaluates this file with
 * an empty context and no Railway variables in `process.env` — `ctx.environment`
 * is always undefined — so there is no way to detect it here. Guessing would be
 * worse than failing: a wrong guess points an environment at the wrong branch,
 * and apply deletes whatever disagrees.
 *
 *   AVA_ENV=stage railway config plan
 *   AVA_ENV=production railway config apply
 */
const ENVIRONMENTS = {
  stage: { branch: "stage" },
  production: { branch: "main" },
} as const;

/**
 * Railway's private network addresses a service by its own name, so these are
 * stable without a reference.
 */
const API_INTERNAL = "http://api.railway.internal:8000";
const BROKER_INTERNAL = "tcp://mosquitto.railway.internal:1883";

export default defineRailway(() => {
  const name = process.env.AVA_ENV;

  if (name !== "stage" && name !== "production") {
    throw new Error(
      `Set AVA_ENV to "stage" or "production" (got ${JSON.stringify(name)}). ` +
        "It cannot be inferred from the linked environment.",
    );
  }

  const { branch } = ENVIRONMENTS[name];
  const repo = "xenomech/ava";

  const db = postgres("postgres");

  /**
   * The broker the hub talks to.
   *
   * `tcp` asks Railway for a public TCP proxy. A hub runs on someone's Pi at
   * home and reaches this from outside the project, so unlike the API it cannot
   * live on the private network alone.
   *
   * Credentials are `preserve()`: set once in the Railway UI, never written
   * here. The broker bootstraps its dynamic-security file on first boot from
   * MQTT_ADMIN_PASSWORD, so changing it later does nothing until the volume is
   * wiped.
   *
   * One volume, mounted where the broker writes. It cannot go on
   * /mosquitto/config: the image bakes mosquitto.conf in there and a mount
   * would hide it. Both the dynamic-security file and the persistence data are
   * therefore under /mosquitto/data.
   */
  const brokerData = volume("mosquitto-data", { sizeMB: 1024 });

  const broker = service("mosquitto", {
    source: github(repo, { branch }),
    build: {
      builder: "DOCKERFILE",
      dockerfilePath: "docker/mosquitto/Dockerfile",
      watchPatterns: ["docker/mosquitto/**"],
    },
    tcp: [1883],
    volumeMounts: { "mosquitto-data": { mountPath: "/mosquitto/data" } },
    env: {
      MQTT_ADMIN_USERNAME: "ava-api",
      MQTT_ADMIN_PASSWORD: preserve(),
    },
  });

  /**
   * Caddy serving the built SPA, and proxying /api to the API service.
   *
   * The proxy is what keeps the browser on a single origin. Session cookies are
   * SameSite=Lax, and two Railway services sit on different hosts under a public
   * suffix, so split across two origins the cookie would never be sent and every
   * request after sign-in would 401.
   *
   * VITE_API_URL is baked in at build time — the Dockerfile fails without it —
   * and is relative on purpose, so no environment-specific host ends up compiled
   * into the bundle.
   *
   * Declared before the API because the two would otherwise reference each other
   * in a cycle. This direction is the one that can be written by hand.
   */
  const web = service("web", {
    source: github(repo, { branch }),
    build: {
      builder: "DOCKERFILE",
      dockerfilePath: "frontend/Dockerfile",
      watchPatterns: ["frontend/**", ".railway/**"],
    },
    healthcheckPath: "/",
    healthcheckTimeout: 30,
    env: {
      VITE_API_URL: "/api/v1",
      API_UPSTREAM: API_INTERNAL,
    },
  });

  const api = service("api", {
    source: github(repo, { branch }),
    build: {
      builder: "DOCKERFILE",
      dockerfilePath: "backend/services/api/Dockerfile",
      watchPatterns: ["backend/pkg/**", "backend/services/api/**", "go.work", ".railway/**"],
    },
    healthcheckPath: "/api/v1/health",
    healthcheckTimeout: 60,
    env: {
      PORT: "8000",
      SERVER_ENV: name,

      DB_HOST: db.env.PGHOST,
      DB_PORT: db.env.PGPORT,
      DB_USER: db.env.PGUSER,
      DB_PASSWORD: db.env.PGPASSWORD,
      DB_DATABASE: db.env.PGDATABASE,

      JWT_SECRET: preserve(),
      JWT_ACCESS_EXPIRY: "15m",
      JWT_REFRESH_EXPIRY: "168h",

      /* Railway's own interpolation, not a JS template. A `service.env.X`
         reference is an object that only survives as a whole value — dropped
         into a template literal it stringifies to "[object Object]", which is
         exactly what the first deploy of this file shipped. The host is
         assigned by Railway and is not name-derivable, so it has to be resolved
         at deploy time rather than written here. */
      CORS_ALLOWED_ORIGINS: "https://${{web.RAILWAY_PUBLIC_DOMAIN}}",
      APP_URL: "https://${{web.RAILWAY_PUBLIC_DOMAIN}}",
      COOKIE_DOMAIN: "",

      MQTT_BROKER_URL: BROKER_INTERNAL,
      MQTT_USERNAME: "ava-api",
      MQTT_PASSWORD: preserve(),

      RESEND_API_KEY: preserve(),
      RESEND_FROM_EMAIL: preserve(),
      RESEND_FROM_NAME: "Ava",
    },
  });

  return project("ava", { resources: [db, brokerData, broker, web, api] });
});
